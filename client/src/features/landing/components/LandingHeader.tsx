import { Link, NavLink, useNavigate } from 'react-router-dom';
import { useAuthStore } from '../../../stores/authStore';
import styles from './LandingHeader.module.css';

export function LandingHeader() {
  const { player, isAuthenticated, logout } = useAuthStore();
  const navigate = useNavigate();
  const isAdmin = player?.role === 'admin';

  const handleLogout = async () => {
    await logout();
    navigate('/');
  };

  return (
    <header className={styles.header}>
      <div className={styles.left}>
        <Link to="/" className={styles.logo}>
          Weapons of Order
        </Link>
        <nav className={styles.nav}>
          <NavLink to="/seasons" className={({ isActive }) => `${styles.navLink} ${isActive ? styles.navLinkActive : ''}`}>Seasons</NavLink>
          <NavLink to="/kingdoms" className={({ isActive }) => `${styles.navLink} ${isActive ? styles.navLinkActive : ''}`}>Kingdoms</NavLink>
          <NavLink to="/leaderboards" className={({ isActive }) => `${styles.navLink} ${isActive ? styles.navLinkActive : ''}`}>Leaderboards</NavLink>
        </nav>
      </div>

      <div className={styles.right}>
        {isAuthenticated && player ? (
          <div className={styles.user}>
            {isAdmin && (
              <Link to="/admin" className={styles.adminLink}>👑 Admin</Link>
            )}
            <Link to="/profile" className={styles.usernameLink}>
              {player.username}
            </Link>
            <button className={styles.logoutBtn} onClick={handleLogout}>
              Logout
            </button>
          </div>
        ) : (
          <>
            <Link to="/login" className={styles.authLink}>Login</Link>
            <Link to="/register" className={styles.authLink}>Register</Link>
          </>
        )}
      </div>
    </header>
  );
}
