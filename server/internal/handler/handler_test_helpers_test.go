package handler_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/luisfpires18/woo/internal/handler"
	"github.com/luisfpires18/woo/internal/middleware"
	"github.com/luisfpires18/woo/internal/repository"
	"github.com/luisfpires18/woo/internal/repository/sqlite"
	"github.com/luisfpires18/woo/internal/service"
	"github.com/luisfpires18/woo/internal/testutil"
)

// testEnv bundles all repos and services needed for handler tests.
type testEnv struct {
	DB              *sql.DB
	AuthService     *service.AuthService
	VillageService  *service.VillageService
	BuildingService *service.BuildingService
	TrainingService *service.TrainingService
	AdminService    *service.AdminService
	MapService      *service.MapService
	TemplateService *service.TemplateService
	SeasonService   *service.SeasonService
	AuthHandler     *handler.AuthHandler
	VillageHandler  *handler.VillageHandler
	TrainingHandler *handler.TrainingHandler
	PlayerHandler   *handler.PlayerHandler
	AdminHandler    *handler.AdminHandler
	TemplateHandler *handler.TemplateHandler
	SeasonHandler   *handler.SeasonHandler
}

// newTestEnv creates a full test environment with an in-memory DB and all services wired.
func newTestEnv(t *testing.T) *testEnv {
	t.Helper()
	db := testutil.NewTestDB(t)

	playerRepo := sqlite.NewPlayerRepo(db)
	villageRepo := sqlite.NewVillageRepo(db)
	buildingRepo := sqlite.NewBuildingRepo(db)
	resourceRepo := sqlite.NewResourceRepo(db)
	refreshTokenRepo := sqlite.NewRefreshTokenRepo(db)
	worldMapRepo := sqlite.NewWorldMapRepo(db)
	queueRepo := sqlite.NewBuildingQueueRepo(db)
	announcementRepo := sqlite.NewAnnouncementRepo(db)
	gameAssetRepo := sqlite.NewGameAssetRepo(db)
	resBuildingConfigRepo := sqlite.NewResourceBuildingConfigRepo(db)

	uow := sqlite.NewUnitOfWork(db)

	authService := service.NewAuthService(playerRepo, refreshTokenRepo, "test-secret", "woo-test")
	mapService := service.NewMapService(worldMapRepo, villageRepo)
	villageService := service.NewVillageService(villageRepo, buildingRepo, resourceRepo, mapService)
	buildingService := service.NewBuildingService(uow, villageRepo, buildingRepo, resourceRepo, queueRepo, playerRepo)
	troopRepo := sqlite.NewTroopRepo(db)
	trainingQueueRepo := sqlite.NewTrainingQueueRepo(db)
	trainingService := service.NewTrainingService(uow, villageRepo, buildingRepo, resourceRepo, troopRepo, trainingQueueRepo, playerRepo)
	adminService := service.NewAdminService(playerRepo, villageRepo, announcementRepo, gameAssetRepo, resBuildingConfigRepo, sqlite.NewBuildingDisplayConfigRepo(db), sqlite.NewTroopDisplayConfigRepo(db))

	// Template repo (file-based, uses temp dir)
	templateRepo, err := repository.NewTemplateRepository(t.TempDir())
	if err != nil {
		t.Fatalf("NewTemplateRepository: %v", err)
	}
	templateService := service.NewTemplateService(templateRepo, worldMapRepo)

	// Generate a world map so village creation works
	if err := mapService.GenerateMap(context.Background()); err != nil {
		t.Fatalf("GenerateMap: %v", err)
	}

	playerService := service.NewPlayerService(playerRepo, villageService)
	seasonRepo := sqlite.NewSeasonRepo(db)
	seasonService := service.NewSeasonService(seasonRepo, playerRepo, villageService)

	return &testEnv{
		DB:              db,
		AuthService:     authService,
		VillageService:  villageService,
		BuildingService: buildingService,
		TrainingService: trainingService,
		AdminService:    adminService,
		MapService:      mapService,
		TemplateService: templateService,
		SeasonService:   seasonService,
		AuthHandler:     handler.NewAuthHandler(authService, villageService),
		VillageHandler:  handler.NewVillageHandler(villageService, buildingService, trainingService),
		TrainingHandler: handler.NewTrainingHandler(trainingService),
		PlayerHandler:   handler.NewPlayerHandler(playerService, seasonService),
		AdminHandler:    handler.NewAdminHandler(adminService, mapService),
		TemplateHandler: handler.NewTemplateHandler(templateService),
		SeasonHandler:   handler.NewSeasonHandler(seasonService),
	}
}

// authCtx returns a context with the given player ID and role injected.
func authCtx(playerID int64, role string) context.Context {
	return middleware.NewPlayerContext(context.Background(), playerID, role)
}

// decodeEnvelope decodes the standard apiResponse envelope from a test response.
func decodeEnvelope(t *testing.T, rec *httptest.ResponseRecorder) (json.RawMessage, string) {
	t.Helper()
	var env struct {
		Data  json.RawMessage `json:"data,omitempty"`
		Error string          `json:"error,omitempty"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&env); err != nil {
		t.Fatalf("decode envelope: %v", err)
	}
	return env.Data, env.Error
}
