// Battle Report modal — choice-first flow: outcome hidden until player decides

import { useCallback, useEffect, useState } from 'react';
import { useExpeditionStore } from '../../../stores/expeditionStore';
import { fetchBattleReport, fetchBattleReplay } from '../../../services/camp';
import type { BattleReportResponse, BattleReplayResponse } from '../../../types/api';
import { Modal } from '../../../components/Modal/Modal';
import { Card } from '../../../components/Card/Card';
import { Button } from '../../../components/Button/Button';
import { BattleReplayCanvas } from './BattleReplayCanvas';
import styles from './BattleReportModal.module.css';

type Phase = 'choice' | 'replay' | 'report';

interface BattleReportModalProps {
  battleId: number;
  onClose: () => void;
}

function resultLabel(result: string): string {
  switch (result) {
    case 'attacker_won': return 'Victory!';
    case 'defender_won': return 'Defeat';
    case 'draw': return 'Draw';
    default: return result;
  }
}

function resultClass(result: string): string {
  switch (result) {
    case 'attacker_won': return styles.victory ?? '';
    case 'defender_won': return styles.defeat ?? '';
    case 'draw': return styles.draw ?? '';
    default: return '';
  }
}

export function BattleReportModal({ battleId, onClose }: BattleReportModalProps) {
  const cachedReport = useExpeditionStore((s) => s.battleReports[battleId]);
  const cacheBattleReport = useExpeditionStore((s) => s.cacheBattleReport);
  const isViewed = useExpeditionStore((s) => s.viewedBattles.has(battleId));
  const markBattleViewed = useExpeditionStore((s) => s.markBattleViewed);
  const [report, setReport] = useState<BattleReportResponse | null>(cachedReport ?? null);
  const [loading, setLoading] = useState(!cachedReport);
  const [error, setError] = useState('');

  // Phase state machine — skip choice if already viewed
  const [phase, setPhase] = useState<Phase>(isViewed ? 'report' : 'choice');

  // Replay state
  const [replay, setReplay] = useState<BattleReplayResponse | null>(null);
  const [replayLoading, setReplayLoading] = useState(false);
  const [replayError, setReplayError] = useState('');

  useEffect(() => {
    if (cachedReport) {
      setReport(cachedReport);
      return;
    }
    setLoading(true);
    fetchBattleReport(battleId)
      .then((r) => {
        setReport(r);
        cacheBattleReport(r);
      })
      .catch((err) => setError(err instanceof Error ? err.message : 'Failed to load report'))
      .finally(() => setLoading(false));
  }, [battleId, cachedReport, cacheBattleReport]);

  // Choice → Replay: load replay data then transition
  const handleSeeBattle = useCallback(async () => {
    setReplayLoading(true);
    setReplayError('');
    try {
      const data = await fetchBattleReplay(battleId);
      setReplay(data);
      markBattleViewed(battleId);
      setPhase('replay');
    } catch (err) {
      setReplayError(err instanceof Error ? err.message : 'Failed to load replay');
    } finally {
      setReplayLoading(false);
    }
  }, [battleId, markBattleViewed]);

  // Choice → Report: skip straight to results
  const handleSkipToResults = useCallback(() => {
    markBattleViewed(battleId);
    setPhase('report');
  }, [battleId, markBattleViewed]);

  // Replay → Report: after watching, see results
  const handleReplayDone = useCallback(() => {
    setPhase('report');
  }, []);

  return (
    <Modal isOpen onClose={onClose} title="Battle Report" size="lg">
      {loading && <div className={styles.loading}>Loading battle report...</div>}
      {error && <div className={styles.loading}>{error}</div>}

      {/* ── Choice screen: outcome hidden ── */}
      {report && phase === 'choice' && (
        <div className={styles.choiceScreen}>
          <div className={styles.choiceIcon}>⚔️</div>
          <div className={styles.choiceText}>
            Your forces clashed with the camp defenders.
          </div>
          <div className={styles.choiceButtons}>
            <Button
              variant="primary"
              onClick={handleSeeBattle}
              loading={replayLoading}
            >
              See Battle
            </Button>
            <Button variant="ghost" onClick={handleSkipToResults}>
              Skip to Results
            </Button>
          </div>
          {replayError && <div className={styles.errorText}>{replayError}</div>}
        </div>
      )}

      {/* ── Replay phase ── */}
      {phase === 'replay' && replay && (
        <BattleReplayCanvas replay={replay} onClose={handleReplayDone} />
      )}

      {/* ── Report phase: full results ── */}
      {report && phase === 'report' && (
        <div className={styles.reportContent}>
          <div className={`${styles.resultBanner} ${resultClass(report.result)}`}>
            {resultLabel(report.result)}
          </div>

          <div className={styles.statsColumns}>
            <Card header="Your Forces">
              <div className={styles.statRow}>
                <span className={styles.statLabel}>Sent</span>
                <span className={styles.statValue}>{report.attacker_losses.total_sent}</span>
              </div>
              <div className={styles.statRow}>
                <span className={styles.statLabel}>Lost</span>
                <span className={`${styles.statValue} ${styles.statLost}`}>
                  {report.attacker_losses.total_lost}
                </span>
              </div>
              <div className={styles.statRow}>
                <span className={styles.statLabel}>Survived</span>
                <span className={`${styles.statValue} ${styles.statSurvived}`}>
                  {report.attacker_losses.total_survived}
                </span>
              </div>
            </Card>

            <Card header="Defenders">
              <div className={styles.statRow}>
                <span className={styles.statLabel}>Total</span>
                <span className={styles.statValue}>{report.defender_losses.total_sent}</span>
              </div>
              <div className={styles.statRow}>
                <span className={styles.statLabel}>Slain</span>
                <span className={`${styles.statValue} ${styles.statSurvived}`}>
                  {report.defender_losses.total_lost}
                </span>
              </div>
              <div className={styles.statRow}>
                <span className={styles.statLabel}>Survived</span>
                <span className={`${styles.statValue} ${styles.statLost}`}>
                  {report.defender_losses.total_survived}
                </span>
              </div>
            </Card>
          </div>

          <Card header="Rewards">
            {report.rewards.length > 0 ? (
              <div className={styles.rewardList}>
                {report.rewards.map((r, i) => (
                  <div className={styles.rewardChip} key={i}>
                    <span className={styles.rewardType}>{r.resource_type}</span>
                    <span className={styles.rewardAmount}>+{r.amount}</span>
                  </div>
                ))}
              </div>
            ) : (
              <div className={styles.noRewards}>No rewards — camp defenders survived.</div>
            )}
          </Card>

          <div className={styles.timestamp}>
            Battle fought at {new Date(report.fought_at).toLocaleString()}
          </div>
        </div>
      )}
    </Modal>
  );
}
