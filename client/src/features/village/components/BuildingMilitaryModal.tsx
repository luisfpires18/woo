import { useState, useEffect } from 'react';
import type { BuildingInfo, ResourcesResponse, BuildingQueueResponse, TrainingQueueResponse } from '../../../types/api';
import type { Kingdom } from '../../../types/game';
import { Modal } from '../../../components/Modal';
import { GameIcon } from '../../../components/GameIcon/GameIcon';
import {
  BUILDING_CONFIGS,
  costAtLevel,
  timeAtLevel,
  checkPrerequisites,
  formatDuration as formatBuildDuration,
} from '../../../config/buildings';
import { getTroopsForBuilding, formatDuration } from '../../../config/troops';
import type { TroopType } from '../../../config/troops';
import { useBuildingDisplayNames } from '../../../hooks/useBuildingDisplayNames';
import { useStartUpgrade } from '../hooks/useBuildingUpgrade';
import { getTrainingCost } from '../../../services/training';
import { useStartTraining } from '../hooks/useTrainingMutations';
import styles from './BuildingMilitaryModal.module.css';

interface BuildingMilitaryModalProps {
  isOpen: boolean;
  onClose: () => void;
  building: BuildingInfo;
  allBuildings: BuildingInfo[];
  villageId: number;
  resources: ResourcesResponse;
  kingdom: Kingdom | '';
  buildQueue: BuildingQueueResponse[];
  trainingQueue: TrainingQueueResponse[];
}

