import { useState, useEffect, useCallback, useRef } from 'react';
import { fetchGameAssets, uploadSprite, deleteSprite, createGameAsset, deleteGameAsset } from '../../../services/admin';
import type { GameAsset, AssetCategory } from '../../../types/api';
import { LoadingSpinner } from '../../../components/LoadingSpinner/LoadingSpinner';
import { useAssetStore } from '../../../stores/assetStore';
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
  unit: { w: 64, h: 64 },
};

/** Categories that support adding / removing variants */
const VARIANT_CATEGORIES: Set<AssetCategory> = new Set(['zone_tile', 'terrain_tile']);

export function AdminAssetsPage() {
  const [assets, setAssets] = useState<GameAsset[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [uploading, setUploading] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [activeUploadId, setActiveUploadId] = useState<string | null>(null);

  const upsertStore = useAssetStore((s) => s.upsert);
  const clearSpriteStore = useAssetStore((s) => s.clearSprite);
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

  const handleUploadClick = (id: string) => {
    setActiveUploadId(id);
    fileInputRef.current?.click();
  };

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file || !activeUploadId) return;

    // Reset the input so the same file can be re-selected
    e.target.value = '';

    setUploading(activeUploadId);
    setError(null);
    try {
      const updated = await uploadSprite(activeUploadId, file);
      setAssets((prev) => prev.map((a) => (a.id === updated.id ? updated : a)));
      upsertStore(updated);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Upload failed.');
    } finally {
      setUploading(null);
      setActiveUploadId(null);
    }
  };

  const handleDelete = async (id: string) => {
    setUploading(id);
    setError(null);
    try {
      await deleteSprite(id);
      setAssets((prev) =>
        prev.map((a) => (a.id === id ? { ...a, sprite_url: null } : a)),
      );
      clearSpriteStore(id);
    } catch {
      setError('Failed to delete sprite.');
    } finally {
      setUploading(null);
    }
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
   * Delete a variant asset entirely (row + sprite).
   * Only allowed for variant IDs (those with _v\d+ suffix).
   */
  const handleDeleteAsset = async (id: string) => {
    if (!confirm(`Delete asset "${id}" permanently?`)) return;
    setUploading(id);
    setError(null);
    try {
      await deleteGameAsset(id);
      setAssets((prev) => prev.filter((a) => a.id !== id));
      removeAssetStore(id);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete asset.');
    } finally {
      setUploading(null);
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
        Upload custom PNG sprites for buildings, resources, and units. Emoji icons are used as fallback.
      </p>

      {error && <div className={styles.error}>{error}</div>}

      {/* Hidden file input shared across all cards */}
      <input
        ref={fileInputRef}
        type="file"
        accept="image/png"
        className={styles.hiddenInput}
        onChange={handleFileChange}
      />

      {grouped.map((group) => (
        <section key={group.category} className={styles.section}>
          <h3 className={styles.sectionTitle}>
            {group.label}
            <span className={styles.dimBadge}>
              {group.dims.w}×{group.dims.h}px
            </span>
          </h3>

          <div className={styles.grid}>
            {group.items.map((asset) => (
              <div key={asset.id} className={styles.card}>
                <div className={styles.preview}>
                  {asset.sprite_url ? (
                    <img
                      src={asset.sprite_url}
                      alt={asset.display_name}
                      width={group.dims.w}
                      height={group.dims.h}
                      className={styles.spriteImg}
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
                  <button
                    className={styles.uploadBtn}
                    disabled={uploading === asset.id}
                    onClick={() => handleUploadClick(asset.id)}
                  >
                    {uploading === asset.id ? '…' : asset.sprite_url ? 'Replace' : 'Upload'}
                  </button>
                  {asset.sprite_url && (
                    <button
                      className={styles.deleteBtn}
                      disabled={uploading === asset.id}
                      onClick={() => handleDelete(asset.id)}
                    >
                      Remove
                    </button>
                  )}
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
                      disabled={uploading === asset.id}
                      onClick={() => handleDeleteAsset(asset.id)}
                      title="Delete this variant permanently"
                    >
                      Delete
                    </button>
                  )}
                </div>
              </div>
            ))}
          </div>
        </section>
      ))}
    </div>
  );
}
