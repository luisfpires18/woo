import { lazy, Suspense, useEffect } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { useAuthStore } from './stores/authStore';
import { useGameStore } from './stores/gameStore';
import { useThemeStore } from './stores/themeStore';
import { PublicLayout } from './components/Layout/PublicLayout';
import { ProtectedLayout } from './components/Layout/ProtectedLayout';
import { LoadingSpinner } from './components/LoadingSpinner/LoadingSpinner';
import type { Kingdom } from './types/game';
import { VALID_KINGDOMS } from './utils/constants';
import appStyles from './App.module.css';

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
const WorldMapPage = lazy(() =>
  import('./features/map/pages/WorldMapPage').then((m) => ({
    default: m.WorldMapPage,
  })),
);
const ProfilePage = lazy(() =>
  import('./features/profile/pages/ProfilePage').then((m) => ({
    default: m.ProfilePage,
  })),
);

// Landing sub-pages
const LandingLayout = lazy(() =>
  import('./features/landing/components/LandingLayout').then((m) => ({
    default: m.LandingLayout,
  })),
);
const LandingHeroPage = lazy(() =>
  import('./features/landing/pages/LandingHeroPage').then((m) => ({
    default: m.LandingHeroPage,
  })),
);
const SeasonsPublicPage = lazy(() =>
  import('./features/landing/pages/SeasonsPublicPage').then((m) => ({
    default: m.SeasonsPublicPage,
  })),
);
const KingdomsShowcasePage = lazy(() =>
  import('./features/landing/pages/KingdomsShowcasePage').then((m) => ({
    default: m.KingdomsShowcasePage,
  })),
);
const LeaderboardsPage = lazy(() =>
  import('./features/landing/pages/LeaderboardsPage').then((m) => ({
    default: m.LeaderboardsPage,
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
const AdminBuildingsPage = lazy(() =>
  import('./features/admin/pages/AdminBuildingsPage').then((m) => ({
    default: m.AdminBuildingsPage,
  })),
);
const AdminMapEditorPage = lazy(() =>
  import('./features/admin/pages/AdminMapEditorPage').then((m) => ({
    default: m.AdminMapEditorPage,
  })),
);
const AdminSeasonsPage = lazy(() =>
  import('./features/admin/pages/AdminSeasonsPage').then((m) => ({
    default: m.AdminSeasonsPage,
  })),
);

function FullPageLoader() {
  return (
    <div className={appStyles.fullPageLoader}>
      <LoadingSpinner size="lg" />
    </div>
  );
}

/** Redirect /game to the first village or show empty state */
function GameRedirect() {
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
      <div className={appStyles.emptyState}>
        <h2 className={appStyles.emptyStateTitle}>
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
  const setKingdom = useThemeStore((s) => s.setKingdom);
  const playerKingdom = useAuthStore((s) => s.player?.kingdom);

  useEffect(() => {
    restore();
    initTheme();
  }, [restore, initTheme]);

  // Sync kingdom theme whenever player.kingdom changes (login, restore, kingdom selection)
  useEffect(() => {
    if (playerKingdom && (VALID_KINGDOMS as readonly string[]).includes(playerKingdom)) {
      setKingdom(playerKingdom as Kingdom);
    } else {
      setKingdom(null);
    }
  }, [playerKingdom, setKingdom]);

  return (
    <BrowserRouter>
      <Suspense fallback={<FullPageLoader />}>
        <Routes>
          {/* Landing layout — public, works for both guests and logged-in users */}
          <Route element={<LandingLayout />}>
            <Route path="/" element={<LandingHeroPage />} />
            <Route path="/seasons" element={<SeasonsPublicPage />} />
            <Route path="/kingdoms" element={<KingdomsShowcasePage />} />
            <Route path="/leaderboards" element={<LeaderboardsPage />} />
            <Route path="/profile" element={<ProfilePage />} />
          </Route>

          {/* Public routes (auth forms) */}
          <Route element={<PublicLayout />}>
            <Route path="/login" element={<LoginPage />} />
            <Route path="/register" element={<RegisterPage />} />
          </Route>

          {/* Protected routes (in-game) */}
          <Route element={<ProtectedLayout />}>
            <Route path="/game" element={<GameRedirect />} />
            <Route path="/village/:id" element={<VillagePage />} />
            <Route path="/map" element={<WorldMapPage />} />
            <Route path="/choose-kingdom" element={<KingdomSelectionPage />} />
          </Route>

          {/* Admin routes */}
          <Route path="/admin" element={<AdminLayout />}>
            <Route index element={<Navigate to="/admin/players" replace />} />
            <Route path="players" element={<AdminPlayersPage />} />
            <Route path="stats" element={<AdminStatsPage />} />
            <Route path="announcements" element={<AdminAnnouncementsPage />} />
            <Route path="assets" element={<AdminAssetsPage />} />
            <Route path="buildings" element={<AdminBuildingsPage />} />
            <Route path="seasons" element={<AdminSeasonsPage />} />
            <Route path="map-editor" element={<AdminMapEditorPage />} />
          </Route>
        </Routes>
      </Suspense>
    </BrowserRouter>
  );
}

export default App;
