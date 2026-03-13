// Admin Camps page — manage beast templates, camp templates, spawn rules, reward tables, and battle tuning

import { useState, useEffect, useCallback } from 'react';
import { LoadingSpinner } from '../../../components/LoadingSpinner/LoadingSpinner';
import {
  fetchBeastTemplates,
  createBeastTemplate,
  updateBeastTemplate,
  deleteBeastTemplate,
  fetchCampTemplates,
  createCampTemplate,
  updateCampTemplate,
  deleteCampTemplate,
  fetchSpawnRules,
  createSpawnRule,
  updateSpawnRule,
  deleteSpawnRule,
  fetchRewardTables,
  createRewardTable,
  updateRewardTable,
  deleteRewardTable,
  fetchBattleTuning,
  updateBattleTuning,
} from '../../../services/camp';
import type {
  BeastTemplateResponse,
  CreateBeastTemplateRequest,
  CampTemplateResponse,
  CreateCampTemplateRequest,
  SpawnRuleResponse,
  CreateSpawnRuleRequest,
  RewardTableResponse,
  CreateRewardTableRequest,
  BattleTuningResponse,
} from '../../../types/api';
import styles from './AdminCampsPage.module.css';

type Tab = 'beasts' | 'camps' | 'spawns' | 'rewards' | 'tuning';

export function AdminCampsPage() {
  const [tab, setTab] = useState<Tab>('beasts');
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const clearMessages = () => {
    setError(null);
    setSuccess(null);
  };

  return (
    <div className={styles.page}>
      <h1 className={styles.heading}>Camp Management</h1>
      <p className={styles.subtitle}>Configure beast templates, camp templates, spawn rules, reward tables, and battle tuning.</p>

      <div className={styles.tabBar}>
        {(['beasts', 'camps', 'spawns', 'rewards', 'tuning'] as Tab[]).map((t) => (
          <button
            key={t}
            className={tab === t ? styles.tabBtnActive : styles.tabBtn}
            onClick={() => { setTab(t); clearMessages(); }}
          >
            {t === 'beasts' ? 'Beast Templates' :
             t === 'camps' ? 'Camp Templates' :
             t === 'spawns' ? 'Spawn Rules' :
             t === 'rewards' ? 'Reward Tables' : 'Battle Tuning'}
          </button>
        ))}
      </div>

      {error && <div className={styles.error}>{error}</div>}
      {success && <div className={styles.success}>{success}</div>}

      {tab === 'beasts' && <BeastTemplatesSection onError={setError} onSuccess={setSuccess} />}
      {tab === 'camps' && <CampTemplatesSection onError={setError} onSuccess={setSuccess} />}
      {tab === 'spawns' && <SpawnRulesSection onError={setError} onSuccess={setSuccess} />}
      {tab === 'rewards' && <RewardTablesSection onError={setError} onSuccess={setSuccess} />}
      {tab === 'tuning' && <BattleTuningSection onError={setError} onSuccess={setSuccess} />}
    </div>
  );
}

// ── Shared types for section props ──────────────────────────────────────────

interface SectionProps {
  onError: (msg: string) => void;
  onSuccess: (msg: string) => void;
}

// ── Beast Templates Section ─────────────────────────────────────────────────

