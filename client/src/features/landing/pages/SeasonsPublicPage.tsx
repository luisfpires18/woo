import { useState, useEffect, useCallback } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuthStore } from '../../../stores/authStore';
import { fetchPublicSeasons, joinSeason, fetchMySeasons } from '../../../services/season';
import { Button } from '../../../components/Button/Button';
import { LoadingSpinner } from '../../../components/LoadingSpinner/LoadingSpinner';
import { ApiRequestError } from '../../../services/api';
import type { SeasonResponse, PlayerInfo } from '../../../types/api';
import { KingdomCard, KINGDOMS } from '../../kingdom/components/KingdomCard';
import styles from './SeasonsPublicPage.module.css';

const PLAYABLE_KINGDOMS = KINGDOMS.filter((k) => k.playable);

function StatusBadge({ status }: { status: string }) {
  const cls =
    status === 'active'
      ? styles.badgeActive
      : status === 'upcoming'
        ? styles.badgeUpcoming
        : styles.badgeEnded;
  return <span className={`${styles.badge} ${cls}`}>{status}</span>;
}

function SeasonCard({
  season,
  joined,
  isAuthenticated,
  onJoinClick,
}: {
  season: SeasonResponse;
  joined: boolean;
  isAuthenticated: boolean;
  onJoinClick: (s: SeasonResponse) => void;
}) {
  return (
    <article className={styles.seasonCard}>
      <div className={styles.seasonCardHeader}>
        <h3 className={styles.seasonName}>{season.name}</h3>
        <StatusBadge status={season.status} />
      </div>

      {season.description && (
        <p className={styles.seasonDesc}>{season.description}</p>
      )}

      <div className={styles.seasonMeta}>
        <span>Speed: {season.game_speed}x</span>
        <span>Resources: {season.resource_multiplier}x</span>
        <span>Map: {season.map_width}&times;{season.map_height}</span>
        {season.start_date && (
          <span>Starts: {new Date(season.start_date).toLocaleDateString()}</span>
        )}
      </div>

      <div className={styles.seasonFooter}>
        <span className={styles.playerCount}>
          {season.player_count} player{season.player_count !== 1 ? 's' : ''}
        </span>

        {isAuthenticated ? (
          joined ? (
            <span className={styles.joinedBadge}>Joined</span>
          ) : season.status === 'active' ? (
            <Button size="sm" onClick={() => onJoinClick(season)}>
              Join Season
            </Button>
          ) : null
        ) : (
          <Link to="/login" className={styles.authLink}>Login to Join</Link>
        )}
      </div>
    </article>
  );
}

