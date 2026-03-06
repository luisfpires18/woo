import { useState } from 'react';
import { useParams } from 'react-router-dom';
import { useVillage } from '../hooks/useVillage';
import { useAuthStore } from '../../../stores/authStore';
import { BuildingGrid } from '../components/BuildingGrid';
import { ResourcePanel } from '../components/ResourcePanel';
import { ResourceFieldsGrid } from '../components/ResourceFieldsGrid';
import { ConstructionQueue } from '../components/ConstructionQueue';
import { BuildModal } from '../components/BuildModal';
import { BuildingDetailModal } from '../components/BuildingDetailModal';
import { LoadingSpinner } from '../../../components/LoadingSpinner/LoadingSpinner';
import type { BuildingInfo } from '../../../types/api';
import styles from './VillagePage.module.css';

export function VillagePage() {
  const { id } = useParams<{ id: string }>();
  const villageId = Number(id);
  const player = useAuthStore((s) => s.player);

  const { data: village, isLoading, error } = useVillage(villageId);

  // Modal state
  const [buildModalOpen, setBuildModalOpen] = useState(false);
  const [selectedBuilding, setSelectedBuilding] = useState<BuildingInfo | null>(null);

  if (isLoading) {
    return (
      <div className={styles.loading}>
        <LoadingSpinner size="lg" />
      </div>
    );
  }

  if (error || !village) {
    return (
      <div className={styles.error}>
        <p>Failed to load village.</p>
      </div>
    );
  }

  const buildQueue = village.build_queue ?? [];

  return (
    <div className={styles.page}>
      <header className={styles.header}>
        <h1 className={styles.title}>{village.name}</h1>
        <span className={styles.coords}>
          ({village.x}, {village.y})
        </span>
      </header>

      <ConstructionQueue queue={buildQueue} villageId={villageId} />

      <div className={styles.content}>
        <div className={styles.main}>
          <section className={styles.section}>
            <h2 className={styles.sectionTitle}>Resource Fields</h2>
            <ResourceFieldsGrid
              buildings={village.buildings}
              onBuildingClick={(b) => setSelectedBuilding(b)}
            />
          </section>

          <section className={styles.section}>
            <h2 className={styles.sectionTitle}>Buildings</h2>
            <BuildingGrid
              buildings={village.buildings}
              onBuildingClick={(b) => setSelectedBuilding(b)}
              onEmptySlotClick={() => setBuildModalOpen(true)}
            />
          </section>
        </div>

        <aside className={styles.sidebar}>
          <ResourcePanel resources={village.resources} />
        </aside>
      </div>

      {/* Build new building modal */}
      <BuildModal
        isOpen={buildModalOpen}
        onClose={() => setBuildModalOpen(false)}
        buildings={village.buildings}
        villageId={villageId}
        resources={village.resources}
        playerKingdom={player?.kingdom ?? ''}
        queue={buildQueue}
      />

      {/* Building detail / upgrade modal */}
      {selectedBuilding && (
        <BuildingDetailModal
          isOpen={true}
          onClose={() => setSelectedBuilding(null)}
          building={selectedBuilding}
          allBuildings={village.buildings}
          villageId={villageId}
          resources={village.resources}
          queue={buildQueue}
        />
      )}
    </div>
  );
}
