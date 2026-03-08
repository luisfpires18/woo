import { useState, useEffect, useRef } from 'react';
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
  const refetchScheduled = useRef(false);

  useEffect(() => {
    if (queue.length === 0) return;
    setNow(Date.now());
    const id = setInterval(() => setNow(Date.now()), 1000);
    return () => clearInterval(id);
  }, [queue.length]);

  // Auto-refetch when earliest training completes
  useEffect(() => {
    if (queue.length === 0) return;
    refetchScheduled.current = false;

    const earliest = Math.min(
      ...queue.map((q) => new Date(q.completes_at).getTime()),
    );
    const delay = earliest - Date.now() + 1500;

    if (delay <= 0) {
      queryClient.invalidateQueries({ queryKey: ['village', villageId] });
      return;
    }

    const timer = setTimeout(() => {
      if (!refetchScheduled.current) {
        refetchScheduled.current = true;
        queryClient.invalidateQueries({ queryKey: ['village', villageId] });
      }
    }, delay);

    return () => clearTimeout(timer);
  }, [queue, villageId, queryClient]);

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
  const startMs = new Date(item.started_at).getTime();
  const endMs = new Date(item.completes_at).getTime();
  const totalMs = endMs - startMs;
  const elapsedMs = now - startMs;
  const remainingMs = Math.max(0, endMs - now);
  const progress = totalMs > 0 ? Math.max(0, Math.min(100, (elapsedMs / totalMs) * 100)) : 0;
  const remainingSec = Math.ceil(remainingMs / 1000);

  return (
    <div className={styles.item}>
      <div className={styles.info}>
        <span className={styles.name}>
          {displayName} × {item.quantity}
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