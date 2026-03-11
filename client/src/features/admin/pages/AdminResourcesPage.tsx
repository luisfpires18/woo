import { useState, useEffect, useCallback } from 'react';
import {
  fetchGameAssets,
  fetchResourceBuildingConfigs,
  updateResourceBuildingConfig,
  fetchKingdomBuildingSprites,
} from '../../../services/admin';
import type { GameAsset, ResourceBuildingConfig, BuildingSpriteInfo } from '../../../types/api';
import { LoadingSpinner } from '../../../components/LoadingSpinner/LoadingSpinner';
import { getSpriteUrl } from '../../../utils/spriteUrl';
import styles from './AdminResourcesPage.module.css';

const KINGDOMS = [
  'arkazia',
  'drakanith',
  'draxys',
  'lumus',
  'nordalh',
  'sylvara',
  'veridor',
  'zandres',
];

interface EditingState {
  display_name: string;
  description: string;
  default_icon: string;
}

/** Build a lookup key for sprite info: "food_1" */
function spriteMapKey(resourceType: string, slot: number): string {
  return `${resourceType}_${slot}`;
}

export function AdminResourcesPage() {
  // --- Resource game assets state ---
  const [assets, setAssets] = useState<GameAsset[]>([]);
  const [assetsLoading, setAssetsLoading] = useState(true);

  // --- Resource building configs state ---
  const [configs, setConfigs] = useState<ResourceBuildingConfig[]>([]);
  const [configsLoading, setConfigsLoading] = useState(true);
  const [kingdom, setKingdom] = useState<string>('arkazia');
  const [editing, setEditing] = useState<Record<number, EditingState>>({});
  const [saving, setSaving] = useState<number | null>(null);

  // --- Sprite listing state (per kingdom) ---
  const [spriteMap, setSpriteMap] = useState<Record<string, BuildingSpriteInfo>>({});

  // --- Shared state ---
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [failedSprites, setFailedSprites] = useState<Set<string>>(new Set());

  // --- Load resource game assets ---
  const loadAssets = useCallback(async () => {
    setAssetsLoading(true);
    try {
      const resp = await fetchGameAssets();
      setAssets(resp.assets.filter((a) => a.category === 'resource'));
    } catch {
      setError('Failed to load resource assets.');
    } finally {
      setAssetsLoading(false);
    }
  }, []);

  // --- Load resource building configs ---
  const loadConfigs = useCallback(async () => {
    setConfigsLoading(true);
    try {
      const resp = await fetchResourceBuildingConfigs();
      setConfigs(resp.configs);
    } catch {
      setError('Failed to load resource building configs.');
    } finally {
      setConfigsLoading(false);
    }
  }, []);

  // --- Load sprite listing for the selected kingdom ---
  const loadSprites = useCallback(async (k: string) => {
    try {
      const resp = await fetchKingdomBuildingSprites(k);
      const map: Record<string, BuildingSpriteInfo> = {};
      for (const s of resp.sprites) {
        map[spriteMapKey(s.resource_type, s.slot)] = s;
      }
      setSpriteMap(map);
    } catch {
      setSpriteMap({});
    }
  }, []);

  useEffect(() => {
    loadAssets();
    loadConfigs();
  }, [loadAssets, loadConfigs]);

  // Re-fetch sprites when kingdom changes
  useEffect(() => {
    loadSprites(kingdom);
    setFailedSprites(new Set());
  }, [kingdom, loadSprites]);

  const loading = assetsLoading || configsLoading;
  const filtered = configs.filter((c) => c.kingdom === kingdom);

  const handleSpriteError = (key: string) => {
    setFailedSprites((prev) => new Set(prev).add(key));
  };

  // ── Config inline editing handlers ──

  const startEdit = (cfg: ResourceBuildingConfig) => {
    setEditing((prev) => ({
      ...prev,
      [cfg.id]: {
        display_name: cfg.display_name,
        description: cfg.description,
        default_icon: cfg.default_icon,
      },
    }));
  };

  const cancelEdit = (id: number) => {
    setEditing((prev) => {
      const next = { ...prev };
      delete next[id];
      return next;
    });
  };

  const handleChange = (id: number, field: keyof EditingState, value: string) => {
    setEditing((prev) => {
      const current = prev[id] ?? { display_name: '', description: '', default_icon: '' };
      return { ...prev, [id]: { ...current, [field]: value } };
    });
  };

  const handleSave = async (id: number) => {
    const edit = editing[id];
    if (!edit) return;

    setSaving(id);
    setError(null);
    setSuccess(null);
    try {
      await updateResourceBuildingConfig(id, {
        display_name: edit.display_name,
        description: edit.description,
        default_icon: edit.default_icon,
      });
      setConfigs((prev) =>
        prev.map((c) =>
          c.id === id
            ? {
                ...c,
                display_name: edit.display_name,
                description: edit.description,
                default_icon: edit.default_icon,
                updated_at: new Date().toISOString(),
              }
            : c,
        ),
      );
      cancelEdit(id);
      setSuccess('Config updated.');
      setTimeout(() => setSuccess(null), 3000);
    } catch {
      setError('Failed to save config.');
    } finally {
      setSaving(null);
    }
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
      <h2 className={styles.heading}>Resources</h2>
      <p className={styles.subtitle}>
        Manage resource icons and per-kingdom resource building display settings.
        Place resource sprites
        in <code>uploads/sprites/resources/&#123;resource_type&#125;.png</code> and building sprites
        in <code>uploads/sprites/kingdoms/&#123;kingdom&#125;/buildings/&#123;kingdom&#125;_&#123;resource&#125;_&#123;slot&#125;_name.png</code>.
      </p>

      {error && <div className={styles.error}>{error}</div>}
      {success && <div className={styles.success}>{success}</div>}

      {/* Section A: Resource game assets (icons) */}
      <section className={styles.section}>
        <h3 className={styles.sectionTitle}>
          Resource Icons <span className={styles.dimHint}>(32×32 PNG)</span>
        </h3>
        <div className={styles.assetGrid}>
          {assets.map((asset) => {
            const spriteUrl = getSpriteUrl({ kind: 'resource', id: asset.id });
            const showSprite = spriteUrl && !failedSprites.has(asset.id);

            return (
              <div key={asset.id} className={styles.assetCard}>
                <div className={styles.assetPreview}>
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
                <div className={styles.assetInfo}>
                  <span className={styles.assetName}>{asset.display_name}</span>
                  <span className={styles.assetId}>{asset.id}</span>
                </div>
              </div>
            );
          })}
        </div>
      </section>

      {/* Section B: Resource building configs (per kingdom) */}
      <section className={styles.section}>
        <h3 className={styles.sectionTitle}>
          Resource Building Names{' '}
          <span className={styles.dimHint}>(per kingdom)</span>
        </h3>

        <div className={styles.filterRow}>
          <span className={styles.filterLabel}>Kingdom:</span>
          {KINGDOMS.map((k) => (
            <button
              key={k}
              className={`${styles.filterBtn} ${kingdom === k ? styles.filterBtnActive : ''}`}
              onClick={() => setKingdom(k)}
            >
              {k.charAt(0).toUpperCase() + k.slice(1)}
            </button>
          ))}
        </div>

        <div className={styles.configGrid}>
          {filtered.map((cfg) => {
            const edit = editing[cfg.id];
            const isEditing = !!edit;
            const isSaving = saving === cfg.id;
            const key = spriteMapKey(cfg.resource_type, cfg.slot);
            const spriteInfo = spriteMap[key];
            const spriteUrl = spriteInfo?.url;
            const showSprite = spriteUrl && !failedSprites.has(key);

            return (
              <div key={cfg.id} className={styles.configCard}>
                {/* Hero sprite area */}
                <div className={styles.cardSpriteHero}>
                  {showSprite ? (
                    <img
                      src={spriteUrl}
                      alt={cfg.display_name}
                      className={styles.heroImg}
                      onError={() => handleSpriteError(key)}
                    />
                  ) : (
                    <div className={styles.spritePlaceholder}>
                      <span className={styles.placeholderEmoji}>
                        {isEditing ? edit.default_icon : cfg.default_icon}
                      </span>
                      <span className={styles.placeholderHint}>
                        {cfg.kingdom}_{cfg.resource_type}_{cfg.slot}_name.png
                      </span>
                    </div>
                  )}
                </div>

                {/* Sprite filename badge */}
                {spriteInfo && (
                  <span className={styles.spriteFilename}>{spriteInfo.filename}</span>
                )}

                <div className={styles.cardBody}>
                  <div className={styles.cardHeader}>
                    <span className={styles.configType}>
                      {cfg.resource_type} #{cfg.slot}
                    </span>
                    <span className={styles.kingdom}>{cfg.kingdom}</span>
                  </div>

                  {isEditing ? (
                    <>
                      <div className={styles.fieldGroup}>
                        <label className={styles.fieldLabel}>Display Name</label>
                        <input
                          type="text"
                          className={styles.input}
                          value={edit.display_name}
                          onChange={(e) => handleChange(cfg.id, 'display_name', e.target.value)}
                          disabled={isSaving}
                        />
                      </div>
                      <div className={styles.fieldGroup}>
                        <label className={styles.fieldLabel}>Description</label>
                        <textarea
                          className={styles.textarea}
                          value={edit.description}
                          onChange={(e) => handleChange(cfg.id, 'description', e.target.value)}
                          disabled={isSaving}
                        />
                      </div>
                      <div className={styles.fieldGroup}>
                        <label className={styles.fieldLabel}>Icon (emoji fallback)</label>
                        <input
                          type="text"
                          className={styles.input}
                          value={edit.default_icon}
                          onChange={(e) => handleChange(cfg.id, 'default_icon', e.target.value)}
                          disabled={isSaving}
                        />
                      </div>
                      <div className={styles.formActions}>
                        <button
                          className={styles.saveBtn}
                          onClick={() => handleSave(cfg.id)}
                          disabled={isSaving || !edit.display_name.trim()}
                        >
                          {isSaving ? 'Saving...' : 'Save'}
                        </button>
                        <button
                          className={styles.cancelBtn}
                          onClick={() => cancelEdit(cfg.id)}
                          disabled={isSaving}
                        >
                          Cancel
                        </button>
                      </div>
                    </>
                  ) : (
                    <>
                      <div className={styles.fieldGroup}>
                        <span className={styles.fieldLabel}>Display Name</span>
                        <span>{cfg.display_name}</span>
                      </div>
                      {cfg.description && (
                        <div className={styles.fieldGroup}>
                          <span className={styles.fieldLabel}>Description</span>
                          <span>{cfg.description}</span>
                        </div>
                      )}
                      <div className={styles.formActions}>
                        <button className={styles.saveBtn} onClick={() => startEdit(cfg)}>
                          Edit
                        </button>
                      </div>
                    </>
                  )}
                </div>
              </div>
            );
          })}
        </div>
      </section>
    </div>
  );
}
