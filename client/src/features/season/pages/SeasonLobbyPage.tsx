import { useState, useEffect, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { useSeasonStore } from '../../../stores/seasonStore';
import { joinSeason } from '../../../services/season';
import { Button } from '../../../components/Button/Button';
import { LoadingSpinner } from '../../../components/LoadingSpinner/LoadingSpinner';
import { ApiRequestError } from '../../../services/api';
import type { SeasonResponse } from '../../../types/api';
import { KINGDOMS } from '../../kingdom/components/KingdomCard';
import styles from './SeasonLobbyPage.module.css';

const PLAYABLE_KINGDOMS = KINGDOMS.filter((k) => k.playable);

const STATUS_ORDER = ['active', 'upcoming', 'ended', 'archived'] as const;

const STATUS_LABELS: Record<string, string> = {
  active: 'Active Seasons',
  upcoming: 'Upcoming Seasons',
  ended: 'Ended Seasons',
  archived: 'Archived Seasons',
};

function StatusBadge({ status }: { status: string }) {
  const badgeClass =
    status === 'active'
      ? styles.badgeActive
      : status === 'upcoming'
        ? styles.badgeUpcoming
        : status === 'ended'
          ? styles.badgeEnded
          : styles.badgeArchived;

  return <span className={`${styles.badge} ${badgeClass}`}>{status}</span>;
}

export function SeasonLobbyPage() {
  const { seasons, mySeasons, loaded, loading, loadSeasons, loadMySeasons } =
    useSeasonStore();
  const navigate = useNavigate();

  const [joining, setJoining] = useState(false);
  const [selectedSeason, setSelectedSeason] = useState<SeasonResponse | null>(null);
  const [selectedKingdom, setSelectedKingdom] = useState('');
  const [error, setError] = useState('');

  const load = useCallback(async () => {
    await Promise.all([loadSeasons(), loadMySeasons()]);
  }, [loadSeasons, loadMySeasons]);

  useEffect(() => {
    load();
  }, [load]);

  const mySeasonIds = new Set(mySeasons.map((s) => s.id));

  const handleJoinClick = (season: SeasonResponse) => {
    setSelectedSeason(season);
    setSelectedKingdom('');
    setError('');
  };

  const handleCloseModal = () => {
    setSelectedSeason(null);
    setSelectedKingdom('');
    setError('');
  };

  const handleJoinConfirm = async () => {
    if (!selectedSeason || !selectedKingdom) return;
    setJoining(true);
    setError('');

    try {
      const resp = await joinSeason(selectedSeason.id, selectedKingdom);
      // Refresh seasons list
      await load();
      setSelectedSeason(null);
      // Navigate to the new village
      if (resp.village_id) {
        navigate(`/village/${resp.village_id}`);
      }
    } catch (err) {
      if (err instanceof ApiRequestError) {
        setError(err.message);
      } else {
        setError('Failed to join season');
      }
    } finally {
      setJoining(false);
    }
  };

  // Group seasons by status
  const grouped = STATUS_ORDER.map((status) => ({
    status,
    label: STATUS_LABELS[status],
    items: seasons.filter((s) => s.status === status),
  })).filter((g) => g.items.length > 0);

  if (!loaded && loading) {
    return (
      <div className={styles.center}>
        <LoadingSpinner size="lg" />
      </div>
    );
  }

  return (
    <div className={styles.page}>
      <header className={styles.header}>
        <h1 className={styles.title}>Seasons</h1>
        <p className={styles.subtitle}>
          Choose your battlefield. Each season is a new world with its own leaderboard.
        </p>
      </header>

      {grouped.length === 0 && (
        <p className={styles.empty}>No seasons available yet. Check back soon!</p>
      )}

      {grouped.map((group) => (
        <section key={group.status} className={styles.section}>
          <h2 className={styles.sectionTitle}>{group.label}</h2>
          <div className={styles.grid}>
            {group.items.map((season) => {
              const joined = mySeasonIds.has(season.id);
              return (
                <article key={season.id} className={styles.card}>
                  <div className={styles.cardHeader}>
                    <h3 className={styles.cardName}>{season.name}</h3>
                    <StatusBadge status={season.status} />
                  </div>

                  {season.description && (
                    <p className={styles.cardDesc}>{season.description}</p>
                  )}

                  <div className={styles.cardMeta}>
                    <span>Speed: {season.game_speed}x</span>
                    <span>Resources: {season.resource_multiplier}x</span>
                    <span>
                      Map: {season.map_width}×{season.map_height}
                    </span>
                    {season.start_date && (
                      <span>
                        Starts: {new Date(season.start_date).toLocaleDateString()}
                      </span>
                    )}
                  </div>

                  <div className={styles.cardFooter}>
                    <span className={styles.playerCount}>
                      {season.player_count} player{season.player_count !== 1 ? 's' : ''}
                    </span>

                    {joined ? (
                      <span className={styles.joinedBadge}>Joined</span>
                    ) : season.status === 'active' ? (
                      <Button
                        size="sm"
                        onClick={() => handleJoinClick(season)}
                      >
                        Join
                      </Button>
                    ) : null}
                  </div>
                </article>
              );
            })}
          </div>
        </section>
      ))}

      {/* Kingdom selection modal for joining */}
      {selectedSeason && (
        <div className={styles.modalOverlay} onClick={handleCloseModal}>
          <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
            <h2 className={styles.modalTitle}>
              Join {selectedSeason.name}
            </h2>
            <p style={{ color: 'var(--text-muted)', fontSize: '0.875rem', margin: 0 }}>
              Choose your kingdom for this season. This cannot be changed later.
            </p>

            {error && <p className={styles.error}>{error}</p>}

            <div className={styles.kingdomOptions}>
              {PLAYABLE_KINGDOMS.map((k) => (
                <button
                  key={k.id}
                  className={`${styles.kingdomOption} ${selectedKingdom === k.id ? styles.kingdomOptionSelected : ''}`}
                  onClick={() => setSelectedKingdom(k.id)}
                >
                  {k.name}
                </button>
              ))}
            </div>

            <div className={styles.modalActions}>
              <Button variant="ghost" size="sm" onClick={handleCloseModal}>
                Cancel
              </Button>
              <Button
                size="sm"
                disabled={!selectedKingdom}
                loading={joining}
                onClick={handleJoinConfirm}
              >
                Confirm
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