function BeastTemplatesSection({ onError, onSuccess }: SectionProps) {
  const [items, setItems] = useState<BeastTemplateResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editId, setEditId] = useState<number | null>(null);

  const load = useCallback(async () => {
    setLoading(true);
    try {
      setItems(await fetchBeastTemplates());
    } catch {
      onError('Failed to load beast templates.');
    } finally {
      setLoading(false);
    }
  }, [onError]);

  useEffect(() => { load(); }, [load]);

  const handleDelete = async (id: number) => {
    if (!confirm('Delete this beast template?')) return;
    try {
      await deleteBeastTemplate(id);
      onSuccess('Beast template deleted.');
      load();
    } catch {
      onError('Failed to delete beast template.');
    }
  };

  const handleSave = async (data: CreateBeastTemplateRequest) => {
    try {
      if (editId) {
        await updateBeastTemplate(editId, data);
        onSuccess('Beast template updated.');
      } else {
        await createBeastTemplate(data);
        onSuccess('Beast template created.');
      }
      setShowForm(false);
      setEditId(null);
      load();
    } catch {
      onError('Failed to save beast template.');
    }
  };

  if (loading) return <div className={styles.center}><LoadingSpinner size="md" /></div>;

  const editItem = editId ? items.find((i) => i.id === editId) : undefined;

  return (
    <>
      <div className={styles.sectionHeader}>
        <span>{items.length} beast template{items.length !== 1 ? 's' : ''}</span>
        <button className={styles.addBtn} onClick={() => { setEditId(null); setShowForm(true); }}>
          + New Beast
        </button>
      </div>

      {items.length === 0 ? (
        <div className={styles.empty}>No beast templates yet. Create one to get started.</div>
      ) : (
        <div className={styles.cardList}>
          {items.map((b) => (
            <div className={styles.card} key={b.id}>
              <div className={styles.cardHeader}>
                <h3 className={styles.cardTitle}>{b.name}</h3>
                <div className={styles.cardActions}>
                  <button className={styles.editBtn} onClick={() => { setEditId(b.id); setShowForm(true); }}>Edit</button>
                  <button className={styles.deleteBtn} onClick={() => handleDelete(b.id)}>Delete</button>
                </div>
              </div>
              <div className={styles.cardRow}>
                <span className={styles.cardLabel}>Sprite</span>
                <span className={styles.cardValue}>{b.sprite_key || '—'}</span>
              </div>
              <div className={styles.cardRow}>
                <span className={styles.cardLabel}>HP / ATK / Interval</span>
                <span className={styles.cardValue}>{b.hp} / {b.attack_power} / {b.attack_interval}</span>
              </div>
              <div className={styles.cardRow}>
                <span className={styles.cardLabel}>DEF% / CRIT%</span>
                <span className={styles.cardValue}>{b.defense_percent}% / {b.crit_chance_percent}%</span>
              </div>
              {b.description && (
                <div className={styles.cardRow}>
                  <span className={styles.cardLabel}>Desc</span>
                  <span className={styles.cardValue}>{b.description}</span>
                </div>
              )}
            </div>
          ))}
        </div>
      )}

      {showForm && (
        <BeastTemplateForm
          initial={editItem}
          onSave={handleSave}
          onCancel={() => { setShowForm(false); setEditId(null); }}
        />
      )}
    </>
  );
}

