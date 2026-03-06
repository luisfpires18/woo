import type { ResourcesResponse } from '../../../types/api';
import { useResourceTicker } from '../../../hooks/useResourceTicker';
import { GameIcon } from '../../../components/GameIcon/GameIcon';
import styles from './ResourcePanel.module.css';

interface ResourcePanelProps {
  resources: ResourcesResponse;
}

const RESOURCE_ROWS: {
  key: keyof ResourcesResponse;
  rateKey: keyof ResourcesResponse;
  label: string;
  assetId: string;
  fallbackIcon: string;
}[] = [
  { key: 'iron', rateKey: 'iron_rate', label: 'Iron', assetId: 'iron', fallbackIcon: '⛏️' },
  { key: 'wood', rateKey: 'wood_rate', label: 'Wood', assetId: 'wood', fallbackIcon: '🪵' },
  { key: 'stone', rateKey: 'stone_rate', label: 'Stone', assetId: 'stone', fallbackIcon: '🪨' },
  { key: 'food', rateKey: 'food_rate', label: 'Food', assetId: 'food', fallbackIcon: '🌾' },
];

export function ResourcePanel({ resources }: ResourcePanelProps) {
  const live = useResourceTicker(resources);

  return (
    <div className={styles.panel}>
      <h3 className={styles.heading}>Resources</h3>

      <div className={styles.rows}>
        {RESOURCE_ROWS.map((r) => (
          <div key={r.key} className={styles.row}>
            <GameIcon assetId={r.assetId} fallback={r.fallbackIcon} size={18} className={styles.icon} />
            <span className={styles.label}>{r.label}</span>
            <span className={styles.amount}>
              {Math.floor(live[r.key])}
            </span>
            <span className={styles.rate}>
              +{Math.floor(live[r.rateKey])}/s
            </span>
          </div>
        ))}
      </div>

      <div className={styles.storage}>
        <span className={styles.storageLabel}>Storage</span>
        <span className={styles.storageValue}>{live.max_storage}</span>
      </div>

      {live.food_consumption > 0 && (
        <div className={styles.consumption}>
          <span className={styles.consumptionLabel}>Food consumption</span>
          <span className={styles.consumptionValue}>
            -{live.food_consumption}/s
          </span>
        </div>
      )}
    </div>
  );
}
