import { Link, useNavigate } from 'react-router-dom';
import { useAuthStore } from '../../stores/authStore';
import styles from './Header.module.css';

export function Header() {
  const { player, logout, isAuthenticated } = useAuthStore();
  const navigate = useNavigate();

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
      </div>

      <div className={styles.center} />

      <div className={styles.right}>
        {isAuthenticated && player && (
          <div className={styles.user}>
            <span className={styles.username}>{player.username}</span>
            <span className={styles.kingdom}>{player.kingdom}</span>
            <button className={styles.logoutBtn} onClick={handleLogout}>
              Logout
            </button>
          </div>
        )}
      </div>
    </header>
  );
}