function BeastTemplateForm({
  initial,
  onSave,
  onCancel,
}: {
  initial?: BeastTemplateResponse;
  onSave: (data: CreateBeastTemplateRequest) => void;
  onCancel: () => void;
}) {
  const [name, setName] = useState(initial?.name ?? '');
  const [spriteKey, setSpriteKey] = useState(initial?.sprite_key ?? '');
  const [hp, setHp] = useState(initial?.hp ?? 100);
  const [attackPower, setAttackPower] = useState(initial?.attack_power ?? 10);
  const [attackInterval, setAttackInterval] = useState(initial?.attack_interval ?? 5);
  const [defensePct, setDefensePct] = useState(initial?.defense_percent ?? 0);
  const [critPct, setCritPct] = useState(initial?.crit_chance_percent ?? 0);
  const [desc, setDesc] = useState(initial?.description ?? '');

  return (
    <div className={styles.formOverlay} onClick={onCancel}>
      <div className={styles.formModal} onClick={(e) => e.stopPropagation()}>
        <h3 className={styles.formTitle}>{initial ? 'Edit Beast Template' : 'New Beast Template'}</h3>
        <div className={styles.formField}>
          <label className={styles.formLabel}>Name</label>
          <input className={styles.formInput} value={name} onChange={(e) => setName(e.target.value)} />
        </div>
        <div className={styles.formField}>
          <label className={styles.formLabel}>Sprite Key</label>
          <input className={styles.formInput} value={spriteKey} onChange={(e) => setSpriteKey(e.target.value)} placeholder="e.g. forest_spider" />
        </div>
        <div className={styles.formField}>
          <label className={styles.formLabel}>HP</label>
          <input className={styles.formInput} type="number" min={1} value={hp} onChange={(e) => setHp(+e.target.value)} />
        </div>
        <div className={styles.formField}>
          <label className={styles.formLabel}>Attack Power</label>
          <input className={styles.formInput} type="number" min={1} value={attackPower} onChange={(e) => setAttackPower(+e.target.value)} />
        </div>
        <div className={styles.formField}>
          <label className={styles.formLabel}>Attack Interval (ticks)</label>
          <input className={styles.formInput} type="number" min={1} value={attackInterval} onChange={(e) => setAttackInterval(+e.target.value)} />
        </div>
        <div className={styles.formField}>
          <label className={styles.formLabel}>Defense %</label>
          <input className={styles.formInput} type="number" min={0} max={100} step={0.1} value={defensePct} onChange={(e) => setDefensePct(+e.target.value)} />
        </div>
        <div className={styles.formField}>
          <label className={styles.formLabel}>Crit Chance %</label>
          <input className={styles.formInput} type="number" min={0} max={100} step={0.1} value={critPct} onChange={(e) => setCritPct(+e.target.value)} />
        </div>
        <div className={styles.formField}>
          <label className={styles.formLabel}>Description</label>
          <input className={styles.formInput} value={desc} onChange={(e) => setDesc(e.target.value)} />
        </div>
        <div className={styles.formActions}>
          <button className={styles.cancelFormBtn} onClick={onCancel}>Cancel</button>
          <button className={styles.saveBtn} disabled={!name} onClick={() => onSave({
            name,
            sprite_key: spriteKey,
            hp,
            attack_power: attackPower,
            attack_interval: attackInterval,
            defense_percent: defensePct,
            crit_chance_percent: critPct,
            description: desc,
          })}>Save</button>
        </div>
      </div>
    </div>
  );
}

// ── Camp Templates Section ──────────────────────────────────────────────────

function CampTemplatesSection({ onError, onSuccess }: SectionProps) {
  const [items, setItems] = useState<CampTemplateResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editId, setEditId] = useState<number | null>(null);

  const load = useCallback(async () => {
    setLoading(true);
    try {
      setItems(await fetchCampTemplates());
    } catch {
      onError('Failed to load camp templates.');
    } finally {
      setLoading(false);
    }
  }, [onError]);

  useEffect(() => { load(); }, [load]);

  const handleDelete = async (id: number) => {
    if (!confirm('Delete this camp template?')) return;
    try {
      await deleteCampTemplate(id);
      onSuccess('Camp template deleted.');
      load();
    } catch {
      onError('Failed to delete camp template.');
    }
  };

  const handleSave = async (data: CreateCampTemplateRequest) => {
    try {
      if (editId) {
        await updateCampTemplate(editId, { name: data.name, tier: data.tier, min_beasts: data.min_beasts, max_beasts: data.max_beasts, reward_table_id: data.reward_table_id, description: data.description });
        onSuccess('Camp template updated.');
      } else {
        await createCampTemplate(data);
        onSuccess('Camp template created.');
      }
      setShowForm(false);
      setEditId(null);
      load();
    } catch {
      onError('Failed to save camp template.');
    }
  };

  if (loading) return <div className={styles.center}><LoadingSpinner size="md" /></div>;

  const editItem = editId ? items.find((i) => i.id === editId) : undefined;

  return (
    <>
      <div className={styles.sectionHeader}>
        <span>{items.length} camp template{items.length !== 1 ? 's' : ''}</span>
        <button className={styles.addBtn} onClick={() => { setEditId(null); setShowForm(true); }}>
          + New Camp
        </button>
      </div>

      {items.length === 0 ? (
        <div className={styles.empty}>No camp templates yet.</div>
      ) : (
        <div className={styles.cardList}>
          {items.map((c) => (
            <div className={styles.card} key={c.id}>
              <div className={styles.cardHeader}>
                <h3 className={styles.cardTitle}>{c.name} (Tier {c.tier})</h3>
                <div className={styles.cardActions}>
                  <button className={styles.editBtn} onClick={() => { setEditId(c.id); setShowForm(true); }}>Edit</button>
                  <button className={styles.deleteBtn} onClick={() => handleDelete(c.id)}>Delete</button>
                </div>
              </div>
              <div className={styles.cardRow}>
                <span className={styles.cardLabel}>Beasts</span>
                <span className={styles.cardValue}>{c.min_beasts}–{c.max_beasts}</span>
              </div>
              <div className={styles.cardRow}>
                <span className={styles.cardLabel}>Reward Table</span>
                <span className={styles.cardValue}>#{c.reward_table_id}</span>
              </div>
              {c.beast_slots.length > 0 && (
                <div className={styles.cardRow}>
                  <span className={styles.cardLabel}>Slots</span>
                  <span className={styles.cardValue}>{c.beast_slots.map((s) => s.beast_name).join(', ')}</span>
                </div>
              )}
            </div>
          ))}
        </div>
      )}

      {showForm && (
        <CampTemplateForm
          initial={editItem}
          onSave={handleSave}
          onCancel={() => { setShowForm(false); setEditId(null); }}
        />
      )}
    </>
  );
}

