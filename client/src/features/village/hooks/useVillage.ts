import { useQuery } from '@tanstack/react-query';
import { fetchVillage } from '../../../services/village';
import { useGameStore } from '../../../stores/gameStore';
import type { VillageResponse } from '../../../types/api';

export function useVillage(villageId: number) {
  const setCurrentVillage = useGameStore((s) => s.setCurrentVillage);

  return useQuery<VillageResponse>({
    queryKey: ['village', villageId],
    queryFn: async () => {
      const village = await fetchVillage(villageId);
      setCurrentVillage(village);
      return village;
    },
    enabled: villageId > 0,
    refetchInterval: 60_000, // refresh resources every 60 s
  });
}
