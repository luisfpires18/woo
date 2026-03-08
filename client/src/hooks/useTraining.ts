import { useState, useCallback } from 'react';
import type { TrainingQueueResponse, TrainingCostResponse, TroopInfo } from '../types/api';
import * as trainingApi from '../services/training';

interface UseTrainingReturn {
  /** Start training troops. Returns the queue item on success. */
  train: (villageId: number, troopType: string, quantity: number) => Promise<TrainingQueueResponse>;
  /** Cancel a training queue item. */
  cancel: (villageId: number, queueId: number) => Promise<void>;
  /** Fetch training cost preview. */
  getCost: (villageId: number, troopType: string, quantity: number) => Promise<TrainingCostResponse>;
  /** Fetch troops stationed in a village. */
  getTroops: (villageId: number) => Promise<TroopInfo[]>;
  /** Whether an API call is in progress. */
  loading: boolean;
  /** Last error message, if any. */
  error: string | null;
}

/**
 * Hook that wraps training API calls with loading/error state.
 */
export function useTraining(): UseTrainingReturn {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const train = useCallback(async (villageId: number, troopType: string, quantity: number) => {
    setLoading(true);
    setError(null);
    try {
      const result = await trainingApi.startTraining(villageId, troopType, quantity);
      return result;
    } catch (err) {
      const msg = err instanceof Error ? err.message : 'Failed to start training';
      setError(msg);
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  const cancel = useCallback(async (villageId: number, queueId: number) => {
    setLoading(true);
    setError(null);
    try {
      await trainingApi.cancelTraining(villageId, queueId);
    } catch (err) {
      const msg = err instanceof Error ? err.message : 'Failed to cancel training';
      setError(msg);
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  const getCost = useCallback(async (villageId: number, troopType: string, quantity: number) => {
    setLoading(true);
    setError(null);
    try {
      const result = await trainingApi.getTrainingCost(villageId, troopType, quantity);
      return result;
    } catch (err) {
      const msg = err instanceof Error ? err.message : 'Failed to get training cost';
      setError(msg);
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  const getTroops = useCallback(async (villageId: number) => {
    setLoading(true);
    setError(null);
    try {
      const result = await trainingApi.fetchTroops(villageId);
      return result.troops;
    } catch (err) {
      const msg = err instanceof Error ? err.message : 'Failed to fetch troops';
      setError(msg);
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  return { train, cancel, getCost, getTroops, loading, error };
}
