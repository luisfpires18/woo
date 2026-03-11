import { useState, useEffect, useCallback } from 'react';
import {
  fetchGameAssets,
  createGameAsset,
  deleteGameAsset,
} from '../../../services/admin';
import type { GameAsset, AssetCategory } from '../../../types/api';
import { LoadingSpinner } from '../../../components/LoadingSpinner/LoadingSpinner';
import { useAssetStore } from '../../../stores/assetStore';
import { getSpriteUrl } from '../../../utils/spriteUrl';
import styles from './AdminMapAssetsPage.module.css';

const CATEGORIES: AssetCategory[] = ['village_marker', 'zone_tile', 'terrain_tile'];

const CATEGORY_LABELS: Record<string, string> = {
  village_marker: 'Village Markers',
  zone_tile: 'Zone Tiles',
  terrain_tile: 'Terrain Tiles',
};

const SPRITE_DIMS: Record<string, string> = {
  village_marker: '256×256',
  zone_tile: '256×256',
  terrain_tile: '256×256',
};

const VARIANT_CATEGORIES = new Set<string>(['zone_tile', 'terrain_tile']);

export function AdminMapAssetsPage() {
  const [assets, setAssets] = useState<GameAsset[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [failedSprites, setFailedSprites] = useState<Set<string>>(new Set());

  const addAssetToStore = useAssetStore((s) => s.addAsset);
  const removeAssetFromStore = useAssetStore((s) => s.removeAsset);

  const load = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const resp = await fetchGameAssets();
      setAssets(
        resp.assets.filter((a) =>
          (CATEGORIES as string[]).includes(a.category),
        ),
      );
    } catch {
      setError('Failed to load map assets.');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  const grouped = CATEGORIES.map((cat) => ({
    category: cat,
    items: assets.filter((a) => a.category === cat),
  }));

  const handleSpriteError = (id: string) => {
    setFailedSprites((prev) => new Set(prev).add(id));
  };

  const handleAddVariant = async (asset: GameAsset) => {
    setError(null);
    const existing = assets.filter(
      (a) => a.category === asset.category && a.id.startsWith(asset.id.replace(/_v\d+$/, '')),
    );
    const nextNum = existing.length + 1;
    const baseId = asset.id.replace(/_v\d+$/, '');
    const variantId = `${baseId}_v${nextNum}`;
    const variantName = `${asset.display_name.replace(/ v\d+$/, '')} v${nextNum}`;

    try {
      const newAsset = await createGameAsset({
        id: variantId,
        category: asset.category,
        display_name: variantName,
        default_icon: asset.default_icon,
      });
      setAssets((prev) => [...prev, newAsset]);
      addAssetToStore(newAsset);
      setSuccess(`Variant "${variantId}" created.`);
      setTimeout(() => setSuccess(null), 3000);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create variant.');
    }
  };

  const handleDeleteAsset = async (id: string) => {
    if (!confirm(`Delete variant "${id}"? This cannot be undone.`)) return;
    setError(null);
    try {
      await deleteGameAsset(id);
      setAssets((prev) => prev.filter((a) => a.id !== id));
      removeAssetFromStore(id);
      setSuccess(`Asset "${id}" deleted.`);
      setTimeout(() => setSuccess(null), 3000);
    } catch {
      setError('Failed to delete asset.');
    }
  };

  const isVariant = (id: string) => /_v\d+$/.test(id);

  if (loading) {
    return (
      <div className={styles.center}>
        <LoadingSpinner size="lg" />
      </div>
    );
  }

  return (
    <div className={styles.page}>
      <h2 className={styles.heading}>Map Assets</h2>
      <p className={styles.subtitle}>
        Manage village markers, zone tiles, and terrain tiles used on the world map.
        Zone and terrain tiles support variants. Drop PNGs into the sprites folder to add sprites.
      </p>

      {error && <div className={styles.error}>{error}</div>}
      {success && <div className={styles.success}>{success}</div>}

      {grouped.map(({ category, items }) => (
        <section key={category} className={styles.section}>
          <h3 className={styles.sectionTitle}>
            {CATEGORY_LABELS[category]}{' '}
            <span className={styles.dimHint}>({SPRITE_DIMS[category]} PNG)</span>
          </h3>
          <div className={styles.grid}>
            {items.map((asset) => {
              const spriteUrl = getSpriteUrl({ kind: asset.category as any, id: asset.id });
              const showSprite = spriteUrl && !failedSprites.has(asset.id);

              return (
                <div key={asset.id} className={styles.card}>
                  <div className={styles.preview}>
                    {showSprite ? (
                      <img
                        src={spriteUrl}
                        alt={asset.display_name}
                        className={styles.spriteImg}
                        onError={() => handleSpriteError(asset.id)}
                      />
                    ) : (
                      <span className={styles.emoji}>{asset.default_icon}</span>
                    )}
                  </div>
                  <div className={styles.info}>
                    <span className={styles.name}>{asset.display_name}</span>
                    <span className={styles.id}>{asset.id}</span>
                  </div>
                  <div className={styles.actions}>
                    {VARIANT_CATEGORIES.has(category) && !isVariant(asset.id) && (
                      <button
                        className={styles.variantBtn}
                        onClick={() => handleAddVariant(asset)}
                      >
                        + Variant
                      </button>
                    )}
                    {VARIANT_CATEGORIES.has(category) && isVariant(asset.id) && (
                      <button
                        className={styles.deleteAssetBtn}
                        onClick={() => handleDeleteAsset(asset.id)}
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
