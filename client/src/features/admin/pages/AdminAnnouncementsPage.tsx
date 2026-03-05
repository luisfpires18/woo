import { useState, useEffect, useCallback, type FormEvent } from 'react';
import {
  fetchAnnouncements,
  createAnnouncement,
  deleteAnnouncement,
} from '../../../services/admin';
import type { AnnouncementResponse } from '../../../types/api';
import { LoadingSpinner } from '../../../components/LoadingSpinner/LoadingSpinner';
import styles from './AdminAnnouncementsPage.module.css';

export function AdminAnnouncementsPage() {
  const [announcements, setAnnouncements] = useState<AnnouncementResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [deleting, setDeleting] = useState<number | null>(null);

  // Form state
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [expiresAt, setExpiresAt] = useState('');
  const [submitting, setSubmitting] = useState(false);

  const load = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await fetchAnnouncements();
      setAnnouncements(data);
    } catch {
      setError('Failed to load announcements.');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  const handleCreate = async (e: FormEvent) => {
    e.preventDefault();
    if (!title.trim() || !content.trim()) return;

    setSubmitting(true);
    setError(null);
    try {
      const newAnn = await createAnnouncement({
        title: title.trim(),
        content: content.trim(),
        expires_at: expiresAt || undefined,
      });
      setAnnouncements((prev) => [newAnn, ...prev]);
      setTitle('');
      setContent('');
      setExpiresAt('');
    } catch {
      setError('Failed to create announcement.');
    } finally {
      setSubmitting(false);
    }
  };

  const handleDelete = async (id: number) => {
    setDeleting(id);
    setError(null);
    try {
      await deleteAnnouncement(id);
      setAnnouncements((prev) => prev.filter((a) => a.id !== id));
    } catch {
      setError('Failed to delete announcement.');
    } finally {
      setDeleting(null);
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
      <h2 className={styles.heading}>Announcements</h2>
      <p className={styles.subtitle}>
        Broadcast messages to all players.
      </p>

      {error && <div className={styles.error}>{error}</div>}

      {/* Create form */}
      <form onSubmit={handleCreate} className={styles.form}>
        <h3 className={styles.formTitle}>New Announcement</h3>
        <input
          type="text"
          placeholder="Title"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          className={styles.input}
          required
          disabled={submitting}
        />
        <textarea
          placeholder="Content"
          value={content}
          onChange={(e) => setContent(e.target.value)}
          className={styles.textarea}
          rows={3}
          required
          disabled={submitting}
        />
        <div className={styles.formRow}>
          <label className={styles.label}>
            <span>Expires at (optional)</span>
            <input
              type="datetime-local"
              value={expiresAt}
              onChange={(e) => setExpiresAt(e.target.value)}
              className={styles.input}
              disabled={submitting}
            />
          </label>
          <button type="submit" disabled={submitting} className={styles.submitBtn}>
            {submitting ? 'Creating…' : 'Create'}
          </button>
        </div>
      </form>

      {/* List */}
      <div className={styles.list}>
        {announcements.length === 0 ? (
          <p className={styles.empty}>No active announcements.</p>
        ) : (
          announcements.map((ann) => (
            <div key={ann.id} className={styles.card}>
              <div className={styles.cardHeader}>
                <h4 className={styles.cardTitle}>{ann.title}</h4>
                <button
                  onClick={() => handleDelete(ann.id)}
                  disabled={deleting === ann.id}
                  className={styles.deleteBtn}
                >
                  {deleting === ann.id ? '…' : '✕'}
                </button>
              </div>
              <p className={styles.cardContent}>{ann.content}</p>
              <div className={styles.cardMeta}>
                <span>Created {new Date(ann.created_at).toLocaleString()}</span>
                {ann.expires_at && (
                  <span>Expires {new Date(ann.expires_at).toLocaleString()}</span>
                )}
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
