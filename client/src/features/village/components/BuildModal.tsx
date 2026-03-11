import type { BuildingInfo, ResourcesResponse, BuildingQueueResponse } from '../../../types/api';
import type { BuildingType } from '../../../types/game';
import { Modal } from '../../../components/Modal';
import { GameIcon } from '../../../components/GameIcon/GameIcon';
import {
  BUILDING_CONFIGS,
  RESOURCE_BUILDING_TYPES,
  costAtLevel,
  timeAtLevel,
  checkPrerequisites,
  formatDuration,
  type PrerequisiteCheck,
} from '../../../config/buildings';
import { useBuildingDisplayNames } from '../../../hooks/useBuildingDisplayNames';
import { useStartUpgrade } from '../hooks/useBuildingUpgrade';
import styles from './BuildModal.module.css';

interface BuildModalProps {
  isOpen: boolean;
  onClose: () => void;
  buildings: BuildingInfo[];
  villageId: number;
  resources: ResourcesResponse;

  queue: BuildingQueueResponse[];
}

interface BuildOption {
  building: BuildingInfo;
  type: BuildingType;
  displayName: string;
  cost: { food: number; water: number; lumber: number; stone: number };
  timeSec: number;
  prereqs: { allMet: boolean; checks: PrerequisiteCheck[] };
  canAfford: boolean;
}

export function BuildModal({
  isOpen,
  onClose,
  buildings,
  villageId,
  resources,
  queue,
}: BuildModalProps) {
  const upgradeMutation = useStartUpgrade(villageId);
  const { getDisplayName } = useBuildingDisplayNames();
  const queueActive = queue.length > 0;

  // Get all village buildings at level 0, excluding resource fields and wrong-kingdom buildings
  const options: BuildOption[] = buildings
    .filter((b) => {
      if (b.level > 0) return false;
      if (RESOURCE_BUILDING_TYPES.has(b.building_type)) return false;
      const cfg = BUILDING_CONFIGS[b.building_type];
      if (!cfg) return false;
      return true;
    })
    .map((b) => {
      const type = b.building_type;
      const cost = costAtLevel(type, 1);
      const prereqs = checkPrerequisites(type, buildings, getDisplayName);
      return {
        building: b,
        type,
        displayName: getDisplayName(type),
        cost,
        timeSec: timeAtLevel(type, 1),
        prereqs,
        canAfford:
          resources.food >= cost.food &&
          resources.water >= cost.water &&
          resources.lumber >= cost.lumber &&
          resources.stone >= cost.stone,
      };
    });

  const available = options.filter((o) => o.prereqs.allMet);
  const locked = options.filter((o) => !o.prereqs.allMet);

  const handleBuild = async (buildingType: string) => {
    try {
      await upgradeMutation.mutateAsync(buildingType);
      onClose();
    } catch {
      // Error handling is done via mutation state
    }
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="Build New Building" size="md">
      {options.length === 0 ? (
        <p className={styles.emptyMsg}>All available buildings have been constructed.</p>
      ) : (
        <div className={styles.sections}>
          {/* ── Available Buildings ── */}
          {available.length > 0 && (
            <section>
              <h3 className={styles.sectionTitle}>
                <span className={styles.sectionIcon}>✓</span> Available Buildings
              </h3>
              <div className={styles.list}>
                {available.map((o) => (
                  <div key={o.building.id} className={styles.row}>
                    <div className={styles.rowHeader}>
                      <GameIcon assetId={o.building.building_type} fallback="🏗️" size={24} />
                      <span className={styles.buildingName}>{o.displayName}</span>
                    </div>

                    <div className={styles.costs}>
                      <CostItem label="Food" value={o.cost.food} available={resources.food} />
                      <CostItem label="Water" value={o.cost.water} available={resources.water} />
                      <CostItem label="Lumber" value={o.cost.lumber} available={resources.lumber} />
                      <CostItem label="Stone" value={o.cost.stone} available={resources.stone} />
                      <span className={styles.timeValue}>⏱ {formatDuration(o.timeSec)}</span>
                    </div>

                    <button
                      type="button"
                      className={styles.buildBtn}
                      disabled={!o.canAfford || queueActive || upgradeMutation.isPending}
                      onClick={() => handleBuild(o.building.building_type)}
                    >
                      {queueActive ? 'Queue busy' : !o.canAfford ? 'Not enough resources' : 'Build'}
                    </button>
                  </div>
                ))}
              </div>
            </section>
          )}

          {/* ── Locked Buildings ── */}
          {locked.length > 0 && (
            <section>
              <h3 className={`${styles.sectionTitle} ${styles.lockedTitle}`}>
                <span className={styles.sectionIcon}>🔒</span> Locked Buildings
              </h3>
              <div className={styles.list}>
                {locked.map((o) => (
                  <div key={o.building.id} className={`${styles.row} ${styles.lockedRow}`}>
                    <div className={styles.rowHeader}>
                      <GameIcon assetId={o.building.building_type} fallback="🏗️" size={24} />
                      <span className={styles.buildingName}>{o.displayName}</span>
                    </div>

                    <div className={styles.prereqs}>
                      {o.prereqs.checks.map((p) => (
                        <span
                          key={p.buildingType}
                          className={p.met ? styles.prereqMet : styles.prereqUnmet}
                        >
                          {p.met ? '✓' : '✗'} {p.displayName} Lv {p.minLevel}
                          {!p.met && (
                            <span className={styles.prereqCurrent}>
                              {' '}(yours: {p.currentLevel})
                            </span>
                          )}
                        </span>
                      ))}
                    </div>
                  </div>
                ))}
              </div>
            </section>
          )}
        </div>
      )}

      {upgradeMutation.isError && (
        <p className={styles.error}>
          {(upgradeMutation.error as Error).message || 'Failed to start construction'}
        </p>
      )}
    </Modal>
  );
}

function CostItem({ label, value, available }: { label: string; value: number; available: number }) {
  if (value === 0) return null;
  const enough = available >= value;
  return (
    <span className={enough ? styles.costOk : styles.costInsufficient}>
      {label}: {value}
    </span>
  );
}
