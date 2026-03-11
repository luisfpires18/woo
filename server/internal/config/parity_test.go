package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// findRepoRoot walks up from the test file's directory until it finds go.mod,
// which marks the server/ root. The repo root is one level above that.
func findRepoRoot(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("unable to determine test file path")
	}
	dir := filepath.Dir(filename)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			// dir is server/, repo root is parent
			return filepath.Dir(dir)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find go.mod walking up from test file")
		}
		dir = parent
	}
}

func loadJSONFile[T any](t *testing.T, path string) map[string]T {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read %s: %v", path, err)
	}
	var result map[string]T
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal %s: %v", path, err)
	}
	return result
}

func TestBuildingsParity(t *testing.T) {
	repoRoot := findRepoRoot(t)
	jsonPath := filepath.Join(repoRoot, "client", "src", "config", "generated", "buildings.json")
	jsonBuildings := loadJSONFile[GeneratedBuildingConfig](t, jsonPath)

	// Key set must match exactly
	goKeys := SortedBuildingKeys()
	if len(jsonBuildings) != len(goKeys) {
		t.Fatalf("key count mismatch: Go has %d buildings, JSON has %d", len(goKeys), len(jsonBuildings))
	}
	for _, key := range goKeys {
		if _, ok := jsonBuildings[key]; !ok {
			t.Errorf("building %q exists in Go but not in JSON", key)
		}
	}
	for key := range jsonBuildings {
		if _, ok := BuildingConfigs[key]; !ok {
			t.Errorf("building %q exists in JSON but not in Go", key)
		}
	}

	// Field-level comparison
	for _, key := range goKeys {
		goConfig := ToGeneratedBuilding(BuildingConfigs[key])
		jsonConfig := jsonBuildings[key]

		// Compare prerequisites order-independently
		goPrereqs := goConfig.Prerequisites
		jsonPrereqs := jsonConfig.Prerequisites
		goConfig.Prerequisites = nil
		jsonConfig.Prerequisites = nil

		if !reflect.DeepEqual(goConfig, jsonConfig) {
			t.Errorf("building %q mismatch (excluding prerequisites):\n  Go:   %+v\n  JSON: %+v", key, goConfig, jsonConfig)
		}

		if len(goPrereqs) != len(jsonPrereqs) {
			t.Errorf("building %q prerequisite count mismatch: Go=%d, JSON=%d", key, len(goPrereqs), len(jsonPrereqs))
			continue
		}
		goPrereqMap := make(map[string]int, len(goPrereqs))
		for _, p := range goPrereqs {
			goPrereqMap[p.BuildingType] = p.MinLevel
		}
		for _, p := range jsonPrereqs {
			if goLevel, ok := goPrereqMap[p.BuildingType]; !ok {
				t.Errorf("building %q: JSON has prerequisite %q not in Go", key, p.BuildingType)
			} else if goLevel != p.MinLevel {
				t.Errorf("building %q: prerequisite %q level mismatch Go=%d JSON=%d", key, p.BuildingType, goLevel, p.MinLevel)
			}
		}
	}
}

func TestTroopsParity(t *testing.T) {
	repoRoot := findRepoRoot(t)
	jsonPath := filepath.Join(repoRoot, "client", "src", "config", "generated", "troops.json")
	jsonTroops := loadJSONFile[GeneratedTroopConfig](t, jsonPath)

	// Key set must match exactly
	goKeys := SortedTroopKeys()
	if len(jsonTroops) != len(goKeys) {
		t.Fatalf("key count mismatch: Go has %d troops, JSON has %d", len(goKeys), len(jsonTroops))
	}
	for _, key := range goKeys {
		if _, ok := jsonTroops[key]; !ok {
			t.Errorf("troop %q exists in Go but not in JSON", key)
		}
	}
	for key := range jsonTroops {
		if _, ok := TroopConfigs[key]; !ok {
			t.Errorf("troop %q exists in JSON but not in Go", key)
		}
	}

	// Field-level comparison
	for _, key := range goKeys {
		goConfig := ToGeneratedTroop(TroopConfigs[key])
		jsonConfig := jsonTroops[key]
		if !reflect.DeepEqual(goConfig, jsonConfig) {
			t.Errorf("troop %q mismatch:\n  Go:   %+v\n  JSON: %+v", key, goConfig, jsonConfig)
		}
	}
}

func TestResourceEconomyParity(t *testing.T) {
	repoRoot := findRepoRoot(t)
	jsonPath := filepath.Join(repoRoot, "client", "src", "config", "generated", "resources.json")

	data, err := os.ReadFile(jsonPath)
	if err != nil {
		t.Fatalf("failed to read %s: %v", jsonPath, err)
	}
	var jsonEconomy GeneratedResourceEconomy
	if err := json.Unmarshal(data, &jsonEconomy); err != nil {
		t.Fatalf("failed to unmarshal %s: %v", jsonPath, err)
	}

	goEconomy := ToGeneratedResourceEconomy()
	if !reflect.DeepEqual(goEconomy, jsonEconomy) {
		t.Errorf("resource economy mismatch:\n  Go:   %+v\n  JSON: %+v", goEconomy, jsonEconomy)
	}
}
