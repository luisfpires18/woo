import { useState, useEffect, useCallback } from 'react';
import { fetchPlayers, updatePlayerRole } from '../../../services/admin';
import type { PlayerListItem } from '../../../types/api';
import { LoadingSpinner } from '../../../components/LoadingSpinner/LoadingSpinner';
import styles from './AdminPlayersPage.module.css';

const PAGE_SIZE = 20;

export function AdminPlayersPage() {
  const [players, setPlayers] = useState<PlayerListItem[]>([]);
  const [total, setTotal] = useState(0);
  const [offset, setOffset] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [updating, setUpdating] = useState<number | null>(null);

  const load = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const resp = await fetchPlayers(offset, PAGE_SIZE);
      setPlayers(resp.players);
      setTotal(resp.total);
    } catch {
      setError('Failed to load players.');
    } finally {
      setLoading(false);
    }
  }, [offset]);

  useEffect(() => {
    load();
  }, [load]);

  const handleRoleChange = async (id: number, newRole: 'player' | 'admin') => {
    setUpdating(id);
    try {
      await updatePlayerRole(id, newRole);
      setPlayers((prev) =>
        prev.map((p) => (p.id === id ? { ...p, role: newRole } : p)),
      );
    } catch {
      setError('Failed to update role.');
    } finally {
      setUpdating(null);
    }
  };

  const totalPages = Math.ceil(total / PAGE_SIZE);
  const currentPage = Math.floor(offset / PAGE_SIZE) + 1;

  if (loading) {
    return (
      <div className={styles.center}>
        <LoadingSpinner size="lg" />
      </div>
    );
  }

  return (
    <div className={styles.page}>
      <h2 className={styles.heading}>Player Management</h2>
      <p className={styles.subtitle}>{total} registered player{total !== 1 ? 's' : ''}</p>

      {error && <div className={styles.error}>{error}</div>}

      <div className={styles.tableWrap}>
        <table className={styles.table}>
          <thead>
            <tr>
              <th>ID</th>
              <th>Username</th>
              <th>Email</th>
              <th>Kingdom</th>
              <th>Role</th>
              <th>Joined</th>
            </tr>
          </thead>
          <tbody>
            {players.map((p) => (
              <tr key={p.id}>
                <td className={styles.id}>{p.id}</td>
                <td>{p.username}</td>
                <td className={styles.email}>{p.email}</td>
                <td className={styles.kingdom}>{p.kingdom}</td>
                <td>
                  <select
                    value={p.role}
                    disabled={updating === p.id}
                    onChange={(e) =>
                      handleRoleChange(p.id, e.target.value as 'player' | 'admin')
                    }
                    className={`${styles.roleSelect} ${p.role === 'admin' ? styles.adminRole : ''}`}
                  >
                    <option value="player">player</option>
                    <option value="admin">admin</option>
                  </select>
                </td>
                <td className={styles.date}>
                  {new Date(p.created_at).toLocaleDateString()}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {totalPages > 1 && (
        <div className={styles.pagination}>
          <button
            disabled={offset === 0}
            onClick={() => setOffset((o) => Math.max(0, o - PAGE_SIZE))}
            className={styles.pageBtn}
          >
            ← Prev
          </button>
          <span className={styles.pageInfo}>
            Page {currentPage} of {totalPages}
          </span>
          <button
            disabled={offset + PAGE_SIZE >= total}
            onClick={() => setOffset((o) => o + PAGE_SIZE)}
            className={styles.pageBtn}
          >
            Next →
          </button>
        </div>
      )}
    </div>
  );
}
