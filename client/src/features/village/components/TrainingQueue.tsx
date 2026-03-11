import { useState, useEffect, useMemo } from 'react';
import { useQueryClient, useMutation } from '@tanstack/react-query';
import type { TrainingQueueResponse } from '../../../types/api';
import { TROOP_CONFIGS, formatDuration } from '../../../config/troops';
import { useCancelTraining } from '../hooks/useTrainingMutations';
import { useAuthStore } from '../../../stores/authStore';
import { instantCompleteTraining } from '../../../services/training';
import styles from './TrainingQueue.module.css';

interface TrainingQueueProps {
  queue: TrainingQueueResponse[];
  villageId: number;
}

export function TrainingQueue({ queue, villageId }: TrainingQueueProps) {
  const cancelMutation = useCancelTraining(villageId);
  const queryClient = useQueryClient();
  const player = useAuthStore((s) => s.player);
  const isAdmin = player?.role === 'admin';
  const [now, setNow] = useState(() => Date.now());

  // Tick `now` frequently while queue is active for smooth progress
  useEffect(() => {
    if (queue.length === 0) return;
    setNow(Date.now());
    const id = setInterval(() => setNow(Date.now()), 250);
    return () => clearInterval(id);
  }, [queue.length]);

  // Stable key: only changes when queue IDs, quantities, or completion times actually change
  const queueKey = useMemo(
    () => queue.map((q) => `${q.id}:${q.quantity}:${q.completes_at}`).join(','),
    [queue],
  );

  // Auto-refetch when earliest item completes, with polling fallback
  useEffect(() => {
    if (queue.length === 0) return;

    const earliest = Math.min(
      ...queue.map((q) => new Date(q.completes_at).getTime()),
    );
    const msUntil = earliest - Date.now();

    // Item already past due — poll every 2s until server catches up
    if (msUntil <= 0) {
      queryClient.invalidateQueries({ queryKey: ['village', villageId] });
      const poll = setInterval(() => {
        queryClient.invalidateQueries({ queryKey: ['village', villageId] });
      }, 2000);
      return () => clearInterval(poll);
    }

    // Schedule refetch for when item should complete (+1.5s buffer for game-loop tick)
    const timer = setTimeout(() => {
      queryClient.invalidateQueries({ queryKey: ['village', villageId] });
    }, msUntil + 1500);

    return () => clearTimeout(timer);
  }, [queueKey, villageId, queryClient]); // eslint-disable-line react-hooks/exhaustive-deps

  if (queue.length === 0) return null;

  return (
    <div className={styles.container}>
      <h3 className={styles.heading}>Training Queue</h3>
      {queue.map((item) => (
        <TrainingQueueItem
          key={item.id}
          item={item}
          now={now}
          villageId={villageId}
          isAdmin={isAdmin}
          cancelMutation={cancelMutation}
        />
      ))}
    </div>
  );
}

interface TrainingQueueItemProps {
  item: TrainingQueueResponse;
  now: number;
  villageId: number;
  isAdmin: boolean;
  cancelMutation: ReturnType<typeof useCancelTraining>;
}

function TrainingQueueItem({ item, now, villageId, isAdmin, cancelMutation }: TrainingQueueItemProps) {
  const queryClient = useQueryClient();

  const instantMutation = useMutation({
    mutationFn: () => instantCompleteTraining(item.id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['village', villageId] });
    },
  });

  const cfg = TROOP_CONFIGS[item.troop_type];
  const displayName = cfg?.displayName ?? item.troop_type;

  // Per-unit progress: current unit started at (completes_at - each_duration_sec)
  const eachMs = item.each_duration_sec * 1000;
  const endMs = new Date(item.completes_at).getTime();
  const unitStartMs = endMs - eachMs;
  const elapsedMs = now - unitStartMs;
  const remainingMs = Math.max(0, endMs - now);
  const progress = eachMs > 0 ? Math.max(0, Math.min(100, (elapsedMs / eachMs) * 100)) : 0;
  const remainingSec = Math.ceil(remainingMs / 1000);

  // X/Y progress: how many completed so far + current
  const completed = item.original_quantity - item.quantity;
  const currentUnit = completed + 1;

  return (
    <div className={styles.item}>
      <div className={styles.info}>
        <span className={styles.name}>
          {displayName} — {currentUnit}/{item.original_quantity}
        </span>
        <span className={styles.time}>
          {remainingSec > 0 ? formatDuration(remainingSec) : 'Completing...'}
        </span>
      </div>
      <div className={styles.barTrack}>
        <div
          className={styles.barFill}
          style={{ width: `${progress}%` }}
        />
      </div>
      <div className={styles.actions}>
        {isAdmin && (
          <button
            className={styles.instantBtn}
            type="button"
            onClick={() => instantMutation.mutate()}
            disabled={instantMutation.isPending}
            title="Admin: Complete instantly"
          >
            ⚡
          </button>
        )}
        <button
          className={styles.cancel}
          type="button"
          onClick={() => cancelMutation.mutate(item.id)}
          disabled={cancelMutation.isPending}
          aria-label={`Cancel ${displayName} training`}
        >
          ✕
        </button>
      </div>
    </div>
  );
}