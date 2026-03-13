// World Map page — displays the Canvas 2D map renderer with tile info, camps, and expeditions

import { useCallback, useEffect, useRef, useState } from 'react';
import { MapRenderer } from '../components/MapRenderer';
import type { MapRendererHandle } from '../components/MapRenderer';
import { TileInfoPanel } from '../components/TileInfoPanel';
import { DispatchExpeditionModal } from '../components/DispatchExpeditionModal';
import { ExpeditionPanel } from '../components/ExpeditionPanel';
import { BattleReportModal } from '../components/BattleReportModal';
import { useMapStore } from '../../../stores/mapStore';
import { useGameStore } from '../../../stores/gameStore';
import { useExpeditionStore } from '../../../stores/expeditionStore';
import { fetchCamps, fetchExpeditions } from '../../../services/camp';
import type { MapTile } from '../../../types/map';
import type { CampResponse } from '../../../types/api';
import styles from './WorldMapPage.module.css';

export function WorldMapPage() {
  const containerRef = useRef<HTMLDivElement>(null);
  const mapRendererRef = useRef<MapRendererHandle>(null);
  const [dimensions, setDimensions] = useState({ width: 800, height: 600 });
  const [hoverCoords, setHoverCoords] = useState<{ x: number; y: number } | null>(null);
  const [gotoInput, setGotoInput] = useState('');
  const selectedTile = useMapStore((s) => s.selectedTile);
  const selectTile = useMapStore((s) => s.selectTile);
  const loading = useMapStore((s) => s.loading);
  const villages = useGameStore((s) => s.villages);

  // Camp & expedition state
  const setCamps = useExpeditionStore((s) => s.setCamps);
  const setExpeditions = useExpeditionStore((s) => s.setExpeditions);
  const campsLoaded = useExpeditionStore((s) => s.campsLoaded);
  const [attackCamp, setAttackCamp] = useState<CampResponse | null>(null);
  const [viewBattleId, setViewBattleId] = useState<number | null>(null);
  const [expeditionPanelOpen, setExpeditionPanelOpen] = useState(true);

  // Get the player's first village as the initial center
  const firstVillage = villages[0];
  const initialX = firstVillage?.x ?? 0;
  const initialY = firstVillage?.y ?? 0;

  // Resize the canvas to fill the container
  useEffect(() => {
    const updateSize = () => {
      if (containerRef.current) {
        const rect = containerRef.current.getBoundingClientRect();
        setDimensions({
          width: Math.floor(rect.width),
          height: Math.floor(rect.height),
        });
      }
    };

    updateSize();
    window.addEventListener('resize', updateSize);
    return () => window.removeEventListener('resize', updateSize);
  }, []);

  // Load camps and expeditions when the map page mounts
  useEffect(() => {
    if (!campsLoaded) {
      fetchCamps().then((c) => setCamps(c ?? [])).catch(() => {});
    }
    fetchExpeditions().then((e) => setExpeditions(e ?? [])).catch(() => {});
    // Refresh expeditions every 5s, camps every 30s
    let campTick = 0;
    const interval = setInterval(() => {
      fetchExpeditions().then((e) => setExpeditions(e ?? [])).catch(() => {});
      campTick++;
      if (campTick % 6 === 0) {
        fetchCamps().then((c) => setCamps(c ?? [])).catch(() => {});
      }
    }, 5_000);
    return () => clearInterval(interval);
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  const handleTileClick = useCallback(
    (tile: MapTile) => {
      selectTile(tile);
    },
    [selectTile],
  );

  const handleClosePanel = useCallback(() => {
    selectTile(null);
  }, [selectTile]);

  const handleTileHover = useCallback((x: number, y: number) => {
    setHoverCoords({ x, y });
  }, []);

  const handleAttackCamp = useCallback((camp: CampResponse) => {
    setAttackCamp(camp);
  }, []);

  const handleCloseDispatch = useCallback(() => {
    setAttackCamp(null);
  }, []);

  const handleViewReport = useCallback((battleId: number) => {
    setViewBattleId(battleId);
  }, []);

  const handleCloseReport = useCallback(() => {
    setViewBattleId(null);
  }, []);

  const handleGotoSubmit = useCallback(
    (e: React.FormEvent) => {
      e.preventDefault();
      // Parse "x, y" or "x y" format
      const match = gotoInput.match(/(-?\d+)\s*[,\s]\s*(-?\d+)/);
      if (match && match[1] && match[2]) {
        const gx = parseInt(match[1], 10);
        const gy = parseInt(match[2], 10);
        mapRendererRef.current?.navigateTo(gx, gy);
        setGotoInput('');
      }
    },
    [gotoInput],
  );

  return (
    <div className={styles.page}>
      <div className={styles.header}>
        <h1 className={styles.title}>World Map</h1>
        <div className={styles.coords}>
          {hoverCoords
            ? `Hover: (${hoverCoords.x}, ${hoverCoords.y})`
            : selectedTile
              ? `Selected: (${selectedTile.x}, ${selectedTile.y})`
              : `Center: (${initialX}, ${initialY})`}
        </div>
        <form className={styles.gotoForm} onSubmit={handleGotoSubmit}>
          <input
            className={styles.gotoInput}
            type="text"
            value={gotoInput}
            onChange={(e) => setGotoInput(e.target.value)}
            placeholder="Go to x, y"
            aria-label="Go to coordinates"
          />
          <button className={styles.gotoBtn} type="submit">Go</button>
        </form>
        {loading && <div className={styles.loading}>Loading...</div>}
      </div>

      <div className={styles.mapContainer} ref={containerRef}>
        {dimensions.width > 0 && dimensions.height > 0 && (
          <MapRenderer
            ref={mapRendererRef}
            initialX={initialX}
            initialY={initialY}
            width={dimensions.width}
            height={dimensions.height}
            selectedTile={selectedTile}
            onTileClick={handleTileClick}
            onTileHover={handleTileHover}
          />
        )}

        <TileInfoPanel tile={selectedTile} onClose={handleClosePanel} onAttackCamp={handleAttackCamp} />
        {expeditionPanelOpen ? (
          <ExpeditionPanel onViewReport={handleViewReport} onClose={() => setExpeditionPanelOpen(false)} />
        ) : (
          <button className={styles.reopenExpeditions} onClick={() => setExpeditionPanelOpen(true)}>
            ⚔ Expeditions
          </button>
        )}
      </div>

      {attackCamp && (
        <DispatchExpeditionModal camp={attackCamp} onClose={handleCloseDispatch} />
      )}

      {viewBattleId !== null && (
        <BattleReportModal battleId={viewBattleId} onClose={handleCloseReport} />
      )}

      <div className={styles.legend}>
        <span className={styles.legendItem}>
          <span className={styles.swatch} style={{ backgroundColor: '#7ec850' }} /> Plains
        </span>
        <span className={styles.legendItem}>
          <span className={styles.swatch} style={{ backgroundColor: '#2d7a3a' }} /> Forest
        </span>
        <span className={styles.legendItem}>
          <span className={styles.swatch} style={{ backgroundColor: '#8b7355' }} /> Mountain
        </span>
        <span className={styles.legendItem}>
          <span className={styles.swatch} style={{ backgroundColor: '#3a7ec8' }} /> Water
        </span>
        <span className={styles.legendItem}>
          <span className={styles.swatch} style={{ backgroundColor: '#d4a843' }} /> Desert
        </span>
        <span className={styles.legendItem}>
          <span className={styles.swatch} style={{ backgroundColor: '#5a6e3a' }} /> Swamp
        </span>
        <span className={styles.legendItem}>
          <span className={styles.swatch} style={{ backgroundColor: '#1a0a2e' }} /> Chasm
        </span>
        <span className={styles.legendItem}>
          <span className={styles.swatch} style={{ backgroundColor: '#8b6914' }} /> Bridge
        </span>
      </div>
    </div>
  );
}
