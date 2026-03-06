import { lazy, Suspense, useEffect } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { useAuthStore } from './stores/authStore';
import { useGameStore } from './stores/gameStore';
import { useThemeStore } from './stores/themeStore';
import { PublicLayout } from './components/Layout/PublicLayout';
import { ProtectedLayout } from './components/Layout/ProtectedLayout';
import { LoadingSpinner } from './components/LoadingSpinner/LoadingSpinner';

const LoginPage = lazy(() =>
  import('./features/auth/pages/LoginPage').then((m) => ({
    default: m.LoginPage,
  })),
);
const RegisterPage = lazy(() =>
  import('./features/auth/pages/RegisterPage').then((m) => ({
    default: m.RegisterPage,
  })),
);
const VillagePage = lazy(() =>
  import('./features/village/pages/VillagePage').then((m) => ({
    default: m.VillagePage,
  })),
);
const KingdomSelectionPage = lazy(() =>
  import('./features/kingdom/pages/KingdomSelectionPage').then((m) => ({
    default: m.KingdomSelectionPage,
  })),
);

// Admin (lazy-loaded)
const AdminLayout = lazy(() =>
  import('./components/Layout/AdminLayout').then((m) => ({
    default: m.AdminLayout,
  })),
);
const AdminPlayersPage = lazy(() =>
  import('./features/admin/pages/AdminPlayersPage').then((m) => ({
    default: m.AdminPlayersPage,
  })),
);
const AdminConfigPage = lazy(() =>
  import('./features/admin/pages/AdminConfigPage').then((m) => ({
    default: m.AdminConfigPage,
  })),
);
const AdminStatsPage = lazy(() =>
  import('./features/admin/pages/AdminStatsPage').then((m) => ({
    default: m.AdminStatsPage,
  })),
);
const AdminAnnouncementsPage = lazy(() =>
  import('./features/admin/pages/AdminAnnouncementsPage').then((m) => ({
    default: m.AdminAnnouncementsPage,
  })),
);
const AdminAssetsPage = lazy(() =>
  import('./features/admin/pages/AdminAssetsPage').then((m) => ({
    default: m.AdminAssetsPage,
  })),
);

function FullPageLoader() {
  return (
    <div style={{ display: 'flex', justifyContent: 'center', marginTop: '40vh' }}>
      <LoadingSpinner size="lg" />
    </div>
  );
}

/** Redirect / to the first village or show empty state */
function HomeRedirect() {
  const player = useAuthStore((s) => s.player);
  const villages = useGameStore((s) => s.villages);
  const villagesLoaded = useGameStore((s) => s.villagesLoaded);
  const first = villages[0];

  // No kingdom chosen yet — go to selection
  if (player && !player.kingdom) {
    return <Navigate to="/choose-kingdom" replace />;
  }

  if (first) {
    return <Navigate to={`/village/${first.id}`} replace />;
  }

  if (villagesLoaded) {
    // Villages fetched but list is empty
    return (
      <div style={{ textAlign: 'center', marginTop: '20vh', color: 'var(--text-muted)' }}>
        <h2 style={{ fontFamily: 'var(--font-heading)', color: 'var(--text-primary)' }}>
          No Village Yet
        </h2>
        <p>Your empire awaits. A village will be created for you shortly.</p>
      </div>
    );
  }

  // Still loading
  return <FullPageLoader />;
}

function App() {
  const restore = useAuthStore((s) => s.restore);
  const initTheme = useThemeStore((s) => s.init);

  useEffect(() => {
    restore();
    initTheme();
  }, [restore, initTheme]);

  return (
    <BrowserRouter>
      <Suspense fallback={<FullPageLoader />}>
        <Routes>
          {/* Public routes */}
          <Route element={<PublicLayout />}>
            <Route path="/login" element={<LoginPage />} />
            <Route path="/register" element={<RegisterPage />} />
          </Route>

          {/* Protected routes */}
          <Route element={<ProtectedLayout />}>
            <Route path="/" element={<HomeRedirect />} />
            <Route path="/village/:id" element={<VillagePage />} />
            <Route path="/choose-kingdom" element={<KingdomSelectionPage />} />
          </Route>

          {/* Admin routes */}
          <Route path="/admin" element={<AdminLayout />}>
            <Route index element={<Navigate to="/admin/players" replace />} />
            <Route path="players" element={<AdminPlayersPage />} />
            <Route path="config" element={<AdminConfigPage />} />
            <Route path="stats" element={<AdminStatsPage />} />
            <Route path="announcements" element={<AdminAnnouncementsPage />} />
            <Route path="assets" element={<AdminAssetsPage />} />
          </Route>
        </Routes>
      </Suspense>
    </BrowserRouter>
  );
}

export default App;
