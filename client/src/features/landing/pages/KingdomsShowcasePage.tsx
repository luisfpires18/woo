import { KingdomCard, KINGDOMS } from '../../kingdom/components/KingdomCard';
import styles from './KingdomsShowcasePage.module.css';

const PLAYABLE_KINGDOMS = KINGDOMS.filter((k) => k.playable);

export function KingdomsShowcasePage() {
  return (
    <div className={styles.page}>
      <div className={styles.section}>
        <h1 className={styles.sectionTitle}>Kingdoms</h1>
        <p className={styles.sectionSubtitle}>
          Seven rival kingdoms compete for dominance. Each offers unique strengths and playstyles.
        </p>
        <div className={styles.kingdomGrid}>
          {PLAYABLE_KINGDOMS.map((k) => (
            <KingdomCard
              key={k.id}
              kingdom={k}
              selected={false}
              onSelect={() => {}}
              displayOnly
            />
          ))}
        </div>
      </div>
    </div>
  );
}
