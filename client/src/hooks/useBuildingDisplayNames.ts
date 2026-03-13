import { useCallback } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useAuthStore } from '../stores/authStore';
import { fetchBuildingDisplayConfigs } from '../services/village';
import { BUILDING_CONFIGS } from '../config/buildings';
import type { BuildingType } from '../types/game';
import type { BuildingDisplayConfig } from '../types/api';

interface BuildingDisplay {
  displayName: string;
  spriteUrl: string | null;
  emoji: string;
  config: BuildingDisplayConfig | null;
}

/**
 * Hook that fetches admin-configured building display names for the player's kingdom
 * and returns resolver functions. Falls back to the hardcoded BUILDING_CONFIGS names
 * if the API data hasn't loaded yet or the building type isn't found.
 *
 * React Query deduplicates calls — multiple components using this hook
 * will share the same cached data.
 */
export function useBuildingDisplayNames() {
  const kingdom = useAuthStore((s) => s.player?.kingdom ?? '');

  const { data: configMap } = useQuery({
    queryKey: ['building-display-configs', kingdom],
    queryFn: async () => {
      const resp = await fetchBuildingDisplayConfigs(kingdom);
      const map: Record<string, BuildingDisplayConfig> = {};
      for (const cfg of resp.configs) {
        map[cfg.building_type] = cfg;
      }
      return map;
    },
    enabled: !!kingdom,
    staleTime: 5 * 60_000, // 5 minutes — display names rarely change
  });

  /** Resolve the display name for a building type. */
  const getDisplayName = useCallback(
    (buildingType: string): string => {
      return (
        configMap?.[buildingType]?.display_name ??
        BUILDING_CONFIGS[buildingType as BuildingType]?.displayName ??
        buildingType
      );
    },
    [configMap],
  );

  /** Resolve full display info (name, sprite, emoji) for a building type. */
  const getDisplay = useCallback(
    (buildingType: string): BuildingDisplay => {
      const cfg = configMap?.[buildingType] ?? null;

      let spriteUrl: string | null = null;
      if (kingdom && buildingType) {
        spriteUrl = `/api/sprites/building-display/${kingdom}/${buildingType}`;
      }

      return {
        displayName:
          cfg?.display_name ??
          BUILDING_CONFIGS[buildingType as BuildingType]?.displayName ??
          buildingType,
        spriteUrl,
        emoji: cfg?.default_icon ?? '🏗️',
        config: cfg,
      };
    },
    [configMap, kingdom],
  );

  return { getDisplayName, getDisplay, kingdom };
}
