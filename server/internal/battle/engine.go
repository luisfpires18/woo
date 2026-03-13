package battle

import (
	"encoding/json"
	"math"
	"math/rand"
	"sort"
)

// ── Combat Unit ──────────────────────────────────────────────────────────────

// Unit represents a single combatant in the battle simulation.
type Unit struct {
	ID                int     `json:"id"`
	Side              string  `json:"side"` // "attacker" or "defender"
	Name              string  `json:"name"`
	SpriteKey         string  `json:"sprite_key"`
	HP                int     `json:"hp"`
	MaxHP             int     `json:"max_hp"`
	AttackPower       int     `json:"attack_power"`
	AttackInterval    int     `json:"attack_interval"`
	DefensePercent    float64 `json:"defense_percent"`
	CritChancePercent float64 `json:"crit_chance_percent"`
	Cooldown          int     `json:"-"` // ticks until next attack
}

// ── Battle Event ─────────────────────────────────────────────────────────────

// Event represents a single action during the battle replay.
type Event struct {
	Tick          int    `json:"tick"`
	Type          string `json:"type"` // "attack", "death"
	SourceID      int    `json:"source_id"`
	TargetID      int    `json:"target_id"`
	Damage        int    `json:"damage,omitempty"`
	IsCrit        bool   `json:"is_crit,omitempty"`
	TargetHPAfter int    `json:"target_hp_after,omitempty"`
	IsKill        bool   `json:"is_kill,omitempty"`
}

// ── Tuning Parameters ────────────────────────────────────────────────────────

// Tuning holds the admin-configurable battle parameters.
type Tuning struct {
	CritDamageMultiplier float64
	MaxDefensePercent    float64
	MaxCritChancePercent float64
	MinAttackInterval    int
	MaxTicks             int
}

// DefaultTuning returns sensible defaults for battle tuning.
func DefaultTuning() Tuning {
	return Tuning{
		CritDamageMultiplier: 2.0,
		MaxDefensePercent:    90.0,
		MaxCritChancePercent: 75.0,
		MinAttackInterval:    1,
		MaxTicks:             10000,
	}
}

// ── Battle Result ────────────────────────────────────────────────────────────

// Result holds the full outcome of a battle simulation.
type Result struct {
	Outcome           string  `json:"outcome"` // "attacker_won", "defender_won", "draw"
	DurationTicks     int     `json:"duration_ticks"`
	Events            []Event `json:"events"`
	AttackerSurvivors []Unit  `json:"attacker_survivors"`
	DefenderSurvivors []Unit  `json:"defender_survivors"`
}

// ── Replay Envelope ──────────────────────────────────────────────────────────

// ReplayData is the top-level structure serialized as the replay blob.
type ReplayData struct {
	Version    int     `json:"version"`
	TickRateMs int     `json:"tick_rate_ms"`
	Attackers  []Unit  `json:"attackers"`
	Defenders  []Unit  `json:"defenders"`
	Events     []Event `json:"events"`
	Result     string  `json:"result"`
	TotalTicks int     `json:"total_ticks"`
}

// ── Simulation ───────────────────────────────────────────────────────────────

