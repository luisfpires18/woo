import styles from './ResourceBar.module.css';
import type { ResourcesResponse } from '../../types/api';
import { useResourceTicker } from '../../hooks/useResourceTicker';

interface ResourceBarProps {
  resources: ResourcesResponse;
}

const RESOURCE_CONFIG = [
  { key: 'food' as const, maxKey: 'max_food' as const, label: 'Food', emoji: '\uD83C\uDF3E' },
  { key: 'water' as const, maxKey: 'max_water' as const, label: 'Water', emoji: '\uD83D\uDCA7' },
  { key: 'lumber' as const, maxKey: 'max_lumber' as const, label: 'Lumber', emoji: '\uD83E\uDEB5' },
  { key: 'stone' as const, maxKey: 'max_stone' as const, label: 'Stone', emoji: '\uD83E\uDEA8' },
];

export function ResourceBar({ resources }: ResourceBarProps) {
  const live = useResourceTicker(resources);

  return (
    <div className={styles.bar}>
      {RESOURCE_CONFIG.map(({ key, maxKey, label, emoji }) => {
        const amount = Math.floor(live[key]);
        const rate = live[`${key}_rate`];
        const max = Math.floor(live[maxKey]);
        return (
          <div key={key} className={styles.resource} title={`${label}: ${amount} / ${max}`}>
            <span className={styles.emoji}>{emoji}</span>
            <span className={styles.amount}>{amount.toLocaleString()}</span>
            <span className={styles.rate}>+{rate}/s</span>
          </div>
        );
      })}
    </div>
  );
}
