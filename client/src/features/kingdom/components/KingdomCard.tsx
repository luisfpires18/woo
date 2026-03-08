import type { Kingdom } from '../../../types/game';
import { useAssetStore } from '../../../stores/assetStore';
import styles from './KingdomCard.module.css';

export interface KingdomInfo {
  id: Kingdom;
  name: string;
  tagline: string;
  description: string;
  traits: string[];
  colorVar: string;         // CSS color value
  glowVar: string;          // CSS glow/shadow color
  playable: boolean;        // Whether players can choose this kingdom
  lockReason?: string;      // Why the kingdom is locked (shown on card)
}

export const KINGDOMS: KingdomInfo[] = [
  {
    id: 'sylvara',
    name: 'Sylvara',
    tagline: 'Wardens of the Wild',
    description: 'The cradle kingdom — where humanity first took root on Bellum. Forests, rivers, and discipline.',
    traits: ['Archers & Scouts', 'Nature / Animal Runes', 'Healing & Survival'],
    colorVar: '#2E7D32',
    glowVar: 'rgba(46, 125, 50, 0.35)',
    playable: true,
  },
  {
    id: 'arkazia',
    name: 'Arkazia',
    tagline: 'Champions of the Arena',
    description: 'Mountain strongholds ruled by warlords. Steel, honour, and the Chapter Fortress.',
    traits: ['Heavy Cavalry & Pikemen', 'Iron / Obsidian Runes', 'Fortress & Arena Culture'],
    colorVar: '#DC143C',
    glowVar: 'rgba(220, 20, 60, 0.35)',
    playable: true,
  },
  {
    id: 'veridor',
    name: 'Veridor',
    tagline: 'Lords of the Tide',
    description: 'Born from betrayal and ambition. The wealthiest realm — trade routes, navy, and empire.',
    traits: ['Navy & Iron Verdict', 'Aquatic / Avian Runes', 'Trade & Expansion'],
    colorVar: '#2196F3',
    glowVar: 'rgba(33, 150, 243, 0.35)',
    playable: true,
  },
  {
    id: 'draxys',
    name: 'Draxys',
    tagline: 'Scions of the Sands',
    description: 'Desert frontier born from Veridor\'s colonies. Heat, survival, and relentless endurance.',
    traits: ['Desert Infantry & Scorpion Riders', 'Scale / Swarm Runes', 'Arena & Frontier Culture'],
    colorVar: '#FDD835',
    glowVar: 'rgba(253, 216, 53, 0.35)',
    playable: false,
    lockReason: 'Desert frontier — coming in a future update',
  },
  {
    id: 'zandres',
    name: 'Zandres',
    tagline: 'Keepers of the Deep',
    description: 'Hidden beneath the earth — a lost city of Thalori technic legacy and bioluminescent caverns.',
    traits: ['Underground Mining & Stonework', 'Technic / Circuit Runes', 'Secrecy & Ancient Systems'],
    colorVar: '#795548',
    glowVar: 'rgba(121, 85, 72, 0.35)',
    playable: false,
    lockReason: 'Underground realm — requires unique mechanics',
  },
  {
    id: 'lumus',
    name: 'Lumus',
    tagline: 'Children of the Sun',
    description: 'Light-soaked western island of ritual and radiance. A surviving fragment of ancient purpose.',
    traits: ['Martial Artists & Staff Users', 'Physical / Solar Runes', 'Ritual & Discipline'],
    colorVar: '#FFFFFF',
    glowVar: 'rgba(255, 255, 255, 0.35)',
    playable: false,
    lockReason: 'Island kingdom — coming in a future update',
  },
  {
    id: 'nordalh',
    name: 'Nordalh',
    tagline: 'Wolves of the North',
    description: 'Cold northern colony forged from Arkazian expansion. Fjords, smithing, and direwolf cavalry.',
    traits: ['Direwolf Cavalry & Smiths', 'Frost / Beast Runes', 'Clan Endurance & Hearth Law'],
    colorVar: '#7B1FA2',
    glowVar: 'rgba(123, 31, 162, 0.35)',
    playable: false,
    lockReason: 'Northern colony — coming in a future update',
  },
  {
    id: 'drakanith',
    name: 'Drakanith',
    tagline: 'Blood of the Volcano',
    description: 'Volcanic homeland of the Drakani — a people whose bodies carry draconic traits.',
    traits: ['Drakani Bloodline Warriors', 'Magma / Primal Runes', 'Volcanic Survival & Heat Forge'],
    colorVar: '#FF6D00',
    glowVar: 'rgba(255, 109, 0, 0.35)',
    playable: false,
    lockReason: 'Draconic bloodline — coming in a future expansion',
  },
];

interface KingdomCardProps {
  kingdom: KingdomInfo;
  selected: boolean;
  onSelect: () => void;
  locked?: boolean;
  displayOnly?: boolean;
}

export function KingdomCard({ kingdom, selected, onSelect, locked = false, displayOnly = false }: KingdomCardProps) {
  const flagAsset = useAssetStore((s) => s.getById(`flag_${kingdom.id}`));
  const flagUrl = flagAsset?.sprite_url;

  const cardClass = [
    styles.card,
    selected ? styles.selected : '',
    locked ? styles.locked : '',
    displayOnly ? styles.displayOnly : '',
  ].filter(Boolean).join(' ');

  if (displayOnly) {
    return (
      <div
        className={cardClass}
        style={{
          '--kingdom-color': kingdom.colorVar,
          '--kingdom-glow': kingdom.glowVar,
        } as React.CSSProperties}
      >
        <div className={styles.flagArea}>
          {flagUrl ? (
            <img src={flagUrl} alt={`${kingdom.name} flag`} className={styles.flagImg} />
          ) : (
            <div className={styles.flagPlaceholder}>
              <span className={styles.flagInitial}>{kingdom.name[0]}</span>
            </div>
          )}
        </div>
        <div className={styles.header}>
          <span className={styles.name}>{kingdom.name}</span>
          <span className={styles.tagline}>{kingdom.tagline}</span>
        </div>
        <p className={styles.description}>{kingdom.description}</p>
        <ul className={styles.traits}>
          {kingdom.traits.map((trait) => (
            <li key={trait} className={styles.trait}>{trait}</li>
          ))}
        </ul>
        {locked && kingdom.lockReason && (
          <p className={styles.lockReason}>{kingdom.lockReason}</p>
        )}
      </div>
    );
  }

  return (
    <button
      type="button"
      className={cardClass}
      onClick={locked ? undefined : onSelect}
      aria-pressed={locked ? undefined : selected}
      aria-disabled={locked || undefined}
      style={{
        '--kingdom-color': kingdom.colorVar,
        '--kingdom-glow': kingdom.glowVar,
      } as React.CSSProperties}
    >
      {locked && <div className={styles.lockBadge}>Coming Soon</div>}

      <div className={styles.flagArea}>
        {flagUrl ? (
          <img
            src={flagUrl}
            alt={`${kingdom.name} flag`}
            className={styles.flagImg}
          />
        ) : (
          <div className={styles.flagPlaceholder}>
            <span className={styles.flagInitial}>{kingdom.name[0]}</span>
          </div>
        )}
      </div>

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

      {locked && kingdom.lockReason && (
        <p className={styles.lockReason}>{kingdom.lockReason}</p>
      )}

      {selected && !locked && <div className={styles.checkmark}>&#10003;</div>}
    </button>
  );
}
