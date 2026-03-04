import { lazy, Suspense, useEffect } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { useAuthStore } from './stores/authStore';
import { useGameStore } from './stores/gameStore';
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

function FullPageLoader() {
  return (
    <div style={{ display: 'flex', justifyContent: 'center', marginTop: '40vh' }}>
      <LoadingSpinner size="lg" />
    </div>
  );
}

/** Redirect / to the first village or a fallback */
function HomeRedirect() {
  const villages = useGameStore((s) => s.villages);
  const first = villages[0];

  if (first) {
    return <Navigate to={`/village/${first.id}`} replace />;
  }

  // villages not loaded yet — ProtectedLayout will fetch them
  return <FullPageLoader />;
}

function App() {
  const restore = useAuthStore((s) => s.restore);

  useEffect(() => {
    restore();
  }, [restore]);

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
          </Route>
        </Routes>
      </Suspense>
    </BrowserRouter>
  );
}

export default App;
