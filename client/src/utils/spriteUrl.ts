/**
 * Convention-based sprite URL builder.
 *
 * Sprites live on disk under `server/uploads/sprites/` and are served via
 * the static file server at `/uploads/sprites/…`. The frontend constructs
 * URLs purely by naming convention — no database column needed.
 *
 * Folder layout (kingdom-first):
 *   sprites/{kingdom}/units/{troop_type}.png
 *   sprites/{kingdom}/buildings/{building_type}.png
 *   sprites/kingdoms/{kingdom}/buildings/{kingdom}_{resource}_{slot}[_name].png  (resolved via /api/sprites/building/)
 *   sprites/flags/{kingdom}.png
 *   sprites/resources/{resource_type}.png
 *   sprites/map/village_markers/{id}.png
 *   sprites/map/zone_tiles/{id}.png
 *   sprites/map/terrain_tiles/{id}.png
 *
 * GameIcon's existing onError fallback handles missing files gracefully.
 */

type SpriteKind =
  | 'unit'
  | 'building'
  | 'resource_building'
  | 'kingdom_flag'
  | 'resource'
  | 'village_marker'
  | 'zone_tile'
  | 'terrain_tile';

interface SpriteUrlOptions {
  kind: SpriteKind;
  /** The asset / type identifier (e.g. "infantry_1", "barracks", "food"). */
  id: string;
  /** Required for kingdom-scoped sprites (unit, building, resource_building). */
  kingdom?: string;
  /** Slot number — only for resource_building. */
  slot?: number;
}

/**
 * Return the convention sprite URL for a given asset kind + id.
 * Returns `null` when required fields are missing (caller should fall back to emoji).
 */
export function getSpriteUrl(opts: SpriteUrlOptions): string | null {
  const { kind, id, kingdom, slot } = opts;

  switch (kind) {
    case 'unit':
      if (!kingdom) return null;
      return `/api/sprites/troop/${kingdom}/${id}`;

    case 'building':
      if (!kingdom) return null;
      return `/api/sprites/building-display/${kingdom}/${id}`;

    case 'resource_building':
      if (!kingdom || slot == null) return null;
      return `/api/sprites/building/${kingdom}/${id}_${slot}`;

    case 'kingdom_flag':
      return `/uploads/sprites/flags/${id}.png`;

    case 'resource':
      return `/uploads/sprites/resources/${id}.png`;

    case 'village_marker':
      return `/uploads/sprites/map/village_markers/${id}.png`;

    case 'zone_tile':
      return `/uploads/sprites/map/zone_tiles/${id}.png`;

    case 'terrain_tile':
      return `/uploads/sprites/map/terrain_tiles/${id}.png`;

    default:
      return null;
  }
}
