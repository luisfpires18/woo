// Tile info panel — shows details about the selected map tile

import type { MapTile } from '../../../types/map';
import { TERRAIN_CONFIG } from '../../../types/map';
import styles from './TileInfoPanel.module.css';

interface TileInfoPanelProps {
  tile: MapTile | null;
  onClose: () => void;
}

/** Format zone name for display */
function formatZone(zone: string): string {
  if (!zone || zone === 'wilderness') return 'Wilderness';
  return zone
    .split('_')
    .map((w) => w.charAt(0).toUpperCase() + w.slice(1))
    .join(' ');
}

export function TileInfoPanel({ tile, onClose }: TileInfoPanelProps) {
  if (!tile) return null;

  const terrain = TERRAIN_CONFIG[tile.terrain] ?? TERRAIN_CONFIG.plains;

  return (
    <div className={styles.panel}>
      <div className={styles.header}>
        <h3 className={styles.title}>
          Tile ({tile.x}, {tile.y})
        </h3>
        <button className={styles.closeBtn} onClick={onClose} aria-label="Close">
          ×
        </button>
      </div>

      <div className={styles.details}>
        <div className={styles.row}>
          <span className={styles.label}>Terrain</span>
          <span className={styles.value}>
            <span
              className={styles.terrainSwatch}
              style={{ backgroundColor: `#${terrain.color.toString(16).padStart(6, '0')}` }}
            />
            {terrain.label}
          </span>
        </div>

        <div className={styles.row}>
          <span className={styles.label}>Zone</span>
          <span className={styles.value}>{formatZone(tile.zone)}</span>
        </div>

        {terrain.passable && (
          <div className={styles.row}>
            <span className={styles.label}>Movement</span>
            <span className={styles.value}>{terrain.movementMod}×</span>
          </div>
        )}

        {!terrain.passable && (
          <div className={styles.row}>
            <span className={styles.label}>Movement</span>
            <span className={styles.impassable}>Impassable</span>
          </div>
        )}

        {tile.village_name && (
          <>
            <div className={styles.divider} />
            <div className={styles.row}>
              <span className={styles.label}>Village</span>
              <span className={styles.value}>{tile.village_name}</span>
            </div>
            {tile.owner_name && (
              <div className={styles.row}>
                <span className={styles.label}>Owner</span>
                <span className={styles.value}>{tile.owner_name}</span>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
}
