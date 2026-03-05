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
	"github.com/luisfpires18/woo/internal/handler"
	"github.com/luisfpires18/woo/internal/middleware"
	"github.com/luisfpires18/woo/internal/repository/sqlite"
	"github.com/luisfpires18/woo/internal/service"
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
	worldConfigRepo := sqlite.NewWorldConfigRepo(db)
	announcementRepo := sqlite.NewAnnouncementRepo(db)
	gameAssetRepo := sqlite.NewGameAssetRepo(db)

	// Ensure uploads directory exists
	for _, dir := range []string{"uploads/sprites/building", "uploads/sprites/resource", "uploads/sprites/unit"} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			slog.Error("failed to create uploads directory", "dir", dir, "error", err)
			os.Exit(1)
		}
	}

	// Wire up services
	authService := service.NewAuthService(playerRepo, refreshTokenRepo, cfg.JWTSecret, cfg.JWTIssuer)
	villageService := service.NewVillageService(villageRepo, buildingRepo, resourceRepo)
	adminService := service.NewAdminService(playerRepo, villageRepo, worldConfigRepo, announcementRepo, gameAssetRepo)

	// Wire up handlers
	authHandler := handler.NewAuthHandler(authService, villageService)
	villageHandler := handler.NewVillageHandler(villageService)
	adminHandler := handler.NewAdminHandler(adminService)

	// Auth middleware for protected routes
	authMiddleware := middleware.Auth(authService.ValidateAccessToken)

	// Set up HTTP router
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Auth routes (public)
	authHandler.RegisterRoutes(mux)

	// Protected routes — wrapped with auth middleware
	protectedMux := http.NewServeMux()
	villageHandler.RegisterRoutes(protectedMux)

	// Mount protected routes under the auth middleware
	mux.Handle("/api/villages", authMiddleware(protectedMux))
	mux.Handle("/api/villages/", authMiddleware(protectedMux))

	// Game assets — read is auth-only (all players need icons), write is admin-only
	mux.Handle("GET /api/assets", authMiddleware(http.HandlerFunc(adminHandler.ListAssets)))

	// Admin routes — wrapped with auth + admin middleware
	adminMux := http.NewServeMux()
	adminHandler.RegisterRoutes(adminMux)
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

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutdown error", "error", err)
	}
	slog.Info("server stopped")
}
