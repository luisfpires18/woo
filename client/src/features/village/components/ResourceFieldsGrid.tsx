import type { BuildingInfo } from '../../../types/api';
import type { BuildingType } from '../../../types/game';
import { GameIcon } from '../../../components/GameIcon/GameIcon';
import { BUILDING_CONFIGS, RESOURCE_BUILDING_TYPES } from '../../../config/buildings';
import styles from './ResourceFieldsGrid.module.css';

interface ResourceFieldsGridProps {
  buildings: BuildingInfo[];
  onBuildingClick: (building: BuildingInfo) => void;
}

/** Order in which resource buildings appear. */
const FIELD_ORDER: BuildingType[] = ['iron_mine', 'lumber_mill', 'quarry', 'farm'];

export function ResourceFieldsGrid({ buildings, onBuildingClick }: ResourceFieldsGridProps) {
  const resourceBuildings = FIELD_ORDER.map((type) =>
    buildings.find((b) => b.building_type === type && RESOURCE_BUILDING_TYPES.has(b.building_type)),
  ).filter(Boolean) as BuildingInfo[];

  return (
    <div className={styles.grid}>
      {resourceBuildings.map((b) => {
        const cfg = BUILDING_CONFIGS[b.building_type as BuildingType];
        const isBuilt = b.level > 0;

        return (
          <button
            key={b.id}
            type="button"
            className={`${styles.card} ${isBuilt ? '' : styles.unbuilt}`}
            onClick={() => onBuildingClick(b)}
          >
            <GameIcon
              assetId={b.building_type}
              fallback="🌿"
              size={28}
              className={styles.icon}
            />
            <span className={styles.name}>{cfg.displayName}</span>
            <span className={styles.level}>
              {isBuilt ? `Lv ${b.level}` : 'Not built'}
            </span>
          </button>
        );
      })}
    </div>
  );
}