export function BuildingMilitaryModal({
  isOpen,
  onClose,
  building,
  allBuildings,
  villageId,
  resources,
  kingdom,
  buildQueue,
  trainingQueue,
}: BuildingMilitaryModalProps) {
  const { getDisplayName } = useBuildingDisplayNames();
  const upgradeMutation = useStartUpgrade(villageId);
  const startTraining = useStartTraining(villageId);

  const type = building.building_type;
  const cfg = BUILDING_CONFIGS[type];
  const displayName = getDisplayName(type);

  // ── Upgrade state ──────────────────────────────────────────────────────────
  const isMaxLevel = building.level >= (cfg?.maxLevel ?? 0);
  const targetLevel = building.level + 1;
  const effectiveTarget = building.level === 0 ? 1 : targetLevel;

  const queueActive = buildQueue.length > 0;
  const isUpgrading = buildQueue.some((q) => q.building_type === type);

  const cost = !isMaxLevel ? costAtLevel(type, effectiveTarget) : null;
  const timeSec = !isMaxLevel ? timeAtLevel(type, effectiveTarget) : 0;
  const prereqs = checkPrerequisites(type, allBuildings, getDisplayName);
  const unmetPrereqs = prereqs.checks.filter((p) => !p.met);

  const canAffordUpgrade = cost
    ? resources.food >= cost.food &&
      resources.water >= cost.water &&
      resources.lumber >= cost.lumber &&
      resources.stone >= cost.stone
    : false;

  const canUpgrade = !isMaxLevel && prereqs.allMet && canAffordUpgrade && !queueActive;

  const handleUpgrade = async () => {
    try {
      await upgradeMutation.mutateAsync(type);
      onClose();
    } catch {
      // Error displayed via mutation state
    }
  };

  // ── Training state ─────────────────────────────────────────────────────────
  const troops = getTroopsForBuilding(type, kingdom as Kingdom);

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

  useEffect(() => {
    setSelectedTroop('');
    setQuantity(1);
    setCostPreview(null);
  }, [building.id]);

  useEffect(() => {
    if (!selectedTroop || quantity < 1) {
      setCostPreview(null);
      setCostLoading(false);
      return;
    }
    setCostPreview(null);
    setCostLoading(true);
    let cancelled = false;
    const timeout = setTimeout(() => {
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
        .catch(() => { if (!cancelled) setCostPreview(null); })
        .finally(() => { if (!cancelled) setCostLoading(false); });
    }, 300);
    return () => { cancelled = true; clearTimeout(timeout); };
  }, [villageId, selectedTroop, quantity]);

  const handleTrain = () => {
    if (!selectedTroop || quantity < 1) return;
    startTraining.mutate(
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

  const canAffordTraining = costPreview
    ? resources.food >= costPreview.totalFood &&
      resources.water >= costPreview.totalWater &&
      resources.lumber >= costPreview.totalLumber &&
      resources.stone >= costPreview.totalStone
    : false;

  const queueCount = trainingQueue.length;

  return (
    <Modal isOpen={isOpen} onClose={onClose} title={displayName} size="md">
      <div className={styles.content}>

        {/* ── Header ─────────────────────────────────────────────────────── */}
        <div className={styles.headerSection}>
          <GameIcon assetId={type} fallback="⚔️" size={40} />
          <div className={styles.headerInfo}>
            <span className={styles.currentLevel}>Level {building.level}</span>
            {isMaxLevel && <span className={styles.maxBadge}>Max Level</span>}
            {queueCount > 0 && (
              <span className={styles.queueBadge}>{queueCount} training</span>
            )}
          </div>
        </div>

        {/* ── Training section ────────────────────────────────────────────── */}
        <div className={styles.section}>
          <h4 className={styles.sectionTitle}>Assemble Units</h4>

          {troops.length === 0 ? (
            <p className={styles.emptyMsg}>No troops available for this building.</p>
          ) : (
            <div className={styles.troopList}>
              {[...troops]
                .sort(([, a], [, b]) => {
                  const aLocked = building.level < a.buildingLevelReq ? 1 : 0;
                  const bLocked = building.level < b.buildingLevelReq ? 1 : 0;
                  return aLocked - bLocked;
                })
                .map(([troopType, troopCfg]) => {
                  const meetsReq = building.level >= troopCfg.buildingLevelReq;
                  const isSelected = selectedTroop === troopType;
                  return (
                    <button
                      key={troopType}
                      type="button"
                      className={`${styles.troopCard} ${isSelected ? styles.selected : ''} ${!meetsReq ? styles.locked : ''}`}
                      onClick={() => meetsReq && setSelectedTroop(troopType)}
                      disabled={!meetsReq}
                      title={meetsReq ? troopCfg.displayName : `Requires ${displayName} level ${troopCfg.buildingLevelReq}`}
                    >
                      {!meetsReq && <span className={styles.lockOverlay}>🔒</span>}
                      <div className={styles.troopIcon}>
                        <GameIcon assetId={troopType} fallback="⚔️" size={48} />
                      </div>
                      <span className={styles.troopName}>{troopCfg.displayName}</span>
                      <div className={styles.troopStats}>
                        <span className={styles.stat}>⚔️ {troopCfg.attack}</span>
                        <span className={styles.stat}>🛡️ {troopCfg.defInfantry}</span>
                        <span className={styles.stat}>🏇 {troopCfg.defCavalry}</span>
                      </div>
                      {!meetsReq && (
                        <span className={styles.troopReq}>Lv {troopCfg.buildingLevelReq} required</span>
                      )}
                    </button>
                  );
                })}
            </div>
          )}
        </div>

        {/* Training controls */}
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

            <div className={`${styles.costGrid} ${costLoading ? styles.costGridLoading : ''}`}>
              {costLoading ? (
                <>
                  <div className={styles.costRow}><span className={styles.costIcon}>🌾</span><span className={styles.costLabel}>Food</span><span className={styles.costSkeleton} /></div>
                  <div className={styles.costRow}><span className={styles.costIcon}>💧</span><span className={styles.costLabel}>Water</span><span className={styles.costSkeleton} /></div>
                  <div className={styles.costRow}><span className={styles.costIcon}>🪵</span><span className={styles.costLabel}>Lumber</span><span className={styles.costSkeleton} /></div>
                  <div className={styles.costRow}><span className={styles.costIcon}>🪨</span><span className={styles.costLabel}>Stone</span><span className={styles.costSkeleton} /></div>
                  <div className={styles.costTime}>⏱ —</div>
                </>
              ) : costPreview ? (
                <>
                  {Math.ceil(costPreview.totalFood) > 0 && (
                    <div className={styles.costRow}>
                      <span className={styles.costIcon}>🌾</span>
                      <span className={styles.costLabel}>Food</span>
                      <span className={resources.food >= costPreview.totalFood ? styles.costOk : styles.costBad}>
                        {Math.ceil(costPreview.totalFood)}
                      </span>
                    </div>
                  )}
                  {Math.ceil(costPreview.totalWater) > 0 && (
                    <div className={styles.costRow}>
                      <span className={styles.costIcon}>💧</span>
                      <span className={styles.costLabel}>Water</span>
                      <span className={resources.water >= costPreview.totalWater ? styles.costOk : styles.costBad}>
                        {Math.ceil(costPreview.totalWater)}
                      </span>
                    </div>
                  )}
                  {Math.ceil(costPreview.totalLumber) > 0 && (
                    <div className={styles.costRow}>
                      <span className={styles.costIcon}>🪵</span>
                      <span className={styles.costLabel}>Lumber</span>
                      <span className={resources.lumber >= costPreview.totalLumber ? styles.costOk : styles.costBad}>
                        {Math.ceil(costPreview.totalLumber)}
                      </span>
                    </div>
                  )}
                  {Math.ceil(costPreview.totalStone) > 0 && (
                    <div className={styles.costRow}>
                      <span className={styles.costIcon}>🪨</span>
                      <span className={styles.costLabel}>Stone</span>
                      <span className={resources.stone >= costPreview.totalStone ? styles.costOk : styles.costBad}>
                        {Math.ceil(costPreview.totalStone)}
                      </span>
                    </div>
                  )}
                  <div className={styles.costTime}>
                    ⏱ {formatDuration(costPreview.eachTimeSec)}/unit · {formatDuration(costPreview.totalTimeSec)} total
                  </div>
                </>
              ) : null}
            </div>

            <button
              type="button"
              className={styles.trainBtn}
              onClick={handleTrain}
              disabled={!canAffordTraining || startTraining.isPending || costLoading}
            >
              {startTraining.isPending
                ? 'Training...'
                : !canAffordTraining && costPreview
                  ? 'Insufficient resources'
                  : 'Train'}
            </button>

            {startTraining.isError && (
              <p className={styles.errorMsg}>
                {startTraining.error instanceof Error
                  ? startTraining.error.message
                  : 'Failed to start training'}
              </p>
            )}
          </div>
        )}

        {/* ── Divider ─────────────────────────────────────────────────────── */}
        <div className={styles.divider} />

        {/* ── Upgrade section ─────────────────────────────────────────────── */}
        {isUpgrading ? (
          <div className={styles.upgradingBanner}>
            Upgrading to Lv {buildQueue.find((q) => q.building_type === type)?.target_level}...
          </div>
        ) : !isMaxLevel ? (
          <div className={styles.sectionCentered}>
            <h4 className={`${styles.sectionTitle} ${styles.upgradeSectionTitle}`}>
              Upgrade to Level {effectiveTarget}
            </h4>
            <div className={styles.costGridSmall}>
              <CostRow label="Food"   icon="🌾" value={cost!.food}   available={resources.food} />
              <CostRow label="Water"  icon="💧" value={cost!.water}  available={resources.water} />
              <CostRow label="Lumber" icon="🪵" value={cost!.lumber} available={resources.lumber} />
              <CostRow label="Stone"  icon="🪨" value={cost!.stone}  available={resources.stone} />
            </div>
            <div className={styles.buildTimeSmall}>
              <span>⏱</span>
              <span>{formatBuildDuration(timeSec)}</span>
            </div>

            {unmetPrereqs.length > 0 && (
              <div className={styles.prereqList}>
                {unmetPrereqs.map((p) => (
                  <span key={p.buildingType} className={styles.prereqUnmet}>
                    ✗ {p.displayName} Lv {p.minLevel}
                    <span className={styles.prereqCurrent}> (current: {p.currentLevel})</span>
                  </span>
                ))}
              </div>
            )}

            <button
              type="button"
              className={styles.upgradeBtn}
              disabled={!canUpgrade || upgradeMutation.isPending}
              onClick={handleUpgrade}
            >
              {upgradeMutation.isPending
                ? 'Starting...'
                : queueActive
                  ? 'Queue busy'
                  : !prereqs.allMet
                    ? 'Prerequisites not met'
                    : !canAffordUpgrade
                      ? 'Insufficient resources'
                      : 'Upgrade'}
            </button>

            {upgradeMutation.isError && (
              <p className={styles.errorMsg}>
                {(upgradeMutation.error as Error).message || 'Failed to start upgrade'}
              </p>
            )}
          </div>
        ) : (
          <div className={styles.maxLevelNote}>Building is at max level.</div>
        )}
      </div>
    </Modal>
  );
}

function CostRow({
  label,
  icon,
  value,
  available,
}: {
  label: string;
  icon: string;
  value: number;
  available: number;
}) {
  if (value === 0) return null;
  const enough = available >= value;
  return (
    <div className={styles.costRow}>
      <span className={styles.costIcon}>{icon}</span>
      <span className={styles.costLabel}>{label}</span>
      <span className={enough ? styles.costValueOk : styles.costValueBad}>{value}</span>
    </div>
  );
}
