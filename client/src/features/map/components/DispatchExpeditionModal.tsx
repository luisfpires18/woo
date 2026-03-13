// Dispatch expedition modal — select troops from active village to attack a camp

import { useState, useMemo } from 'react';
import { useGameStore } from '../../../stores/gameStore';
import { useExpeditionStore } from '../../../stores/expeditionStore';
import { dispatchExpedition } from '../../../services/camp';
import type { CampResponse, TroopInfo } from '../../../types/api';
import styles from './DispatchExpeditionModal.module.css';

interface DispatchExpeditionModalProps {
  camp: CampResponse;
  onClose: () => void;
}

export function DispatchExpeditionModal({ camp, onClose }: DispatchExpeditionModalProps) {
  const village = useGameStore((s) => s.currentVillage);
  const addExpedition = useExpeditionStore((s) => s.addExpedition);
  const [troopAmounts, setTroopAmounts] = useState<Record<string, number>>({});
  const [error, setError] = useState('');
  const [dispatching, setDispatching] = useState(false);

  // Filter troops that are available (idle status, quantity > 0)
  const availableTroops: TroopInfo[] = useMemo(() => {
    if (!village?.troops) return [];
    return village.troops.filter((t) => t.quantity > 0 && t.status === 'stationed');
  }, [village?.troops]);

  const totalSelected = Object.values(troopAmounts).reduce((s, n) => s + (n || 0), 0);

  const handleAmountChange = (troopType: string, value: string) => {
    const num = Math.max(0, parseInt(value, 10) || 0);
    const available = availableTroops.find((t) => t.type === troopType)?.quantity ?? 0;
    setTroopAmounts((prev) => ({
      ...prev,
      [troopType]: Math.min(num, available),
    }));
  };

  const handleDispatch = async () => {
    if (!village) return;
    if (totalSelected === 0) {
      setError('Select at least one troop to send.');
      return;
    }

    setError('');
    setDispatching(true);

    const troops = Object.entries(troopAmounts)
      .filter(([, qty]) => qty > 0)
      .map(([troop_type, quantity]) => ({ troop_type, quantity }));

    try {
      const expedition = await dispatchExpedition(village.id, {
        camp_id: camp.id,
        troops,
      });
      addExpedition(expedition);
      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to dispatch expedition');
    } finally {
      setDispatching(false);
    }
  };

  if (!village) return null;

  return (
    <div className={styles.overlay} onClick={onClose}>
      <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
        <h2 className={styles.title}>Attack Camp</h2>
        <p className={styles.subtitle}>
          Send troops from <strong>{village.name}</strong>
        </p>

        <div className={styles.campInfo}>
          <div className={styles.campRow}>
            <span className={styles.campLabel}>Target</span>
            <span className={styles.campValue}>{camp.template_name}</span>
          </div>
          <div className={styles.campRow}>
            <span className={styles.campLabel}>Tier</span>
            <span className={styles.campValue}>{camp.tier}</span>
          </div>
          <div className={styles.campRow}>
            <span className={styles.campLabel}>Location</span>
            <span className={styles.campValue}>({camp.tile_x}, {camp.tile_y})</span>
          </div>
          <div className={styles.campRow}>
            <span className={styles.campLabel}>Defenders</span>
            <span className={styles.campValue}>
              {(camp.beasts ?? []).reduce((s, b) => s + b.count, 0)} beasts
            </span>
          </div>
        </div>

        <h3 className={styles.sectionTitle}>Select Troops</h3>

        {availableTroops.length === 0 ? (
          <div className={styles.noTroops}>No troops available in this village.</div>
        ) : (
          <div className={styles.troopList}>
            {availableTroops.map((troop) => (
              <div className={styles.troopRow} key={troop.type}>
                <span className={styles.troopName}>
                  {troop.type.replace(/_/g, ' ').replace(/\b\w/g, (c) => c.toUpperCase())}
                </span>
                <span className={styles.troopAvailable}>/{troop.quantity}</span>
                <input
                  className={styles.troopInput}
                  type="number"
                  min={0}
                  max={troop.quantity}
                  value={troopAmounts[troop.type] || 0}
                  onChange={(e) => handleAmountChange(troop.type, e.target.value)}
                />
              </div>
            ))}
          </div>
        )}

        <div className={styles.totalRow}>
          <span className={styles.totalLabel}>Total troops</span>
          <span className={styles.totalValue}>{totalSelected}</span>
        </div>

        {error && <div className={styles.errorMsg}>{error}</div>}

        <div className={styles.actions}>
          <button className={styles.cancelBtn} onClick={onClose}>
            Cancel
          </button>
          <button
            className={styles.dispatchBtn}
            onClick={handleDispatch}
            disabled={dispatching || totalSelected === 0}
          >
            {dispatching ? 'Dispatching...' : '⚔ Dispatch'}
          </button>
        </div>
      </div>
    </div>
  );
}
