import type { Kingdom } from '../../../types/game';
import styles from './KingdomCard.module.css';

export interface KingdomInfo {
  id: Kingdom;
  name: string;
  tagline: string;
  description: string;
  traits: string[];
  colorVar: string;         // CSS color value
  glowVar: string;          // CSS glow/shadow color
}

export const KINGDOMS: KingdomInfo[] = [
  {
    id: 'sylvara',
    name: 'Sylvara',
    tagline: 'Wardens of the Wild',
    description: 'The cradle kingdom — where humanity first took root on Bellum. Forests, rivers, and discipline.',
    traits: ['Archers & Scouts', 'Nature / Animal Runes', 'Healing & Survival'],
    colorVar: '#4CAF50',
    glowVar: 'rgba(76, 175, 80, 0.35)',
  },
  {
    id: 'arkazia',
    name: 'Arkazia',
    tagline: 'Champions of the Arena',
    description: 'Mountain strongholds ruled by warlords. Steel, honour, and the glory of the Colosseum.',
    traits: ['Heavy Cavalry & Pikemen', 'Iron / Obsidian Runes', 'Fortress & Arena Culture'],
    colorVar: '#FF9800',
    glowVar: 'rgba(255, 152, 0, 0.35)',
  },
  {
    id: 'veridor',
    name: 'Veridor',
    tagline: 'Lords of the Tide',
    description: 'Born from betrayal and ambition. The wealthiest realm — trade routes, navy, and empire.',
    traits: ['Navy & Iron Verdict', 'Aquatic / Avian Runes', 'Trade & Expansion'],
    colorVar: '#2196F3',
    glowVar: 'rgba(33, 150, 243, 0.35)',
  },
  {
    id: 'draxys',
    name: 'Draxys',
    tagline: 'Scions of the Sands',
    description: 'Desert frontier born from Veridor\'s colonies. Heat, survival, and relentless endurance.',
    traits: ['Desert Infantry & Scorpion Riders', 'Scale / Swarm Runes', 'Arena & Frontier Culture'],
    colorVar: '#E6A817',
    glowVar: 'rgba(230, 168, 23, 0.35)',
  },
  {
    id: 'zandres',
    name: 'Zandres',
    tagline: 'Keepers of the Deep',
    description: 'Hidden beneath the earth — a lost city of Thalori technic legacy and bioluminescent caverns.',
    traits: ['Underground Mining & Stonework', 'Technic / Circuit Runes', 'Secrecy & Ancient Systems'],
    colorVar: '#9C27B0',
    glowVar: 'rgba(156, 39, 176, 0.35)',
  },
  {
    id: 'lumus',
    name: 'Lumus',
    tagline: 'Children of the Sun',
    description: 'Light-soaked western island of ritual and radiance. A surviving fragment of ancient purpose.',
    traits: ['Martial Artists & Staff Users', 'Physical / Solar Runes', 'Ritual & Discipline'],
    colorVar: '#FDD835',
    glowVar: 'rgba(253, 216, 53, 0.35)',
  },
  {
    id: 'nordalh',
    name: 'Nordalh',
    tagline: 'Wolves of the North',
    description: 'Cold northern colony forged from Arkazian expansion. Fjords, smithing, and direwolf cavalry.',
    traits: ['Direwolf Cavalry & Smiths', 'Frost / Beast Runes', 'Clan Endurance & Hearth Law'],
    colorVar: '#78909C',
    glowVar: 'rgba(120, 144, 156, 0.35)',
  },
  {
    id: 'drakanith',
    name: 'Drakanith',
    tagline: 'Blood of the Volcano',
    description: 'Volcanic homeland of the Drakani — a people whose bodies carry draconic traits.',
    traits: ['Drakani Bloodline Warriors', 'Magma / Primal Runes', 'Volcanic Survival & Heat Forge'],
    colorVar: '#F44336',
    glowVar: 'rgba(244, 67, 54, 0.35)',
  },
];

interface KingdomCardProps {
  kingdom: KingdomInfo;
  selected: boolean;
  onSelect: () => void;
}

export function KingdomCard({ kingdom, selected, onSelect }: KingdomCardProps) {
  return (
    <button
      type="button"
      className={`${styles.card} ${selected ? styles.selected : ''}`}
      onClick={onSelect}
      aria-pressed={selected}
      style={{
        '--kingdom-color': kingdom.colorVar,
        '--kingdom-glow': kingdom.glowVar,
      } as React.CSSProperties}
    >
      <div className={styles.header}>
        <span className={styles.name}>{kingdom.name}</span>
        <span className={styles.tagline}>{kingdom.tagline}</span>
      </div>

      <p className={styles.description}>{kingdom.description}</p>

      <ul className={styles.traits}>
        {kingdom.traits.map((trait) => (
          <li key={trait} className={styles.trait}>
            {trait}
          </li>
        ))}
      </ul>

      {selected && <div className={styles.checkmark}>&#10003;</div>}
    </button>
  );
}
