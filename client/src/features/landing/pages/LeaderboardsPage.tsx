import styles from './LeaderboardsPage.module.css';

export function LeaderboardsPage() {
  return (
    <div className={styles.page}>
      <div className={styles.section}>
        <h1 className={styles.sectionTitle}>Leaderboards</h1>
        <div className={styles.placeholder}>
          <span className={styles.placeholderIcon}>🏆</span>
          <span className={styles.placeholderText}>Coming Soon</span>
          <p className={styles.placeholderDesc}>
            Compete against other players and kingdoms. Leaderboards will be available once seasons are live.
          </p>
        </div>
      </div>
    </div>
  );
}
