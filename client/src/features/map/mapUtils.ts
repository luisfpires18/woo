// Shared utility functions for map canvas rendering.
// Used by both MapRenderer and AdminMapEditorPage.

/** Tile size in pixels (world units). */
export const TILE_SIZE = 128;

/** Convert a 0xRRGGBB number to a CSS hex string. */
export function hexColor(n: number): string {
  return `#${n.toString(16).padStart(6, '0')}`;
}

/** Convert a 0xRRGGBB number to an rgba() CSS string. */
export function hexColorAlpha(n: number, a: number): string {
  const r = (n >> 16) & 0xff;
  const g = (n >> 8) & 0xff;
  const b = n & 0xff;
  return `rgba(${r},${g},${b},${a})`;
}

/** Convert screen coordinates to map tile coordinates. */
export function screenToTile(
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

/**
 * Simple deterministic hash for (x, y) → unsigned 32-bit integer.
 * Used to pick a sprite variant for a tile position so the choice is
 * stable across renders.
 */
export function tileHash(x: number, y: number): number {
  let h = (x * 374761393 + y * 668265263) | 0;
  h = (h ^ (h >> 13)) * 1274126177;
  h = h ^ (h >> 16);
  return h >>> 0;
}

/**
 * Extract the base zone/terrain name from an asset ID.
 * "zone_veridor"     → "veridor"
 * "zone_veridor_v2"  → "veridor"
 * "terrain_plains"   → "plains"
 * "terrain_forest_v2"→ "forest"
 */
export function extractBaseName(id: string, prefix: string): string {
  const stripped = id.startsWith(prefix) ? id.slice(prefix.length) : id;
  return stripped.replace(/_v\d+$/, '');
}
