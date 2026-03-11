import type { ResourcesResponse } from '../../../types/api';
import { useResourceTicker } from '../../../hooks/useResourceTicker';
import { ResourceCard } from './ResourceCard';
import styles from './ResourcePanel.module.css';

interface ResourcePanelProps {
  resources: ResourcesResponse;
}

const RESOURCE_ITEMS: {
  key: keyof ResourcesResponse;
  rateKey: keyof ResourcesResponse;
  maxKey: keyof ResourcesResponse;
  label: string;
  assetId: string;
  fallbackIcon: string;
}[] = [
  { key: 'food', rateKey: 'food_rate', maxKey: 'max_food', label: 'Food', assetId: 'food', fallbackIcon: '🌾' },
  { key: 'water', rateKey: 'water_rate', maxKey: 'max_water', label: 'Water', assetId: 'water', fallbackIcon: '💧' },
  { key: 'lumber', rateKey: 'lumber_rate', maxKey: 'max_lumber', label: 'Lumber', assetId: 'lumber', fallbackIcon: '🪵' },
  { key: 'stone', rateKey: 'stone_rate', maxKey: 'max_stone', label: 'Stone', assetId: 'stone', fallbackIcon: '🪨' },
];

export function ResourcePanel({ resources }: ResourcePanelProps) {
  const live = useResourceTicker(resources);

  const popCap = live.pop_cap ?? 0;
  const popUsed = live.pop_used ?? 0;
  const popPct = popCap > 0 ? Math.min(100, (popUsed / popCap) * 100) : 0;
  const popWarning = popPct >= 90;

  return (
    <div className={styles.panel}>
      <h3 className={styles.heading}>Resources</h3>

      <div className={styles.popSection}>
        <div className={styles.popHeader}>
          <span className={styles.popIcon}>👥</span>
          <span className={styles.popLabel}>Population</span>
          <span className={`${styles.popCount} ${popWarning ? styles.popWarning : ''}`}>
            {popUsed} / {popCap}
          </span>
        </div>
        <div className={styles.popBarTrack}>
          <div
            className={`${styles.popBarFill} ${popWarning ? styles.popBarWarning : ''}`}
            style={{ width: `${popPct}%` }}
          />
        </div>
        {popWarning && (
          <p className={styles.popWarnMsg}>Population almost full — build more to increase capacity.</p>
        )}
      </div>

      <div className={styles.grid}>
        {RESOURCE_ITEMS.map((r) => {
          // For food, show net rate (production minus troop consumption)
          const netRate = r.key === 'food'
            ? Math.floor(live[r.rateKey] - (live.food_consumption ?? 0))
            : Math.floor(live[r.rateKey]);

          return (
            <ResourceCard
              key={r.key}
              assetId={r.assetId}
              fallbackIcon={r.fallbackIcon}
              label={r.label}
              current={Math.floor(live[r.key])}
              max={live[r.maxKey]}
              rate={netRate}
            />
          );
        })}
      </div>
    </div>
  );
}