export function SeasonsPublicPage() {
  const { isAuthenticated } = useAuthStore();
  const setPlayer = useAuthStore((s) => s.setPlayer);
  const player = useAuthStore((s) => s.player);
  const navigate = useNavigate();

  const [seasons, setSeasons] = useState<SeasonResponse[]>([]);
  const [mySeasonIds, setMySeasonIds] = useState<Set<number>>(new Set());
  const [loading, setLoading] = useState(true);

  // Join-season modal state
  const [selectedSeason, setSelectedSeason] = useState<SeasonResponse | null>(null);
  const [selectedKingdom, setSelectedKingdom] = useState('');
  const [joining, setJoining] = useState(false);
  const [error, setError] = useState('');

  const load = useCallback(async (retries = 1) => {
    setLoading(true);
    setError('');
    try {
      const all = await fetchPublicSeasons();
      setSeasons(all);
      if (isAuthenticated) {
        try {
          const mine = await fetchMySeasons();
          setMySeasonIds(new Set(mine.map((s) => s.id)));
        } catch {
          // My-seasons can fail without blocking the page
        }
      }
    } catch (err) {
      // Retry once after a short delay (handles startup race condition)
      if (retries > 0) {
        await new Promise((r) => setTimeout(r, 1500));
        return load(retries - 1);
      }
      console.error('[SeasonsPublicPage] Failed to load seasons:', err);
      setError('Failed to load seasons. Please refresh the page.');
    } finally {
      setLoading(false);
    }
  }, [isAuthenticated]);

  useEffect(() => {
    load();
  }, [load]);

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
      setSelectedSeason(null);

      // Update player kingdom in auth store so ProtectedLayout works correctly
      if (player && !player.kingdom) {
        setPlayer({ ...player, kingdom: selectedKingdom } as PlayerInfo);
      }

      if (resp.village_id) {
        navigate(`/village/${resp.village_id}`);
      } else {
        navigate('/game');
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

  const liveSeasons = seasons.filter((s) => s.status === 'active');
  const upcomingSeasons = seasons.filter((s) => s.status === 'upcoming');

  return (
    <div className={styles.page}>
      <div className={styles.section}>
        <h1 className={styles.sectionTitle}>Seasons</h1>
        <p className={styles.sectionSubtitle}>
          Each season is a new world with its own map, leaderboard, and timeline.
          Join an active season to begin your conquest.
        </p>

        {error && !selectedSeason && (
          <p className={styles.error}>{error}</p>
        )}

        {loading ? (
          <div className={styles.center}>
            <LoadingSpinner size="md" />
          </div>
        ) : liveSeasons.length === 0 && upcomingSeasons.length === 0 ? (
          <p className={styles.empty}>No seasons available yet. Check back soon!</p>
        ) : (
          <>
            {/* ── Live Seasons ─────────────────────────────── */}
            <h2 className={styles.gridTitle}>
              <span className={styles.gridTitleDot + ' ' + styles.dotLive} />
              Live Seasons
            </h2>
            {liveSeasons.length === 0 ? (
              <p className={styles.gridEmpty}>No live seasons right now.</p>
            ) : (
              <div className={styles.seasonGrid}>
                {liveSeasons.map((s) => (
                  <SeasonCard
                    key={s.id}
                    season={s}
                    joined={mySeasonIds.has(s.id)}
                    isAuthenticated={isAuthenticated}
                    onJoinClick={handleJoinClick}
                  />
                ))}
              </div>
            )}

            {/* ── Upcoming Seasons ─────────────────────────── */}
            <h2 className={`${styles.gridTitle} ${styles.gridTitleSpaced}`}>
              <span className={styles.gridTitleDot + ' ' + styles.dotUpcoming} />
              Upcoming Seasons
            </h2>
            {upcomingSeasons.length === 0 ? (
              <p className={styles.gridEmpty}>No upcoming seasons announced.</p>
            ) : (
              <div className={styles.seasonGrid}>
                {upcomingSeasons.map((s) => (
                  <SeasonCard
                    key={s.id}
                    season={s}
                    joined={mySeasonIds.has(s.id)}
                    isAuthenticated={isAuthenticated}
                    onJoinClick={handleJoinClick}
                  />
                ))}
              </div>
            )}
          </>
        )}
      </div>

      {/* Kingdom selection modal */}
      {selectedSeason && (
        <div className={styles.modalOverlay} onClick={handleCloseModal}>
          <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
            <header className={styles.modalHeader}>
              <h2 className={styles.modalTitle}>Join {selectedSeason.name}</h2>
              <p className={styles.modalSubtitle}>
                This choice is permanent and shapes your buildings, troops, and destiny.
              </p>
            </header>

            {error && <p className={styles.error}>{error}</p>}

            <div className={styles.kingdomGrid}>
              {PLAYABLE_KINGDOMS.map((k) => (
                <KingdomCard
                  key={k.id}
                  kingdom={k}
                  selected={selectedKingdom === k.id}
                  onSelect={() => setSelectedKingdom(k.id)}
                />
              ))}
            </div>

            <div className={styles.modalActions}>
              <Button variant="ghost" size="sm" onClick={handleCloseModal}>
                Cancel
              </Button>
              <Button
                size="lg"
                disabled={!selectedKingdom}
                loading={joining}
                onClick={handleJoinConfirm}
              >
                {selectedKingdom
                  ? `Pledge to ${KINGDOMS.find((k) => k.id === selectedKingdom)?.name}`
                  : 'Select a Kingdom'}
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
