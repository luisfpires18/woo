import type { BuildingInfo } from '../../../types/api';
import { GameIcon } from '../../../components/GameIcon/GameIcon';
import { useBuildingDisplayNames } from '../../../hooks/useBuildingDisplayNames';
import styles from './BuildingCard.module.css';

interface BuildingCardProps {
  building: BuildingInfo;
  onClick?: () => void;
  isMilitary?: boolean;
}

export function BuildingCard({ building, onClick, isMilitary }: BuildingCardProps) {
  const { getDisplayName } = useBuildingDisplayNames();
  const label = getDisplayName(building.building_type);

  return (
    <button
      type="button"
      className={`${styles.card} ${isMilitary ? styles.military : ''}`}
      onClick={onClick}
    >
      <GameIcon assetId={building.building_type} fallback="🏗️" size={28} className={styles.icon} />
      <span className={styles.name}>{label}</span>
      <span className={styles.level}>Lv {building.level}</span>
    </button>
  );
}