// Simulate runs a deterministic tick-based battle.
// Given the same attackers, defenders, tuning, and seed,
// it always produces the same result.
func Simulate(attackers []Unit, defenders []Unit, tuning Tuning, seed int64) Result {
	rng := rand.New(rand.NewSource(seed))

	// Clamp stats and initialize cooldowns
	allUnits := make([]*Unit, 0, len(attackers)+len(defenders))
	for i := range attackers {
		u := &attackers[i]
		u.Side = "attacker"
		clampUnit(u, tuning)
		u.Cooldown = u.AttackInterval // first attack after one full interval
		allUnits = append(allUnits, u)
	}
	for i := range defenders {
		u := &defenders[i]
		u.Side = "defender"
		clampUnit(u, tuning)
		u.Cooldown = u.AttackInterval
		allUnits = append(allUnits, u)
	}

	var events []Event
	tick := 0

	for tick < tuning.MaxTicks {
		tick++

		// Collect living units sorted by ID for deterministic ordering
		living := getLiving(allUnits)
		if len(living) == 0 {
			break
		}

		// Check if one side is fully eliminated before processing
		attackersAlive := filterBySide(living, "attacker")
		defendersAlive := filterBySide(living, "defender")
		if len(attackersAlive) == 0 || len(defendersAlive) == 0 {
			break
		}

		// Process each living unit in deterministic order (ascending ID)
		sort.Slice(living, func(i, j int) bool { return living[i].ID < living[j].ID })

		// Decrement cooldowns and collect units ready to attack, grouped by side
		var readyAttackers, readyDefenders []*Unit
		for _, unit := range living {
			unit.Cooldown--
			if unit.Cooldown > 0 {
				continue
			}
			if unit.Side == "attacker" {
				readyAttackers = append(readyAttackers, unit)
			} else {
				readyDefenders = append(readyDefenders, unit)
			}
		}

		// Assign targets round-robin: spread attacks across enemies 1-on-1
		assignments := assignTargets(readyAttackers, defendersAlive)
		assignments = append(assignments, assignTargets(readyDefenders, attackersAlive)...)

		// Sort assignments by attacker ID for deterministic event order
		sort.Slice(assignments, func(i, j int) bool {
			return assignments[i].source.ID < assignments[j].source.ID
		})

		// Execute all attacks for this tick
		for _, a := range assignments {
			if a.source.HP <= 0 {
				continue // killed earlier this tick
			}
			target := a.target
			if target.HP <= 0 {
				// Original target died; find another enemy if available
				var enemies []*Unit
				if a.source.Side == "attacker" {
					enemies = filterBySide(getLiving(allUnits), "defender")
				} else {
					enemies = filterBySide(getLiving(allUnits), "attacker")
				}
				if len(enemies) == 0 {
					a.source.Cooldown = a.source.AttackInterval
					continue
				}
				target = enemies[0] // pick first available
			}

			// Roll crit
			isCrit := rng.Float64() < (a.source.CritChancePercent / 100.0)

			// Calculate damage: crit → then defense
			baseDamage := float64(a.source.AttackPower)
			if isCrit {
				baseDamage *= tuning.CritDamageMultiplier
			}
			defReduction := clampFloat(target.DefensePercent, 0, tuning.MaxDefensePercent) / 100.0
			finalDamage := int(math.Floor(baseDamage * (1.0 - defReduction)))
			if finalDamage < 1 {
				finalDamage = 1
			}

			// Apply damage
			target.HP -= finalDamage
			isKill := target.HP <= 0
			if target.HP < 0 {
				target.HP = 0
			}

			events = append(events, Event{
				Tick:          tick,
				Type:          "attack",
				SourceID:      a.source.ID,
				TargetID:      target.ID,
				Damage:        finalDamage,
				IsCrit:        isCrit,
				TargetHPAfter: target.HP,
				IsKill:        isKill,
			})

			// Reset cooldown
			a.source.Cooldown = a.source.AttackInterval
		}
	}

	// Determine outcome
	attackerSurvivors := getSurvivors(allUnits, "attacker")
	defenderSurvivors := getSurvivors(allUnits, "defender")

	outcome := "draw"
	if len(attackerSurvivors) > 0 && len(defenderSurvivors) == 0 {
		outcome = "attacker_won"
	} else if len(defenderSurvivors) > 0 && len(attackerSurvivors) == 0 {
		outcome = "defender_won"
	}

	return Result{
		Outcome:           outcome,
		DurationTicks:     tick,
		Events:            events,
		AttackerSurvivors: toValues(attackerSurvivors),
		DefenderSurvivors: toValues(defenderSurvivors),
	}
}

// BuildReplayJSON creates the replay payload as JSON bytes.
func BuildReplayJSON(attackers []Unit, defenders []Unit, result Result, tickRateMs int) ([]byte, error) {
	replay := ReplayData{
		Version:    1,
		TickRateMs: tickRateMs,
		Attackers:  attackers,
		Defenders:  defenders,
		Events:     result.Events,
		Result:     result.Outcome,
		TotalTicks: result.DurationTicks,
	}
	return json.Marshal(replay)
}

// ── Internal helpers ─────────────────────────────────────────────────────────

func clampUnit(u *Unit, t Tuning) {
	if u.AttackInterval < t.MinAttackInterval {
		u.AttackInterval = t.MinAttackInterval
	}
	u.DefensePercent = clampFloat(u.DefensePercent, 0, t.MaxDefensePercent)
	u.CritChancePercent = clampFloat(u.CritChancePercent, 0, t.MaxCritChancePercent)
	if u.HP < 1 {
		u.HP = 1
	}
	if u.MaxHP < u.HP {
		u.MaxHP = u.HP
	}
	if u.AttackPower < 1 {
		u.AttackPower = 1
	}
}

func clampFloat(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func getLiving(units []*Unit) []*Unit {
	var alive []*Unit
	for _, u := range units {
		if u.HP > 0 {
			alive = append(alive, u)
		}
	}
	return alive
}

func filterBySide(units []*Unit, side string) []*Unit {
	var result []*Unit
	for _, u := range units {
		if u.Side == side {
			result = append(result, u)
		}
	}
	return result
}

// assignment pairs a source unit with its assigned target for one tick.
type assignment struct {
	source *Unit
	target *Unit
}

// assignTargets distributes ready attackers across living enemies round-robin.
// In a 5v5, each attacker gets a unique defender. In a 5v3, two defenders
// will each face two attackers while one faces one.
func assignTargets(ready []*Unit, enemies []*Unit) []assignment {
	if len(ready) == 0 || len(enemies) == 0 {
		return nil
	}
	// Sort enemies by ID for deterministic assignment
	sorted := make([]*Unit, len(enemies))
	copy(sorted, enemies)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].ID < sorted[j].ID })

	assignments := make([]assignment, len(ready))
	for i, u := range ready {
		assignments[i] = assignment{source: u, target: sorted[i%len(sorted)]}
	}
	return assignments
}

func getSurvivors(units []*Unit, side string) []*Unit {
	var survivors []*Unit
	for _, u := range units {
		if u.Side == side && u.HP > 0 {
			survivors = append(survivors, u)
		}
	}
	return survivors
}

func toValues(ptrs []*Unit) []Unit {
	result := make([]Unit, len(ptrs))
	for i, p := range ptrs {
		result[i] = *p
	}
	return result
}
