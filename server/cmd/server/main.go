package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/luisfpires18/woo/internal/config"
	"github.com/luisfpires18/woo/internal/gameloop"
	"github.com/luisfpires18/woo/internal/handler"
	"github.com/luisfpires18/woo/internal/middleware"
	"github.com/luisfpires18/woo/internal/repository"
	"github.com/luisfpires18/woo/internal/repository/sqlite"
	"github.com/luisfpires18/woo/internal/service"
	wws "github.com/luisfpires18/woo/internal/websocket"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Set up structured logging
	logLevel := slog.LevelInfo
	if cfg.Debug {
		logLevel = slog.LevelDebug
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(logger)

	// Initialize database
	db := sqlite.NewConnection(cfg.DatabasePath)
	defer db.Close()

	// Run migrations
	if err := sqlite.RunMigrations(db, cfg.MigrationsPath); err != nil {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	// Wire up repositories
	playerRepo := sqlite.NewPlayerRepo(db)
	refreshTokenRepo := sqlite.NewRefreshTokenRepo(db)
	villageRepo := sqlite.NewVillageRepo(db)
	buildingRepo := sqlite.NewBuildingRepo(db)
	resourceRepo := sqlite.NewResourceRepo(db)
	buildingQueueRepo := sqlite.NewBuildingQueueRepo(db)
	troopRepo := sqlite.NewTroopRepo(db)
	trainingQueueRepo := sqlite.NewTrainingQueueRepo(db)
	worldConfigRepo := sqlite.NewWorldConfigRepo(db)
	announcementRepo := sqlite.NewAnnouncementRepo(db)
	gameAssetRepo := sqlite.NewGameAssetRepo(db)
	resBuildingConfigRepo := sqlite.NewResourceBuildingConfigRepo(db)
	buildingDisplayConfigRepo := sqlite.NewBuildingDisplayConfigRepo(db)
	worldMapRepo := sqlite.NewWorldMapRepo(db)
	kingdomRelationRepo := sqlite.NewKingdomRelationRepo(db)
	_ = kingdomRelationRepo // used later for diplomacy features

	// Ensure uploads directory exists
	for _, dir := range []string{"uploads/sprites/building", "uploads/sprites/building_display", "uploads/sprites/resource", "uploads/sprites/unit", "uploads/sprites/kingdom_flag", "uploads/sprites/village_marker", "uploads/sprites/zone_tile", "uploads/sprites/terrain_tile"} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			slog.Error("failed to create uploads directory", "dir", dir, "error", err)
			os.Exit(1)
		}
	}

	// Wire up template repository (file-system backed)
	templateRepo, err := repository.NewTemplateRepository("data/templates")
	if err != nil {
		slog.Error("failed to initialise template repository", "error", err)
		os.Exit(1)
	}

	// Wire up unit of work (transactional operations)
	uow := sqlite.NewUnitOfWork(db)

	// Wire up services
	authService := service.NewAuthService(playerRepo, refreshTokenRepo, cfg.JWTSecret, cfg.JWTIssuer)
	mapService := service.NewMapService(worldMapRepo, villageRepo)
	villageService := service.NewVillageService(villageRepo, buildingRepo, resourceRepo, mapService)
	buildingService := service.NewBuildingService(uow, villageRepo, buildingRepo, resourceRepo, buildingQueueRepo, playerRepo)
	trainingService := service.NewTrainingService(uow, villageRepo, buildingRepo, resourceRepo, troopRepo, trainingQueueRepo, playerRepo)
	adminService := service.NewAdminService(playerRepo, villageRepo, worldConfigRepo, announcementRepo, gameAssetRepo, resBuildingConfigRepo, buildingDisplayConfigRepo)
	templateService := service.NewTemplateService(templateRepo, worldMapRepo)

	// Season repository + service
	seasonRepo := sqlite.NewSeasonRepo(db)
	seasonService := service.NewSeasonService(seasonRepo, playerRepo, villageService)

	// Generate world map on startup (idempotent — skips if already done)
	if err := mapService.GenerateMap(context.Background()); err != nil {
		slog.Error("failed to generate world map", "error", err)
		os.Exit(1)
	}

	playerService := service.NewPlayerService(playerRepo, villageService)

	// Wire up handlers
	authHandler := handler.NewAuthHandler(authService, villageService)
	villageHandler := handler.NewVillageHandler(villageService, buildingService, trainingService)
	trainingHandler := handler.NewTrainingHandler(trainingService)
	mapHandler := handler.NewMapHandler(mapService)
	adminHandler := handler.NewAdminHandler(adminService, mapService)
	templateHandler := handler.NewTemplateHandler(templateService)
	playerHandler := handler.NewPlayerHandler(playerService, seasonService)
	seasonHandler := handler.NewSeasonHandler(seasonService)

	// Auth middleware for protected routes
	authMiddleware := middleware.Auth(authService.ValidateAccessToken)
	optionalAuthMiddleware := middleware.OptionalAuth(authService.ValidateAccessToken)

	// Set up HTTP router
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Auth routes (public)
	authHandler.RegisterRoutes(mux)

	// Public seasons list (no auth required, optional auth for context)
	publicSeasonMux := http.NewServeMux()
	seasonHandler.RegisterPlayerRoutes(publicSeasonMux)
	publicSeasonHandler := optionalAuthMiddleware(http.StripPrefix("/api/public/seasons", publicSeasonMux))
	mux.Handle("GET /api/public/seasons", optionalAuthMiddleware(http.HandlerFunc(seasonHandler.ListSeasons)))
	mux.Handle("/api/public/seasons/", publicSeasonHandler)

	// Protected routes — wrapped with auth middleware
	protectedMux := http.NewServeMux()
	villageHandler.RegisterRoutes(protectedMux)
	trainingHandler.RegisterRoutes(protectedMux)

	// Mount protected routes under the auth middleware
	mux.Handle("/api/villages", authMiddleware(protectedMux))
	mux.Handle("/api/villages/", authMiddleware(protectedMux))

	// Map routes (protected)
	mapMux := http.NewServeMux()
	mapHandler.RegisterRoutes(mapMux)
	mux.Handle("/api/map", authMiddleware(mapMux))
	mux.Handle("/api/map/", authMiddleware(mapMux))

	// Player routes (protected)
	mux.Handle("GET /api/player/me", authMiddleware(http.HandlerFunc(playerHandler.GetMe)))
	mux.Handle("GET /api/player/profile", authMiddleware(http.HandlerFunc(playerHandler.GetProfile)))
	mux.Handle("PUT /api/player/kingdom", authMiddleware(http.HandlerFunc(playerHandler.ChooseKingdom)))

	// Season routes (protected — player-facing)
	seasonMux := http.NewServeMux()
	seasonHandler.RegisterPlayerRoutes(seasonMux)
	mux.Handle("/api/seasons", authMiddleware(http.StripPrefix("/api/seasons", seasonMux)))
	mux.Handle("/api/seasons/", authMiddleware(http.StripPrefix("/api/seasons", seasonMux)))

	// Game assets — read is auth-only (all players need icons), write is admin-only
	mux.Handle("GET /api/assets", authMiddleware(http.HandlerFunc(adminHandler.ListAssets)))

	// Resource building configs — read is auth-only (players need display names per kingdom)
	mux.Handle("GET /api/resource-building-configs", authMiddleware(http.HandlerFunc(adminHandler.ListResourceBuildingConfigs)))

	// Building display configs — read is auth-only (players need display names per kingdom)
	mux.Handle("GET /api/building-display-configs", authMiddleware(http.HandlerFunc(adminHandler.ListBuildingDisplayConfigs)))

	// Admin routes — wrapped with auth + admin middleware
	adminMux := http.NewServeMux()
	adminHandler.RegisterRoutes(adminMux)
	templateHandler.RegisterRoutes(adminMux)
	trainingHandler.RegisterAdminRoutes(adminMux)
	seasonHandler.RegisterAdminRoutes(adminMux)
	mux.Handle("/api/admin/", authMiddleware(middleware.RequireAdmin(http.StripPrefix("/api/admin", adminMux))))

	// Serve uploaded sprites with caching
	fileServer := http.FileServer(http.Dir("uploads"))
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=86400")
		fileServer.ServeHTTP(w, r)
	})))

	// Apply middleware stack
	handler := middleware.Chain(
		mux,
		middleware.Logging(logger),
		middleware.CORS(cfg.CORSOrigin),
		middleware.RateLimit(30),
	)

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start game loop (processes building/training completions every second)
	gl := gameloop.New(buildingService, trainingService, 1*time.Second)

	// Start WebSocket hub
	wsHub := wws.NewHub()
	go wsHub.Run()
	wsHandler := wws.NewHandler(wsHub, authService.ValidateAccessToken)

	// Connect game loop → WebSocket hub for build/train completion notifications
	gl.SetNotifier(wsHub)
	gl.SetTrainNotifier(wsHub)
	gl.Start()

	// WebSocket endpoint (public — auth is via token query param)
	mux.Handle("/ws", wsHandler)

	// Graceful shutdown context
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Start server
	go func() {
		slog.Info("server starting", "port", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	slog.Info("shutting down...")

	gl.Stop()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutdown error", "error", err)
	}
	slog.Info("server stopped")
}