function CampTemplateForm({
  initial,
  onSave,
  onCancel,
}: {
  initial?: CampTemplateResponse;
  onSave: (data: CreateCampTemplateRequest) => void;
  onCancel: () => void;
}) {
  const [name, setName] = useState(initial?.name ?? '');
  const [tier, setTier] = useState(initial?.tier ?? 1);
  const [minBeasts, setMinBeasts] = useState(initial?.min_beasts ?? 1);
  const [maxBeasts, setMaxBeasts] = useState(initial?.max_beasts ?? 3);
  const [rewardTableId, setRewardTableId] = useState(initial?.reward_table_id ?? 1);
  const [desc, setDesc] = useState(initial?.description ?? '');

  return (
    <div className={styles.formOverlay} onClick={onCancel}>
      <div className={styles.formModal} onClick={(e) => e.stopPropagation()}>
        <h3 className={styles.formTitle}>{initial ? 'Edit Camp Template' : 'New Camp Template'}</h3>
        <div className={styles.formField}>
          <label className={styles.formLabel}>Name</label>
          <input className={styles.formInput} value={name} onChange={(e) => setName(e.target.value)} />
        </div>
        <div className={styles.formField}>
          <label className={styles.formLabel}>Tier (1-10)</label>
          <input className={styles.formInput} type="number" min={1} max={10} value={tier} onChange={(e) => setTier(+e.target.value)} />
        </div>
        <div className={styles.formField}>
          <label className={styles.formLabel}>Min Beasts</label>
          <input className={styles.formInput} type="number" min={1} value={minBeasts} onChange={(e) => setMinBeasts(+e.target.value)} />
        </div>
        <div className={styles.formField}>
          <label className={styles.formLabel}>Max Beasts</label>
          <input className={styles.formInput} type="number" min={1} value={maxBeasts} onChange={(e) => setMaxBeasts(+e.target.value)} />
        </div>
        <div className={styles.formField}>
          <label className={styles.formLabel}>Reward Table ID</label>
          <input className={styles.formInput} type="number" min={1} value={rewardTableId} onChange={(e) => setRewardTableId(+e.target.value)} />
        </div>
        <div className={styles.formField}>
          <label className={styles.formLabel}>Description</label>
          <input className={styles.formInput} value={desc} onChange={(e) => setDesc(e.target.value)} />
        </div>
        <div className={styles.formActions}>
          <button className={styles.cancelFormBtn} onClick={onCancel}>Cancel</button>
          <button className={styles.saveBtn} disabled={!name} onClick={() => onSave({
            name,
            tier,
            min_beasts: minBeasts,
            max_beasts: maxBeasts,
            reward_table_id: rewardTableId,
            description: desc,
            beast_slots: [],
          })}>Save</button>
        </div>
      </div>
    </div>
  );
}

// ── Spawn Rules Section ─────────────────────────────────────────────────────

