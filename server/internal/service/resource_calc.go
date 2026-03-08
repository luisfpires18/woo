package service

import (
	"math"
	"time"

	"github.com/luisfpires18/woo/internal/model"
)

// FlushResources recalculates resource amounts based on the elapsed time since
// the last update. Rates are per-second. The resource snapshot is mutated in
// place and LastUpdated is advanced to now.
//
// Returns false (no-op) when the clock has not advanced since the last flush.
func FlushResources(res *model.Resources, now time.Time) bool {
	elapsed := now.Sub(res.LastUpdated).Seconds()
	if elapsed <= 0 {
		return false
	}

	res.Food = clampResource(res.Food+(res.FoodRate-res.FoodConsumption)*elapsed, res.MaxStorage)
	res.Water = clampResource(res.Water+res.WaterRate*elapsed, res.MaxStorage)
	res.Lumber = clampResource(res.Lumber+res.LumberRate*elapsed, res.MaxStorage)
	res.Stone = clampResource(res.Stone+res.StoneRate*elapsed, res.MaxStorage)
	res.LastUpdated = now

	return true
}

// clampResource caps a resource value between 0 and the storage maximum.
func clampResource(val, max float64) float64 {
	if val < 0 {
		return 0
	}
	return math.Min(val, max)
}
