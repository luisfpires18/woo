import { useState, useEffect, useCallback } from 'react';
import {
  fetchBuildingDisplayConfigs,
  updateBuildingDisplayConfig,
  fetchKingdomDisplayBuildingSprites,
} from '../../../services/admin';
import type { BuildingDisplayConfig, DisplayBuildingSpriteInfo } from '../../../types/api';
import { LoadingSpinner } from '../../../components/LoadingSpinner/LoadingSpinner';
import { useAssetStore, buildingConfigToAsset } from '../../../stores/assetStore';
import styles from './AdminBuildingsPage.module.css';

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

export function AdminBuildingsPage() {
  const [configs, setConfigs] = useState<BuildingDisplayConfig[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [kingdom, setKingdom] = useState<string>('arkazia');
  const [editing, setEditing] = useState<Record<number, EditingState>>({});
  const [saving, setSaving] = useState<number | null>(null);
  const [failedSprites, setFailedSprites] = useState<Set<string>>(new Set());
  const [spriteMap, setSpriteMap] = useState<Record<string, DisplayBuildingSpriteInfo>>({});
  const upsertAsset = useAssetStore((s) => s.upsert);
  const addOrUpdateAsset = useAssetStore((s) => s.addAsset);
  const getAsset = useAssetStore((s) => s.getById);

  /** Push a building config change into the asset store so GameIcon updates everywhere. */
  const syncToAssetStore = useCallback((cfg: BuildingDisplayConfig) => {
    const asset = buildingConfigToAsset(cfg);
    if (getAsset(asset.id)) {
      upsertAsset(asset);
    } else {
      addOrUpdateAsset(asset);
    }
  }, [upsertAsset, addOrUpdateAsset, getAsset]);

  const load = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const resp = await fetchBuildingDisplayConfigs();
      setConfigs(resp.configs);
    } catch {
      setError('Failed to load building display configs.');
    } finally {
      setLoading(false);
    }
  }, []);

  const loadSprites = useCallback(async (k: string) => {
    try {
      const resp = await fetchKingdomDisplayBuildingSprites(k);
      const map: Record<string, DisplayBuildingSpriteInfo> = {};
      for (const s of resp.sprites) {
        map[s.building_type] = s;
      }
      setSpriteMap(map);
    } catch {
      setSpriteMap({});
    }
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  useEffect(() => {
    loadSprites(kingdom);
    setFailedSprites(new Set());
  }, [kingdom, loadSprites]);

  const filtered = configs.filter((c) => c.kingdom === kingdom);

  const startEdit = (cfg: BuildingDisplayConfig) => {
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

  const handleChange = (
    id: number,
    field: keyof EditingState,
    value: string,
  ) => {
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
      await updateBuildingDisplayConfig(id, {
        display_name: edit.display_name,
        description: edit.description,
        default_icon: edit.default_icon,
      });
      const updatedAt = new Date().toISOString();
      setConfigs((prev) => {
        const updated = prev.map((c) =>
          c.id === id
            ? {
                ...c,
                display_name: edit.display_name,
                description: edit.description,
                default_icon: edit.default_icon,
                updated_at: updatedAt,
              }
            : c,
        );
        const cfg = updated.find((c) => c.id === id);
        if (cfg) syncToAssetStore(cfg);
        return updated;
      });
      cancelEdit(id);
      setSuccess('Building config updated.');
      setTimeout(() => setSuccess(null), 3000);
    } catch {
      setError('Failed to save building config.');
    } finally {
      setSaving(null);
    }
  };

  const handleSpriteError = (key: string) => {
    setFailedSprites((prev) => new Set(prev).add(key));
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
      <h2 className={styles.heading}>Building Display Names</h2>
      <p className={styles.subtitle}>
        Customise building names, descriptions, and icons per kingdom. Place sprites
        in <code>uploads/sprites/kingdoms/&#123;kingdom&#125;/buildings/&#123;kingdom&#125;_&#123;building_type&#125;_name.png</code>.
      </p>

      {error && <div className={styles.error}>{error}</div>}
      {success && <div className={styles.success}>{success}</div>}

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
          const spriteInfo = spriteMap[cfg.building_type];
          const spriteUrl = spriteInfo?.url;
          const showSprite = spriteUrl && !failedSprites.has(cfg.building_type);

          return (
            <div key={cfg.id} className={styles.configCard}>
              {/* Hero sprite area */}
              <div className={styles.cardSpriteHero}>
                {showSprite ? (
                  <img
                    src={spriteUrl}
                    alt={cfg.display_name}
                    className={styles.heroImg}
                    onError={() => handleSpriteError(cfg.building_type)}
                  />
                ) : (
                  <div className={styles.spritePlaceholder}>
                    <span className={styles.placeholderEmoji}>
                      {isEditing ? edit.default_icon : cfg.default_icon}
                    </span>
                    <span className={styles.placeholderHint}>
                      {cfg.kingdom}_{cfg.building_type}_name.png
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
                  <span className={styles.buildingType}>
                    {cfg.building_type.replace(/_/g, ' ')}
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
                        onChange={(e) =>
                          handleChange(cfg.id, 'display_name', e.target.value)
                        }
                        disabled={isSaving}
                      />
                    </div>
                    <div className={styles.fieldGroup}>
                      <label className={styles.fieldLabel}>Description</label>
                      <textarea
                        className={styles.textarea}
                        value={edit.description}
                        onChange={(e) =>
                          handleChange(cfg.id, 'description', e.target.value)
                        }
                        disabled={isSaving}
                      />
                    </div>
                    <div className={styles.fieldGroup}>
                      <label className={styles.fieldLabel}>Icon (emoji fallback)</label>
                      <input
                        type="text"
                        className={styles.input}
                        value={edit.default_icon}
                        onChange={(e) =>
                          handleChange(cfg.id, 'default_icon', e.target.value)
                        }
                        disabled={isSaving}
                      />
                    </div>
                    <div className={styles.formActions}>
                      <button
                        className={styles.saveBtn}
                        onClick={() => handleSave(cfg.id)}
                        disabled={isSaving || !edit.display_name.trim()}
                      >
                        {isSaving ? 'Saving…' : 'Save'}
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
                      <button
                        className={styles.saveBtn}
                        onClick={() => startEdit(cfg)}
                      >
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
    </div>
  );
}
