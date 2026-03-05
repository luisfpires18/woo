import { useState, useEffect, useCallback } from 'react';
import { fetchWorldConfig, setWorldConfig } from '../../../services/admin';
import type { WorldConfigEntry } from '../../../types/api';
import { LoadingSpinner } from '../../../components/LoadingSpinner/LoadingSpinner';
import styles from './AdminConfigPage.module.css';

export function AdminConfigPage() {
  const [configs, setConfigs] = useState<WorldConfigEntry[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [editing, setEditing] = useState<Record<string, string>>({});
  const [saving, setSaving] = useState<string | null>(null);

  const load = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const resp = await fetchWorldConfig();
      setConfigs(resp.configs);
    } catch {
      setError('Failed to load configuration.');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  const handleEdit = (key: string, value: string) => {
    setEditing((prev) => ({ ...prev, [key]: value }));
  };

  const handleSave = async (key: string) => {
    const newValue = editing[key];
    if (newValue === undefined) return;

    setSaving(key);
    setError(null);
    try {
      await setWorldConfig(key, newValue);
      setConfigs((prev) =>
        prev.map((c) =>
          c.key === key
            ? { ...c, value: newValue, updated_at: new Date().toISOString() }
            : c,
        ),
      );
      setEditing((prev) => {
        const next = { ...prev };
        delete next[key];
        return next;
      });
    } catch {
      setError(`Failed to save "${key}".`);
    } finally {
      setSaving(null);
    }
  };

  const handleCancel = (key: string) => {
    setEditing((prev) => {
      const next = { ...prev };
      delete next[key];
      return next;
    });
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
      <h2 className={styles.heading}>World Configuration</h2>
      <p className={styles.subtitle}>
        Adjust game-world settings. Changes take effect immediately.
      </p>

      {error && <div className={styles.error}>{error}</div>}

      <div className={styles.grid}>
        {configs.map((cfg) => {
          const isEditing = editing[cfg.key] !== undefined;
          const isSaving = saving === cfg.key;

          return (
            <div key={cfg.key} className={styles.card}>
              <div className={styles.cardHeader}>
                <span className={styles.configKey}>{cfg.key}</span>
                <span className={styles.updatedAt}>
                  Updated {new Date(cfg.updated_at).toLocaleDateString()}
                </span>
              </div>

              {cfg.description && (
                <p className={styles.description}>{cfg.description}</p>
              )}

              <div className={styles.valueRow}>
                {isEditing ? (
                  <>
                    <input
                      type="text"
                      value={editing[cfg.key]}
                      onChange={(e) => handleEdit(cfg.key, e.target.value)}
                      className={styles.input}
                      disabled={isSaving}
                      autoFocus
                    />
                    <button
                      onClick={() => handleSave(cfg.key)}
                      disabled={isSaving}
                      className={styles.saveBtn}
                    >
                      {isSaving ? 'Saving…' : 'Save'}
                    </button>
                    <button
                      onClick={() => handleCancel(cfg.key)}
                      disabled={isSaving}
                      className={styles.cancelBtn}
                    >
                      Cancel
                    </button>
                  </>
                ) : (
                  <>
                    <span className={styles.value}>{cfg.value}</span>
                    <button
                      onClick={() => handleEdit(cfg.key, cfg.value)}
                      className={styles.editBtn}
                    >
                      Edit
                    </button>
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
