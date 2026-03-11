import { useState, useEffect, useCallback } from 'react';
import { fetchGameAssets, createGameAsset, deleteGameAsset } from '../../../services/admin';
import type { GameAsset, AssetCategory } from '../../../types/api';
import { LoadingSpinner } from '../../../components/LoadingSpinner/LoadingSpinner';
import { useAssetStore } from '../../../stores/assetStore';
import { getSpriteUrl } from '../../../utils/spriteUrl';
import styles from './AdminAssetsPage.module.css';

const CATEGORY_ORDER: AssetCategory[] = ['kingdom_flag', 'village_marker', 'zone_tile', 'terrain_tile', 'building', 'resource', 'unit'];

const CATEGORY_LABELS: Record<AssetCategory, string> = {
  kingdom_flag: 'Kingdom Flags',
  village_marker: 'Village Markers',
  zone_tile: 'Zone Tiles',
  terrain_tile: 'Terrain Tiles',
  building: 'Buildings',
  resource: 'Resources',
  unit: 'Units',
};

const SPRITE_DIMENSIONS: Record<AssetCategory, { w: number; h: number }> = {
  kingdom_flag: { w: 256, h: 256 },
  village_marker: { w: 256, h: 256 },
  zone_tile: { w: 256, h: 256 },
  terrain_tile: { w: 256, h: 256 },
  building: { w: 96, h: 96 },
  resource: { w: 32, h: 32 },
  unit: { w: 256, h: 256 },
};

/** Categories that support adding / removing variants */
const VARIANT_CATEGORIES: Set<AssetCategory> = new Set(['zone_tile', 'terrain_tile']);

export function AdminAssetsPage() {
  const [assets, setAssets] = useState<GameAsset[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [deleting, setDeleting] = useState<string | null>(null);
  const [failedSprites, setFailedSprites] = useState<Set<string>>(new Set());

  const addAssetStore = useAssetStore((s) => s.addAsset);
  const removeAssetStore = useAssetStore((s) => s.removeAsset);

  const load = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const resp = await fetchGameAssets();
      setAssets(resp.assets);
    } catch {
      setError('Failed to load game assets.');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  const handleSpriteError = (id: string) => {
    setFailedSprites((prev) => new Set(prev).add(id));
  };

  /**
   * Add a new variant for a zone_tile or terrain_tile.
   * Auto-generates the next variant ID: e.g. zone_veridor → zone_veridor_v2, zone_veridor_v3 …
   */
  const handleAddVariant = async (baseAsset: GameAsset) => {
    setError(null);
    // Find existing variants: base ID + _v\d+
    const baseId = baseAsset.id;
    const existing = assets.filter(
      (a) => a.id === baseId || (a.id.startsWith(baseId + '_v') && a.category === baseAsset.category),
    );
    // Next variant number
    let maxV = 1;
    for (const a of existing) {
      const match = a.id.match(/_v(\d+)$/);
      if (match) {
        maxV = Math.max(maxV, parseInt(match[1]!, 10));
      }
    }
    const newId = `${baseId}_v${maxV + 1}`;
    const newName = `${baseAsset.display_name} v${maxV + 1}`;

    try {
      const created = await createGameAsset({
        id: newId,
        category: baseAsset.category,
        display_name: newName,
        default_icon: baseAsset.default_icon,
      });
      setAssets((prev) => [...prev, created]);
      addAssetStore(created);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create variant.');
    }
  };

  /**
   * Delete a variant asset entirely.
   * Only allowed for variant IDs (those with _v\d+ suffix).
   */
  const handleDeleteAsset = async (id: string) => {
    if (!confirm(`Delete asset "${id}" permanently?`)) return;
    setDeleting(id);
    setError(null);
    try {
      await deleteGameAsset(id);
      setAssets((prev) => prev.filter((a) => a.id !== id));
      removeAssetStore(id);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete asset.');
    } finally {
      setDeleting(null);
    }
  };

  /** Check if an asset is a variant (has _v\d+ suffix) — only variants can be deleted */
  const isVariant = (id: string) => /_v\d+$/.test(id);

  if (loading) {
    return (
      <div className={styles.center}>
        <LoadingSpinner size="lg" />
      </div>
    );
  }

  const grouped = CATEGORY_ORDER.map((cat) => ({
    category: cat,
    label: CATEGORY_LABELS[cat],
    dims: SPRITE_DIMENSIONS[cat],
    items: assets.filter((a) => a.category === cat),
  })).filter((g) => g.items.length > 0);

  return (
    <div className={styles.page}>
      <h2 className={styles.heading}>Game Assets</h2>
      <p className={styles.subtitle}>
        Game asset sprites are loaded by convention from the filesystem. Emoji icons are used as fallback.
      </p>

      {error && <div className={styles.error}>{error}</div>}

      {grouped.map((group) => (
        <section key={group.category} className={styles.section}>
          <h3 className={styles.sectionTitle}>
            {group.label}
            <span className={styles.dimBadge}>
              {group.dims.w}×{group.dims.h}px
            </span>
          </h3>

          <div className={styles.grid}>
            {group.items.map((asset) => {
              const spriteUrl = getSpriteUrl({ kind: asset.category as any, id: asset.id });
              const showSprite = spriteUrl && !failedSprites.has(asset.id);

              return (
                <div key={asset.id} className={styles.card}>
                  <div className={styles.preview}>
                    {showSprite ? (
                      <img
                        src={spriteUrl}
                        alt={asset.display_name}
                        width={group.dims.w}
                        height={group.dims.h}
                        className={styles.spriteImg}
                        onError={() => handleSpriteError(asset.id)}
                      />
                    ) : (
                      <span className={styles.emoji}>{asset.default_icon}</span>
                    )}
                  </div>

                  <div className={styles.info}>
                    <span className={styles.assetName}>{asset.display_name}</span>
                    <span className={styles.assetId}>{asset.id}</span>
                  </div>

                  <div className={styles.actions}>
                    {/* Variant management for zone_tile / terrain_tile */}
                    {VARIANT_CATEGORIES.has(group.category) && !isVariant(asset.id) && (
                      <button
                        className={styles.variantBtn}
                        onClick={() => handleAddVariant(asset)}
                        title="Add a new sprite variant for this tile"
                      >
                        + Variant
                      </button>
                    )}
                    {VARIANT_CATEGORIES.has(group.category) && isVariant(asset.id) && (
                      <button
                        className={styles.deleteAssetBtn}
                        disabled={deleting === asset.id}
                        onClick={() => handleDeleteAsset(asset.id)}
                        title="Delete this variant permanently"
                      >
                        Delete
                      </button>
                    )}
                  </div>
                </div>
              );
            })}
          </div>
        </section>
      ))}
    </div>
  );
}
