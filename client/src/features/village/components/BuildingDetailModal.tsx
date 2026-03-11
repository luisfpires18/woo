import { useState } from 'react';
import type { BuildingInfo, ResourcesResponse, BuildingQueueResponse } from '../../../types/api';
import { Modal } from '../../../components/Modal';
import { GameIcon } from '../../../components/GameIcon/GameIcon';
import {
  BUILDING_CONFIGS,
  RESOURCE_BUILDING_TYPES,
  costAtLevel,
  timeAtLevel,
  checkPrerequisites,
  formatDuration,
} from '../../../config/buildings';
import { useBuildingDisplayNames } from '../../../hooks/useBuildingDisplayNames';
import { useResourceBuildingDisplay } from '../../../hooks/useResourceBuildingDisplay';
import { useStartUpgrade } from '../hooks/useBuildingUpgrade';
import { useGameStore } from '../../../stores/gameStore';
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
  const { getDisplayName } = useBuildingDisplayNames();
  const { getDisplay } = useResourceBuildingDisplay();
  const [spriteFailed, setSpriteFailed] = useState(false);
  const playerGold = useGameStore((s) => s.playerGold);
  const type = building.building_type;
  const cfg = BUILDING_CONFIGS[type];
  const isResourceBuilding = RESOURCE_BUILDING_TYPES.has(type);
  const resourceDisplay = isResourceBuilding ? getDisplay(type) : null;
  const displayName = isResourceBuilding ? resourceDisplay!.displayName : getDisplayName(type);

  const isMaxLevel = building.level >= (cfg?.maxLevel ?? 0);
  const targetLevel = building.level + 1;

  const queueActive = queue.length > 0;
  const isUpgrading = queue.some((q) => q.building_type === building.building_type);

  // For level 0 buildings (not built yet), target level 1
  const effectiveTarget = building.level === 0 ? 1 : targetLevel;

  const cost = !isMaxLevel ? costAtLevel(type, effectiveTarget) : null;
  const timeSec = !isMaxLevel ? timeAtLevel(type, effectiveTarget) : 0;
  const prereqs = checkPrerequisites(type, allBuildings, getDisplayName);

  const canAfford = cost
    ? resources.food >= cost.food &&
      resources.water >= cost.water &&
      resources.lumber >= cost.lumber &&
      resources.stone >= cost.stone &&
      playerGold >= cost.gold
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
          {isResourceBuilding && resourceDisplay?.spriteUrl && !spriteFailed ? (
            <img
              src={resourceDisplay.spriteUrl}
              alt={displayName}
              className={styles.headerSprite}
              onError={() => setSpriteFailed(true)}
              draggable={false}
            />
          ) : (
            <GameIcon assetId={building.building_type} fallback={resourceDisplay?.emoji ?? '🏗️'} size={64} />
          )}
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
                <CostRow label="Gold" icon="🪙" value={cost!.gold} available={playerGold} />
                <CostRow label="Food" icon="🌾" value={cost!.food} available={resources.food} />
                <CostRow label="Water" icon="💧" value={cost!.water} available={resources.water} />
                <CostRow label="Lumber" icon="🪵" value={cost!.lumber} available={resources.lumber} />
                <CostRow label="Stone" icon="🪨" value={cost!.stone} available={resources.stone} />
              </div>
              <div className={styles.buildTime}>
                <span>⏱ Build time:</span>
                <span className={styles.timeValue}>{formatDuration(timeSec)}</span>
              </div>
            </div>

            {prereqs.checks.some((p) => !p.met) && (
              <div className={styles.section}>
                <h4 className={styles.sectionTitle}>Prerequisites</h4>
                <div className={styles.prereqList}>
                  {prereqs.checks.filter((p) => !p.met).map((p) => (
                    <span key={p.buildingType} className={styles.prereqUnmet}>
                      ✗ {p.displayName} Lv {p.minLevel}
                      <span className={styles.prereqCurrent}> (current: {p.currentLevel})</span>
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
