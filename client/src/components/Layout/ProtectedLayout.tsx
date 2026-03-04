import { useEffect } from 'react';
import { Navigate, Outlet } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { useAuthStore } from '../../stores/authStore';
import { useGameStore } from '../../stores/gameStore';
import { fetchVillages } from '../../services/village';
import { Header } from './Header';
import { Sidebar } from './Sidebar';
import { LoadingSpinner } from '../LoadingSpinner/LoadingSpinner';
import styles from './ProtectedLayout.module.css';

export function ProtectedLayout() {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);
  const setVillages = useGameStore((s) => s.setVillages);

  // Fetch the player's village list on mount
  const { data: villages, isLoading } = useQuery({
    queryKey: ['villages'],
    queryFn: fetchVillages,
    enabled: isAuthenticated,
  });

  useEffect(() => {
    if (villages) {
      setVillages(villages);
    }
  }, [villages, setVillages]);

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
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
