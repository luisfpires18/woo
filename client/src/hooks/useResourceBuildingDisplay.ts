import { useCallback } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useAuthStore } from '../stores/authStore';
import { fetchResourceBuildingConfigs } from '../services/village';
import type { ResourceBuildingConfig } from '../types/api';

interface ResourceBuildingDisplay {
  displayName: string;
  spriteUrl: string | null;
  emoji: string;
  config: ResourceBuildingConfig | null;
}

/**
 * Hook that fetches resource building configs for the player's kingdom
 * and returns a resolver function that provides display names and sprite URLs.
 *
 * Resource buildings use the sprite resolver endpoint which does prefix-based
 * file matching: /api/sprites/building/{kingdom}/{resource}_{slot}
 *
 * Falls back to the building_type key (e.g. "food_1") if no config is found.
 */
export function useResourceBuildingDisplay() {
  const kingdom = useAuthStore((s) => s.player?.kingdom ?? '');

  const { data: configMap } = useQuery({
    queryKey: ['resource-building-configs', kingdom],
    queryFn: async () => {
      const resp = await fetchResourceBuildingConfigs(kingdom);
      const map: Record<string, ResourceBuildingConfig> = {};
      for (const cfg of resp.configs) {
        // Key by building_type pattern: "food_1", "water_2", etc.
        map[`${cfg.resource_type}_${cfg.slot}`] = cfg;
      }
      return map;
    },
    enabled: !!kingdom,
    staleTime: 5 * 60_000,
  });

  /** Resolve display info for a resource building type (e.g. "food_1"). */
  const getDisplay = useCallback(
    (buildingType: string): ResourceBuildingDisplay => {
      const cfg = configMap?.[buildingType] ?? null;

      // Build sprite URL via the resolver endpoint
      let spriteUrl: string | null = null;
      if (kingdom && buildingType) {
        spriteUrl = `/api/sprites/building/${kingdom}/${buildingType}`;
      }

      return {
        displayName: cfg?.display_name ?? buildingType,
        spriteUrl,
        emoji: cfg?.default_icon ?? '🏗️',
        config: cfg,
      };
    },
    [configMap, kingdom],
  );

  return { getDisplay, kingdom };
}
