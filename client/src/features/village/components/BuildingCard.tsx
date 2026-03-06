import type { BuildingInfo } from '../../../types/api';
import type { BuildingType } from '../../../types/game';
import { GameIcon } from '../../../components/GameIcon/GameIcon';
import { BUILDING_CONFIGS } from '../../../config/buildings';
import styles from './BuildingCard.module.css';

interface BuildingCardProps {
  building: BuildingInfo;
  onClick?: () => void;
}

export function BuildingCard({ building, onClick }: BuildingCardProps) {
  const cfg = BUILDING_CONFIGS[building.building_type as BuildingType];
  const label = cfg?.displayName ?? building.building_type;

  return (
    <button
      type="button"
      className={styles.card}
      onClick={onClick}
    >
      <GameIcon assetId={building.building_type} fallback="🏗️" size={28} className={styles.icon} />
      <span className={styles.name}>{label}</span>
      <span className={styles.level}>Lv {building.level}</span>
    </button>
  );
}
