// Expeditions history page — full expedition log with battle report access

import { useEffect, useState } from 'react';
import { useExpeditionStore } from '../../../stores/expeditionStore';
import { fetchExpeditions } from '../../../services/camp';
import { BattleReportModal } from '../../map/components/BattleReportModal';
import { Card } from '../../../components/Card/Card';
import { Button } from '../../../components/Button/Button';
import { LoadingSpinner } from '../../../components/LoadingSpinner/LoadingSpinner';
import type { ExpeditionResponse } from '../../../types/api';
import styles from './ExpeditionsPage.module.css';

function statusLabel(status: string): string {
  return status.charAt(0).toUpperCase() + status.slice(1);
}

function statusClass(status: string): string {
  switch (status) {
    case 'marching':
      return styles.statusMarching ?? '';
    case 'battling':
      return styles.statusBattling ?? '';
    case 'returning':
      return styles.statusReturning ?? '';
    case 'completed':
      return styles.statusCompleted ?? '';
    default:
      return '';
  }
}

function timeRemaining(targetIso: string): string {
  const diff = new Date(targetIso).getTime() - Date.now();
  if (diff <= 0) return 'arrived';
  const secs = Math.ceil(diff / 1000);
  if (secs < 60) return `${secs}s`;
  const mins = Math.floor(secs / 60);
  const remainSecs = secs % 60;
  return `${mins}m ${remainSecs}s`;
}

function formatDate(iso: string): string {
  return new Date(iso).toLocaleString();
}

export function ExpeditionsPage() {
  const expeditions = useExpeditionStore((s) => s.expeditions);
  const setExpeditions = useExpeditionStore((s) => s.setExpeditions);
  const [loading, setLoading] = useState(true);
  const [viewBattleId, setViewBattleId] = useState<number | null>(null);
  const [, setTick] = useState(0);

  // Fetch expeditions on mount and every 5s
  useEffect(() => {
    fetchExpeditions()
      .then((e) => setExpeditions(e ?? []))
      .catch(() => {})
      .finally(() => setLoading(false));
    const interval = setInterval(() => {
      fetchExpeditions()
        .then((e) => setExpeditions(e ?? []))
        .catch(() => {});
    }, 5_000);
    return () => clearInterval(interval);
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  // Re-render every second for timers
  useEffect(() => {
    const id = setInterval(() => setTick((t) => t + 1), 1000);
    return () => clearInterval(id);
  }, []);

  // Sort: active first (marching, battling, returning), then completed newest first
  const active = expeditions
    .filter((e) => e.status !== 'completed')
    .sort((a, b) => new Date(b.dispatched_at).getTime() - new Date(a.dispatched_at).getTime());
  const completed = expeditions
    .filter((e) => e.status === 'completed')
    .sort((a, b) => new Date(b.completed_at ?? b.dispatched_at).getTime() - new Date(a.completed_at ?? a.dispatched_at).getTime());

  if (loading) {
    return (
      <div className={styles.loading}>
        <LoadingSpinner size="lg" />
      </div>
    );
  }

  const renderExpedition = (exp: ExpeditionResponse) => {
    const totalTroops = exp.troops.reduce((s, t) => s + t.quantity_sent, 0);
    const totalSurvived = exp.troops.reduce((s, t) => s + t.quantity_survived, 0);

    return (
      <Card className={styles.expCard} key={exp.id}>
        <div className={styles.expHeader}>
          <span className={styles.expId}>Expedition #{exp.id}</span>
          <span className={`${styles.expStatus} ${statusClass(exp.status)}`}>
            {statusLabel(exp.status)}
          </span>
        </div>

        <div className={styles.expDetails}>
          <div className={styles.detailRow}>
            <span className={styles.detailLabel}>Camp</span>
            <span className={styles.detailValue}>#{exp.camp_id}</span>
          </div>
          <div className={styles.detailRow}>
            <span className={styles.detailLabel}>Troops sent</span>
            <span className={styles.detailValue}>{totalTroops}</span>
          </div>
          {exp.status === 'completed' && (
            <div className={styles.detailRow}>
              <span className={styles.detailLabel}>Survived</span>
              <span className={styles.detailValue}>{totalSurvived}</span>
            </div>
          )}
          <div className={styles.detailRow}>
            <span className={styles.detailLabel}>Dispatched</span>
            <span className={styles.detailValue}>{formatDate(exp.dispatched_at)}</span>
          </div>
          {exp.status === 'marching' && (
            <div className={styles.detailRow}>
              <span className={styles.detailLabel}>Arrives in</span>
              <span className={styles.detailValue}>{timeRemaining(exp.arrives_at)}</span>
            </div>
          )}
          {exp.status === 'returning' && exp.return_at && (
            <div className={styles.detailRow}>
              <span className={styles.detailLabel}>Returns in</span>
              <span className={styles.detailValue}>{timeRemaining(exp.return_at)}</span>
            </div>
          )}
          {exp.completed_at && (
            <div className={styles.detailRow}>
              <span className={styles.detailLabel}>Completed</span>
              <span className={styles.detailValue}>{formatDate(exp.completed_at)}</span>
            </div>
          )}
        </div>

        {exp.status === 'completed' && (
          <div className={styles.troopBreakdown}>
            {exp.troops.map((t) => (
              <span className={styles.troopChip} key={t.troop_type}>
                {t.troop_type}: {t.quantity_survived}/{t.quantity_sent}
              </span>
            ))}
          </div>
        )}

        {exp.battle_id && (
          <Button
            variant="secondary"
            size="sm"
            onClick={() => setViewBattleId(exp.battle_id!)}
            style={{ marginTop: 'var(--spacing-xs)' }}
          >
            View Battle Report
          </Button>
        )}
      </Card>
    );
  };

  return (
    <div className={styles.page}>
      <h1 className={styles.title}>Expeditions</h1>

      {active.length === 0 && completed.length === 0 && (
        <p className={styles.empty}>No expeditions yet. Send troops from the world map!</p>
      )}

      {active.length > 0 && (
        <section>
          <h2 className={styles.sectionTitle}>Active</h2>
          <div className={styles.list}>{active.map(renderExpedition)}</div>
        </section>
      )}

      {completed.length > 0 && (
        <section>
          <h2 className={styles.sectionTitle}>Completed</h2>
          <div className={styles.list}>{completed.map(renderExpedition)}</div>
        </section>
      )}

      {viewBattleId !== null && (
        <BattleReportModal
          battleId={viewBattleId}
          onClose={() => setViewBattleId(null)}
        />
      )}
    </div>
  );
}
