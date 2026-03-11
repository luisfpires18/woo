import { TROOP_CONFIGS, type TroopType } from '../../../config/troops';
import { GameIcon } from '../../../components/GameIcon/GameIcon';
import styles from './TroopCard.module.css';

interface TroopCardProps {
  /** Troop type key (e.g. "iron_legionary"). */
  troopType: TroopType;
  /** Number of troops of this type. */
  quantity: number;
}

export function TroopCard({ troopType, quantity }: TroopCardProps) {
  const cfg = TROOP_CONFIGS[troopType];
  const displayName = cfg?.displayName ?? troopType;

  return (
    <div className={styles.card}>
      <div className={styles.icon}>
        <GameIcon assetId={troopType} fallback="⚔️" size={28} />
      </div>
      <span className={styles.name}>{displayName}</span>
      <span className={styles.quantity}>{quantity.toLocaleString()}</span>
    </div>
  );
}
