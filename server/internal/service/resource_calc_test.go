package service

import (
	"testing"
	"time"

	"github.com/luisfpires18/woo/internal/model"
)

func TestFlushResources_AdvancesTime(t *testing.T) {
	res := &model.Resources{
		Food: 100, Water: 100, Lumber: 100, Stone: 100,
		FoodRate: 10, WaterRate: 10, LumberRate: 10, StoneRate: 10,
		FoodConsumption: 2,
		MaxFood: 5000, MaxWater: 5000, MaxLumber: 5000, MaxStone: 5000,
		LastUpdated: time.Now().UTC().Add(-10 * time.Second),
	}
	now := time.Now().UTC()
	changed := FlushResources(res, now)
	if !changed {
		t.Fatal("expected FlushResources to return true")
	}

	// Food: 100 + (10-2)*~10 = ~180
	if res.Food < 170 || res.Food > 190 {
		t.Errorf("food = %f, want ~180", res.Food)
	}
	// Other resources: 100 + 10*~10 = ~200
	if res.Water < 190 || res.Water > 210 {
		t.Errorf("water = %f, want ~200", res.Water)
	}
	if !res.LastUpdated.Equal(now) {
		t.Errorf("LastUpdated not advanced to now")
	}
}

func TestFlushResources_NoOpWhenNotAdvanced(t *testing.T) {
	now := time.Now().UTC()
	res := &model.Resources{
		Food: 100, Water: 100, Lumber: 100, Stone: 100,
		FoodRate: 10, WaterRate: 10, LumberRate: 10, StoneRate: 10,
		MaxFood: 5000, MaxWater: 5000, MaxLumber: 5000, MaxStone: 5000,
		LastUpdated: now,
	}
	changed := FlushResources(res, now)
	if changed {
		t.Fatal("expected FlushResources to return false for zero elapsed time")
	}
	if res.Food != 100 {
		t.Errorf("food should be unchanged, got %f", res.Food)
	}
}

func TestFlushResources_ClampsToMaxStorage(t *testing.T) {
	res := &model.Resources{
		Food: 990, Water: 990, Lumber: 990, Stone: 990,
		FoodRate: 100, WaterRate: 100, LumberRate: 100, StoneRate: 100,
		MaxFood: 1000, MaxWater: 1000, MaxLumber: 1000, MaxStone: 1000,
		LastUpdated: time.Now().UTC().Add(-10 * time.Second),
	}
	FlushResources(res, time.Now().UTC())

	if res.Food != 1000 {
		t.Errorf("food = %f, want 1000 (capped at max)", res.Food)
	}
	if res.Water != 1000 {
		t.Errorf("water = %f, want 1000 (capped at max)", res.Water)
	}
}

func TestFlushResources_FoodFloorsAtZero(t *testing.T) {
	res := &model.Resources{
		Food:            10,
		FoodRate:        1,
		FoodConsumption: 100, // consumption >> rate
		MaxFood: 5000, MaxWater: 5000, MaxLumber: 5000, MaxStone: 5000,
		LastUpdated: time.Now().UTC().Add(-10 * time.Second),
	}
	FlushResources(res, time.Now().UTC())

	if res.Food != 0 {
		t.Errorf("food = %f, want 0 (should floor at zero)", res.Food)
	}
}

func TestClampResource(t *testing.T) {
	tests := []struct {
		val, max, want float64
	}{
		{500, 1000, 500},
		{1500, 1000, 1000},
		{-50, 1000, 0},
		{0, 1000, 0},
	}
	for _, tc := range tests {
		got := clampResource(tc.val, tc.max)
		if got != tc.want {
			t.Errorf("clampResource(%f, %f) = %f, want %f", tc.val, tc.max, got, tc.want)
		}
	}
}
