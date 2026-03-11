import { Navigate, NavLink, Outlet } from 'react-router-dom';
import { useAuthStore } from '../../stores/authStore';
import styles from './AdminLayout.module.css';

const ADMIN_TABS = [
  { label: 'Players', path: '/admin/players' },
  { label: 'Seasons', path: '/admin/seasons' },
  { label: 'Stats', path: '/admin/stats' },
  { label: 'Announcements', path: '/admin/announcements' },
  { label: 'Kingdoms', path: '/admin/kingdoms' },
  { label: 'Buildings', path: '/admin/buildings' },
  { label: 'Units', path: '/admin/units' },
  { label: 'Resources', path: '/admin/resources' },
  { label: 'Map Assets', path: '/admin/map-assets' },
  { label: 'Map Editor', path: '/admin/map-editor' },
];

export function AdminLayout() {
  const player = useAuthStore((s) => s.player);
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  if (player?.role !== 'admin') {
    return <Navigate to="/" replace />;
  }

  return (
    <div className={styles.layout}>
      <header className={styles.header}>
        <h1 className={styles.title}>Admin Panel</h1>
        <NavLink to="/" className={styles.backLink}>
          ← Back to game
        </NavLink>
      </header>

      <nav className={styles.tabs}>
        {ADMIN_TABS.map((tab) => (
          <NavLink
            key={tab.path}
            to={tab.path}
            className={({ isActive }) =>
              `${styles.tab} ${isActive ? styles.activeTab : ''}`
            }
          >
            {tab.label}
          </NavLink>
        ))}
      </nav>

      <main className={styles.content}>
        <Outlet />
      </main>
    </div>
  );
}
