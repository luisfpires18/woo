import type { ResourcesResponse } from '../../../types/api';
import styles from './ResourcePanel.module.css';

interface ResourcePanelProps {
  resources: ResourcesResponse;
}

const RESOURCE_ROWS: {
  key: keyof ResourcesResponse;
  rateKey: keyof ResourcesResponse;
  label: string;
  icon: string;
}[] = [
  { key: 'iron', rateKey: 'iron_rate', label: 'Iron', icon: '⛏️' },
  { key: 'wood', rateKey: 'wood_rate', label: 'Wood', icon: '🪵' },
  { key: 'stone', rateKey: 'stone_rate', label: 'Stone', icon: '🪨' },
  { key: 'food', rateKey: 'food_rate', label: 'Food', icon: '🌾' },
];

export function ResourcePanel({ resources }: ResourcePanelProps) {
  return (
    <div className={styles.panel}>
      <h3 className={styles.heading}>Resources</h3>

      <div className={styles.rows}>
        {RESOURCE_ROWS.map((r) => (
          <div key={r.key} className={styles.row}>
            <span className={styles.icon}>{r.icon}</span>
            <span className={styles.label}>{r.label}</span>
            <span className={styles.amount}>
              {Math.floor(resources[r.key])}
            </span>
            <span className={styles.rate}>
              +{Math.floor(resources[r.rateKey])}/h
            </span>
          </div>
        ))}
      </div>

      <div className={styles.storage}>
        <span className={styles.storageLabel}>Storage</span>
        <span className={styles.storageValue}>{resources.max_storage}</span>
      </div>

      {resources.food_consumption > 0 && (
        <div className={styles.consumption}>
          <span className={styles.consumptionLabel}>Food consumption</span>
          <span className={styles.consumptionValue}>
            -{resources.food_consumption}/h
          </span>
        </div>
      )}
    </div>
  );
}
