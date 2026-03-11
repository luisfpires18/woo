import { GameIcon } from '../../../components/GameIcon/GameIcon';
import styles from './ResourceCard.module.css';

interface ResourceCardProps {
  /** Asset id for GameIcon lookup (e.g. "food", "water"). */
  assetId: string;
  /** Emoji fallback when no sprite uploaded. */
  fallbackIcon: string;
  /** Display name (e.g. "Food"). */
  label: string;
  /** Current amount (floored). */
  current: number;
  /** Maximum storage capacity. */
  max: number;
  /** Production rate per second. */
  rate: number;
}

export function ResourceCard({
  assetId,
  fallbackIcon,
  label,
  current,
  max,
  rate,
}: ResourceCardProps) {
  const ratio = max > 0 ? Math.min(current / max, 1) : 0;
  const pct = Math.round(ratio * 100);

  // Color thresholds for progress bar
  let fillClass = styles.progressFill;
  if (ratio >= 0.95) {
    fillClass = `${styles.progressFill} ${styles.danger}`;
  } else if (ratio >= 0.8) {
    fillClass = `${styles.progressFill} ${styles.warning}`;
  }

  return (
    <div className={styles.card}>
      <div className={styles.iconRow}>
        <GameIcon assetId={assetId} fallback={fallbackIcon} size={22} />
        <span className={styles.name}>{label}</span>
      </div>

      <div className={styles.capacity}>
        <span className={styles.current}>{current.toLocaleString()}</span>
        <span className={styles.separator}>/</span>
        <span className={styles.max}>{max.toLocaleString()}</span>
      </div>

      <div className={styles.progressTrack}>
        <div className={fillClass} style={{ width: `${pct}%` }} />
      </div>

      <span className={rate < 0 ? styles.rateNegative : styles.rate}>
        {rate >= 0 ? `+${rate}` : `${rate}`}/s
      </span>
    </div>
  );
}
