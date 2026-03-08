import { useState, useEffect, useCallback, useRef } from 'react';
import {
  listTemplates,
  getTemplate,
  createTemplate,
  deleteTemplate,
  updateTemplateTerrain,
  updateTemplateZones,
  applyTemplate,
  exportTemplate,
  importTemplate,
  resizeTemplate,
} from '../../../services/template';
import type { TemplateInfo, TemplateTile, TerrainType } from '../../../types/map';
import { TERRAIN_CONFIG, ZONE_CONFIG, KINGDOM_ZONES } from '../../../types/map';
import { TILE_SIZE, hexColor, hexColorAlpha, screenToTile } from '../../map/mapUtils';
import styles from './AdminMapEditorPage.module.css';

/** All paintable terrain types */
const TERRAIN_TYPES: TerrainType[] = ['plains', 'forest', 'mountain', 'water', 'desert', 'swamp'];

/** Paint mode toggle */
type PaintMode = 'terrain' | 'zone';

export function AdminMapEditorPage() {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const offsetRef = useRef({ x: 0, y: 0 });
  const scaleRef = useRef(0.6);
  const draggingRef = useRef(false);
  const dragStartRef = useRef({ x: 0, y: 0 });
  const lastOffsetRef = useRef({ x: 0, y: 0 });
  const rafRef = useRef(0);
  const hoverTileRef = useRef<{ x: number; y: number } | null>(null);

  // Template state
  const [templates, setTemplates] = useState<TemplateInfo[]>([]);
  const [activeTemplate, setActiveTemplate] = useState<string | null>(null);
  const [tiles, setTiles] = useState<Map<string, TemplateTile>>(new Map());
  const [mapSize, setMapSize] = useState(51);
  const mapHalf = Math.floor((mapSize - 1) / 2);

  // Paint state
  const [paintMode, setPaintMode] = useState<PaintMode>('terrain');
  const [selectedTerrain, setSelectedTerrain] = useState<TerrainType>('plains');
  const [selectedZone, setSelectedZone] = useState<string>('wilderness');
  const [brushSize, setBrushSize] = useState(1);
  const [pendingTerrain, setPendingTerrain] = useState<Map<string, { x: number; y: number; terrain_type: string }>>(new Map());
  const [pendingZones, setPendingZones] = useState<Map<string, { x: number; y: number; kingdom_zone: string }>>(new Map());

  // UI state
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [status, setStatus] = useState<string | null>(null);
  const [showNewDialog, setShowNewDialog] = useState(false);
  const [newName, setNewName] = useState('');
  const [newDesc, setNewDesc] = useState('');
  const [newSize, setNewSize] = useState(51);
  const [loadingTemplates, setLoadingTemplates] = useState(true);
  const [loadingTiles, setLoadingTiles] = useState(false);

  // Refs for draw callback (avoids stale closure)
  const tilesRef = useRef(tiles);
  tilesRef.current = tiles;
  const pendingTerrainRef = useRef(pendingTerrain);
  pendingTerrainRef.current = pendingTerrain;
  const pendingZonesRef = useRef(pendingZones);
  pendingZonesRef.current = pendingZones;
  const paintModeRef = useRef(paintMode);
  paintModeRef.current = paintMode;

  const canvasWidth = 900;
  const canvasHeight = 600;

  // ── Load template list on mount ──────────────────────────────────────
  useEffect(() => {
    loadTemplateList();
  }, []);

  const loadTemplateList = async () => {
    setLoadingTemplates(true);
    try {
      const list = await listTemplates();
      setTemplates(list ?? []);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load templates');
    } finally {
      setLoadingTemplates(false);
    }
  };

  // ── Load a template's tiles ──────────────────────────────────────────
  const loadTemplate = async (name: string) => {
    setLoadingTiles(true);
    setError(null);
    setStatus(null);
    setPendingTerrain(new Map());
    setPendingZones(new Map());
    try {
      const tmpl = await getTemplate(name);
      const tileMap = new Map<string, TemplateTile>();
      for (const t of tmpl.tiles) {
        tileMap.set(`${t.x},${t.y}`, t);
      }
      setTiles(tileMap);
      setActiveTemplate(name);
      setMapSize(tmpl.map_size || 51);
      setStatus(`Loaded template "${name}" — ${tmpl.map_size}×${tmpl.map_size} (${tmpl.tiles.length} tiles)`);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load template');
    } finally {
      setLoadingTiles(false);
    }
  };

  // Center on (0,0)
  useEffect(() => {
    offsetRef.current = { x: canvasWidth / 2, y: canvasHeight / 2 };
    scaleRef.current = 0.6;
  }, []);

  // ── Canvas draw function ─────────────────────────────────────────────
  const drawMap = useCallback(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    const dpr = window.devicePixelRatio || 1;
    const scale = scaleRef.current;
    const ox = offsetRef.current.x;
    const oy = offsetRef.current.y;

    ctx.setTransform(dpr, 0, 0, dpr, 0, 0);
    ctx.fillStyle = '#111122';
    ctx.fillRect(0, 0, canvasWidth, canvasHeight);

    ctx.setTransform(dpr * scale, 0, 0, dpr * scale, dpr * ox, dpr * oy);

    const invScale = 1 / scale;
    const worldLeft = -ox * invScale;
    const worldTop = -oy * invScale;
    const worldRight = worldLeft + canvasWidth * invScale;
    const worldBottom = worldTop + canvasHeight * invScale;

    const currentTiles = tilesRef.current;
    const pendingT = pendingTerrainRef.current;
    const pendingZ = pendingZonesRef.current;
    const mode = paintModeRef.current;

    // Draw tiles
    currentTiles.forEach((tile) => {
      const px = tile.x * TILE_SIZE;
      const py = -tile.y * TILE_SIZE;

      if (px + TILE_SIZE < worldLeft || px > worldRight) return;
      if (py + TILE_SIZE < worldTop || py > worldBottom) return;

      const key = `${tile.x},${tile.y}`;

      // Resolve effective terrain and zone (with pending changes)
      const pTerrain = pendingT.get(key);
      const pZone = pendingZ.get(key);
      const terrain = pTerrain ? pTerrain.terrain_type : tile.terrain_type;
      const zone = pZone ? pZone.kingdom_zone : tile.kingdom_zone;

      // Draw terrain base
      const cfg = TERRAIN_CONFIG[terrain as TerrainType] ?? TERRAIN_CONFIG.plains;
      ctx.fillStyle = hexColor(cfg.color);
      ctx.fillRect(px, py, TILE_SIZE, TILE_SIZE);

      // Draw zone overlay (semi-transparent)
      if (zone && zone !== 'wilderness') {
        const zoneCfg = ZONE_CONFIG[zone];
        if (zoneCfg) {
          ctx.fillStyle = hexColorAlpha(zoneCfg.color, 0.35);
          ctx.fillRect(px, py, TILE_SIZE, TILE_SIZE);
        }
      }

      // Grid lines
      ctx.strokeStyle = 'rgba(255,255,255,0.08)';
      ctx.lineWidth = 0.5;
      ctx.strokeRect(px, py, TILE_SIZE, TILE_SIZE);

      // Mark pending changes with a yellow dot
      const hasPending = (mode === 'terrain' && pTerrain) || (mode === 'zone' && pZone);
      if (hasPending) {
        ctx.fillStyle = 'rgba(255, 204, 0, 0.6)';
        ctx.beginPath();
        ctx.arc(px + TILE_SIZE - 10, py + 10, 5, 0, Math.PI * 2);
        ctx.fill();
      }

      // Zone label in small text when zoomed in enough
      if (zone && zone !== 'wilderness' && scale > 0.4) {
        ctx.font = '10px sans-serif';
        ctx.textAlign = 'center';
        ctx.textBaseline = 'bottom';
        ctx.fillStyle = 'rgba(255,255,255,0.5)';
        ctx.fillText(zone, px + TILE_SIZE / 2, py + TILE_SIZE - 2);
      }
    });

    // Hover highlight with brush size indicator
    const hover = hoverTileRef.current;
    if (hover) {
      const half = Math.floor(brushSize / 2);
      for (let dy = -half; dy <= half; dy++) {
        for (let dx = -half; dx <= half; dx++) {
          const hpx = (hover.x + dx) * TILE_SIZE;
          const hpy = -(hover.y + dy) * TILE_SIZE;
          ctx.fillStyle = 'rgba(255, 255, 255, 0.2)';
          ctx.fillRect(hpx, hpy, TILE_SIZE, TILE_SIZE);
          ctx.strokeStyle = 'rgba(255, 255, 255, 0.8)';
          ctx.lineWidth = 1.5;
          ctx.strokeRect(hpx + 0.5, hpy + 0.5, TILE_SIZE - 1, TILE_SIZE - 1);
        }
      }
    }

    ctx.setTransform(1, 0, 0, 1, 0, 0);
  }, [canvasWidth, canvasHeight, brushSize]);

  const drawMapRef = useRef(drawMap);
  drawMapRef.current = drawMap;

  const requestDraw = useCallback(() => {
    cancelAnimationFrame(rafRef.current);
    rafRef.current = requestAnimationFrame(drawMapRef.current);
  }, []);

  // Redraw on state changes
  useEffect(() => {
    drawMap();
  }, [tiles, pendingTerrain, pendingZones, paintMode, drawMap]);

  // ── Paint tiles at a position with current brush ─────────────────────
  const paintAt = useCallback(
    (tileX: number, tileY: number) => {
      const half = Math.floor(brushSize / 2);
      if (paintMode === 'terrain') {
        setPendingTerrain((prev) => {
          const next = new Map(prev);
          for (let dy = -half; dy <= half; dy++) {
            for (let dx = -half; dx <= half; dx++) {
              const x = tileX + dx;
              const y = tileY + dy;
              if (x < -mapHalf || x > mapHalf || y < -mapHalf || y > mapHalf) continue;
              next.set(`${x},${y}`, { x, y, terrain_type: selectedTerrain });
            }
          }
          return next;
        });
      } else {
        setPendingZones((prev) => {
          const next = new Map(prev);
          for (let dy = -half; dy <= half; dy++) {
            for (let dx = -half; dx <= half; dx++) {
              const x = tileX + dx;
              const y = tileY + dy;
              if (x < -mapHalf || x > mapHalf || y < -mapHalf || y > mapHalf) continue;
              next.set(`${x},${y}`, { x, y, kingdom_zone: selectedZone });
            }
          }
          return next;
        });
      }
    },
    [brushSize, paintMode, selectedTerrain, selectedZone, mapHalf],
  );

  // ── Canvas event handlers ────────────────────────────────────────────
  useEffect(() => {
    const el = canvasRef.current;
    if (!el) return;

    const dpr = window.devicePixelRatio || 1;
    el.width = canvasWidth * dpr;
    el.height = canvasHeight * dpr;
    el.style.width = `${canvasWidth}px`;
    el.style.height = `${canvasHeight}px`;

    let painting = false;

    const onMouseDown = (e: MouseEvent) => {
      if (e.button === 2 || e.shiftKey) {
        draggingRef.current = true;
        dragStartRef.current = { x: e.clientX, y: e.clientY };
        lastOffsetRef.current = { ...offsetRef.current };
        el.style.cursor = 'grabbing';
      } else {
        painting = true;
        const rect = el.getBoundingClientRect();
        const { tileX, tileY } = screenToTile(
          e.clientX - rect.left,
          e.clientY - rect.top,
          offsetRef.current.x,
          offsetRef.current.y,
          scaleRef.current,
        );
        paintAt(tileX, tileY);
      }
    };

    const onMouseMove = (e: MouseEvent) => {
      const rect = el.getBoundingClientRect();
      const screenX = e.clientX - rect.left;
      const screenY = e.clientY - rect.top;

      if (draggingRef.current) {
        const dx = e.clientX - dragStartRef.current.x;
        const dy = e.clientY - dragStartRef.current.y;
        offsetRef.current = {
          x: lastOffsetRef.current.x + dx,
          y: lastOffsetRef.current.y + dy,
        };
        requestDraw();
        return;
      }

      const { tileX, tileY } = screenToTile(screenX, screenY, offsetRef.current.x, offsetRef.current.y, scaleRef.current);

      const prev = hoverTileRef.current;
      if (!prev || prev.x !== tileX || prev.y !== tileY) {
        hoverTileRef.current = { x: tileX, y: tileY };
        requestDraw();
        if (painting) {
          paintAt(tileX, tileY);
        }
      }
    };

    const onMouseUp = () => {
      if (draggingRef.current) {
        draggingRef.current = false;
        el.style.cursor = 'crosshair';
      }
      painting = false;
    };

    const onMouseLeave = () => {
      hoverTileRef.current = null;
      painting = false;
      requestDraw();
    };

    const onWheel = (e: WheelEvent) => {
      e.preventDefault();
      const oldScale = scaleRef.current;
      const delta = e.deltaY > 0 ? -0.05 : 0.05;
      const newScale = Math.max(0.15, Math.min(2, oldScale + delta));

      const rect = el.getBoundingClientRect();
      const mouseX = e.clientX - rect.left;
      const mouseY = e.clientY - rect.top;
      const worldX = (mouseX - offsetRef.current.x) / oldScale;
      const worldY = (mouseY - offsetRef.current.y) / oldScale;

      scaleRef.current = newScale;
      offsetRef.current = {
        x: mouseX - worldX * newScale,
        y: mouseY - worldY * newScale,
      };
      requestDraw();
    };

    const onContextMenu = (e: Event) => e.preventDefault();

    el.addEventListener('mousedown', onMouseDown);
    window.addEventListener('mousemove', onMouseMove);
    window.addEventListener('mouseup', onMouseUp);
    el.addEventListener('mouseleave', onMouseLeave);
    el.addEventListener('wheel', onWheel, { passive: false });
    el.addEventListener('contextmenu', onContextMenu);
    el.style.cursor = 'crosshair';

    drawMapRef.current();

    return () => {
      el.removeEventListener('mousedown', onMouseDown);
      window.removeEventListener('mousemove', onMouseMove);
      window.removeEventListener('mouseup', onMouseUp);
      el.removeEventListener('mouseleave', onMouseLeave);
      el.removeEventListener('wheel', onWheel);
      el.removeEventListener('contextmenu', onContextMenu);
      cancelAnimationFrame(rafRef.current);
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [requestDraw, paintAt]);

  // ── Save pending changes to template ─────────────────────────────────
  const handleSave = async () => {
    if (!activeTemplate) return;
    const totalPending = pendingTerrain.size + pendingZones.size;
    if (totalPending === 0) return;

    setSaving(true);
    setError(null);
    setStatus(null);
    try {
      if (pendingTerrain.size > 0) {
        await updateTemplateTerrain(activeTemplate, Array.from(pendingTerrain.values()));
      }
      if (pendingZones.size > 0) {
        await updateTemplateZones(activeTemplate, Array.from(pendingZones.values()));
      }
      setPendingTerrain(new Map());
      setPendingZones(new Map());
      // Reload to see updated data
      await loadTemplate(activeTemplate);
      setStatus(`Saved ${totalPending} change(s)`);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save changes');
    } finally {
      setSaving(false);
    }
  };

  const handleClear = () => {
    setPendingTerrain(new Map());
    setPendingZones(new Map());
    setStatus(null);
  };

  // ── Template management ──────────────────────────────────────────────
  const handleCreate = async () => {
    if (!newName.trim()) return;
    setError(null);
    try {
      await createTemplate(newName.trim(), newDesc.trim(), newSize);
      setShowNewDialog(false);
      setNewName('');
      setNewDesc('');
      setNewSize(51);
      await loadTemplateList();
      await loadTemplate(newName.trim());
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create template');
    }
  };

  const handleDelete = async () => {
    if (!activeTemplate) return;
    if (!confirm(`Delete template "${activeTemplate}"? This cannot be undone.`)) return;
    try {
      await deleteTemplate(activeTemplate);
      setActiveTemplate(null);
      setTiles(new Map());
      setPendingTerrain(new Map());
      setPendingZones(new Map());
      await loadTemplateList();
      setStatus('Template deleted');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete template');
    }
  };

  const handleResize = async () => {
    if (!activeTemplate) return;
    const input = prompt(`Current size: ${mapSize}×${mapSize}\nEnter new map size (odd number, 3–201):`, String(mapSize));
    if (!input) return;
    const size = parseInt(input, 10);
    if (isNaN(size) || size < 3 || size > 201) {
      setError('Invalid size. Must be an odd number between 3 and 201.');
      return;
    }
    setError(null);
    try {
      await resizeTemplate(activeTemplate, size);
      await loadTemplate(activeTemplate);
      await loadTemplateList();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to resize template');
    }
  };

  const handleApply = async () => {
    if (!activeTemplate) return;
    if (!confirm(
      `Apply template "${activeTemplate}" to the live game?\n\nThis will overwrite the current world map terrain and zones.\nVillage ownership is preserved.`,
    )) return;
    setError(null);
    try {
      await applyTemplate(activeTemplate);
      setStatus('Template applied to live map!');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to apply template');
    }
  };

  const handleExport = async () => {
    if (!activeTemplate) return;
    try {
      await exportTemplate(activeTemplate);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to export template');
    }
  };

  const handleImport = () => {
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = '.json';
    input.onchange = async (e) => {
      const file = (e.target as HTMLInputElement).files?.[0];
      if (!file) return;
      try {
        await importTemplate(file);
        await loadTemplateList();
        setStatus(`Imported template from ${file.name}`);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to import template');
      }
    };
    input.click();
  };

  const totalPending = pendingTerrain.size + pendingZones.size;
  const hover = hoverTileRef.current;

  return (
    <div className={styles.page}>
      <h2 className={styles.heading}>Map Template Editor</h2>
      <p className={styles.subtitle}>
        Design map templates offline. Paint terrain & zones, then apply to the live game.
        Left-click to paint, shift+drag or right-drag to pan, scroll to zoom.
      </p>

      {error && <div className={styles.error}>{error}</div>}
      {status && <div className={styles.status}>{status}</div>}

      {/* ── Template selector bar ─────────────────────────────────────── */}
      <div className={styles.toolbar}>
        <div className={styles.templateSelector}>
          <select
            className={styles.templateSelect}
            value={activeTemplate ?? ''}
            onChange={(e) => {
              if (e.target.value) loadTemplate(e.target.value);
            }}
            disabled={loadingTiles}
          >
            <option value="">— Select template —</option>
            {templates.map((t) => (
              <option key={t.name} value={t.name}>
                {t.name} ({t.map_size}×{t.map_size})
              </option>
            ))}
          </select>
          <button className={styles.clearBtn} onClick={() => setShowNewDialog(true)}>
            + New
          </button>
          <button className={styles.clearBtn} onClick={handleImport}>
            Import
          </button>
          {activeTemplate && (
            <>
              <button className={styles.clearBtn} onClick={handleResize}>
                Resize
              </button>
              <button className={styles.clearBtn} onClick={handleExport}>
                Export
              </button>
              <button className={styles.clearBtn} onClick={handleDelete} style={{ color: '#dc3545' }}>
                Delete
              </button>
            </>
          )}
        </div>

        {activeTemplate && (
          <button className={styles.applyBtn} onClick={handleApply}>
            Apply to Game ({mapSize}×{mapSize})
          </button>
        )}
      </div>

      {/* ── New template dialog ───────────────────────────────────────── */}
      {showNewDialog && (
        <div className={styles.toolbar}>
          <input
            className={styles.templateInput}
            placeholder="Template name"
            value={newName}
            onChange={(e) => setNewName(e.target.value)}
            autoFocus
          />
          <input
            className={styles.templateInput}
            placeholder="Description (optional)"
            value={newDesc}
            onChange={(e) => setNewDesc(e.target.value)}
          />
          <div className={styles.sizeControl}>
            <label className={styles.brushLabel}>Size:</label>
            <input
              className={styles.sizeInput}
              type="number"
              min={3}
              max={201}
              step={2}
              value={newSize}
              onChange={(e) => {
                let v = parseInt(e.target.value, 10);
                if (!isNaN(v)) {
                  if (v % 2 === 0) v++;
                  setNewSize(Math.max(3, Math.min(201, v)));
                }
              }}
            />
            <span className={styles.brushLabel}>{newSize}×{newSize}</span>
          </div>
          <button className={styles.saveBtn} onClick={handleCreate} disabled={!newName.trim()}>
            Create
          </button>
          <button className={styles.clearBtn} onClick={() => setShowNewDialog(false)}>
            Cancel
          </button>
        </div>
      )}

      {/* ── Paint toolbar ─────────────────────────────────────────────── */}
      {activeTemplate && tiles.size > 0 && (
        <div className={styles.toolbar}>
          {/* Mode toggle */}
          <div className={styles.modeToggle}>
            <button
              className={`${styles.modeBtn} ${paintMode === 'terrain' ? styles.modeBtnActive : ''}`}
              onClick={() => setPaintMode('terrain')}
            >
              Terrain
            </button>
            <button
              className={`${styles.modeBtn} ${paintMode === 'zone' ? styles.modeBtnActive : ''}`}
              onClick={() => setPaintMode('zone')}
            >
              Zone
            </button>
          </div>

          {/* Terrain picker or zone picker depending on mode */}
          {paintMode === 'terrain' ? (
            <div className={styles.terrainPicker}>
              {TERRAIN_TYPES.map((t) => {
                const cfg = TERRAIN_CONFIG[t];
                return (
                  <button
                    key={t}
                    className={`${styles.terrainBtn} ${selectedTerrain === t ? styles.terrainBtnActive : ''}`}
                    onClick={() => setSelectedTerrain(t)}
                    style={{ '--terrain-color': hexColor(cfg.color) } as React.CSSProperties}
                    title={cfg.label}
                  >
                    <span className={styles.terrainSwatch} />
                    <span className={styles.terrainLabel}>{cfg.label}</span>
                  </button>
                );
              })}
            </div>
          ) : (
            <div className={styles.terrainPicker}>
              {KINGDOM_ZONES.map((z) => {
                const cfg = ZONE_CONFIG[z];
                if (!cfg) return null;
                return (
                  <button
                    key={z}
                    className={`${styles.terrainBtn} ${selectedZone === z ? styles.terrainBtnActive : ''}`}
                    onClick={() => setSelectedZone(z)}
                    style={{ '--terrain-color': hexColor(cfg.color) } as React.CSSProperties}
                    title={cfg.label}
                  >
                    <span className={styles.terrainSwatch} />
                    <span className={styles.terrainLabel}>{cfg.label}</span>
                  </button>
                );
              })}
            </div>
          )}

          <div className={styles.brushControl}>
            <label className={styles.brushLabel}>
              Brush: {brushSize}×{brushSize}
            </label>
            <input
              type="range"
              min="1"
              max="5"
              value={brushSize}
              onChange={(e) => setBrushSize(parseInt(e.target.value, 10))}
              className={styles.brushSlider}
            />
          </div>

          <div className={styles.saveControls}>
            <span className={styles.pendingCount}>
              {totalPending} pending
            </span>
            <button className={styles.clearBtn} onClick={handleClear} disabled={totalPending === 0}>
              Clear
            </button>
            <button
              className={styles.saveBtn}
              onClick={handleSave}
              disabled={totalPending === 0 || saving}
            >
              {saving ? 'Saving…' : 'Save'}
            </button>
          </div>
        </div>
      )}

      {/* ── Canvas ────────────────────────────────────────────────────── */}
      {activeTemplate && (
        <div className={styles.canvasWrapper}>
          {loadingTiles && (
            <div className={styles.loadingOverlay}>Loading template…</div>
          )}
          <canvas ref={canvasRef} />
          {hover && (
            <div className={styles.coords}>
              ({hover.x}, {hover.y})
            </div>
          )}
        </div>
      )}

      {/* ── Empty state ───────────────────────────────────────────────── */}
      {!activeTemplate && !loadingTemplates && (
        <div className={styles.emptyState}>
          <p>Select a template above or create a new one to start editing.</p>
        </div>
      )}
    </div>
  );
}
