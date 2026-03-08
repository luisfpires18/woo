import { useMutation, useQueryClient } from '@tanstack/react-query';
import { startTraining, cancelTraining } from '../../../services/training';

interface TrainVariables {
  troopType: string;
  quantity: number;
}

/**
 * Mutation hook for starting troop training.
 * Invalidates the village query on success to refresh data.
 */
export function useStartTraining(villageId: number) {
  const qc = useQueryClient();

  return useMutation({
    mutationFn: ({ troopType, quantity }: TrainVariables) =>
      startTraining(villageId, troopType, quantity),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['village', villageId] });
    },
  });
}

/**
 * Mutation hook for cancelling a queued training item.
 * Invalidates the village query on success to refresh data.
 */
export function useCancelTraining(villageId: number) {
  const qc = useQueryClient();

  return useMutation({
    mutationFn: (queueId: number) => cancelTraining(villageId, queueId),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['village', villageId] });
    },
  });
}
