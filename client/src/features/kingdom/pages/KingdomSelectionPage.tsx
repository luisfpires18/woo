import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '../../../stores/authStore';
import { useThemeStore } from '../../../stores/themeStore';
import { useAssetStore } from '../../../stores/assetStore';
import { KingdomCard, KINGDOMS } from '../components/KingdomCard';
import { Button } from '../../../components/Button/Button';
import { chooseKingdom } from '../../../services/player';
import { ApiRequestError } from '../../../services/api';
import type { Kingdom } from '../../../types/game';
import styles from './KingdomSelectionPage.module.css';

export function KingdomSelectionPage() {
  const player = useAuthStore((s) => s.player);
  const setPlayer = useAuthStore((s) => s.setPlayer);
  const setKingdom = useThemeStore((s) => s.setKingdom);
  const loadAssets = useAssetStore((s) => s.load);
  const navigate = useNavigate();

  const [selected, setSelected] = useState<Kingdom | null>(null);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const playableKingdoms = KINGDOMS.filter((k) => k.playable);
  const lockedKingdoms = KINGDOMS.filter((k) => !k.playable);

  // Load game assets so kingdom flags are available
  useEffect(() => {
    loadAssets();
  }, [loadAssets]);

  // If somehow the player already has a kingdom, redirect
  if (player?.kingdom) {
    navigate('/', { replace: true });
    return null;
  }

  const handleConfirm = async () => {
    if (!selected) return;
    setError('');
    setLoading(true);

    try {
      const resp = await chooseKingdom(selected);
      // Update player in authStore with the new kingdom & apply theme immediately
      setPlayer(resp.player);
      setKingdom(selected);
      navigate('/', { replace: true });
    } catch (err) {
      if (err instanceof ApiRequestError) {
        setError(err.message);
      } else {
        setError('An unexpected error occurred');
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.page}>
      <div className={styles.container}>
        <header className={styles.header}>
          <h1 className={styles.title}>Choose Your Kingdom</h1>
          <p className={styles.subtitle}>
            This choice is permanent and shapes your buildings, troops, and destiny on Bellum.
          </p>
        </header>

        {error && <p className={styles.error}>{error}</p>}

        <div className={styles.grid}>
          {playableKingdoms.map((k) => (
            <KingdomCard
              key={k.id}
              kingdom={k}
              selected={selected === k.id}
              onSelect={() => setSelected(k.id)}
            />
          ))}
        </div>

        <div className={styles.footer}>
          <Button
            size="lg"
            loading={loading}
            disabled={!selected}
            onClick={handleConfirm}
          >
            {selected
              ? `Pledge to ${KINGDOMS.find((k) => k.id === selected)?.name}`
              : 'Select a Kingdom'}
          </Button>
        </div>

        {lockedKingdoms.length > 0 && (
          <>
            <div className={styles.comingSoonDivider}>
              <span className={styles.comingSoonLabel}>Coming Soon</span>
            </div>

            <div className={styles.gridLocked}>
              {lockedKingdoms.map((k) => (
                <KingdomCard
                  key={k.id}
                  kingdom={k}
                  selected={false}
                  onSelect={() => {}}
                  locked
                />
              ))}
            </div>
          </>
        )}
      </div>
    </div>
  );
}
