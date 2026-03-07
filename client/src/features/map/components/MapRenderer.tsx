// Canvas 2D map renderer — renders the world map grid using native canvas API
// Uses a <canvas> element with setTransform for panning/zoom
// Features: adaptive chunk loading, hover highlight, selected tile highlight,
//           chunk loading during drag & zoom, keyboard navigation (arrows/WASD/+/-)

import { useEffect, useRef, useCallback, useImperativeHandle, forwardRef } from 'react';
import { useMapStore } from '../../../stores/mapStore';
import { useAssetStore } from '../../../stores/assetStore';
import type { MapTile } from '../../../types/map';
import { TERRAIN_CONFIG, ZONE_COLORS } from '../../../types/map';

const TILE_SIZE = 128;
/** Maximum range the server accepts */
const MAX_SERVER_RANGE = 40;
/** Minimum tile distance the viewport center must move before re-loading during drag */
const DRAG_LOAD_THRESHOLD = 5;

/** Convert a 0xRRGGBB number to a CSS hex string */
function hexColor(n: number): string {
  return `#${n.toString(16).padStart(6, '0')}`;
}

/** Convert a 0xRRGGBB number + alpha (0-1) to an rgba() string */
function hexColorAlpha(n: number, alpha: number): string {
  const r = (n >> 16) & 0xff;
  const g = (n >> 8) & 0xff;
  const b = n & 0xff;
  return `rgba(${r},${g},${b},${alpha})`;
}



/**
 * Compute the range needed to fill the viewport at the current zoom level.
 * Returns the radius in tiles (clamped to MAX_SERVER_RANGE).
 */
function computeAdaptiveRange(
  viewportWidth: number,
  viewportHeight: number,
  scale: number,
): number {
  const tilesAcross = viewportWidth / (TILE_SIZE * scale);
  const tilesDown = viewportHeight / (TILE_SIZE * scale);
  // Use the larger dimension, halved to get radius, + a buffer of 3 tiles
  const radius = Math.ceil(Math.max(tilesAcross, tilesDown) / 2) + 3;
  return Math.min(radius, MAX_SERVER_RANGE);
}

/** Convert screen coordinates to map tile coordinates */
function screenToTile(
  screenX: number,
  screenY: number,
  ox: number,
  oy: number,
  scale: number,
): { tileX: number; tileY: number } {
  const worldX = (screenX - ox) / scale;
  const worldY = (screenY - oy) / scale;
  return {
    tileX: Math.floor(worldX / TILE_SIZE),
    tileY: -Math.floor(worldY / TILE_SIZE),
  };
}

interface MapRendererProps {
  initialX?: number;
  initialY?: number;
  width?: number;
  height?: number;
  selectedTile?: MapTile | null;
  onTileClick?: (tile: MapTile) => void;
  onTileHover?: (x: number, y: number) => void;
}

/** World bounds for minimap (map is -25 to +25) */
const WORLD_MIN = -25;
const WORLD_MAX = 25;
const WORLD_SIZE = WORLD_MAX - WORLD_MIN + 1; // 51
const MINIMAP_SIZE = 140;
const MINIMAP_MARGIN = 10;

/**
 * Draw a minimap overlay in the bottom-left corner (screen-space).
 * Shows all loaded tiles as colored dots and a white rect for the viewport.
 */