function SpawnRulesSection({ onError, onSuccess }: SectionProps) {
  const [items, setItems] = useState<SpawnRuleResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editId, setEditId] = useState<number | null>(null);

  const load = useCallback(async () => {
    setLoading(true);
    try {
      setItems(await fetchSpawnRules());
    } catch {
      onError('Failed to load spawn rules.');
    } finally {
      setLoading(false);
    }
  }, [onError]);

  useEffect(() => { load(); }, [load]);

  const handleDelete = async (id: number) => {
    if (!confirm('Delete this spawn rule?')) return;
    try {
      await deleteSpawnRule(id);
      onSuccess('Spawn rule deleted.');
      load();
    } catch {
      onError('Failed to delete spawn rule.');
    }
  };

  const handleToggle = async (rule: SpawnRuleResponse) => {
    try {
      await updateSpawnRule(rule.id, { enabled: !rule.enabled });
      onSuccess(`Spawn rule ${!rule.enabled ? 'enabled' : 'disabled'}.`);
      load();
    } catch {
      onError('Failed to toggle spawn rule.');
    }
  };

  const handleSave = async (data: CreateSpawnRuleRequest) => {
    try {
      if (editId) {
        await updateSpawnRule(editId, data);
        onSuccess('Spawn rule updated.');
      } else {
        await createSpawnRule(data);
        onSuccess('Spawn rule created.');
      }
      setShowForm(false);
      setEditId(null);
      load();
    } catch {
      onError('Failed to save spawn rule.');
    }
  };

  if (loading) return <div className={styles.center}><LoadingSpinner size="md" /></div>;

  const editItem = editId ? items.find((i) => i.id === editId) : undefined;

  return (
    <>
      <div className={styles.sectionHeader}>
        <span>{items.length} spawn rule{items.length !== 1 ? 's' : ''}</span>
        <button className={styles.addBtn} onClick={() => { setEditId(null); setShowForm(true); }}>
          + New Rule
        </button>
      </div>

      {items.length === 0 ? (
        <div className={styles.empty}>No spawn rules yet.</div>
      ) : (
        <div className={styles.cardList}>
          {items.map((r) => (
            <div className={styles.card} key={r.id}>
              <div className={styles.cardHeader}>
                <h3 className={styles.cardTitle}>
                  {r.name}{' '}
                  <span className={`${styles.badge} ${r.enabled ? styles.badgeEnabled : styles.badgeDisabled}`}>
                    {r.enabled ? 'Enabled' : 'Disabled'}
                  </span>
                </h3>
                <div className={styles.cardActions}>
                  <button className={styles.editBtn} onClick={() => handleToggle(r)}>
                    {r.enabled ? 'Disable' : 'Enable'}
                  </button>
                  <button className={styles.editBtn} onClick={() => { setEditId(r.id); setShowForm(true); }}>Edit</button>
                  <button className={styles.deleteBtn} onClick={() => handleDelete(r.id)}>Delete</button>
                </div>
              </div>
              <div className={styles.cardRow}>
                <span className={styles.cardLabel}>Max Camps</span>
                <span className={styles.cardValue}>{r.max_camps}</span>
              </div>
              <div className={styles.cardRow}>
                <span className={styles.cardLabel}>Spawn Interval</span>
                <span className={styles.cardValue}>{r.spawn_interval_sec}s</span>
              </div>
              <div className={styles.cardRow}>
                <span className={styles.cardLabel}>Despawn After</span>
                <span className={styles.cardValue}>{r.despawn_after_sec}s</span>
              </div>
              <div className={styles.cardRow}>
                <span className={styles.cardLabel}>Terrain</span>
                <span className={styles.cardValue}>{r.terrain_types.join(', ') || 'any'}</span>
              </div>
              <div className={styles.cardRow}>
                <span className={styles.cardLabel}>Zones</span>
                <span className={styles.cardValue}>{r.zone_types.join(', ') || 'any'}</span>
              </div>
              <div className={styles.cardRow}>
                <span className={styles.cardLabel}>Templates</span>
                <span className={styles.cardValue}>{r.camp_template_pool.length} in pool</span>
              </div>
            </div>
          ))}
        </div>
      )}

      {showForm && (
        <SpawnRuleForm
          initial={editItem}
          onSave={handleSave}
          onCancel={() => { setShowForm(false); setEditId(null); }}
        />
      )}
    </>
  );
}

