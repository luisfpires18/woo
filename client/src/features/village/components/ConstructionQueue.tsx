import { useState, useEffect, useRef } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import type { BuildingQueueResponse } from '../../../types/api';
import { formatDuration } from '../../../config/buildings';
import { useBuildingDisplayNames } from '../../../hooks/useBuildingDisplayNames';
import { useCancelUpgrade } from '../hooks/useBuildingUpgrade';
import styles from './ConstructionQueue.module.css';

interface ConstructionQueueProps {
  queue: BuildingQueueResponse[];
  villageId: number;
}

export function ConstructionQueue({ queue, villageId }: ConstructionQueueProps) {
  const cancelMutation = useCancelUpgrade(villageId);
  const { getDisplayName } = useBuildingDisplayNames();
  const queryClient = useQueryClient();
  const [now, setNow] = useState(() => Date.now());
  const refetchScheduled = useRef(false);

  // Refresh `now` immediately when queue transitions from empty → non-empty
  // so the first render uses a current timestamp, not the stale mount-time value.
  useEffect(() => {
    if (queue.length === 0) return;
    setNow(Date.now());
    const id = setInterval(() => setNow(Date.now()), 1000);
    return () => clearInterval(id);
  }, [queue.length]);

  // Auto-refetch village data when the earliest build completes
  useEffect(() => {
    if (queue.length === 0) return;
    refetchScheduled.current = false;

    const earliest = Math.min(
      ...queue.map((q) => new Date(q.completes_at).getTime()),
    );
    const delay = earliest - Date.now() + 1500; // +1.5 s buffer for server game-loop tick

    if (delay <= 0) {
      // Already past — refetch immediately
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
      <h3 className={styles.heading}>Construction Queue</h3>
      {queue.map((item) => {
        const displayName = getDisplayName(item.building_type);
        const startMs = new Date(item.started_at).getTime();
        const endMs = new Date(item.completes_at).getTime();
        const totalMs = endMs - startMs;
        const elapsedMs = now - startMs;
        const remainingMs = Math.max(0, endMs - now);
        const progress = totalMs > 0 ? Math.max(0, Math.min(100, (elapsedMs / totalMs) * 100)) : 0;
        const remainingSec = Math.ceil(remainingMs / 1000);

        return (
          <div key={item.id} className={styles.item}>
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
        );
      })}
    </div>
  );
}
