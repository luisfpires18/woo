import type { TroopInfo } from '../../../types/api';
import type { TroopType } from '../../../config/troops';
import { TroopCard } from './TroopCard';
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
        {troops.map((troop) => (
          <TroopCard
            key={troop.type}
            troopType={troop.type as TroopType}
            quantity={troop.quantity}
          />
        ))}
      </div>
    </div>
  );
}
