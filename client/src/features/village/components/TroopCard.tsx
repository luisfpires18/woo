import { useState } from 'react';
import type { TroopType } from '../../../config/troops';
import { useTroopDisplay } from '../../../hooks/useTroopDisplay';
import styles from './TroopCard.module.css';

interface TroopCardProps {
  /** Troop type key (e.g. "iron_legionary"). */
  troopType: TroopType;
  /** Number of troops of this type. */
  quantity: number;
}

export function TroopCard({ troopType, quantity }: TroopCardProps) {
  const { getDisplay } = useTroopDisplay();
  const { displayName, spriteUrl, emoji } = getDisplay(troopType);
  const [imgError, setImgError] = useState(false);

  return (
    <div className={styles.card}>
      <div className={styles.icon}>
        {spriteUrl && !imgError ? (
          <img
            src={spriteUrl}
            alt={displayName}
            width={28}
            height={28}
            onError={() => setImgError(true)}
            draggable={false}
            style={{ objectFit: 'contain' }}
          />
        ) : (
          <span style={{ fontSize: '28px', lineHeight: 1 }}>{emoji}</span>
        )}
      </div>
      <span className={styles.name}>{displayName}</span>
      <span className={styles.quantity}>{quantity.toLocaleString()}</span>
    </div>
  );
}
