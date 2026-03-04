import type { BuildingInfo } from '../../../types/api';
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

const BUILDING_ICONS: Record<string, string> = {
  town_hall: '🏛️',
  iron_mine: '⛏️',
  lumber_mill: '🪓',
  quarry: '🪨',
  farm: '🌾',
  warehouse: '📦',
  barracks: '⚔️',
  stable: '🐴',
  forge: '🔨',
  rune_altar: '🔮',
  walls: '🏰',
  marketplace: '🏪',
  embassy: '📜',
  watchtower: '👁️',
  dock: '⚓',
  grove_sanctum: '🌿',
  colosseum: '🏟️',
};

export function BuildingCard({ building }: BuildingCardProps) {
  const label = BUILDING_LABELS[building.building_type] ?? building.building_type;
  const icon = BUILDING_ICONS[building.building_type] ?? '🏗️';

  return (
    <div className={styles.card}>
      <span className={styles.icon}>{icon}</span>
      <span className={styles.name}>{label}</span>
      <span className={styles.level}>Lv {building.level}</span>
    </div>
  );
}
