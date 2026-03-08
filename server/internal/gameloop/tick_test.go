package gameloop

import (
	"context"
	"testing"
	"time"

	"github.com/luisfpires18/woo/internal/service"
)

// stubBuildCompleter records whether CompleteBuilds was called.
type stubBuildCompleter struct {
	called int
	events []service.BuildCompletionEvent
	err    error
}

func (s *stubBuildCompleter) CompleteBuilds(_ context.Context) ([]service.BuildCompletionEvent, error) {
	s.called++
	return s.events, s.err
}

// stubTrainCompleter records whether CompleteTraining was called.
type stubTrainCompleter struct {
	called int
	events []service.TrainCompletionEvent
	err    error
}

func (s *stubTrainCompleter) CompleteTraining(_ context.Context) ([]service.TrainCompletionEvent, error) {
	s.called++
	return s.events, s.err
}

// stubNotifier records build and train completion broadcasts.
type stubNotifier struct {
	buildCalls int
	trainCalls int
}

func (s *stubNotifier) BroadcastBuildComplete(_, _ int64, _ string, _ int) {
	s.buildCalls++
}

func (s *stubNotifier) BroadcastTrainComplete(_, _ int64, _ string, _ int) {
	s.trainCalls++
}

func TestGameLoop_StartStop(t *testing.T) {
	bc := &stubBuildCompleter{}
	tc := &stubTrainCompleter{}
	gl := New(bc, tc, 50*time.Millisecond)
	gl.Start()

	// Let at least one tick run
	time.Sleep(120 * time.Millisecond)
	gl.Stop()

	if bc.called == 0 {
		t.Error("expected at least one call to CompleteBuilds")
	}
	if tc.called == 0 {
		t.Error("expected at least one call to CompleteTraining")
	}
}

func TestGameLoop_TickNotifiesBuild(t *testing.T) {
	bc := &stubBuildCompleter{
		events: []service.BuildCompletionEvent{
			{PlayerID: 1, VillageID: 1, BuildingType: "barracks", NewLevel: 2},
		},
	}
	tc := &stubTrainCompleter{}
	n := &stubNotifier{}

	gl := New(bc, tc, 50*time.Millisecond)
	gl.SetNotifier(n)
	gl.SetTrainNotifier(n)

	// Manually tick to test synchronously
	gl.tick(context.Background())

	if n.buildCalls != 1 {
		t.Errorf("buildCalls = %d, want 1", n.buildCalls)
	}
}

func TestGameLoop_TickNotifiesTrain(t *testing.T) {
	bc := &stubBuildCompleter{}
	tc := &stubTrainCompleter{
		events: []service.TrainCompletionEvent{
			{PlayerID: 1, VillageID: 1, TroopType: "gladiator", NewTotal: 5},
		},
	}
	n := &stubNotifier{}

	gl := New(bc, tc, 50*time.Millisecond)
	gl.SetNotifier(n)
	gl.SetTrainNotifier(n)

	gl.tick(context.Background())

	if n.trainCalls != 1 {
		t.Errorf("trainCalls = %d, want 1", n.trainCalls)
	}
}
