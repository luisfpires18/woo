import { useMutation, useQueryClient } from '@tanstack/react-query';
import { startUpgrade, cancelUpgrade } from '../../../services/village';

/**
 * Mutation hook for starting a building upgrade.
 * Automatically invalidates the village query on success to refresh data.
 */
export function useStartUpgrade(villageId: number) {
  const qc = useQueryClient();

  return useMutation({
    mutationFn: (buildingType: string) => startUpgrade(villageId, buildingType),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['village', villageId] });
    },
  });
}

/**
 * Mutation hook for cancelling a queued building upgrade.
 * Automatically invalidates the village query on success to refresh data.
 */
export function useCancelUpgrade(villageId: number) {
  const qc = useQueryClient();

  return useMutation({
    mutationFn: (queueId: number) => cancelUpgrade(villageId, queueId),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['village', villageId] });
    },
  });
}
