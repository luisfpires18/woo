package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/luisfpires18/woo/internal/config"
)

func main() {
	outDir := flag.String("out", filepath.Join("..", "client", "src", "config", "generated"), "output directory for generated JSON files")
	flag.Parse()

	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "error creating output dir: %v\n", err)
		os.Exit(1)
	}

	// --- Buildings ---
	buildingKeys := config.SortedBuildingKeys()
	buildings := make(map[string]config.GeneratedBuildingConfig, len(buildingKeys))
	for _, k := range buildingKeys {
		buildings[k] = config.ToGeneratedBuilding(config.BuildingConfigs[k])
	}
	if err := writeJSON(filepath.Join(*outDir, "buildings.json"), buildings); err != nil {
		fmt.Fprintf(os.Stderr, "error writing buildings.json: %v\n", err)
		os.Exit(1)
	}

	// --- Troops ---
	troopKeys := config.SortedTroopKeys()
	troops := make(map[string]config.GeneratedTroopConfig, len(troopKeys))
	for _, k := range troopKeys {
		troops[k] = config.ToGeneratedTroop(config.TroopConfigs[k])
	}
	if err := writeJSON(filepath.Join(*outDir, "troops.json"), troops); err != nil {
		fmt.Fprintf(os.Stderr, "error writing troops.json: %v\n", err)
		os.Exit(1)
	}

	// --- Resources ---
	resources := config.ToGeneratedResourceEconomy()
	if err := writeJSON(filepath.Join(*outDir, "resources.json"), resources); err != nil {
		fmt.Fprintf(os.Stderr, "error writing resources.json: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated %d buildings, %d troops, and resource economy → %s\n", len(buildings), len(troops), *outDir)
}

// writeJSON serialises v to a JSON file with sorted keys and 2-space indentation.
func writeJSON(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}
