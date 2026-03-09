import { useState, useEffect } from 'react';
import type { BuildingInfo, ResourcesResponse, TrainingQueueResponse } from '../../../types/api';
import type { Kingdom } from '../../../types/game';
import { Modal } from '../../../components/Modal';
import { GameIcon } from '../../../components/GameIcon/GameIcon';
import { getTroopsForBuilding, formatDuration } from '../../../config/troops';
import type { TroopType } from '../../../config/troops';
import { useBuildingDisplayNames } from '../../../hooks/useBuildingDisplayNames';
import { getTrainingCost } from '../../../services/training';
import { useStartTraining } from '../hooks/useTrainingMutations';
import styles from './BuildingTrainingModal.module.css';

interface BuildingTrainingModalProps {
  isOpen: boolean;
  onClose: () => void;
  building: BuildingInfo;
  villageId: number;
  resources: ResourcesResponse;
  kingdom: Kingdom | '';
  trainingQueue: TrainingQueueResponse[];
}

export function BuildingTrainingModal({
  isOpen,
  onClose,
  building,
  villageId,
  resources,
  kingdom,
  trainingQueue,
}: BuildingTrainingModalProps) {
  const troops = getTroopsForBuilding(building.building_type, kingdom as Kingdom);
  const startMutation = useStartTraining(villageId);
  const { getDisplayName } = useBuildingDisplayNames();

  const [selectedTroop, setSelectedTroop] = useState<TroopType | ''>('');
  const [quantity, setQuantity] = useState(1);
  const [costPreview, setCostPreview] = useState<{
    totalFood: number;
    totalWater: number;
    totalLumber: number;
    totalStone: number;
    eachTimeSec: number;
    totalTimeSec: number;
  } | null>(null);
  const [costLoading, setCostLoading] = useState(false);

  const displayName = getDisplayName(building.building_type);

  // Reset selection when modal opens for a different building
  useEffect(() => {
    setSelectedTroop('');
    setQuantity(1);
    setCostPreview(null);
  }, [building.id]);

  // Fetch cost preview when selection changes (debounced 300ms to reduce API calls)
  useEffect(() => {
    if (!selectedTroop || quantity < 1) {
      setCostPreview(null);
      return;
    }

    let cancelled = false;
    let timeout: number;

    // Debounce the API call by 300ms
    timeout = setTimeout(() => {
      setCostLoading(true);

      getTrainingCost(villageId, selectedTroop, quantity)
        .then((resp) => {
          if (!cancelled) {
            setCostPreview({
              totalFood: resp.total_food,
              totalWater: resp.total_water,
              totalLumber: resp.total_lumber,
              totalStone: resp.total_stone,
              eachTimeSec: resp.each_time_sec,
              totalTimeSec: resp.total_time_sec,
            });
          }
        })
        .catch(() => {
          if (!cancelled) setCostPreview(null);
        })
        .finally(() => {
          if (!cancelled) setCostLoading(false);
        });
    }, 300);

    return () => {
      cancelled = true;
      clearTimeout(timeout);
    };
  }, [villageId, selectedTroop, quantity]);

  const handleTrain = () => {
    if (!selectedTroop || quantity < 1) return;
    startMutation.mutate(
      { troopType: selectedTroop, quantity },
      {
        onSuccess: () => {
          setQuantity(1);
          setCostPreview(null);
          setSelectedTroop('');
        },
      },
    );
  };

  const canAfford = costPreview
    ? resources.food >= costPreview.totalFood &&
      resources.water >= costPreview.totalWater &&
      resources.lumber >= costPreview.totalLumber &&
      resources.stone >= costPreview.totalStone
    : false;

  const queueCount = trainingQueue.length;

  return (
    <Modal isOpen={isOpen} onClose={onClose} title={displayName} size="md">
      <div className={styles.content}>
        {/* Building header */}
        <div className={styles.headerSection}>
          <GameIcon assetId={building.building_type} fallback="⚔️" size={40} />
          <div className={styles.headerInfo}>
            <span className={styles.currentLevel}>Level {building.level}</span>
            {queueCount > 0 && (
              <span className={styles.queueBadge}>
                {queueCount} in queue
              </span>
            )}
          </div>
        </div>

        {/* Troop list */}
        <div className={styles.section}>
          <h4 className={styles.sectionTitle}>Assemble Units</h4>

          {troops.length === 0 ? (
            <p className={styles.emptyMsg}>No troops available for this building.</p>
          ) : (
            <div className={styles.troopList}>
              {troops.map(([troopType, troopCfg]) => {
                const meetsReq = building.level >= troopCfg.buildingLevelReq;
                const isSelected = selectedTroop === troopType;

                return (
                  <button
                    key={troopType}
                    type="button"
                    className={`${styles.troopCard} ${isSelected ? styles.selected : ''} ${!meetsReq ? styles.locked : ''}`}
                    onClick={() => meetsReq && setSelectedTroop(troopType)}
                    disabled={!meetsReq}
                    title={
                      meetsReq
                        ? troopCfg.displayName
                        : `Requires ${displayName} level ${troopCfg.buildingLevelReq}`
                    }
                  >
                    <div className={styles.troopHeader}>
                      <span className={styles.troopName}>{troopCfg.displayName}</span>
                      {!meetsReq && (
                        <span className={styles.lockIcon}>🔒</span>
                      )}
                    </div>
                    <div className={styles.troopStats}>
                      <span className={styles.stat}>⚔️ {troopCfg.attack}</span>
                      <span className={styles.stat}>🛡️ {troopCfg.defInfantry}</span>
                      <span className={styles.stat}>🏇 {troopCfg.defCavalry}</span>
                    </div>
                    {!meetsReq && (
                      <span className={styles.troopReq}>
                        Lv {troopCfg.buildingLevelReq} required
                      </span>
                    )}
                  </button>
                );
              })}
            </div>
          )}
        </div>

        {/* Training controls — shown when a troop is selected */}
        {selectedTroop && (
          <div className={styles.controls}>
            <label className={styles.qtyLabel}>
              Quantity
              <input
                type="number"
                className={styles.qtyInput}
                min={1}
                max={999}
                value={quantity}
                onChange={(e) => setQuantity(Math.max(1, parseInt(e.target.value) || 1))}
              />
            </label>

            {costLoading && <p className={styles.costNote}>Calculating cost...</p>}

            {costPreview && !costLoading && (
              <div className={styles.costGrid}>
                <div className={styles.costRow}>
                  <span className={styles.costIcon}>🌾</span>
                  <span className={styles.costLabel}>Food</span>
                  <span className={resources.food >= costPreview.totalFood ? styles.costOk : styles.costBad}>
                    {Math.ceil(costPreview.totalFood)}
                  </span>
                </div>
                <div className={styles.costRow}>
                  <span className={styles.costIcon}>💧</span>
                  <span className={styles.costLabel}>Water</span>
                  <span className={resources.water >= costPreview.totalWater ? styles.costOk : styles.costBad}>
                    {Math.ceil(costPreview.totalWater)}
                  </span>
                </div>
                <div className={styles.costRow}>
                  <span className={styles.costIcon}>🪵</span>
                  <span className={styles.costLabel}>Lumber</span>
                  <span className={resources.lumber >= costPreview.totalLumber ? styles.costOk : styles.costBad}>
                    {Math.ceil(costPreview.totalLumber)}
                  </span>
                </div>
                <div className={styles.costRow}>
                  <span className={styles.costIcon}>🪨</span>
                  <span className={styles.costLabel}>Stone</span>
                  <span className={resources.stone >= costPreview.totalStone ? styles.costOk : styles.costBad}>
                    {Math.ceil(costPreview.totalStone)}
                  </span>
                </div>
                <div className={styles.costTime}>
                  ⏱ {formatDuration(costPreview.eachTimeSec)}/unit · {formatDuration(costPreview.totalTimeSec)} total
                </div>
              </div>
            )}

            <button
              type="button"
              className={styles.trainBtn}
              onClick={handleTrain}
              disabled={!canAfford || startMutation.isPending || costLoading}
            >
              {startMutation.isPending
                ? 'Training...'
                : !canAfford && costPreview
                  ? 'Insufficient resources'
                  : 'Train'}
            </button>

            {startMutation.isError && (
              <p className={styles.errorMsg}>
                {startMutation.error instanceof Error
                  ? startMutation.error.message
                  : 'Failed to start training'}
              </p>
            )}
          </div>
        )}
      </div>
    </Modal>
  );
}
