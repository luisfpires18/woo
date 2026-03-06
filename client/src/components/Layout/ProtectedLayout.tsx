import { useEffect } from 'react';
import { Navigate, Outlet } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { useAuthStore } from '../../stores/authStore';
import { useGameStore } from '../../stores/gameStore';
import { useAssetStore } from '../../stores/assetStore';
import { fetchVillages } from '../../services/village';
import { Header } from './Header';
import { Sidebar } from './Sidebar';
import { LoadingSpinner } from '../LoadingSpinner/LoadingSpinner';
import styles from './ProtectedLayout.module.css';

export function ProtectedLayout() {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);
  const hydrated = useAuthStore((s) => s.hydrated);
  const player = useAuthStore((s) => s.player);
  const setVillages = useGameStore((s) => s.setVillages);
  const loadAssets = useAssetStore((s) => s.load);

  const hasKingdom = !!player?.kingdom;

  // Fetch the player's village list on mount (only if kingdom is chosen)
  const { data: villages, isLoading } = useQuery({
    queryKey: ['villages'],
    queryFn: fetchVillages,
    enabled: isAuthenticated && hasKingdom,
  });

  useEffect(() => {
    if (villages) {
      setVillages(villages);
    }
  }, [villages, setVillages]);

  // Pre-load game asset sprites on login
  useEffect(() => {
    if (isAuthenticated) {
      loadAssets();
    }
  }, [isAuthenticated, loadAssets]);

  if (!hydrated) {
    return <LoadingSpinner size="lg" />;
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  // Player hasn't chosen a kingdom yet — show kingdom selection
  // (render directly, skip the full layout with Header/Sidebar)
  if (player && !player.kingdom) {
    return <Outlet />;
  }

  if (isLoading) {
    return <LoadingSpinner size="lg" />;
  }

  return (
    <div className={styles.layout}>
      <Header />
      <div className={styles.body}>
        <Sidebar />
        <main className={styles.main}>
          <Outlet />
        </main>
      </div>
    </div>
  );
}
