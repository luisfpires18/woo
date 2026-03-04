import type { BuildingInfo } from '../../../types/api';
import { BuildingCard } from './BuildingCard';
import styles from './BuildingGrid.module.css';

interface BuildingGridProps {
  buildings: BuildingInfo[];
}

export function BuildingGrid({ buildings }: BuildingGridProps) {
  if (buildings.length === 0) {
    return <p className={styles.empty}>No buildings yet.</p>;
  }

  return (
    <div className={styles.grid}>
      {buildings.map((b) => (
        <BuildingCard key={b.id} building={b} />
      ))}
    </div>
  );
}