function SpawnRuleForm({
  initial,
  onSave,
  onCancel,
}: {
  initial?: SpawnRuleResponse;
  onSave: (data: CreateSpawnRuleRequest) => void;
  onCancel: () => void;
}) {
  const [name, setName] = useState(initial?.name ?? '');
  const [terrainTypes, setTerrainTypes] = useState(initial?.terrain_types.join(', ') ?? '');
  const [zoneTypes, setZoneTypes] = useState(initial?.zone_types.join(', ') ?? '');
  const [maxCamps, setMaxCamps] = useState(initial?.max_camps ?? 10);
  const [spawnInterval, setSpawnInterval] = useState(initial?.spawn_interval_sec ?? 60);
  const [despawnAfter, setDespawnAfter] = useState(initial?.despawn_after_sec ?? 3600);
  const [minCampDist, setMinCampDist] = useState(initial?.min_camp_distance ?? 2);
  const [minVillageDist, setMinVillageDist] = useState(initial?.min_village_distance ?? 3);
  const [templatePool, setTemplatePool] = useState(
    initial?.camp_template_pool.map((p) => `${p.camp_template_id}:${p.weight}`).join(', ') ?? ''
  );

  const parsePool = () =>
    templatePool
      .split(',')
      .map((s) => s.trim())
      .filter(Boolean)
      .map((s) => {
        const parts = s.split(':');
        return { camp_template_id: parseInt(parts[0] ?? '', 10), weight: parseInt(parts[1] ?? '1', 10) };
      })
      .filter((p) => !isNaN(p.camp_template_id));

  return (
    <div className={styles.formOverlay} onClick={onCancel}>
      <div className={styles.formModal} onClick={(e) => e.stopPropagation()}>
        <h3 className={styles.formTitle}>{initial ? 'Edit Spawn Rule' : 'New Spawn Rule'}</h3>
        <div className={styles.formField}>
          <label className={styles.formLabel}>Name</label>
          <input className={styles.formInput} value={name} onChange={(e) => setName(e.target.value)} />
        </div>
        <div className={styles.formField}>
          <label className={styles.formLabel}>Terrain Types (comma-separated)</label>
          <input className={styles.formInput} value={terrainTypes} onChange={(e) => setTerrainTypes(e.target.value)} placeholder="forest, plains" />
        </div>
        <div className={styles.formField}>
          <label className={styles.formLabel}>Zone Types (comma-separated)</label>
          <input className={styles.formInput} value={zoneTypes} onChange={(e) => setZoneTypes(e.target.value)} placeholder="wilderness, sylvara" />
        </div>
        <div className={styles.formField}>
          <label className={styles.formLabel}>Max Camps</label>
          <input className={styles.formInput} type="number" min={1} value={maxCamps} onChange={(e) => setMaxCamps(+e.target.value)} />
        </div>
        <div className={styles.formField}>
          <label className={styles.formLabel}>Spawn Interval (seconds)</label>
          <input className={styles.formInput} type="number" min={1} value={spawnInterval} onChange={(e) => setSpawnInterval(+e.target.value)} />
        </div>
        <div className={styles.formField}>
          <label className={styles.formLabel}>Despawn After (seconds)</label>
          <input className={styles.formInput} type="number" min={1} value={despawnAfter} onChange={(e) => setDespawnAfter(+e.target.value)} />
        </div>
        <div className={styles.formField}>
          <label className={styles.formLabel}>Min Camp Distance</label>
          <input className={styles.formInput} type="number" min={0} value={minCampDist} onChange={(e) => setMinCampDist(+e.target.value)} />
        </div>
        <div className={styles.formField}>
          <label className={styles.formLabel}>Min Village Distance</label>
          <input className={styles.formInput} type="number" min={0} value={minVillageDist} onChange={(e) => setMinVillageDist(+e.target.value)} />
        </div>
        <div className={styles.formField}>
          <label className={styles.formLabel}>Template Pool (id:weight, comma-separated)</label>
          <input className={styles.formInput} value={templatePool} onChange={(e) => setTemplatePool(e.target.value)} placeholder="1:50, 2:30, 3:20" />
        </div>
        <div className={styles.formActions}>
          <button className={styles.cancelFormBtn} onClick={onCancel}>Cancel</button>
          <button className={styles.saveBtn} disabled={!name} onClick={() => onSave({
            name,
            terrain_types: terrainTypes.split(',').map((s) => s.trim()).filter(Boolean),
            zone_types: zoneTypes.split(',').map((s) => s.trim()).filter(Boolean),
            camp_template_pool: parsePool(),
            max_camps: maxCamps,
            spawn_interval_sec: spawnInterval,
            despawn_after_sec: despawnAfter,
            min_camp_distance: minCampDist,
            min_village_distance: minVillageDist,
          })}>Save</button>
        </div>
      </div>
    </div>
  );
}