function drawMinimap(
  ctx: CanvasRenderingContext2D,
  tiles: Map<string, MapTile>,
  scale: number,
  ox: number,
  oy: number,
  viewWidth: number,
  viewHeight: number,
) {
  const mx = MINIMAP_MARGIN;
  const my = viewHeight - MINIMAP_SIZE - MINIMAP_MARGIN;
  const pixPerTile = MINIMAP_SIZE / WORLD_SIZE;

  // Background
  ctx.fillStyle = 'rgba(0,0,0,0.6)';
  ctx.fillRect(mx, my, MINIMAP_SIZE, MINIMAP_SIZE);
  ctx.strokeStyle = 'rgba(255,255,255,0.3)';
  ctx.lineWidth = 1;
  ctx.strokeRect(mx, my, MINIMAP_SIZE, MINIMAP_SIZE);

  // Draw tiles as dots
  tiles.forEach((tile) => {
    const tx = mx + (tile.x - WORLD_MIN) * pixPerTile;
    const ty = my + (-tile.y - WORLD_MIN) * pixPerTile; // invert Y
    const cfg = TERRAIN_CONFIG[tile.terrain] ?? TERRAIN_CONFIG.plains;
    ctx.fillStyle = hexColor(cfg.color);
    ctx.fillRect(tx, ty, Math.max(1, pixPerTile), Math.max(1, pixPerTile));
  });

  // Viewport rectangle
  const invScale = 1 / scale;
  const worldLeft = -ox * invScale;
  const worldTop = -oy * invScale;
  const vpW = viewWidth * invScale;
  const vpH = viewHeight * invScale;

  // Convert world pixel coords to tile coords for minimap
  const vpTileLeft = worldLeft / TILE_SIZE;
  const vpTileTop = -worldTop / TILE_SIZE; // invert Y
  const vpTileW = vpW / TILE_SIZE;
  const vpTileH = vpH / TILE_SIZE;

  const rx = mx + (vpTileLeft - WORLD_MIN) * pixPerTile;
  const ry = my + (-vpTileTop - WORLD_MIN) * pixPerTile;
  const rw = vpTileW * pixPerTile;
  const rh = vpTileH * pixPerTile;

  ctx.strokeStyle = '#ffcc00';
  ctx.lineWidth = 1.5;
  ctx.strokeRect(rx, ry, rw, rh);
}

/** Imperative handle exposed to parent via ref */
export interface MapRendererHandle {
  navigateTo: (x: number, y: number) => void;
}

