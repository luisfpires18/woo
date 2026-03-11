import { useState, useEffect, useMemo } from 'react';
import { useQueryClient, useMutation } from '@tanstack/react-query';
import type { BuildingQueueResponse } from '../../../types/api';
import { formatDuration } from '../../../config/buildings';
import { useBuildingDisplayNames } from '../../../hooks/useBuildingDisplayNames';
import { useCancelUpgrade } from '../hooks/useBuildingUpgrade';
import { useAuthStore } from '../../../stores/authStore';
import { instantCompleteBuild } from '../../../services/village';
import styles from './ConstructionQueue.module.css';

interface ConstructionQueueProps {
  queue: BuildingQueueResponse[];
  villageId: number;
}

export function ConstructionQueue({ queue, villageId }: ConstructionQueueProps) {
  const cancelMutation = useCancelUpgrade(villageId);
  const queryClient = useQueryClient();
  const player = useAuthStore((s) => s.player);
  const isAdmin = player?.role === 'admin';
  const [now, setNow] = useState(() => Date.now());

  // Tick `now` every second while queue is active
  useEffect(() => {
    if (queue.length === 0) return;
    setNow(Date.now());
    const id = setInterval(() => setNow(Date.now()), 1000);
    return () => clearInterval(id);
  }, [queue.length]);

  // Stable key: only changes when queue IDs or completion times actually change
  const queueKey = useMemo(
    () => queue.map((q) => `${q.id}:${q.completes_at}`).join(','),
    [queue],
  );

  // Auto-refetch when earliest build completes, with polling fallback
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
      <h3 className={styles.heading}>Construction Queue</h3>
      {queue.map((item) => (
        <ConstructionQueueItem
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

interface ConstructionQueueItemProps {
  item: BuildingQueueResponse;
  now: number;
  villageId: number;
  isAdmin: boolean;
  cancelMutation: ReturnType<typeof useCancelUpgrade>;
}

function ConstructionQueueItem({ item, now, villageId, isAdmin, cancelMutation }: ConstructionQueueItemProps) {
  const queryClient = useQueryClient();
  const { getDisplayName } = useBuildingDisplayNames();

  const instantMutation = useMutation({
    mutationFn: () => instantCompleteBuild(item.id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['village', villageId] });
    },
  });

  const displayName = getDisplayName(item.building_type);
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
          {displayName} → Lv {item.target_level}
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
          aria-label={`Cancel ${displayName} upgrade`}
        >
          ✕
        </button>
      </div>
    </div>
  );
}
