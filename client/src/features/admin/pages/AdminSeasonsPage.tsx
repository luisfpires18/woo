import { useState, useEffect, useCallback } from 'react';
import {
  adminFetchSeasons,
  adminCreateSeason,
  adminDeleteSeason,
  adminLaunchSeason,
  adminEndSeason,
  adminArchiveSeason,
} from '../../../services/season';
import type { SeasonResponse, CreateSeasonRequest } from '../../../types/api';
import { LoadingSpinner } from '../../../components/LoadingSpinner/LoadingSpinner';
import { Button } from '../../../components/Button/Button';
import styles from './AdminSeasonsPage.module.css';

function StatusBadge({ status }: { status: string }) {
  const cls =
    status === 'active'
      ? styles.badgeActive
      : status === 'upcoming'
        ? styles.badgeUpcoming
        : status === 'ended'
          ? styles.badgeEnded
          : styles.badgeArchived;
  return <span className={`${styles.badge} ${cls}`}>{status}</span>;
}

const EMPTY_FORM: CreateSeasonRequest = {
  name: '',
  description: '',
  map_width: 101,
  map_height: 101,
  game_speed: 1,
  resource_multiplier: 1,
  max_villages_per_player: 5,
  weapons_of_chaos_count: 7,
};

export function AdminSeasonsPage() {
  const [seasons, setSeasons] = useState<SeasonResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [acting, setActing] = useState<number | null>(null);

  // Create form
  const [showCreate, setShowCreate] = useState(false);
  const [form, setForm] = useState<CreateSeasonRequest>(EMPTY_FORM);
  const [creating, setCreating] = useState(false);

  const load = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await adminFetchSeasons();
      setSeasons(data);
    } catch {
      setError('Failed to load seasons.');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  const clearMessages = () => {
    setError(null);
    setSuccess(null);
  };

  const handleCreate = async () => {
    if (!form.name.trim()) {
      setError('Season name is required.');
      return;
    }
    setCreating(true);
    clearMessages();
    try {
      await adminCreateSeason(form);
      setShowCreate(false);
      setForm(EMPTY_FORM);
      setSuccess('Season created.');
      await load();
    } catch {
      setError('Failed to create season.');
    } finally {
      setCreating(false);
    }
  };

  const handleAction = async (
    id: number,
    action: 'launch' | 'end' | 'archive' | 'delete',
  ) => {
    setActing(id);
    clearMessages();
    try {
      const fns = { launch: adminLaunchSeason, end: adminEndSeason, archive: adminArchiveSeason, delete: adminDeleteSeason };
      await fns[action](id);
      setSuccess(`Season ${action}${action.endsWith('e') ? 'd' : 'ed'} successfully.`);
      await load();
    } catch {
      setError(`Failed to ${action} season.`);
    } finally {
      setActing(null);
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
      <div className={styles.topRow}>
        <div>
          <h2 className={styles.heading}>Season Management</h2>
          <p className={styles.subtitle}>
            {seasons.length} season{seasons.length !== 1 ? 's' : ''}
          </p>
        </div>
        <Button size="sm" onClick={() => { clearMessages(); setShowCreate(true); }}>
          + New Season
        </Button>
      </div>

      {error && <div className={styles.error}>{error}</div>}
      {success && <div className={styles.success}>{success}</div>}

      {seasons.length === 0 ? (
        <p className={styles.empty}>No seasons yet. Create one to get started.</p>
      ) : (
        <div className={styles.tableWrap}>
          <table className={styles.table}>
            <thead>
              <tr>
                <th>ID</th>
                <th>Name</th>
                <th>Status</th>
                <th>Players</th>
                <th>Speed</th>
                <th>Map</th>
                <th>Created</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {seasons.map((s) => (
                <tr key={s.id}>
                  <td>{s.id}</td>
                  <td>{s.name}</td>
                  <td><StatusBadge status={s.status} /></td>
                  <td>{s.player_count}</td>
                  <td>{s.game_speed}x</td>
                  <td>{s.map_width}×{s.map_height}</td>
                  <td>{new Date(s.created_at).toLocaleDateString()}</td>
                  <td>
                    <div className={styles.actions}>
                      {s.status === 'upcoming' && (
                        <>
                          <button
                            className={styles.actionBtn}
                            disabled={acting === s.id}
                            onClick={() => handleAction(s.id, 'launch')}
                          >
                            Launch
                          </button>
                          <button
                            className={styles.dangerBtn}
                            disabled={acting === s.id}
                            onClick={() => handleAction(s.id, 'delete')}
                          >
                            Delete
                          </button>
                        </>
                      )}
                      {s.status === 'active' && (
                        <button
                          className={styles.actionBtn}
                          disabled={acting === s.id}
                          onClick={() => handleAction(s.id, 'end')}
                        >
                          End
                        </button>
                      )}
                      {s.status === 'ended' && (
                        <button
                          className={styles.actionBtn}
                          disabled={acting === s.id}
                          onClick={() => handleAction(s.id, 'archive')}
                        >
                          Archive
                        </button>
                      )}
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* Create season modal */}
      {showCreate && (
        <div className={styles.formOverlay} onClick={() => setShowCreate(false)}>
          <div className={styles.form} onClick={(e) => e.stopPropagation()}>
            <h3 className={styles.formTitle}>Create Season</h3>

            <div className={styles.fieldGroup}>
              <label>Name *</label>
              <input
                type="text"
                value={form.name}
                onChange={(e) => setForm({ ...form, name: e.target.value })}
                placeholder="Season 1"
              />
            </div>

            <div className={styles.fieldGroup}>
              <label>Description</label>
              <textarea
                value={form.description ?? ''}
                onChange={(e) => setForm({ ...form, description: e.target.value })}
                placeholder="Optional description..."
              />
            </div>

            <div className={styles.fieldRow}>
              <div className={styles.fieldGroup}>
                <label>Map Width</label>
                <input
                  type="number"
                  value={form.map_width}
                  onChange={(e) => setForm({ ...form, map_width: Number(e.target.value) })}
                />
              </div>
              <div className={styles.fieldGroup}>
                <label>Map Height</label>
                <input
                  type="number"
                  value={form.map_height}
                  onChange={(e) => setForm({ ...form, map_height: Number(e.target.value) })}
                />
              </div>
            </div>

            <div className={styles.fieldRow}>
              <div className={styles.fieldGroup}>
                <label>Game Speed</label>
                <input
                  type="number"
                  step="0.5"
                  value={form.game_speed ?? 1}
                  onChange={(e) => setForm({ ...form, game_speed: Number(e.target.value) })}
                />
              </div>
              <div className={styles.fieldGroup}>
                <label>Resource Multiplier</label>
                <input
                  type="number"
                  step="0.5"
                  value={form.resource_multiplier ?? 1}
                  onChange={(e) => setForm({ ...form, resource_multiplier: Number(e.target.value) })}
                />
              </div>
            </div>

            <div className={styles.fieldRow}>
              <div className={styles.fieldGroup}>
                <label>Max Villages/Player</label>
                <input
                  type="number"
                  value={form.max_villages_per_player ?? 5}
                  onChange={(e) => setForm({ ...form, max_villages_per_player: Number(e.target.value) })}
                />
              </div>
              <div className={styles.fieldGroup}>
                <label>Weapons of Chaos</label>
                <input
                  type="number"
                  value={form.weapons_of_chaos_count ?? 7}
                  onChange={(e) => setForm({ ...form, weapons_of_chaos_count: Number(e.target.value) })}
                />
              </div>
            </div>

            <div className={styles.fieldGroup}>
              <label>Planned Start Date</label>
              <input
                type="date"
                value={form.start_date ?? ''}
                onChange={(e) => setForm({ ...form, start_date: e.target.value || undefined })}
              />
            </div>

            <div className={styles.formActions}>
              <Button variant="ghost" size="sm" onClick={() => setShowCreate(false)}>
                Cancel
              </Button>
              <Button size="sm" loading={creating} onClick={handleCreate}>
                Create
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
