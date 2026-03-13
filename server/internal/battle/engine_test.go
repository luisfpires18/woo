package battle

import (
	"encoding/json"
	"testing"
)

func TestSimulate_Deterministic(t *testing.T) {
	attackers := []Unit{
		{ID: 1, Name: "Soldier A", HP: 100, MaxHP: 100, AttackPower: 20, AttackInterval: 2, DefensePercent: 10, CritChancePercent: 10},
		{ID: 2, Name: "Soldier B", HP: 80, MaxHP: 80, AttackPower: 30, AttackInterval: 3, DefensePercent: 5, CritChancePercent: 15},
	}
	defenders := []Unit{
		{ID: 3, Name: "Beast A", HP: 150, MaxHP: 150, AttackPower: 15, AttackInterval: 2, DefensePercent: 20, CritChancePercent: 5},
	}
	tuning := DefaultTuning()
	seed := int64(42)

	result1 := Simulate(
		cloneUnits(attackers),
		cloneUnits(defenders),
		tuning,
		seed,
	)
	result2 := Simulate(
		cloneUnits(attackers),
		cloneUnits(defenders),
		tuning,
		seed,
	)

	if result1.Outcome != result2.Outcome {
		t.Fatalf("outcomes differ: %s vs %s", result1.Outcome, result2.Outcome)
	}
	if result1.DurationTicks != result2.DurationTicks {
		t.Fatalf("duration differs: %d vs %d", result1.DurationTicks, result2.DurationTicks)
	}
	if len(result1.Events) != len(result2.Events) {
		t.Fatalf("event count differs: %d vs %d", len(result1.Events), len(result2.Events))
	}
	for i := range result1.Events {
		e1 := result1.Events[i]
		e2 := result2.Events[i]
		if e1.Tick != e2.Tick || e1.SourceID != e2.SourceID || e1.TargetID != e2.TargetID ||
			e1.Damage != e2.Damage || e1.IsCrit != e2.IsCrit || e1.IsKill != e2.IsKill {
			t.Fatalf("event %d differs: %+v vs %+v", i, e1, e2)
		}
	}
}

func TestSimulate_AttackerWins(t *testing.T) {
	attackers := []Unit{
		{ID: 1, Name: "Strong", HP: 200, MaxHP: 200, AttackPower: 50, AttackInterval: 1, CritChancePercent: 0},
	}
	defenders := []Unit{
		{ID: 2, Name: "Weak", HP: 10, MaxHP: 10, AttackPower: 1, AttackInterval: 5, CritChancePercent: 0},
	}
	result := Simulate(attackers, defenders, DefaultTuning(), 1)

	if result.Outcome != "attacker_won" {
		t.Fatalf("expected attacker_won, got %s", result.Outcome)
	}
	if len(result.AttackerSurvivors) != 1 {
		t.Fatalf("expected 1 attacker survivor, got %d", len(result.AttackerSurvivors))
	}
	if len(result.DefenderSurvivors) != 0 {
		t.Fatalf("expected 0 defender survivors, got %d", len(result.DefenderSurvivors))
	}
}

func TestSimulate_DefenderWins(t *testing.T) {
	attackers := []Unit{
		{ID: 1, Name: "Weak", HP: 10, MaxHP: 10, AttackPower: 1, AttackInterval: 5, CritChancePercent: 0},
	}
	defenders := []Unit{
		{ID: 2, Name: "Strong", HP: 200, MaxHP: 200, AttackPower: 50, AttackInterval: 1, CritChancePercent: 0},
	}
	result := Simulate(attackers, defenders, DefaultTuning(), 1)

	if result.Outcome != "defender_won" {
		t.Fatalf("expected defender_won, got %s", result.Outcome)
	}
}

func TestSimulate_MinimumDamage(t *testing.T) {
	// Attacker has 1 attack vs defender with max defense → should still deal 1 dmg
	attackers := []Unit{
		{ID: 1, Name: "Puny", HP: 1000, MaxHP: 1000, AttackPower: 1, AttackInterval: 1, CritChancePercent: 0},
	}
	defenders := []Unit{
		{ID: 2, Name: "Tank", HP: 5, MaxHP: 5, AttackPower: 1, AttackInterval: 100, DefensePercent: 90, CritChancePercent: 0},
	}
	result := Simulate(attackers, defenders, DefaultTuning(), 1)

	if result.Outcome != "attacker_won" {
		t.Fatalf("expected attacker_won, got %s", result.Outcome)
	}
	// With 1 base damage and 90% def: floor(1 * 0.1) = 0 → clamped to 1
	// So 5 HP should take 5 attacks to kill
	killEvents := 0
	for _, e := range result.Events {
		if e.IsKill {
			killEvents++
		}
	}
	if killEvents != 1 {
		t.Fatalf("expected 1 kill event, got %d", killEvents)
	}
}

func TestSimulate_DistributesTargets(t *testing.T) {
	attackers := []Unit{
		{ID: 1, Name: "Attacker", HP: 100, MaxHP: 100, AttackPower: 10, AttackInterval: 1, CritChancePercent: 0},
	}
	defenders := []Unit{
		{ID: 2, Name: "HighHP", HP: 100, MaxHP: 100, AttackPower: 5, AttackInterval: 3, CritChancePercent: 0},
		{ID: 3, Name: "LowHP", HP: 20, MaxHP: 20, AttackPower: 5, AttackInterval: 3, CritChancePercent: 0},
	}
	result := Simulate(attackers, defenders, DefaultTuning(), 1)

	// With one attacker vs two defenders, the attacker gets assigned to
	// the first defender by ID (round-robin with index 0 → ID 2)
	if len(result.Events) == 0 {
		t.Fatal("expected at least one event")
	}
	firstAttack := result.Events[0]
	if firstAttack.SourceID != 1 {
		t.Fatalf("expected first event source 1, got %d", firstAttack.SourceID)
	}
	if firstAttack.TargetID != 2 {
		t.Fatalf("expected first target to be ID 2 (first by ID, round-robin), got %d", firstAttack.TargetID)
	}
}

