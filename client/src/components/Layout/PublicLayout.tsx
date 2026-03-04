import { Outlet } from 'react-router-dom';
import styles from './PublicLayout.module.css';

export function PublicLayout() {
  return (
    <div className={styles.layout}>
      <div className={styles.container}>
        <div className={styles.brand}>
          <h1 className={styles.title}>Weapons of Order</h1>
          <p className={styles.tagline}>The forge awaits.</p>
        </div>
        <Outlet />
      </div>
    </div>
  );
}
