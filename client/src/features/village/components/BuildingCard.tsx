import { useState } from 'react';
import type { BuildingInfo } from '../../../types/api';
import { useBuildingDisplayNames } from '../../../hooks/useBuildingDisplayNames';
import styles from './BuildingCard.module.css';

interface BuildingCardProps {
  building: BuildingInfo;
  onClick?: () => void;
  isMilitary?: boolean;
}

export function BuildingCard({ building, onClick, isMilitary }: BuildingCardProps) {
  const { getDisplay } = useBuildingDisplayNames();
  const { displayName, spriteUrl, emoji } = getDisplay(building.building_type);
  const [imgError, setImgError] = useState(false);

  return (
    <button
      type="button"
      className={`${styles.card} ${isMilitary ? styles.military : ''}`}
      onClick={onClick}
    >
      {spriteUrl && !imgError ? (
        <img
          src={spriteUrl}
          alt={displayName}
          width={28}
          height={28}
          className={styles.icon}
          onError={() => setImgError(true)}
          draggable={false}
        />
      ) : (
        <span className={styles.icon} style={{ fontSize: '28px', lineHeight: 1 }}>{emoji}</span>
      )}
      <span className={styles.name}>{displayName}</span>
      <span className={styles.level}>Lv {building.level}</span>
    </button>
  );
}
