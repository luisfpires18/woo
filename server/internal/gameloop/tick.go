package gameloop

import (
	"context"
	"log/slog"
	"time"
)

// BuildCompleter is the interface the game loop needs to complete builds.
type BuildCompleter interface {
	CompleteBuilds(ctx context.Context) error
}

// GameLoop runs periodic game ticks to process building completions and other time-based events.
type GameLoop struct {
	buildingService BuildCompleter
	interval        time.Duration
	cancel          context.CancelFunc
	done            chan struct{}
}

// New creates a new GameLoop with the given tick interval.
func New(buildingService BuildCompleter, interval time.Duration) *GameLoop {
	return &GameLoop{
		buildingService: buildingService,
		interval:        interval,
		done:            make(chan struct{}),
	}
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
	if err := gl.buildingService.CompleteBuilds(ctx); err != nil {
		slog.Error("game tick: complete builds failed", "error", err)
	}
}
