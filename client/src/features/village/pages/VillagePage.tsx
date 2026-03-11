import { useState, useRef, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { useQueryClient } from '@tanstack/react-query';
import { useVillage } from '../hooks/useVillage';
import { useAuthStore } from '../../../stores/authStore';
import { renameVillage } from '../../../services/village';
import { BuildingGrid } from '../components/BuildingGrid';
import { ResourcePanel } from '../components/ResourcePanel';
import { ResourceFieldsGrid } from '../components/ResourceFieldsGrid';
import { ConstructionQueue } from '../components/ConstructionQueue';
import { TrainingQueue } from '../components/TrainingQueue';
import { TroopRoster } from '../components/TroopRoster';
import { BuildModal } from '../components/BuildModal';
import { BuildingDetailModal } from '../components/BuildingDetailModal';
import { BuildingMilitaryModal } from '../components/BuildingMilitaryModal';
import { isMilitaryBuilding } from '../../../config/troops';
import { LoadingSpinner } from '../../../components/LoadingSpinner/LoadingSpinner';
import type { BuildingInfo } from '../../../types/api';
import styles from './VillagePage.module.css';

export function VillagePage() {
  const { id } = useParams<{ id: string }>();
  const villageId = Number(id);
  const player = useAuthStore((s) => s.player);

  const { data: village, isLoading, error } = useVillage(villageId);

  // Modal state
  const [buildModalOpen, setBuildModalOpen] = useState(false);
  const [selectedBuilding, setSelectedBuilding] = useState<BuildingInfo | null>(null);
  const [selectedMilitary, setSelectedMilitary] = useState<BuildingInfo | null>(null);

  // Rename state
  const [isRenaming, setIsRenaming] = useState(false);
  const [renameValue, setRenameValue] = useState('');
  const [renameError, setRenameError] = useState<string | null>(null);
  const [renameSaving, setRenameSaving] = useState(false);
  const renameInputRef = useRef<HTMLInputElement>(null);
  const queryClient = useQueryClient();

  useEffect(() => {
    if (isRenaming && renameInputRef.current) {
      renameInputRef.current.focus();
      renameInputRef.current.select();
    }
  }, [isRenaming]);

  const handleStartRename = () => {
    if (!village) return;
    setRenameValue(village.name);
    setRenameError(null);
    setIsRenaming(true);
  };

  const handleCancelRename = () => {
    setIsRenaming(false);
    setRenameError(null);
  };

  const handleSaveRename = async () => {
    const trimmed = renameValue.trim();
    if (trimmed.length < 2 || trimmed.length > 30) {
      setRenameError('Name must be 2–30 characters.');
      return;
    }
    setRenameSaving(true);
    setRenameError(null);
    try {
      await renameVillage(villageId, trimmed);
      await queryClient.invalidateQueries({ queryKey: ['village', villageId] });
      setIsRenaming(false);
    } catch {
      setRenameError('Failed to rename village.');
    } finally {
      setRenameSaving(false);
    }
  };

  const handleRenameKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') handleSaveRename();
    if (e.key === 'Escape') handleCancelRename();
  };

  if (isLoading) {
    return (
      <div className={styles.loading}>
        <LoadingSpinner size="lg" />
      </div>
    );
  }

  if (error || !village) {
    return (
      <div className={styles.error}>
        <p>Failed to load village.</p>
      </div>
    );
  }

  const buildQueue = village.build_queue ?? [];
  const trainingQueue = village.training_queue ?? [];
  const troops = village.troops ?? [];

  const handleBuildingClick = (b: BuildingInfo) => {
    if (isMilitaryBuilding(b.building_type)) {
      setSelectedMilitary(b);
    } else {
      setSelectedBuilding(b);
    }
  };

  return (
    <div className={styles.page}>
      <header className={styles.header}>
        {isRenaming ? (
          <div className={styles.renameRow}>
            <input
              ref={renameInputRef}
              className={styles.renameInput}
              value={renameValue}
              onChange={(e) => setRenameValue(e.target.value)}
              onKeyDown={handleRenameKeyDown}
              maxLength={30}
              disabled={renameSaving}
            />
            <button
              className={styles.renameSave}
              onClick={handleSaveRename}
              disabled={renameSaving}
              aria-label="Confirm rename"
            >
              ✓
            </button>
            <button
              className={styles.renameCancel}
              onClick={handleCancelRename}
              disabled={renameSaving}
              aria-label="Cancel rename"
            >
              ✗
            </button>
            {renameError && <span className={styles.renameError}>{renameError}</span>}
          </div>
        ) : (
          <div className={styles.titleRow}>
            <h1 className={styles.title}>{village.name}</h1>
            <button
              className={styles.renameBtn}
              onClick={handleStartRename}
              title="Rename village"
              aria-label="Rename village"
            >
              ✏️
            </button>
          </div>
        )}
        <span className={styles.coords}>
          ({village.x}, {village.y})
        </span>
      </header>

      <ConstructionQueue queue={buildQueue} villageId={villageId} />
      <TrainingQueue queue={trainingQueue} villageId={villageId} />

      <div className={styles.content}>
        <div className={styles.main}>
          <section className={styles.section}>
            <h2 className={styles.sectionTitle}>Resource Fields</h2>
            <ResourceFieldsGrid
              buildings={village.buildings}
              onBuildingClick={(b) => handleBuildingClick(b)}
            />
          </section>

          <section className={styles.section}>
            <h2 className={styles.sectionTitle}>Buildings</h2>
            <BuildingGrid
              buildings={village.buildings}
              onBuildingClick={(b) => handleBuildingClick(b)}
              onEmptySlotClick={() => setBuildModalOpen(true)}
            />
          </section>
        </div>

        <aside className={styles.sidebar}>
          <ResourcePanel resources={village.resources} />
          <TroopRoster troops={troops} />
        </aside>
      </div>

      {/* Build new building modal */}
      <BuildModal
        isOpen={buildModalOpen}
        onClose={() => setBuildModalOpen(false)}
        buildings={village.buildings}
        villageId={villageId}
        resources={village.resources}
        queue={buildQueue}
      />

      {/* Building detail / upgrade modal */}
      {selectedBuilding && (
        <BuildingDetailModal
          isOpen={true}
          onClose={() => setSelectedBuilding(null)}
          building={selectedBuilding}
          allBuildings={village.buildings}
          villageId={villageId}
          resources={village.resources}
          queue={buildQueue}
        />
      )}

      {/* Military building modal (upgrade + training) */}
      {selectedMilitary && (
        <BuildingMilitaryModal
          isOpen={true}
          onClose={() => setSelectedMilitary(null)}
          building={selectedMilitary}
          allBuildings={village.buildings}
          villageId={villageId}
          resources={village.resources}
          kingdom={player?.kingdom ?? ''}
          buildQueue={buildQueue}
          trainingQueue={trainingQueue}
        />
      )}
    </div>
  );
}
