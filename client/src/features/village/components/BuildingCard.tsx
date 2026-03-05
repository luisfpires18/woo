import type { BuildingInfo } from '../../../types/api';
import { GameIcon } from '../../../components/GameIcon/GameIcon';
import styles from './BuildingCard.module.css';

interface BuildingCardProps {
  building: BuildingInfo;
}

const BUILDING_LABELS: Record<string, string> = {
  town_hall: 'Town Hall',
  iron_mine: 'Iron Mine',
  lumber_mill: 'Lumber Mill',
  quarry: 'Quarry',
  farm: 'Farm',
  warehouse: 'Warehouse',
  barracks: 'Barracks',
  stable: 'Stable',
  forge: 'Forge',
  rune_altar: 'Rune Altar',
  walls: 'Walls',
  marketplace: 'Marketplace',
  embassy: 'Embassy',
  watchtower: 'Watchtower',
  dock: 'Dock',
  grove_sanctum: 'Grove Sanctum',
  colosseum: 'Colosseum',
};

export function BuildingCard({ building }: BuildingCardProps) {
  const label = BUILDING_LABELS[building.building_type] ?? building.building_type;

  return (
    <div className={styles.card}>
      <GameIcon assetId={building.building_type} fallback="🏗️" size={28} className={styles.icon} />
      <span className={styles.name}>{label}</span>
      <span className={styles.level}>Lv {building.level}</span>
    </div>
  );
}
