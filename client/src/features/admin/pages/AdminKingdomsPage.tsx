import { useState, useEffect, useCallback } from 'react';
import { fetchGameAssets } from '../../../services/admin';
import type { GameAsset } from '../../../types/api';
import { LoadingSpinner } from '../../../components/LoadingSpinner/LoadingSpinner';
import { getSpriteUrl } from '../../../utils/spriteUrl';
import styles from './AdminKingdomsPage.module.css';

export function AdminKingdomsPage() {
  const [assets, setAssets] = useState<GameAsset[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [failedSprites, setFailedSprites] = useState<Set<string>>(new Set());

  const load = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const resp = await fetchGameAssets();
      setAssets(resp.assets.filter((a) => a.category === 'kingdom_flag'));
    } catch {
      setError('Failed to load kingdom flag assets.');
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

  if (loading) {
    return (
      <div className={styles.center}>
        <LoadingSpinner size="lg" />
      </div>
    );
  }

  return (
    <div className={styles.page}>
      <h2 className={styles.heading}>Kingdom Flags</h2>
      <p className={styles.subtitle}>
        Kingdom flag sprites (256×256 PNG). Place files in{' '}
        <code>uploads/sprites/flags/&#123;kingdom&#125;.png</code> to add flags.
      </p>

      {error && <div className={styles.error}>{error}</div>}

      <div className={styles.grid}>
        {assets.map((asset) => {
          const spriteUrl = getSpriteUrl({ kind: 'kingdom_flag', id: asset.id });
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
            </div>
          );
        })}
      </div>
    </div>
  );
}
