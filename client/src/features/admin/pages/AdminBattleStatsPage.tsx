import { useState, useEffect, useCallback } from 'react';
import { TROOP_CONFIGS, type TroopType, type TroopConfig } from '../../../config/troops';
import { fetchBeastTemplates } from '../../../services/camp';
import type { BeastTemplateResponse } from '../../../types/api';
import { LoadingSpinner } from '../../../components/LoadingSpinner/LoadingSpinner';
import styles from './AdminBattleStatsPage.module.css';

const KINGDOMS = [
  'all',
  'arkazia',
  'draxys',
  'lumus',
  'nordalh',
  'sylvara',
  'veridor',
  'zandres',
] as const;

type Tab = 'troops' | 'beasts' | 'compare';

export function AdminBattleStatsPage() {
  const [tab, setTab] = useState<Tab>('troops');
  const [kingdom, setKingdom] = useState<string>('all');
  const [beasts, setBeasts] = useState<BeastTemplateResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const loadBeasts = useCallback(async () => {
    setLoading(true);
    try {
      setBeasts(await fetchBeastTemplates());
    } catch {
      setError('Failed to load beast templates.');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadBeasts();
  }, [loadBeasts]);

  const troops = (Object.entries(TROOP_CONFIGS) as [TroopType, TroopConfig][]).filter(
    ([, cfg]) => kingdom === 'all' || cfg.kingdom === kingdom,
  );

  return (
    <div className={styles.page}>
      <h2 className={styles.heading}>Battle Stats Comparison</h2>
      <p className={styles.subtitle}>Compare troop and beast combat stats side by side.</p>

      {error && <div className={styles.error}>{error}</div>}

      {/* Tab bar */}
      <div className={styles.tabBar}>
        {(['troops', 'beasts', 'compare'] as Tab[]).map((t) => (
          <button
            key={t}
            className={`${styles.tabBtn} ${tab === t ? styles.tabBtnActive : ''}`}
            onClick={() => setTab(t)}
          >
            {t === 'troops' ? 'Troops' : t === 'beasts' ? 'Beasts' : 'Side-by-Side'}
          </button>
        ))}
      </div>

      {/* Kingdom filter (troops & compare tabs) */}
      {(tab === 'troops' || tab === 'compare') && (
        <div className={styles.filterRow}>
          <span className={styles.filterLabel}>Kingdom:</span>
          {KINGDOMS.map((k) => (
            <button
              key={k}
              className={`${styles.filterBtn} ${kingdom === k ? styles.filterBtnActive : ''}`}
              onClick={() => setKingdom(k)}
            >
              {k === 'all' ? 'All' : k.charAt(0).toUpperCase() + k.slice(1)}
            </button>
          ))}
        </div>
      )}

      {loading ? (
        <div className={styles.center}><LoadingSpinner size="md" /></div>
      ) : (
        <>
          {tab === 'troops' && <TroopsTable troops={troops} />}
          {tab === 'beasts' && <BeastsTable beasts={beasts} />}
          {tab === 'compare' && <CompareView troops={troops} beasts={beasts} />}
        </>
      )}
    </div>
  );
}

/* ── Troops Table ──────────────────────────────────────────────────────── */

function TroopsTable({ troops }: { troops: [TroopType, TroopConfig][] }) {
  const [sortKey, setSortKey] = useState<string>('attack');
  const [sortAsc, setSortAsc] = useState(false);

  const handleSort = (key: string) => {
    if (sortKey === key) {
      setSortAsc(!sortAsc);
    } else {
      setSortKey(key);
      setSortAsc(false);
    }
  };

  const sorted = [...troops].sort(([, a], [, b]) => {
    const av = a[sortKey as keyof TroopConfig] as number;
    const bv = b[sortKey as keyof TroopConfig] as number;
    return sortAsc ? av - bv : bv - av;
  });

  const cols: { key: string; label: string }[] = [
    { key: 'attack', label: 'ATK' },
    { key: 'defInfantry', label: 'DEF Inf' },
    { key: 'defCavalry', label: 'DEF Cav' },
    { key: 'speed', label: 'Speed' },
    { key: 'carry', label: 'Carry' },
    { key: 'foodUpkeep', label: 'Upkeep' },
  ];

  return (
    <div className={styles.tableWrap}>
      <table className={styles.table}>
        <thead>
          <tr>
            <th className={styles.th}>Name</th>
            <th className={styles.th}>Kingdom</th>
            {cols.map((c) => (
              <th
                key={c.key}
                className={`${styles.th} ${styles.sortable}`}
                onClick={() => handleSort(c.key)}
              >
                {c.label} {sortKey === c.key ? (sortAsc ? '▲' : '▼') : ''}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {sorted.map(([key, cfg]) => (
            <tr key={key} className={styles.row}>
              <td className={styles.td}>{cfg.displayName}</td>
              <td className={styles.td}>{cfg.kingdom}</td>
              <td className={styles.tdNum}>{cfg.attack}</td>
              <td className={styles.tdNum}>{cfg.defInfantry}</td>
              <td className={styles.tdNum}>{cfg.defCavalry}</td>
              <td className={styles.tdNum}>{cfg.speed}</td>
              <td className={styles.tdNum}>{cfg.carry}</td>
              <td className={styles.tdNum}>{cfg.foodUpkeep}</td>
            </tr>
          ))}
        </tbody>
      </table>
      <div className={styles.countLabel}>{sorted.length} troop{sorted.length !== 1 ? 's' : ''}</div>
    </div>
  );
}

/* ── Beasts Table ──────────────────────────────────────────────────────── */

function BeastsTable({ beasts }: { beasts: BeastTemplateResponse[] }) {
  const [sortKey, setSortKey] = useState<string>('attack_power');
  const [sortAsc, setSortAsc] = useState(false);

  const handleSort = (key: string) => {
    if (sortKey === key) {
      setSortAsc(!sortAsc);
    } else {
      setSortKey(key);
      setSortAsc(false);
    }
  };

  const sorted = [...beasts].sort((a, b) => {
    const av = a[sortKey as keyof BeastTemplateResponse] as number;
    const bv = b[sortKey as keyof BeastTemplateResponse] as number;
    return sortAsc ? av - bv : bv - av;
  });

  const cols: { key: string; label: string }[] = [
    { key: 'hp', label: 'HP' },
    { key: 'attack_power', label: 'ATK' },
    { key: 'attack_interval', label: 'Interval' },
    { key: 'defense_percent', label: 'DEF %' },
    { key: 'crit_chance_percent', label: 'Crit %' },
  ];

  return (
    <div className={styles.tableWrap}>
      <table className={styles.table}>
        <thead>
          <tr>
            <th className={styles.th}>Name</th>
            {cols.map((c) => (
              <th
                key={c.key}
                className={`${styles.th} ${styles.sortable}`}
                onClick={() => handleSort(c.key)}
              >
                {c.label} {sortKey === c.key ? (sortAsc ? '▲' : '▼') : ''}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {sorted.map((b) => (
            <tr key={b.id} className={styles.row}>
              <td className={styles.td}>{b.name}</td>
              <td className={styles.tdNum}>{b.hp}</td>
              <td className={styles.tdNum}>{b.attack_power}</td>
              <td className={styles.tdNum}>{b.attack_interval}</td>
              <td className={styles.tdNum}>{b.defense_percent}%</td>
              <td className={styles.tdNum}>{b.crit_chance_percent}%</td>
            </tr>
          ))}
        </tbody>
      </table>
      <div className={styles.countLabel}>{sorted.length} beast{sorted.length !== 1 ? 's' : ''}</div>
    </div>
  );
}

/* ── Side-by-Side Compare ──────────────────────────────────────────────── */

function CompareView({
  troops,
  beasts,
}: {
  troops: [TroopType, TroopConfig][];
  beasts: BeastTemplateResponse[];
}) {
  // Compute troop aggregate stats
  const troopStats = troops.reduce(
    (acc, [, cfg]) => ({
      count: acc.count + 1,
      totalAtk: acc.totalAtk + cfg.attack,
      minAtk: Math.min(acc.minAtk, cfg.attack),
      maxAtk: Math.max(acc.maxAtk, cfg.attack),
      totalDefInf: acc.totalDefInf + cfg.defInfantry,
      totalDefCav: acc.totalDefCav + cfg.defCavalry,
      totalSpeed: acc.totalSpeed + cfg.speed,
    }),
    { count: 0, totalAtk: 0, minAtk: Infinity, maxAtk: 0, totalDefInf: 0, totalDefCav: 0, totalSpeed: 0 },
  );

  const beastStats = beasts.reduce(
    (acc, b) => ({
      count: acc.count + 1,
      totalHp: acc.totalHp + b.hp,
      minHp: Math.min(acc.minHp, b.hp),
      maxHp: Math.max(acc.maxHp, b.hp),
      totalAtk: acc.totalAtk + b.attack_power,
      minAtk: Math.min(acc.minAtk, b.attack_power),
      maxAtk: Math.max(acc.maxAtk, b.attack_power),
      totalDef: acc.totalDef + b.defense_percent,
      totalCrit: acc.totalCrit + b.crit_chance_percent,
    }),
    { count: 0, totalHp: 0, minHp: Infinity, maxHp: 0, totalAtk: 0, minAtk: Infinity, maxAtk: 0, totalDef: 0, totalCrit: 0 },
  );

  const avgTroopAtk = troopStats.count ? (troopStats.totalAtk / troopStats.count).toFixed(1) : '—';
  const avgTroopDefInf = troopStats.count ? (troopStats.totalDefInf / troopStats.count).toFixed(1) : '—';
  const avgTroopDefCav = troopStats.count ? (troopStats.totalDefCav / troopStats.count).toFixed(1) : '—';
  const avgTroopSpeed = troopStats.count ? (troopStats.totalSpeed / troopStats.count).toFixed(1) : '—';

  const avgBeastHp = beastStats.count ? (beastStats.totalHp / beastStats.count).toFixed(1) : '—';
  const avgBeastAtk = beastStats.count ? (beastStats.totalAtk / beastStats.count).toFixed(1) : '—';
  const avgBeastDef = beastStats.count ? (beastStats.totalDef / beastStats.count).toFixed(1) : '—';
  const avgBeastCrit = beastStats.count ? (beastStats.totalCrit / beastStats.count).toFixed(1) : '—';

  return (
    <div className={styles.compareGrid}>
      {/* Troops summary */}
      <div className={styles.compareCard}>
        <h3 className={styles.compareTitle}>Troops ({troopStats.count})</h3>
        <div className={styles.statRow}>
          <span className={styles.statLabel}>Avg Attack</span>
          <span className={styles.statValue}>{avgTroopAtk}</span>
        </div>
        <div className={styles.statRow}>
          <span className={styles.statLabel}>Attack Range</span>
          <span className={styles.statValue}>
            {troopStats.count ? `${troopStats.minAtk} – ${troopStats.maxAtk}` : '—'}
          </span>
        </div>
        <div className={styles.statRow}>
          <span className={styles.statLabel}>Avg DEF (Inf)</span>
          <span className={styles.statValue}>{avgTroopDefInf}</span>
        </div>
        <div className={styles.statRow}>
          <span className={styles.statLabel}>Avg DEF (Cav)</span>
          <span className={styles.statValue}>{avgTroopDefCav}</span>
        </div>
        <div className={styles.statRow}>
          <span className={styles.statLabel}>Avg Speed</span>
          <span className={styles.statValue}>{avgTroopSpeed}</span>
        </div>
      </div>

      {/* Beasts summary */}
      <div className={styles.compareCard}>
        <h3 className={styles.compareTitle}>Beasts ({beastStats.count})</h3>
        <div className={styles.statRow}>
          <span className={styles.statLabel}>Avg HP</span>
          <span className={styles.statValue}>{avgBeastHp}</span>
        </div>
        <div className={styles.statRow}>
          <span className={styles.statLabel}>HP Range</span>
          <span className={styles.statValue}>
            {beastStats.count ? `${beastStats.minHp} – ${beastStats.maxHp}` : '—'}
          </span>
        </div>
        <div className={styles.statRow}>
          <span className={styles.statLabel}>Avg Attack</span>
          <span className={styles.statValue}>{avgBeastAtk}</span>
        </div>
        <div className={styles.statRow}>
          <span className={styles.statLabel}>Attack Range</span>
          <span className={styles.statValue}>
            {beastStats.count ? `${beastStats.minAtk} – ${beastStats.maxAtk}` : '—'}
          </span>
        </div>
        <div className={styles.statRow}>
          <span className={styles.statLabel}>Avg DEF %</span>
          <span className={styles.statValue}>{avgBeastDef}%</span>
        </div>
        <div className={styles.statRow}>
          <span className={styles.statLabel}>Avg Crit %</span>
          <span className={styles.statValue}>{avgBeastCrit}%</span>
        </div>
      </div>
    </div>
  );
}
