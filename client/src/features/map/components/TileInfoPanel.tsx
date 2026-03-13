// Tile info panel — shows details about the selected map tile (including camps)

import { useEffect, useState } from 'react';
import type { MapTile } from '../../../types/map';
import { TERRAIN_CONFIG } from '../../../types/map';
import { useExpeditionStore } from '../../../stores/expeditionStore';
import { fetchCamp } from '../../../services/camp';
import type { CampResponse } from '../../../types/api';
import styles from './TileInfoPanel.module.css';

interface TileInfoPanelProps {
  tile: MapTile | null;
  onClose: () => void;
  onAttackCamp?: (camp: CampResponse) => void;
}

/** Format zone name for display */
function formatZone(zone: string): string {
  if (!zone || zone === 'wilderness') return 'Wilderness';
  return zone
    .split('_')
    .map((w) => w.charAt(0).toUpperCase() + w.slice(1))
    .join(' ');
}

export function TileInfoPanel({ tile, onClose, onAttackCamp }: TileInfoPanelProps) {
  const camps = useExpeditionStore((s) => s.camps);
  const [campDetail, setCampDetail] = useState<CampResponse | null>(null);
  const [loadingCamp, setLoadingCamp] = useState(false);

  // Find camp on selected tile
  const tileCamp = tile ? camps.find((c) => c.tile_x === tile.x && c.tile_y === tile.y) : null;

  // Fetch full camp detail when a camp tile is selected
  useEffect(() => {
    if (!tileCamp) {
      setCampDetail(null);
      return;
    }
    setLoadingCamp(true);
    fetchCamp(tileCamp.id)
      .then(setCampDetail)
      .catch(() => setCampDetail(null))
      .finally(() => setLoadingCamp(false));
  }, [tileCamp?.id]); // eslint-disable-line react-hooks/exhaustive-deps

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

        {tileCamp && (
          <>
            <div className={styles.divider} />
            <div className={styles.row}>
              <span className={styles.label}>Camp</span>
              <span className={styles.value}>
                {campDetail?.template_name ?? tileCamp.template_name ?? 'Loading...'}
              </span>
            </div>
            <div className={styles.row}>
              <span className={styles.label}>Tier</span>
              <span className={styles.value}>{tileCamp.tier}</span>
            </div>
            <div className={styles.row}>
              <span className={styles.label}>Status</span>
              <span className={styles.value} style={{
                color: tileCamp.status === 'active' ? '#cc2222' : '#888888',
                fontWeight: 'bold',
              }}>
                {tileCamp.status.charAt(0).toUpperCase() + tileCamp.status.slice(1).replace('_', ' ')}
              </span>
            </div>

            {loadingCamp && (
              <div className={styles.row}>
                <span className={styles.label}>Beasts</span>
                <span className={styles.value}>Loading...</span>
              </div>
            )}

            {campDetail && campDetail.beasts.length > 0 && (
              <>
                <div className={styles.divider} />
                <div className={styles.row}>
                  <span className={styles.label}>Defenders</span>
                  <span className={styles.value}>{campDetail.beasts.reduce((s, b) => s + b.count, 0)} total</span>
                </div>
                {campDetail.beasts.map((b, i) => (
                  <div className={styles.row} key={i}>
                    <span className={styles.label} style={{ paddingLeft: 8 }}>{b.name}</span>
                    <span className={styles.value}>
                      ×{b.count} — HP {b.hp}, ATK {b.attack_power}
                    </span>
                  </div>
                ))}
              </>
            )}

            {tileCamp.status === 'active' && onAttackCamp && campDetail && (
              <button
                className={styles.attackBtn}
                onClick={() => onAttackCamp(campDetail)}
              >
                ⚔ Attack Camp
              </button>
            )}
          </>
        )}
      </div>
    </div>
  );
}