export const MapRenderer = forwardRef<MapRendererHandle, MapRendererProps>(function MapRenderer(
  {
    initialX = 0,
    initialY = 0,
    width = 800,
    height = 600,
    selectedTile,
    onTileClick,
    onTileHover,
  },
  ref,
) {
  const canvasRef = useRef<HTMLCanvasElement>(null);

  // View state stored in refs so drawing doesn't trigger re-renders
  const offsetRef = useRef({ x: 0, y: 0 });
  const scaleRef = useRef(1);
  const draggingRef = useRef(false);
  const dragStartRef = useRef({ x: 0, y: 0 });
  const lastOffsetRef = useRef({ x: 0, y: 0 });
  const loadingRef = useRef(false);
  const lastChunkRef = useRef('');
  const rafRef = useRef(0);
  const hoverTileRef = useRef<{ x: number; y: number } | null>(null);
  const selectedTileRef = useRef<MapTile | null>(null);

  // Keep selectedTile ref in sync with prop
  selectedTileRef.current = selectedTile ?? null;

  const tiles = useMapStore((s) => s.tiles);
  const loadChunk = useMapStore((s) => s.loadChunk);

  // Asset store for village marker sprites
  const assetStoreAssets = useAssetStore((s) => s.assets);
  const assetStoreLoaded = useAssetStore((s) => s.loaded);
  const assetStoreLoad = useAssetStore((s) => s.load);

  // Image cache for village marker sprites: maps zone name → HTMLImageElement
  const markerImagesRef = useRef<Map<string, HTMLImageElement>>(new Map());
  // Image cache for zone tile sprites: maps zone name → HTMLImageElement
  const zoneTileImagesRef = useRef<Map<string, HTMLImageElement>>(new Map());

  // Preload village marker + zone tile sprites when asset store is ready
  useEffect(() => {
    if (!assetStoreLoaded) {
      assetStoreLoad();
      return;
    }

    const markerAssets = assetStoreAssets.filter((a) => a.category === 'village_marker' && a.sprite_url);
    const markerCache = markerImagesRef.current;

    for (const asset of markerAssets) {
      const zone = asset.id.replace('marker_', '');
      const cached = markerCache.get(zone);
      if (cached && cached.src === new URL(asset.sprite_url!, window.location.origin).href) continue;

      const img = new Image();
      img.src = asset.sprite_url!;
      img.onload = () => {
        markerCache.set(zone, img);
        cancelAnimationFrame(rafRef.current);
        rafRef.current = requestAnimationFrame(drawMapRef.current);
      };
    }

    const zoneTileAssets = assetStoreAssets.filter((a) => a.category === 'zone_tile' && a.sprite_url);
    const zoneCache = zoneTileImagesRef.current;

    for (const asset of zoneTileAssets) {
      // Extract zone from id: "zone_veridor" → "veridor", "zone_default" → "default"
      const zone = asset.id.replace('zone_', '');
      const cached = zoneCache.get(zone);
      if (cached && cached.src === new URL(asset.sprite_url!, window.location.origin).href) continue;

      const img = new Image();
      img.src = asset.sprite_url!;
      img.onload = () => {
        zoneCache.set(zone, img);
        cancelAnimationFrame(rafRef.current);
        rafRef.current = requestAnimationFrame(drawMapRef.current);
      };
    }
  }, [assetStoreAssets, assetStoreLoaded, assetStoreLoad]);

  // We need a stable ref to drawMap so the sprite onload can call it
  const drawMapRef = useRef(() => {});

  // Store tiles in a ref so the draw function always has the latest without re-binding events
  const tilesRef = useRef(tiles);
  tilesRef.current = tiles;

  // ---------- Drawing ----------

  const drawMap = useCallback(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    const dpr = window.devicePixelRatio || 1;
    const scale = scaleRef.current;
    const ox = offsetRef.current.x;
    const oy = offsetRef.current.y;

    // Clear
    ctx.setTransform(dpr, 0, 0, dpr, 0, 0);
    ctx.fillStyle = '#111122';
    ctx.fillRect(0, 0, width, height);

    // Apply pan + zoom transform
    ctx.setTransform(dpr * scale, 0, 0, dpr * scale, dpr * ox, dpr * oy);

    // Compute visible tile range (in world coordinates)
    const invScale = 1 / scale;
    const worldLeft = -ox * invScale;
    const worldTop = -oy * invScale;
    const worldRight = worldLeft + width * invScale;
    const worldBottom = worldTop + height * invScale;

    const currentTiles = tilesRef.current;

    // --- Layer 1: Terrain fills + zone tile sprites ---
    const zoneCache = zoneTileImagesRef.current;
    currentTiles.forEach((tile) => {
      const px = tile.x * TILE_SIZE;
      const py = -tile.y * TILE_SIZE; // Invert Y: map-Y-up → screen-Y-down

      // Frustum cull
      if (px + TILE_SIZE < worldLeft || px > worldRight) return;
      if (py + TILE_SIZE < worldTop || py > worldBottom) return;

      // Try zone-specific tile sprite first, then 'default' tile sprite, then flat color fallback
      const zone = tile.zone || '';
      const zoneTileImg = zoneCache.get(zone) || (zone ? zoneCache.get('default') : undefined);
      if (zoneTileImg && zoneTileImg.complete && zoneTileImg.naturalWidth > 0) {
        ctx.drawImage(zoneTileImg, px, py, TILE_SIZE, TILE_SIZE);
      } else {
        const cfg = TERRAIN_CONFIG[tile.terrain] ?? TERRAIN_CONFIG.plains;
        ctx.fillStyle = hexColor(cfg.color);
        ctx.fillRect(px, py, TILE_SIZE, TILE_SIZE);

        // Subtle zone tint overlay
        const zoneColor = tile.zone ? ZONE_COLORS[tile.zone] : undefined;
        if (zoneColor !== undefined) {
          ctx.fillStyle = hexColorAlpha(zoneColor, 0.08);
          ctx.fillRect(px, py, TILE_SIZE, TILE_SIZE);
        }
      }
    });

    // --- Layer 2: Terrain-edge outlines + zone borders (cel-shaded style) ---
    currentTiles.forEach((tile) => {
      const px = tile.x * TILE_SIZE;
      const py = -tile.y * TILE_SIZE;
      if (px + TILE_SIZE < worldLeft || px > worldRight) return;
      if (py + TILE_SIZE < worldTop || py > worldBottom) return;

      const terrain = tile.terrain || 'plains';
      const zone = tile.zone || '';

      // Check all 4 neighbors for terrain and zone boundaries
      const edges = [
        { dx: 1, dy: 0, x1: px + TILE_SIZE, y1: py, x2: px + TILE_SIZE, y2: py + TILE_SIZE },  // right
        { dx: -1, dy: 0, x1: px, y1: py, x2: px, y2: py + TILE_SIZE },                          // left
        { dx: 0, dy: 1, x1: px, y1: py, x2: px + TILE_SIZE, y2: py },                            // top (map-up)
        { dx: 0, dy: -1, x1: px, y1: py + TILE_SIZE, x2: px + TILE_SIZE, y2: py + TILE_SIZE },  // bottom (map-down)
      ];

      for (const { dx, dy, x1, y1, x2, y2 } of edges) {
        const nKey = `${tile.x + dx},${tile.y + dy}`;
        const neighbor = currentTiles.get(nKey);
        const nTerrain = neighbor?.terrain || '';
        const nZone = neighbor?.zone || '';

        // Terrain edge: thick dark outline where biome changes
        if (neighbor && nTerrain !== terrain) {
          ctx.strokeStyle = '#1a1a2e';
          ctx.lineWidth = 2.5;
          ctx.beginPath();
          ctx.moveTo(x1, y1);
          ctx.lineTo(x2, y2);
          ctx.stroke();
        }
        // Zone border: white line where kingdom changes (on top of terrain outlines)
        else if (neighbor && nZone !== zone) {
          ctx.strokeStyle = 'rgba(255,255,255,0.4)';
          ctx.lineWidth = 1.5;
          ctx.beginPath();
          ctx.moveTo(x1, y1);
          ctx.lineTo(x2, y2);
          ctx.stroke();
        }
      }
    });

    // --- Layer 3: Village markers + labels (all inside the tile) ---
    const markerCache = markerImagesRef.current;
    currentTiles.forEach((tile) => {
      if (!tile.village_id) return;

      const px = tile.x * TILE_SIZE;
      const py = -tile.y * TILE_SIZE;
      if (px + TILE_SIZE < worldLeft || px > worldRight) return;
      if (py + TILE_SIZE < worldTop || py > worldBottom) return;

      const cx = px + TILE_SIZE / 2;

      // Village marker — sprite centered in tile, otherwise circle fallback
      const zone = tile.zone || '';
      const markerImg = markerCache.get(zone);
      if (markerImg && markerImg.complete && markerImg.naturalWidth > 0) {
        // Draw sprite at 100x100 centered in upper portion of tile
        const spriteSize = 100;
        const spriteX = px + (TILE_SIZE - spriteSize) / 2;
        const spriteY = py + 2;
        ctx.drawImage(markerImg, spriteX, spriteY, spriteSize, spriteSize);
      } else {
        // Fallback: white circle with dark outline
        ctx.beginPath();
        ctx.arc(cx, py + 40, 24, 0, Math.PI * 2);
        ctx.fillStyle = '#ffffff';
        ctx.fill();
        ctx.strokeStyle = '#333333';
        ctx.lineWidth = 2;
        ctx.stroke();
      }

      // Village name + coordinates — text with shadow, no black bar
      if (tile.village_name) {
        const nameFont = 'bold 13px Cinzel, serif';
        const coordFont = '10px "EB Garamond", serif';
        const coordText = `(${tile.x}, ${tile.y})`;
        const textY = py + TILE_SIZE - 32;

        // Village name — should fit at 128px width
        ctx.font = nameFont;
        ctx.textAlign = 'center';
        ctx.textBaseline = 'top';
        let displayName = tile.village_name;
        const maxTextW = TILE_SIZE - 10;
        while (ctx.measureText(displayName).width > maxTextW && displayName.length > 3) {
          displayName = displayName.slice(0, -1);
        }
        if (displayName !== tile.village_name) displayName += '\u2026';

        // Text shadow for readability on any background
        ctx.shadowColor = 'rgba(0, 0, 0, 0.9)';
        ctx.shadowBlur = 4;
        ctx.shadowOffsetX = 1;
        ctx.shadowOffsetY = 1;
        ctx.fillStyle = '#f0e6d2';
        ctx.fillText(displayName, cx, textY);

        // Coordinates
        ctx.font = coordFont;
        ctx.fillStyle = '#cccccc';
        ctx.fillText(coordText, cx, textY + 16);

        // Reset shadow
        ctx.shadowColor = 'transparent';
        ctx.shadowBlur = 0;
        ctx.shadowOffsetX = 0;
        ctx.shadowOffsetY = 0;
      }
    });

    // --- Layer 4: Hover highlight (no coord text — coords shown in village label) ---
    const hover = hoverTileRef.current;
    if (hover) {
      const hpx = hover.x * TILE_SIZE;
      const hpy = -hover.y * TILE_SIZE;
      ctx.fillStyle = 'rgba(255, 255, 255, 0.15)';
      ctx.fillRect(hpx, hpy, TILE_SIZE, TILE_SIZE);
      ctx.strokeStyle = 'rgba(255, 255, 255, 0.7)';
      ctx.lineWidth = 1.5;
      ctx.strokeRect(hpx + 0.5, hpy + 0.5, TILE_SIZE - 1, TILE_SIZE - 1);
    }

    // --- Layer 5: Selected tile highlight ---
    const sel = selectedTileRef.current;
    if (sel) {
      const spx = sel.x * TILE_SIZE;
      const spy = -sel.y * TILE_SIZE;
      ctx.strokeStyle = '#ffcc00';
      ctx.lineWidth = 2.5;
      ctx.strokeRect(spx + 1, spy + 1, TILE_SIZE - 2, TILE_SIZE - 2);
    }

    // Reset transform
    ctx.setTransform(1, 0, 0, 1, 0, 0);

    // --- Layer 6: Minimap (bottom-left, screen-space) ---
    drawMinimap(ctx, currentTiles, scale, ox, oy, width, height);
  }, [width, height]);

  // Keep drawMapRef in sync so sprite onload callbacks can trigger redraws
  drawMapRef.current = drawMap;

  // ---------- Chunk loading ----------

  const checkAndLoadChunk = useCallback(() => {
    if (loadingRef.current) return;
    const scale = scaleRef.current;
    const ox = offsetRef.current.x;
    const oy = offsetRef.current.y;
    const viewCenterX = Math.round((width / 2 - ox) / scale / TILE_SIZE);
    const viewCenterY = -Math.round((height / 2 - oy) / scale / TILE_SIZE);
    const range = computeAdaptiveRange(width, height, scale);

    // Include zoom-derived range in dedup key so zoom changes trigger new loads
    const chunkKey = `${viewCenterX},${viewCenterY},${range}`;
    if (chunkKey === lastChunkRef.current) return;
    lastChunkRef.current = chunkKey;

    loadingRef.current = true;
    loadChunk(viewCenterX, viewCenterY, range).finally(() => {
      loadingRef.current = false;
    });
  }, [width, height, loadChunk]);

  // ---------- Initialization ----------

  // Set up canvas sizing + initial offset
  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const dpr = window.devicePixelRatio || 1;
    canvas.width = width * dpr;
    canvas.height = height * dpr;
    canvas.style.width = `${width}px`;
    canvas.style.height = `${height}px`;

    // Center on initial tile
    offsetRef.current = {
      x: width / 2 - initialX * TILE_SIZE,
      y: height / 2 + initialY * TILE_SIZE,
    };
    scaleRef.current = 1;

    // Load initial chunk with adaptive range
    const range = computeAdaptiveRange(width, height, 1);
    loadChunk(initialX, initialY, range);
    lastChunkRef.current = `${initialX},${initialY},${range}`;
    drawMap();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [width, height]);

  // Redraw when tiles or selection change
  useEffect(() => {
    drawMap();
  }, [tiles, selectedTile, drawMap]);

  // ---------- Imperative handle for parent (Go-to navigation) ----------

  useImperativeHandle(
    ref,
    () => ({
      navigateTo(x: number, y: number) {
        // Center the viewport on (x, y) and load a chunk there
        offsetRef.current = {
          x: width / 2 - x * TILE_SIZE * scaleRef.current,
          y: height / 2 + y * TILE_SIZE * scaleRef.current,
        };
        lastChunkRef.current = ''; // force reload
        checkAndLoadChunk();
        drawMap();
      },
    }),
    [width, height, checkAndLoadChunk, drawMap],
  );

  // ---------- Interaction: pan, zoom, click, hover ----------

  const requestDraw = useCallback(() => {
    cancelAnimationFrame(rafRef.current);
    rafRef.current = requestAnimationFrame(drawMap);
  }, [drawMap]);

  const handleTileClick = useCallback(
    (e: MouseEvent) => {
      const canvas = canvasRef.current;
      if (!canvas) return;
      const rect = canvas.getBoundingClientRect();
      const screenX = e.clientX - rect.left;
      const screenY = e.clientY - rect.top;

      const { tileX, tileY } = screenToTile(
        screenX,
        screenY,
        offsetRef.current.x,
        offsetRef.current.y,
        scaleRef.current,
      );

      const key = `${tileX},${tileY}`;
      const tile = tilesRef.current.get(key);
      if (tile && onTileClick) {
        onTileClick(tile);
      }
    },
    [onTileClick],
  );

  // Attach mouse / wheel listeners
  useEffect(() => {
    const el = canvasRef.current;
    if (!el) return;

    /** Track the last chunk-load center during drag (in tile coords) */
    let lastDragLoadX = 0;
    let lastDragLoadY = 0;

    const onMouseDown = (e: MouseEvent) => {
      draggingRef.current = true;
      dragStartRef.current = { x: e.clientX, y: e.clientY };
      lastOffsetRef.current = { ...offsetRef.current };
      el.style.cursor = 'grabbing';

      // Record current center so we can compare during drag
      const scale = scaleRef.current;
      const ox = offsetRef.current.x;
      const oy = offsetRef.current.y;
      lastDragLoadX = Math.round((width / 2 - ox) / scale / TILE_SIZE);
      lastDragLoadY = -Math.round((height / 2 - oy) / scale / TILE_SIZE);
    };

    const onMouseMove = (e: MouseEvent) => {
      if (draggingRef.current) {
        // --- Drag panning ---
        const dx = e.clientX - dragStartRef.current.x;
        const dy = e.clientY - dragStartRef.current.y;
        offsetRef.current = {
          x: lastOffsetRef.current.x + dx,
          y: lastOffsetRef.current.y + dy,
        };
        requestDraw();

        // Load chunks during drag if center moved far enough
        const scale = scaleRef.current;
        const ox = offsetRef.current.x;
        const oy = offsetRef.current.y;
        const nowCenterX = Math.round((width / 2 - ox) / scale / TILE_SIZE);
        const nowCenterY = -Math.round((height / 2 - oy) / scale / TILE_SIZE);
        const dTile = Math.abs(nowCenterX - lastDragLoadX) + Math.abs(nowCenterY - lastDragLoadY);
        if (dTile >= DRAG_LOAD_THRESHOLD) {
          lastDragLoadX = nowCenterX;
          lastDragLoadY = nowCenterY;
          checkAndLoadChunk();
        }
      } else {
        // --- Hover tracking ---
        const rect = el.getBoundingClientRect();
        const screenX = e.clientX - rect.left;
        const screenY = e.clientY - rect.top;
        const { tileX, tileY } = screenToTile(
          screenX,
          screenY,
          offsetRef.current.x,
          offsetRef.current.y,
          scaleRef.current,
        );

        const prev = hoverTileRef.current;
        if (!prev || prev.x !== tileX || prev.y !== tileY) {
          hoverTileRef.current = { x: tileX, y: tileY };
          onTileHover?.(tileX, tileY);
          requestDraw();
        }
      }
    };

    const onMouseLeave = () => {
      if (hoverTileRef.current) {
        hoverTileRef.current = null;
        requestDraw();
      }
    };

    const onMouseUp = (e: MouseEvent) => {
      if (!draggingRef.current) return;
      const dx = Math.abs(e.clientX - dragStartRef.current.x);
      const dy = Math.abs(e.clientY - dragStartRef.current.y);
      draggingRef.current = false;
      el.style.cursor = 'grab';

      // Tiny move = click
      if (dx < 3 && dy < 3) {
        handleTileClick(e);
      }

      checkAndLoadChunk();
    };

    const onWheel = (e: WheelEvent) => {
      e.preventDefault();
      const oldScale = scaleRef.current;
      const delta = e.deltaY > 0 ? -0.1 : 0.1;
      const newScale = Math.max(0.3, Math.min(3, oldScale + delta));

      // Zoom toward mouse position
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
      checkAndLoadChunk();
    };

    el.addEventListener('mousedown', onMouseDown);
    window.addEventListener('mousemove', onMouseMove);
    window.addEventListener('mouseup', onMouseUp);
    el.addEventListener('mouseleave', onMouseLeave);
    el.addEventListener('wheel', onWheel, { passive: false });
    el.style.cursor = 'grab';

    // --- Keyboard navigation ---
    const PAN_SPEED = TILE_SIZE * 3; // pixels per key press

    const onKeyDown = (e: KeyboardEvent) => {
      // Only handle keys when canvas (or body) is focused, not inside inputs
      if (
        e.target instanceof HTMLInputElement ||
        e.target instanceof HTMLTextAreaElement ||
        e.target instanceof HTMLSelectElement
      ) {
        return;
      }

      let handled = true;
      switch (e.key) {
        case 'ArrowLeft':
        case 'a':
          offsetRef.current = { ...offsetRef.current, x: offsetRef.current.x + PAN_SPEED * scaleRef.current };
          break;
        case 'ArrowRight':
        case 'd':
          offsetRef.current = { ...offsetRef.current, x: offsetRef.current.x - PAN_SPEED * scaleRef.current };
          break;
        case 'ArrowUp':
        case 'w':
          offsetRef.current = { ...offsetRef.current, y: offsetRef.current.y + PAN_SPEED * scaleRef.current };
          break;
        case 'ArrowDown':
        case 's':
          offsetRef.current = { ...offsetRef.current, y: offsetRef.current.y - PAN_SPEED * scaleRef.current };
          break;
        case '+':
        case '=':
          scaleRef.current = Math.min(3, scaleRef.current + 0.15);
          break;
        case '-':
        case '_':
          scaleRef.current = Math.max(0.3, scaleRef.current - 0.15);
          break;
        default:
          handled = false;
      }

      if (handled) {
        e.preventDefault();
        requestDraw();
        checkAndLoadChunk();
      }
    };

    window.addEventListener('keydown', onKeyDown);

    return () => {
      el.removeEventListener('mousedown', onMouseDown);
      window.removeEventListener('mousemove', onMouseMove);
      window.removeEventListener('mouseup', onMouseUp);
      el.removeEventListener('mouseleave', onMouseLeave);
      el.removeEventListener('wheel', onWheel);
      window.removeEventListener('keydown', onKeyDown);
      cancelAnimationFrame(rafRef.current);
    };
    // onTileHover is intentionally captured in closure — stable ref via parent
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [width, height, requestDraw, handleTileClick, checkAndLoadChunk]);

  return (
    <canvas
      ref={canvasRef}
      style={{
        width: `${width}px`,
        height: `${height}px`,
        borderRadius: 'var(--radius-md)',
        border: '1px solid var(--border)',
        display: 'block',
      }}
    />
  );
});
