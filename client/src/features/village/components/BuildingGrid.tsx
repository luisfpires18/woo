import type { BuildingInfo } from '../../../types/api';
import { RESOURCE_BUILDING_TYPES } from '../../../config/buildings';
import { isMilitaryBuilding } from '../../../config/troops';
import { BuildingCard } from './BuildingCard';
import styles from './BuildingGrid.module.css';

interface BuildingGridProps {
  buildings: BuildingInfo[];
  onBuildingClick: (building: BuildingInfo) => void;
  onEmptySlotClick: () => void;
}

export function BuildingGrid({ buildings, onBuildingClick, onEmptySlotClick }: BuildingGridProps) {
  // Filter out resource field buildings (they go in ResourceFieldsGrid)
  const villageBuildings = buildings.filter(
    (b) => !RESOURCE_BUILDING_TYPES.has(b.building_type),
  );

  const built = villageBuildings.filter((b) => b.level > 0);
  const unbuiltCount = villageBuildings.filter((b) => b.level === 0).length;

  if (built.length === 0 && unbuiltCount === 0) {
    return <p className={styles.empty}>No buildings yet.</p>;
  }

  return (
    <div className={styles.grid}>
      {built.map((b) => (
        <BuildingCard
          key={b.id}
          building={b}
          onClick={() => onBuildingClick(b)}
          isMilitary={isMilitaryBuilding(b.building_type)}
        />
      ))}
      {unbuiltCount > 0 && (
        <button
          type="button"
          className={styles.emptySlot}
          onClick={onEmptySlotClick}
          aria-label="Build new building"
        >
          <span className={styles.plusIcon}>+</span>
          <span className={styles.emptyLabel}>Build</span>
        </button>
      )}
    </div>
  );
}
