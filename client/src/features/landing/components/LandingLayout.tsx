import { useEffect } from 'react';
import { Outlet } from 'react-router-dom';
import { LandingHeader } from './LandingHeader';
import { useThemeStore } from '../../../stores/themeStore';
import { useAuthStore } from '../../../stores/authStore';
import { VALID_KINGDOMS } from '../../../utils/constants';
import type { Kingdom } from '../../../types/game';
import styles from './LandingLayout.module.css';

export function LandingLayout() {
  const setKingdom = useThemeStore((s) => s.setKingdom);

  // Landing pages always use the default (neutral) theme.
  // On unmount, read the *current* player kingdom from authStore
  // so we never restore a stale value captured at mount time.
  useEffect(() => {
    setKingdom(null);
    return () => {
      const raw = useAuthStore.getState().player?.kingdom;
      const k = raw && (VALID_KINGDOMS as readonly string[]).includes(raw)
        ? (raw as Kingdom)
        : null;
      setKingdom(k);
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <div className={styles.layout}>
      <LandingHeader />
      <main className={styles.content}>
        <Outlet />
      </main>
    </div>
  );
}
