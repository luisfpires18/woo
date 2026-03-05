import { useState, useEffect, useCallback } from 'react';
import { fetchStats } from '../../../services/admin';
import type { StatsResponse } from '../../../types/api';
import { LoadingSpinner } from '../../../components/LoadingSpinner/LoadingSpinner';
import styles from './AdminStatsPage.module.css';

const STAT_CARDS = [
  { key: 'total_players' as const, label: 'Total Players', icon: '👥' },
  { key: 'total_villages' as const, label: 'Total Villages', icon: '🏘️' },
];

export function AdminStatsPage() {
  const [stats, setStats] = useState<StatsResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const load = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const resp = await fetchStats();
      setStats(resp);
    } catch {
      setError('Failed to load stats.');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  if (loading) {
    return (
      <div className={styles.center}>
        <LoadingSpinner size="lg" />
      </div>
    );
  }

  if (error || !stats) {
    return (
      <div className={styles.page}>
        <h2 className={styles.heading}>Server Statistics</h2>
        <div className={styles.error}>{error ?? 'Unknown error.'}</div>
      </div>
    );
  }

  return (
    <div className={styles.page}>
      <h2 className={styles.heading}>Server Statistics</h2>
      <p className={styles.subtitle}>Real-time overview of the game world.</p>

      <div className={styles.grid}>
        {STAT_CARDS.map((card) => (
          <div key={card.key} className={styles.card}>
            <span className={styles.cardIcon}>{card.icon}</span>
            <div className={styles.cardBody}>
              <span className={styles.cardValue}>{stats[card.key]}</span>
              <span className={styles.cardLabel}>{card.label}</span>
            </div>
          </div>
        ))}
      </div>

      <div className={styles.actions}>
        <button onClick={load} className={styles.refreshBtn}>
          Refresh
        </button>
      </div>
    </div>
  );
}
