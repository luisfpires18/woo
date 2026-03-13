// Expedition tracker panel — displays the player's active and recent expeditions

import { useEffect, useState } from 'react';
import { useExpeditionStore } from '../../../stores/expeditionStore';
import { Button } from '../../../components/Button/Button';
import type { ExpeditionResponse } from '../../../types/api';
import styles from './ExpeditionPanel.module.css';

interface ExpeditionPanelProps {
  onViewReport?: (battleId: number) => void;
  onClose?: () => void;
}

function statusClass(status: string): string {
  switch (status) {
    case 'marching': return styles.statusMarching ?? '';
    case 'battling': return styles.statusBattling ?? '';
    case 'returning': return styles.statusReturning ?? '';
    case 'completed': return styles.statusCompleted ?? '';
    default: return '';
  }
}

function formatStatus(status: string): string {
  return status.charAt(0).toUpperCase() + status.slice(1);
}

function timeRemaining(targetIso: string): string {
  const diff = new Date(targetIso).getTime() - Date.now();
  if (diff <= 0) return 'now';
  const secs = Math.ceil(diff / 1000);
  if (secs < 60) return `${secs}s`;
  const mins = Math.floor(secs / 60);
  const remainSecs = secs % 60;
  return `${mins}m ${remainSecs}s`;
}

export function ExpeditionPanel({ onViewReport, onClose }: ExpeditionPanelProps) {
  const expeditions = useExpeditionStore((s) => s.expeditions);
  const dismissExpedition = useExpeditionStore((s) => s.dismissExpedition);
  const [collapsed, setCollapsed] = useState(false);
  const [, setTick] = useState(0);

  // Re-render every second so countdown timers update
  useEffect(() => {
    const id = setInterval(() => setTick((t) => t + 1), 1000);
    return () => clearInterval(id);
  }, []);

  // Show only non-completed or recently completed expeditions
  const active = expeditions.filter((e) => e.status !== 'completed');
  const recent = expeditions
    .filter((e) => e.status === 'completed')
    .slice(-5)
    .reverse();

  if (active.length === 0 && recent.length === 0) return null;

  const renderExpedition = (exp: ExpeditionResponse) => (
    <div className={styles.expeditionCard} key={exp.id}>
      {exp.status === 'completed' && (
        <button
          className={styles.dismissBtn}
          onClick={() => dismissExpedition(exp.id)}
          aria-label="Dismiss expedition"
        >
          ✕
        </button>
      )}
      <div className={styles.expeditionRow}>
        <span className={styles.expLabel}>Camp</span>
        <span className={styles.expValue}>ID #{exp.camp_id}</span>
      </div>
      <div className={styles.expeditionRow}>
        <span className={styles.expLabel}>Status</span>
        <span className={statusClass(exp.status)}>{formatStatus(exp.status)}</span>
      </div>
      {exp.status === 'marching' && (
        <div className={styles.expeditionRow}>
          <span className={styles.expLabel}>Arrives</span>
          <span className={styles.expValue}>{timeRemaining(exp.arrives_at)}</span>
        </div>
      )}
      {exp.status === 'returning' && exp.return_at && (
        <div className={styles.expeditionRow}>
          <span className={styles.expLabel}>Returns</span>
          <span className={styles.expValue}>{timeRemaining(exp.return_at)}</span>
        </div>
      )}
      <div className={styles.expeditionRow}>
        <span className={styles.expLabel}>Troops</span>
        <span className={styles.expValue}>
          {exp.troops.reduce((s, t) => s + t.quantity_sent, 0)} sent
        </span>
      </div>
      {exp.battle_id && onViewReport && (
        <Button
          variant="secondary"
          size="sm"
          style={{ width: '100%', marginTop: 'var(--spacing-xs)' }}
          onClick={() => onViewReport(exp.battle_id!)}
        >
          View Battle Report
        </Button>
      )}
    </div>
  );

  return (
    <div className={styles.panel}>
      <div className={styles.title}>
        <span className={styles.titleLeft}>
          Expeditions
          {active.length > 0 && <span className={styles.badge}>{active.length}</span>}
        </span>
        <span className={styles.titleActions}>
          <button
            className={styles.collapseBtn}
            onClick={() => setCollapsed((c) => !c)}
            aria-label={collapsed ? 'Expand' : 'Collapse'}
          >
            {collapsed ? '▲' : '▼'}
          </button>
          {onClose && (
            <button
              className={styles.closeBtn}
              onClick={onClose}
              aria-label="Close expeditions panel"
            >
              ✕
            </button>
          )}
        </span>
      </div>

      {!collapsed && (
        <div className={styles.expeditionList}>
          {active.map(renderExpedition)}
          {recent.length > 0 && active.length > 0 && (
            <hr className={styles.divider} />
          )}
          {recent.map(renderExpedition)}
        </div>
      )}
    </div>
  );
}
