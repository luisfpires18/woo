import styles from './ResourceBar.module.css';
import type { ResourcesResponse } from '../../types/api';
import { useResourceTicker } from '../../hooks/useResourceTicker';

interface ResourceBarProps {
  resources: ResourcesResponse;
}

const RESOURCE_CONFIG = [
  { key: 'food' as const, label: 'Food', emoji: '\uD83C\uDF3E' },
  { key: 'water' as const, label: 'Water', emoji: '\uD83D\uDCA7' },
  { key: 'lumber' as const, label: 'Lumber', emoji: '\uD83E\uDEB5' },
  { key: 'stone' as const, label: 'Stone', emoji: '\uD83E\uDEA8' },
];

export function ResourceBar({ resources }: ResourceBarProps) {
  const live = useResourceTicker(resources);

  return (
    <div className={styles.bar}>
      {RESOURCE_CONFIG.map(({ key, label, emoji }) => {
        const amount = Math.floor(live[key]);
        const rate = live[`${key}_rate`];
        return (
          <div key={key} className={styles.resource} title={label}>
            <span className={styles.emoji}>{emoji}</span>
            <span className={styles.amount}>{amount.toLocaleString()}</span>
            <span className={styles.rate}>+{rate}/s</span>
          </div>
        );
      })}
      <div className={styles.resource} title="Storage">
        <span className={styles.emoji}>{'\uD83C\uDFE0'}</span>
        <span className={styles.amount}>
          {Math.floor(live.max_storage).toLocaleString()}
        </span>
      </div>
    </div>
  );
}
