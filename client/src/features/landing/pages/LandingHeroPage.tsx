import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '../../../stores/authStore';
import { fetchMySeasons } from '../../../services/season';
import { Button } from '../../../components/Button/Button';
import type { SeasonDetailResponse } from '../../../types/api';
import styles from './LandingHeroPage.module.css';

export function LandingHeroPage() {
  const { isAuthenticated } = useAuthStore();
  const navigate = useNavigate();

  const [myActiveSeasons, setMyActiveSeasons] = useState<SeasonDetailResponse[]>([]);

  // Clear stale data when the user logs out
  useEffect(() => {
    if (!isAuthenticated) {
      setMyActiveSeasons([]);
      return;
    }
    let cancelled = false;
    fetchMySeasons()
      .then((seasons) => {
        if (!cancelled) {
          setMyActiveSeasons(seasons.filter((s) => s.status === 'active'));
        }
      })
      .catch(() => {
        // Silently fail
      });
    return () => { cancelled = true; };
  }, [isAuthenticated]);

  return (
    <section className={styles.hero}>
      <h1 className={styles.heroTitle}>
        <span className={styles.heroAccent}>Weapons</span> of Order
      </h1>
      <p className={styles.heroSubtitle}>
        Build your empire, forge alliances, and craft legendary weapons to defeat the
        forces of chaos. Choose your kingdom and write your legacy.
      </p>

      {/* Currently playing cards */}
      {myActiveSeasons.length > 0 && (
        <div className={styles.playingSection}>
          <h2 className={styles.playingTitle}>Currently Playing</h2>
          <div className={styles.playingCards}>
            {myActiveSeasons.map((season) => (
              <article key={season.id} className={styles.playingCard}>
                <div className={styles.playingCardLeft}>
                  <span className={styles.playingDot} />
                  <div className={styles.playingInfo}>
                    <h3 className={styles.playingName}>{season.name}</h3>
                    <p className={styles.playingMeta}>
                      {season.kingdom && (
                        <span className={styles.playingKingdom}>
                          {season.kingdom.charAt(0).toUpperCase() + season.kingdom.slice(1)}
                        </span>
                      )}
                      <span>{season.player_count} player{season.player_count !== 1 ? 's' : ''}</span>
                    </p>
                  </div>
                </div>
                <Button size="sm" onClick={() => navigate('/game')}>
                  Enter Game
                </Button>
              </article>
            ))}
          </div>
        </div>
      )}

      <div className={styles.heroCta}>
        {isAuthenticated ? (
          <Button size="lg" onClick={() => navigate('/seasons')}>
            Browse Seasons
          </Button>
        ) : (
          <>
            <Button size="lg" onClick={() => navigate('/register')}>
              Create Account
            </Button>
            <Button variant="secondary" size="lg" onClick={() => navigate('/login')}>
              Login
            </Button>
          </>
        )}
      </div>
    </section>
  );
}
