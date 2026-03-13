package gameloop

import (
	"context"
	"log/slog"
	"time"

	"github.com/luisfpires18/woo/internal/service"
)

// BuildCompleter is the interface the game loop needs to complete builds.
type BuildCompleter interface {
	CompleteBuilds(ctx context.Context) ([]service.BuildCompletionEvent, error)
}

// TrainCompleter is the interface the game loop needs to complete training.
type TrainCompleter interface {
	CompleteTraining(ctx context.Context) ([]service.TrainCompletionEvent, error)
}

// CampSpawner is the interface the game loop needs to spawn/despawn camps.
type CampSpawner interface {
	SpawnCamps(ctx context.Context) (int, error)
	DespawnExpiredCamps(ctx context.Context) (int, error)
}

// ExpeditionResolver is the interface the game loop needs to resolve expeditions.
type ExpeditionResolver interface {
	ResolveArrivedExpeditions(ctx context.Context) ([]service.ExpeditionCompletionEvent, error)
	ReturnCompletedExpeditions(ctx context.Context) ([]service.ExpeditionReturnEvent, error)
}

// BuildCompletionNotifier receives notifications when buildings complete.
// If nil, no notifications are sent.
type BuildCompletionNotifier interface {
	BroadcastBuildComplete(playerID, villageID int64, buildingType string, newLevel int)
}

// TrainCompletionNotifier receives notifications when troop training completes.
type TrainCompletionNotifier interface {
	BroadcastTrainComplete(playerID, villageID int64, troopType string, newTotal int)
}

// ExpeditionNotifier receives notifications about expedition events.
type ExpeditionNotifier interface {
	BroadcastExpeditionComplete(playerID, villageID, expeditionID, campID int64, result string)
	BroadcastExpeditionReturn(playerID, villageID, expeditionID int64)
}

// GameLoop runs periodic game ticks to process building completions and other time-based events.
type GameLoop struct {
	buildingService    BuildCompleter
	trainingService    TrainCompleter
	campService        CampSpawner
	expeditionService  ExpeditionResolver
	notifier           BuildCompletionNotifier
	trainNotifier      TrainCompletionNotifier
	expeditionNotifier ExpeditionNotifier
	interval           time.Duration
	campTickCounter    int
	cancel             context.CancelFunc
	done               chan struct{}
}

// New creates a new GameLoop with the given tick interval.
func New(buildingService BuildCompleter, trainingService TrainCompleter, interval time.Duration) *GameLoop {
	return &GameLoop{
		buildingService: buildingService,
		trainingService: trainingService,
		interval:        interval,
		done:            make(chan struct{}),
	}
}

// SetNotifier sets the WebSocket notifier for build completion events.
func (gl *GameLoop) SetNotifier(n BuildCompletionNotifier) {
	gl.notifier = n
}

// SetTrainNotifier sets the WebSocket notifier for train completion events.
func (gl *GameLoop) SetTrainNotifier(n TrainCompletionNotifier) {
	gl.trainNotifier = n
}

// SetCampService sets the camp spawner for camp spawn/despawn ticks.
func (gl *GameLoop) SetCampService(s CampSpawner) {
	gl.campService = s
}

// SetExpeditionService sets the expedition resolver for expedition ticks.
func (gl *GameLoop) SetExpeditionService(s ExpeditionResolver) {
	gl.expeditionService = s
}

// SetExpeditionNotifier sets the WebSocket notifier for expedition events.
func (gl *GameLoop) SetExpeditionNotifier(n ExpeditionNotifier) {
	gl.expeditionNotifier = n
}

// Start begins the game loop in a background goroutine.
func (gl *GameLoop) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	gl.cancel = cancel

	go gl.run(ctx)
	slog.Info("game loop started", "interval", gl.interval.String())
}

// Stop gracefully stops the game loop and waits for it to finish.
func (gl *GameLoop) Stop() {
	if gl.cancel != nil {
		gl.cancel()
		<-gl.done
		slog.Info("game loop stopped")
	}
}

func (gl *GameLoop) run(ctx context.Context) {
	defer close(gl.done)

	ticker := time.NewTicker(gl.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			gl.tick(ctx)
		}
	}
}

func (gl *GameLoop) tick(ctx context.Context) {
	// Process building completions.
	events, err := gl.buildingService.CompleteBuilds(ctx)
	if err != nil {
		slog.Error("game tick: complete builds failed", "error", err)
	} else if gl.notifier != nil {
		for _, ev := range events {
			gl.notifier.BroadcastBuildComplete(ev.PlayerID, ev.VillageID, ev.BuildingType, ev.NewLevel)
		}
	}

	// Process training completions.
	if gl.trainingService != nil {
		trainEvents, terr := gl.trainingService.CompleteTraining(ctx)
		if terr != nil {
			slog.Error("game tick: complete training failed", "error", terr)
		} else if gl.trainNotifier != nil {
			for _, ev := range trainEvents {
				gl.trainNotifier.BroadcastTrainComplete(ev.PlayerID, ev.VillageID, ev.TroopType, ev.NewTotal)
			}
		}
	}

	// Process arrived expeditions (resolve battles).
	if gl.expeditionService != nil {
		expEvents, eerr := gl.expeditionService.ResolveArrivedExpeditions(ctx)
		if eerr != nil {
			slog.Error("game tick: resolve expeditions failed", "error", eerr)
		} else if gl.expeditionNotifier != nil {
			for _, ev := range expEvents {
				gl.expeditionNotifier.BroadcastExpeditionComplete(ev.PlayerID, ev.VillageID, ev.ExpeditionID, ev.CampID, ev.Result)
			}
		}

		// Process returning expeditions (troops arrive home).
		retEvents, rerr := gl.expeditionService.ReturnCompletedExpeditions(ctx)
		if rerr != nil {
			slog.Error("game tick: return expeditions failed", "error", rerr)
		} else if gl.expeditionNotifier != nil {
			for _, ev := range retEvents {
				gl.expeditionNotifier.BroadcastExpeditionReturn(ev.PlayerID, ev.VillageID, ev.ExpeditionID)
			}
		}
	}

	// Camp spawn/despawn (every 30 ticks to reduce overhead).
	if gl.campService != nil {
		gl.campTickCounter++
		if gl.campTickCounter >= 30 {
			gl.campTickCounter = 0
			if spawned, err := gl.campService.SpawnCamps(ctx); err != nil {
				slog.Error("game tick: camp spawn failed", "error", err)
			} else if spawned > 0 {
				slog.Info("game tick: spawned camps", "count", spawned)
			}
			if despawned, err := gl.campService.DespawnExpiredCamps(ctx); err != nil {
				slog.Error("game tick: camp despawn failed", "error", err)
			} else if despawned > 0 {
				slog.Info("game tick: despawned camps", "count", despawned)
			}
		}
	}
}
