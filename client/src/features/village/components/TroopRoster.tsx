import type { TroopInfo } from '../../../types/api';
import { TROOP_CONFIGS } from '../../../config/troops';
import styles from './TroopRoster.module.css';

interface TroopRosterProps {
  troops: TroopInfo[];
}

export function TroopRoster({ troops }: TroopRosterProps) {
  if (!troops || troops.length === 0) {
    return (
      <div className={styles.container}>
        <h3 className={styles.heading}>Troops</h3>
        <p className={styles.empty}>No troops stationed.</p>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <h3 className={styles.heading}>Troops</h3>
      <div className={styles.grid}>
        {troops.map((troop) => {
          const cfg = TROOP_CONFIGS[troop.type];
          const displayName = cfg?.displayName ?? troop.type;
          return (
            <div key={troop.type} className={styles.card}>
              <span className={styles.name}>{displayName}</span>
              <span className={styles.qty}>{troop.quantity}</span>
              {cfg && (
                <div className={styles.stats}>
                  <span title="Attack">⚔️ {cfg.attack}</span>
                  <span title="Infantry Def">🛡️ {cfg.defInfantry}</span>
                  <span title="Speed">🏃 {cfg.speed}</span>
                </div>
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
}