func TestSimulate_1v1Pairing(t *testing.T) {
	// 3v3 — each attacker should hit a different defender
	attackers := []Unit{
		{ID: 1, Name: "A1", HP: 200, MaxHP: 200, AttackPower: 10, AttackInterval: 1, CritChancePercent: 0},
		{ID: 2, Name: "A2", HP: 200, MaxHP: 200, AttackPower: 10, AttackInterval: 1, CritChancePercent: 0},
		{ID: 3, Name: "A3", HP: 200, MaxHP: 200, AttackPower: 10, AttackInterval: 1, CritChancePercent: 0},
	}
	defenders := []Unit{
		{ID: 4, Name: "D1", HP: 200, MaxHP: 200, AttackPower: 10, AttackInterval: 1, CritChancePercent: 0},
		{ID: 5, Name: "D2", HP: 200, MaxHP: 200, AttackPower: 10, AttackInterval: 1, CritChancePercent: 0},
		{ID: 6, Name: "D3", HP: 200, MaxHP: 200, AttackPower: 10, AttackInterval: 1, CritChancePercent: 0},
	}
	result := Simulate(attackers, defenders, DefaultTuning(), 42)

	// Collect first tick's attacker→defender assignments
	targetsHit := map[int]int{} // sourceID → targetID
	for _, e := range result.Events {
		if e.Type != "attack" || e.Tick != 1 {
			continue
		}
		// Only look at attacker-side hits
		if e.SourceID >= 1 && e.SourceID <= 3 {
			targetsHit[e.SourceID] = e.TargetID
		}
	}
	if len(targetsHit) != 3 {
		t.Fatalf("expected 3 attacker hits on tick 1, got %d", len(targetsHit))
	}

	// All three targets should be unique (1-on-1)
	seen := map[int]bool{}
	for src, tgt := range targetsHit {
		if seen[tgt] {
			t.Fatalf("attacker %d and another attacker both targeted defender %d — not 1-on-1", src, tgt)
		}
		seen[tgt] = true
	}
}

func TestSimulate_MaxTicksDraw(t *testing.T) {
	// Both have 0 attack power (clamped to 1) but very high HP → should hit max ticks
	tiny := Tuning{
		CritDamageMultiplier: 2.0,
		MaxDefensePercent:    90.0,
		MaxCritChancePercent: 75.0,
		MinAttackInterval:    1,
		MaxTicks:             5,
	}
	attackers := []Unit{
		{ID: 1, HP: 99999, MaxHP: 99999, AttackPower: 1, AttackInterval: 1},
	}
	defenders := []Unit{
		{ID: 2, HP: 99999, MaxHP: 99999, AttackPower: 1, AttackInterval: 1},
	}
	result := Simulate(attackers, defenders, tiny, 1)

	if result.Outcome != "draw" {
		t.Fatalf("expected draw with MaxTicks=5, got %s", result.Outcome)
	}
	if result.DurationTicks != 5 {
		t.Fatalf("expected 5 ticks, got %d", result.DurationTicks)
	}
}

func TestSimulate_CritDamage(t *testing.T) {
	// 100% crit chance → every attack is a crit
	attackers := []Unit{
		{ID: 1, Name: "Crit", HP: 100, MaxHP: 100, AttackPower: 10, AttackInterval: 1, CritChancePercent: 100},
	}
	defenders := []Unit{
		{ID: 2, Name: "Target", HP: 50, MaxHP: 50, AttackPower: 1, AttackInterval: 100, CritChancePercent: 0},
	}
	result := Simulate(attackers, defenders, DefaultTuning(), 1)

	for _, e := range result.Events {
		if e.Type == "attack" && e.SourceID == 1 {
			if !e.IsCrit {
				t.Fatalf("expected all attacks to be crits with 100%% crit chance")
			}
			// Damage should be 10 * 2.0 = 20
			if e.Damage != 20 {
				t.Fatalf("expected crit damage 20, got %d", e.Damage)
			}
			break
		}
	}
}

func TestBuildReplayJSON(t *testing.T) {
	attackers := []Unit{
		{ID: 1, Name: "Soldier", HP: 50, MaxHP: 50, AttackPower: 10, AttackInterval: 2},
	}
	defenders := []Unit{
		{ID: 2, Name: "Beast", HP: 30, MaxHP: 30, AttackPower: 8, AttackInterval: 3},
	}
	result := Result{
		Outcome:       "attacker_won",
		DurationTicks: 10,
		Events:        []Event{{Tick: 1, Type: "attack", SourceID: 1, TargetID: 2, Damage: 10}},
	}

	data, err := BuildReplayJSON(attackers, defenders, result, 10)
	if err != nil {
		t.Fatalf("BuildReplayJSON error: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected non-empty JSON")
	}

	// Verify it round-trips
	var replay ReplayData
	if err := json.Unmarshal(data, &replay); err != nil {
		t.Fatalf("failed to unmarshal replay: %v", err)
	}
	if replay.Version != 1 {
		t.Fatalf("expected version 1, got %d", replay.Version)
	}
	if replay.Result != "attacker_won" {
		t.Fatalf("expected attacker_won, got %s", replay.Result)
	}
	if replay.TotalTicks != 10 {
		t.Fatalf("expected 10 total ticks, got %d", replay.TotalTicks)
	}
}

// ── Helpers ──────────────────────────────────────────────────────────────────

func cloneUnits(units []Unit) []Unit {
	c := make([]Unit, len(units))
	copy(c, units)
	return c
}
