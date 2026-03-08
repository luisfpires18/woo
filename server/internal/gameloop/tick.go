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

// BuildCompletionNotifier receives notifications when buildings complete.
// If nil, no notifications are sent.
type BuildCompletionNotifier interface {
	BroadcastBuildComplete(playerID, villageID int64, buildingType string, newLevel int)
}

// TrainCompletionNotifier receives notifications when troop training completes.
type TrainCompletionNotifier interface {
	BroadcastTrainComplete(playerID, villageID int64, troopType string, newTotal int)
}

// GameLoop runs periodic game ticks to process building completions and other time-based events.
type GameLoop struct {
	buildingService BuildCompleter
	trainingService TrainCompleter
	notifier        BuildCompletionNotifier
	trainNotifier   TrainCompletionNotifier
	interval        time.Duration
	cancel          context.CancelFunc
	done            chan struct{}
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
}
