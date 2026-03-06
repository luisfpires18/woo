import type { BuildingInfo, ResourcesResponse, BuildingQueueResponse } from '../../../types/api';
import type { BuildingType } from '../../../types/game';
import { Modal } from '../../../components/Modal';
import { GameIcon } from '../../../components/GameIcon/GameIcon';
import {
  BUILDING_CONFIGS,
  costAtLevel,
  timeAtLevel,
  checkPrerequisites,
  formatDuration,
} from '../../../config/buildings';
import { useStartUpgrade } from '../hooks/useBuildingUpgrade';
import styles from './BuildingDetailModal.module.css';

interface BuildingDetailModalProps {
  isOpen: boolean;
  onClose: () => void;
  building: BuildingInfo;
  allBuildings: BuildingInfo[];
  villageId: number;
  resources: ResourcesResponse;
  queue: BuildingQueueResponse[];
}

export function BuildingDetailModal({
  isOpen,
  onClose,
  building,
  allBuildings,
  villageId,
  resources,
  queue,
}: BuildingDetailModalProps) {
  const upgradeMutation = useStartUpgrade(villageId);
  const type = building.building_type as BuildingType;
  const cfg = BUILDING_CONFIGS[type];
  const displayName = cfg?.displayName ?? building.building_type;

  const isMaxLevel = building.level >= (cfg?.maxLevel ?? 0);
  const targetLevel = building.level + 1;

  const queueActive = queue.length > 0;
  const isUpgrading = queue.some((q) => q.building_type === building.building_type);

  // For level 0 buildings (not built yet), target level 1
  const effectiveTarget = building.level === 0 ? 1 : targetLevel;

  const cost = !isMaxLevel ? costAtLevel(type, effectiveTarget) : null;
  const timeSec = !isMaxLevel ? timeAtLevel(type, effectiveTarget) : 0;
  const prereqs = checkPrerequisites(type, allBuildings);

  const canAfford = cost
    ? resources.iron >= cost.iron &&
      resources.wood >= cost.wood &&
      resources.stone >= cost.stone &&
      resources.food >= cost.food
    : false;

  const canUpgrade = !isMaxLevel && prereqs.allMet && canAfford && !queueActive;

  const handleUpgrade = async () => {
    try {
      await upgradeMutation.mutateAsync(building.building_type);
      onClose();
    } catch {
      // Error handled via mutation state
    }
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} title={displayName} size="sm">
      <div className={styles.content}>
        <div className={styles.headerSection}>
          <GameIcon assetId={building.building_type} fallback="🏗️" size={40} />
          <div className={styles.headerInfo}>
            <span className={styles.currentLevel}>
              {building.level > 0 ? `Level ${building.level}` : 'Not built'}
            </span>
            {isMaxLevel && <span className={styles.maxBadge}>Max Level</span>}
          </div>
        </div>

        {isUpgrading && (
          <div className={styles.upgradingBanner}>
            Upgrading to Lv {queue.find((q) => q.building_type === building.building_type)?.target_level}...
          </div>
        )}

        {!isMaxLevel && !isUpgrading && (
          <>
            <div className={styles.section}>
              <h4 className={styles.sectionTitle}>
                {building.level === 0 ? 'Build to Level 1' : `Upgrade to Level ${effectiveTarget}`}
              </h4>
              <div className={styles.costGrid}>
                <CostRow label="Iron" icon="⛏️" value={cost!.iron} available={resources.iron} />
                <CostRow label="Wood" icon="🪵" value={cost!.wood} available={resources.wood} />
                <CostRow label="Stone" icon="🪨" value={cost!.stone} available={resources.stone} />
                <CostRow label="Food" icon="🌾" value={cost!.food} available={resources.food} />
              </div>
              <div className={styles.buildTime}>
                <span>⏱ Build time:</span>
                <span className={styles.timeValue}>{formatDuration(timeSec)}</span>
              </div>
            </div>

            {prereqs.checks.length > 0 && (
              <div className={styles.section}>
                <h4 className={styles.sectionTitle}>Prerequisites</h4>
                <div className={styles.prereqList}>
                  {prereqs.checks.map((p) => (
                    <span
                      key={p.buildingType}
                      className={p.met ? styles.prereqMet : styles.prereqUnmet}
                    >
                      {p.met ? '✓' : '✗'} {p.displayName} Lv {p.minLevel}
                      {!p.met && (
                        <span className={styles.prereqCurrent}> (current: {p.currentLevel})</span>
                      )}
                    </span>
                  ))}
                </div>
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
                    : !canAfford
                      ? 'Insufficient resources'
                      : building.level === 0
                        ? 'Build'
                        : 'Upgrade'}
            </button>
          </>
        )}

        {upgradeMutation.isError && (
          <p className={styles.error}>
            {(upgradeMutation.error as Error).message || 'Failed to start upgrade'}
          </p>
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
  const enough = available >= value;
  return (
    <div className={styles.costRow}>
      <span className={styles.costIcon}>{icon}</span>
      <span className={styles.costLabel}>{label}</span>
      <span className={enough ? styles.costValueOk : styles.costValueBad}>{value}</span>
    </div>
  );
}
