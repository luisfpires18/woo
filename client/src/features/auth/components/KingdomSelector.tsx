import type { Kingdom } from '../../../types/game';
import styles from './KingdomSelector.module.css';

interface KingdomSelectorProps {
  selected: Kingdom | null;
  onSelect: (kingdom: Kingdom) => void;
}

const KINGDOMS: {
  id: Kingdom;
  name: string;
  tagline: string;
  description: string;
  colorClass: string;
}[] = [
  {
    id: 'veridor',
    name: 'Veridor',
    tagline: 'Lords of the Tide',
    description: 'Naval power. Sea trade routes. Dock-based warfare.',
    colorClass: 'veridor',
  },
  {
    id: 'sylvara',
    name: 'Sylvara',
    tagline: 'Wardens of the Wild',
    description: 'Forest mastery. Nature magic. Grove Sanctum rituals.',
    colorClass: 'sylvara',
  },
  {
    id: 'arkazia',
    name: 'Arkazia',
    tagline: 'Champions of the Arena',
    description: 'Mountain strength. Gladiator combat. Colosseum glory.',
    colorClass: 'arkazia',
  },
];

export function KingdomSelector({ selected, onSelect }: KingdomSelectorProps) {
  return (
    <div className={styles.grid}>
      {KINGDOMS.map((k) => (
        <button
          key={k.id}
          type="button"
          className={`${styles.card} ${styles[k.colorClass]} ${
            selected === k.id ? styles.selected : ''
          }`}
          onClick={() => onSelect(k.id)}
          aria-pressed={selected === k.id}
        >
          <span className={styles.name}>{k.name}</span>
          <span className={styles.tagline}>{k.tagline}</span>
          <span className={styles.description}>{k.description}</span>
        </button>
      ))}
    </div>
  );
}
