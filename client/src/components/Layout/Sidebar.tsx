import { NavLink } from 'react-router-dom';
import { useAuthStore } from '../../stores/authStore';
import { useGameStore } from '../../stores/gameStore';
import styles from './Sidebar.module.css';

const NAV_ITEMS = [
  { label: 'Village', path: '/village', icon: '\uD83C\uDFE0', enabled: true },
  { label: 'Map', path: '/map', icon: '\uD83D\uDDFA\uFE0F', enabled: true },
  { label: 'Forge', path: '/forge', icon: '\u2694\uFE0F', enabled: false },
  { label: 'Alliance', path: '/alliance', icon: '\uD83D\uDEE1\uFE0F', enabled: false },
];

export function Sidebar() {
  const villages = useGameStore((s) => s.villages);
  const player = useAuthStore((s) => s.player);
  const firstVillageId = villages[0]?.id;
  const isAdmin = player?.role === 'admin';

  return (
    <nav className={styles.sidebar}>
      <ul className={styles.navList}>
        {NAV_ITEMS.map((item) => {
          // For village, link to the actual first village
          const path =
            item.path === '/village' && firstVillageId
              ? `/village/${firstVillageId}`
              : item.path;

          return (
            <li key={item.label}>
              {item.enabled ? (
                <NavLink
                  to={path}
                  className={({ isActive }) =>
                    `${styles.navLink} ${isActive ? styles.active : ''}`
                  }
                >
                  <span className={styles.icon}>{item.icon}</span>
                  <span className={styles.label}>{item.label}</span>
                </NavLink>
              ) : (
                <span className={`${styles.navLink} ${styles.disabled}`}>
                  <span className={styles.icon}>{item.icon}</span>
                  <span className={styles.label}>{item.label}</span>
                </span>
              )}
            </li>
          );
        })}

        {isAdmin && (
          <li className={styles.adminDivider}>
            <NavLink
              to="/admin"
              className={({ isActive }) =>
                `${styles.navLink} ${styles.adminLink} ${isActive ? styles.active : ''}`
              }
            >
              <span className={styles.icon}>👑</span>
              <span className={styles.label}>Admin</span>
            </NavLink>
          </li>
        )}
      </ul>
    </nav>
  );
}
