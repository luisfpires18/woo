import { useState, useEffect } from 'react';
import { getProfile } from '../../../services/player';
import type { PlayerProfileResponse } from '../../../types/api';
import { LoadingSpinner } from '../../../components/LoadingSpinner/LoadingSpinner';
import styles from './ProfilePage.module.css';

function StatusBadge({ status }: { status: string }) {
  const cls =
    status === 'active'
      ? styles.badgeActive
      : status === 'ended'
        ? styles.badgeEnded
        : styles.badgeArchived;
  return <span className={`${styles.badge} ${cls}`}>{status}</span>;
}

export function ProfilePage() {
  const [profile, setProfile] = useState<PlayerProfileResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        const data = await getProfile();
        if (!cancelled) setProfile(data);
      } catch {
        if (!cancelled) setError('Failed to load profile.');
      } finally {
        if (!cancelled) setLoading(false);
      }
    })();
    return () => { cancelled = true; };
  }, []);

  if (loading) {
    return (
      <div className={styles.center}>
        <LoadingSpinner size="lg" />
      </div>
    );
  }

  if (error || !profile) {
    return (
      <div className={styles.page}>
        <div className={styles.error}>{error ?? 'Profile not found.'}</div>
      </div>
    );
  }

  return (
    <div className={styles.page}>
      {/* Profile card */}
      <div className={styles.profileCard}>
        <div className={styles.avatar}>{profile.username[0]}</div>
        <div className={styles.profileInfo}>
          <h1 className={styles.username}>{profile.username}</h1>
          <div className={styles.meta}>
            <span>{profile.email}</span>
            {profile.role === 'admin' && (
              <span className={styles.roleBadge}>Admin</span>
            )}
            <span>Joined {new Date(profile.created_at).toLocaleDateString()}</span>
          </div>
        </div>
      </div>

      {/* Stats */}
      <div className={styles.stats}>
        <div className={styles.statCard}>
          <div className={styles.statValue}>{profile.total_seasons}</div>
          <div className={styles.statLabel}>Seasons Played</div>
        </div>
        <div className={styles.statCard}>
          <div className={styles.statValue}>
            {profile.season_history.reduce((sum, e) => sum + e.village_count, 0)}
          </div>
          <div className={styles.statLabel}>Total Villages</div>
        </div>
      </div>

      {/* Season history */}
      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>Season History</h2>

        {profile.season_history.length === 0 ? (
          <p className={styles.empty}>You haven&apos;t participated in any season yet.</p>
        ) : (
          <div className={styles.tableWrap}>
            <table className={styles.table}>
              <thead>
                <tr>
                  <th>Season</th>
                  <th>Status</th>
                  <th>Kingdom</th>
                  <th>Villages</th>
                  <th>Joined</th>
                </tr>
              </thead>
              <tbody>
                {profile.season_history.map((entry) => (
                  <tr key={entry.season_id}>
                    <td>{entry.season_name}</td>
                    <td><StatusBadge status={entry.season_status} /></td>
                    <td className={styles.kingdom}>{entry.kingdom}</td>
                    <td>{entry.village_count}</td>
                    <td>{new Date(entry.joined_at).toLocaleDateString()}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>
    </div>
  );
}
