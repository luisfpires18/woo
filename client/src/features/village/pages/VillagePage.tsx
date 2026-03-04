import { useParams } from 'react-router-dom';
import { useVillage } from '../hooks/useVillage';
import { BuildingGrid } from '../components/BuildingGrid';
import { ResourcePanel } from '../components/ResourcePanel';
import { LoadingSpinner } from '../../../components/LoadingSpinner/LoadingSpinner';
import styles from './VillagePage.module.css';

export function VillagePage() {
  const { id } = useParams<{ id: string }>();
  const villageId = Number(id);

  const { data: village, isLoading, error } = useVillage(villageId);

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

  return (
    <div className={styles.page}>
      <header className={styles.header}>
        <h1 className={styles.title}>{village.name}</h1>
        <span className={styles.coords}>
          ({village.x}, {village.y})
        </span>
      </header>

      <div className={styles.content}>
        <section className={styles.buildings}>
          <h2 className={styles.sectionTitle}>Buildings</h2>
          <BuildingGrid buildings={village.buildings} />
        </section>

        <aside className={styles.sidebar}>
          <ResourcePanel resources={village.resources} />
        </aside>
      </div>
    </div>
  );
}
