import { useCallback } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useAuthStore } from '../stores/authStore';
import { fetchTroopDisplayConfigs } from '../services/village';
import type { TroopDisplayConfig } from '../types/api';

interface TroopDisplay {
  displayName: string;
  spriteUrl: string | null;
  emoji: string;
  config: TroopDisplayConfig | null;
}

/**
 * Hook that fetches admin-configured troop display configs for the player's kingdom
 * and returns a resolver function providing display name, sprite URL, emoji, and
 * the full config object.
 *
 * Sprite URLs point to the resolver endpoint which does prefix-based file matching:
 * /api/sprites/troop/{kingdom}/{troopType}
 *
 * Falls back to the troop_type key if no config is found.
 */
export function useTroopDisplay() {
  const kingdom = useAuthStore((s) => s.player?.kingdom ?? '');

  const { data: configMap } = useQuery({
    queryKey: ['troop-display-configs', kingdom],
    queryFn: async () => {
      const resp = await fetchTroopDisplayConfigs(kingdom);
      const map: Record<string, TroopDisplayConfig> = {};
      for (const cfg of resp.configs) {
        map[cfg.troop_type] = cfg;
      }
      return map;
    },
    enabled: !!kingdom,
    staleTime: 5 * 60_000,
  });

  /** Resolve display info for a troop type (e.g. "arkazia_militia"). */
  const getDisplay = useCallback(
    (troopType: string): TroopDisplay => {
      const cfg = configMap?.[troopType] ?? null;

      let spriteUrl: string | null = null;
      if (kingdom && troopType) {
        spriteUrl = `/api/sprites/troop/${kingdom}/${troopType}`;
      }

      return {
        displayName: cfg?.display_name ?? troopType.replace(/_/g, ' '),
        spriteUrl,
        emoji: cfg?.default_icon ?? '⚔️',
        config: cfg,
      };
    },
    [configMap, kingdom],
  );

  return { getDisplay, kingdom };
}
