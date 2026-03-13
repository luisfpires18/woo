// Troop movement banner — shows outgoing expeditions and incoming attacks

import { useEffect, useState } from 'react';
import { useExpeditionStore } from '../../../stores/expeditionStore';
import { Card } from '../../../components/Card/Card';
import styles from './TroopMovements.module.css';

interface TroopMovementsProps {
  villageId: number;
}

function timeRemaining(targetIso: string): string {
  const diff = new Date(targetIso).getTime() - Date.now();
  if (diff <= 0) return 'now';
  const secs = Math.ceil(diff / 1000);
  if (secs < 60) return `${secs}s`;
  const mins = Math.floor(secs / 60);
  const remainSecs = secs % 60;
  return `${mins}m ${remainSecs}s`;
}

export function TroopMovements({ villageId }: TroopMovementsProps) {
  const expeditions = useExpeditionStore((s) => s.expeditions);
  const incomingAttacks = useExpeditionStore((s) => s.incomingAttacks);
  const [, setTick] = useState(0);

  // Re-render every second for countdown timers
  useEffect(() => {
    const id = setInterval(() => setTick((t) => t + 1), 1000);
    return () => clearInterval(id);
  }, []);

  const outgoing = expeditions.filter(
    (e) => e.village_id === villageId && (e.status === 'marching' || e.status === 'returning'),
  );
  const incoming = incomingAttacks.filter((a) => a.village_id === villageId);

  if (outgoing.length === 0 && incoming.length === 0) return null;

  return (
    <Card className={styles.banner}>
      {incoming.map((atk) => (
        <div className={`${styles.movement} ${styles.incoming}`} key={`atk-${atk.id}`}>
          <span className={styles.icon}>⚠️</span>
          <span className={styles.label}>Incoming attack!</span>
          <span className={styles.timer}>{timeRemaining(atk.arrives_at)}</span>
        </div>
      ))}
      {outgoing.map((exp) => (
        <div
          className={`${styles.movement} ${exp.status === 'marching' ? styles.marching : styles.returning}`}
          key={`exp-${exp.id}`}
        >
          <span className={styles.icon}>{exp.status === 'marching' ? '⚔️' : '🔙'}</span>
          <span className={styles.label}>
            {exp.status === 'marching'
              ? `Troops marching to Camp #${exp.camp_id}`
              : `Troops returning from Camp #${exp.camp_id}`}
          </span>
          <span className={styles.detail}>
            {exp.troops.reduce((s, t) => s + t.quantity_sent, 0)} troops
          </span>
          <span className={styles.timer}>
            {exp.status === 'marching'
              ? timeRemaining(exp.arrives_at)
              : exp.return_at
                ? timeRemaining(exp.return_at)
                : '—'}
          </span>
        </div>
      ))}
    </Card>
  );
}
