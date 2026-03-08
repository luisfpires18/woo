import type { BuildingInfo } from '../../../types/api';
import { GameIcon } from '../../../components/GameIcon/GameIcon';
import { RESOURCE_BUILDING_GROUPS } from '../../../config/buildings';
import { useBuildingDisplayNames } from '../../../hooks/useBuildingDisplayNames';
import styles from './ResourceFieldsGrid.module.css';

interface ResourceFieldsGridProps {
  buildings: BuildingInfo[];
  onBuildingClick: (building: BuildingInfo) => void;
}

export function ResourceFieldsGrid({ buildings, onBuildingClick }: ResourceFieldsGridProps) {
  const { getDisplayName } = useBuildingDisplayNames();
  const buildingMap = new Map<string, BuildingInfo>();
  for (const b of buildings) {
    buildingMap.set(b.building_type, b);
  }

  return (
    <div className={styles.groups}>
      {RESOURCE_BUILDING_GROUPS.map((group) => (
        <div key={group.resource} className={styles.group}>
          <h4 className={styles.groupTitle}>
            <span>{group.emoji}</span> {group.label}
          </h4>
          <div className={styles.grid}>
            {group.types.map((type) => {
              const b = buildingMap.get(type);
              if (!b) return null;
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
                    fallback={group.emoji}
                    size={28}
                    className={styles.icon}
                  />
                  <span className={styles.name}>{getDisplayName(type)}</span>
                  <span className={styles.level}>
                    {isBuilt ? `Lv ${b.level}` : 'Not built'}
                  </span>
                </button>
              );
            })}
          </div>
        </div>
      ))}
    </div>
  );
}