// ── Reward Tables Section ───────────────────────────────────────────────────

function RewardTablesSection({ onError, onSuccess }: SectionProps) {
  const [items, setItems] = useState<RewardTableResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);

  const load = useCallback(async () => {
    setLoading(true);
    try {
      setItems(await fetchRewardTables());
    } catch {
      onError('Failed to load reward tables.');
    } finally {
      setLoading(false);
    }
  }, [onError]);

  useEffect(() => { load(); }, [load]);

  const handleDelete = async (id: number) => {
    if (!confirm('Delete this reward table?')) return;
    try {
      await deleteRewardTable(id);
      onSuccess('Reward table deleted.');
      load();
    } catch {
      onError('Failed to delete reward table.');
    }
  };

  const handleCreate = async (data: CreateRewardTableRequest) => {
    try {
      await createRewardTable(data);
      onSuccess('Reward table created.');
      setShowForm(false);
      load();
    } catch {
      onError('Failed to create reward table.');
    }
  };

  const handleRename = async (id: number) => {
    const newName = prompt('New name:');
    if (!newName) return;
    try {
      await updateRewardTable(id, { name: newName });
      onSuccess('Reward table renamed.');
      load();
    } catch {
      onError('Failed to rename reward table.');
    }
  };

  if (loading) return <div className={styles.center}><LoadingSpinner size="md" /></div>;

  return (
    <>
      <div className={styles.sectionHeader}>
        <span>{items.length} reward table{items.length !== 1 ? 's' : ''}</span>
        <button className={styles.addBtn} onClick={() => setShowForm(true)}>+ New Table</button>
      </div>

      {items.length === 0 ? (
        <div className={styles.empty}>No reward tables yet.</div>
      ) : (
        <div className={styles.cardList}>
          {items.map((t) => (
            <div className={styles.card} key={t.id}>
              <div className={styles.cardHeader}>
                <h3 className={styles.cardTitle}>{t.name}</h3>
                <div className={styles.cardActions}>
                  <button className={styles.editBtn} onClick={() => handleRename(t.id)}>Rename</button>
                  <button className={styles.deleteBtn} onClick={() => handleDelete(t.id)}>Delete</button>
                </div>
              </div>
              {t.entries.length > 0 ? (
                t.entries.map((e) => (
                  <div className={styles.cardRow} key={e.id}>
                    <span className={styles.cardLabel}>{e.reward_type}</span>
                    <span className={styles.cardValue}>
                      {e.min_amount}–{e.max_amount} ({e.drop_chance_pct}%)
                    </span>
                  </div>
                ))
              ) : (
                <div className={styles.cardRow}>
                  <span className={styles.cardLabel}>No entries</span>
                  <span className={styles.cardValue}>—</span>
                </div>
              )}
            </div>
          ))}
        </div>
      )}

      {showForm && (
        <RewardTableForm
          onSave={handleCreate}
          onCancel={() => setShowForm(false)}
        />
      )}
    </>
  );
}

