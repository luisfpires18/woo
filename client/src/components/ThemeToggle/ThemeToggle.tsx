import { useThemeStore } from '../../stores/themeStore';
import styles from './ThemeToggle.module.css';

export function ThemeToggle() {
  const theme = useThemeStore((s) => s.theme);
  const toggle = useThemeStore((s) => s.toggle);

  const isDark = theme === 'dark';

  return (
    <button
      className={styles.toggle}
      onClick={toggle}
      aria-label={`Switch to ${isDark ? 'light' : 'dark'} mode`}
      title={`Switch to ${isDark ? 'light' : 'dark'} mode`}
    >
      <span className={`${styles.icon} ${isDark ? styles.sun : styles.moon}`}>
        {isDark ? '☀' : '☾'}
      </span>
    </button>
  );
}