function RewardTableForm({
  onSave,
  onCancel,
}: {
  onSave: (data: CreateRewardTableRequest) => void;
  onCancel: () => void;
}) {
  const [name, setName] = useState('');

  return (
    <div className={styles.formOverlay} onClick={onCancel}>
      <div className={styles.formModal} onClick={(e) => e.stopPropagation()}>
        <h3 className={styles.formTitle}>New Reward Table</h3>
        <div className={styles.formField}>
          <label className={styles.formLabel}>Name</label>
          <input className={styles.formInput} value={name} onChange={(e) => setName(e.target.value)} />
        </div>
        <div className={styles.formActions}>
          <button className={styles.cancelFormBtn} onClick={onCancel}>Cancel</button>
          <button className={styles.saveBtn} disabled={!name} onClick={() => onSave({ name, entries: [] })}>
            Create
          </button>
        </div>
      </div>
    </div>
  );
}

// ── Battle Tuning Section ───────────────────────────────────────────────────

function BattleTuningSection({ onError, onSuccess }: SectionProps) {
  const [tuning, setTuning] = useState<BattleTuningResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  const load = useCallback(async () => {
    setLoading(true);
    try {
      setTuning(await fetchBattleTuning());
    } catch {
      onError('Failed to load battle tuning.');
    } finally {
      setLoading(false);
    }
  }, [onError]);

  useEffect(() => { load(); }, [load]);

  const handleSave = async () => {
    if (!tuning) return;
    setSaving(true);
    try {
      const updated = await updateBattleTuning(tuning);
      setTuning(updated);
      onSuccess('Battle tuning updated.');
    } catch {
      onError('Failed to update battle tuning.');
    } finally {
      setSaving(false);
    }
  };

  if (loading) return <div className={styles.center}><LoadingSpinner size="md" /></div>;
  if (!tuning) return <div className={styles.empty}>No tuning data available.</div>;

  const update = (key: keyof BattleTuningResponse, value: number) => {
    setTuning((prev) => prev ? { ...prev, [key]: value } : prev);
  };

  return (
    <>
      <h3 className={styles.formTitle}>Battle Tuning</h3>
      <div className={styles.tuningGrid}>
        <div className={styles.tuningField}>
          <label className={styles.tuningLabel}>Tick Duration (ms)</label>
          <input className={styles.tuningInput} type="number" min={50} value={tuning.tick_duration_ms} onChange={(e) => update('tick_duration_ms', +e.target.value)} />
        </div>
        <div className={styles.tuningField}>
          <label className={styles.tuningLabel}>Crit Damage Multiplier</label>
          <input className={styles.tuningInput} type="number" min={1} step={0.1} value={tuning.crit_damage_multiplier} onChange={(e) => update('crit_damage_multiplier', +e.target.value)} />
        </div>
        <div className={styles.tuningField}>
          <label className={styles.tuningLabel}>Max Defense %</label>
          <input className={styles.tuningInput} type="number" min={0} max={100} step={1} value={tuning.max_defense_percent} onChange={(e) => update('max_defense_percent', +e.target.value)} />
        </div>
        <div className={styles.tuningField}>
          <label className={styles.tuningLabel}>Max Crit %</label>
          <input className={styles.tuningInput} type="number" min={0} max={100} step={1} value={tuning.max_crit_chance_percent} onChange={(e) => update('max_crit_chance_percent', +e.target.value)} />
        </div>
        <div className={styles.tuningField}>
          <label className={styles.tuningLabel}>Min Attack Interval</label>
          <input className={styles.tuningInput} type="number" min={1} value={tuning.min_attack_interval} onChange={(e) => update('min_attack_interval', +e.target.value)} />
        </div>
        <div className={styles.tuningField}>
          <label className={styles.tuningLabel}>March Speed (tiles/min)</label>
          <input className={styles.tuningInput} type="number" min={0.1} step={0.1} value={tuning.march_speed_tiles_per_min} onChange={(e) => update('march_speed_tiles_per_min', +e.target.value)} />
        </div>
        <div className={styles.tuningField}>
          <label className={styles.tuningLabel}>Max Ticks</label>
          <input className={styles.tuningInput} type="number" min={100} value={tuning.max_ticks} onChange={(e) => update('max_ticks', +e.target.value)} />
        </div>
      </div>
      <div className={styles.tuningActions}>
        <button className={styles.saveBtn} onClick={handleSave} disabled={saving}>
          {saving ? 'Saving...' : 'Save Tuning'}
        </button>
      </div>
    </>
  );
}
