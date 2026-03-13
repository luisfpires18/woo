import { create } from 'zustand';
import type { CampResponse, ExpeditionResponse, BattleReportResponse } from '../types/api';

export interface IncomingAttack {
  id: number;
  village_id: number;
  arrives_at: string;
}

interface ExpeditionState {
  /** Active camps on the world map */
  camps: CampResponse[];
  /** Currently selected camp (for detail panel / attack modal) */
  selectedCamp: CampResponse | null;
  /** Player's expeditions */
  expeditions: ExpeditionResponse[];
  /** Cached battle reports keyed by battle ID */
  battleReports: Record<number, BattleReportResponse>;
  /** Battle IDs that have been viewed (replay consumed) */
  viewedBattles: Set<number>;
  /** Whether camps have been fetched at least once */
  campsLoaded: boolean;
  /** Incoming attacks on player's villages (stub — wired when backend supports it) */
  incomingAttacks: IncomingAttack[];
  /** IDs of dismissed expeditions — persists across refetches */
  dismissedIds: Set<number>;

  setCamps: (camps: CampResponse[]) => void;
  setSelectedCamp: (camp: CampResponse | null) => void;
  setExpeditions: (expeditions: ExpeditionResponse[]) => void;
  addExpedition: (expedition: ExpeditionResponse) => void;
  updateExpedition: (id: number, updates: Partial<ExpeditionResponse>) => void;
  cacheBattleReport: (report: BattleReportResponse) => void;
  /** Mark a battle as viewed so its replay can't be watched again */
  markBattleViewed: (battleId: number) => void;
  /** Remove a camp when it's cleared or despawned */
  removeCamp: (campId: number) => void;
  /** Dismiss (hide) a completed expedition from the panel */
  dismissExpedition: (id: number) => void;
  /** Set all incoming attacks (from API/WS refresh) */
  setIncomingAttacks: (attacks: IncomingAttack[]) => void;
  /** Add a single incoming attack (from WS push) */
  addIncomingAttack: (attack: IncomingAttack) => void;
}

export const useExpeditionStore = create<ExpeditionState>((set) => ({
  camps: [],
  selectedCamp: null,
  expeditions: [],
  battleReports: {},
  viewedBattles: new Set(),
  campsLoaded: false,
  incomingAttacks: [],
  dismissedIds: new Set(),

  setCamps: (camps) => set({ camps, campsLoaded: true }),
  setSelectedCamp: (camp) => set({ selectedCamp: camp }),
  setExpeditions: (expeditions) =>
    set((state) => ({
      expeditions: expeditions.filter((e) => !state.dismissedIds.has(e.id)),
    })),
  addExpedition: (expedition) =>
    set((state) => ({ expeditions: [...state.expeditions, expedition] })),
  updateExpedition: (id, updates) =>
    set((state) => ({
      expeditions: state.expeditions.map((e) =>
        e.id === id ? { ...e, ...updates } : e,
      ),
    })),
  cacheBattleReport: (report) =>
    set((state) => ({
      battleReports: { ...state.battleReports, [report.id]: report },
    })),
  markBattleViewed: (battleId) =>
    set((state) => {
      const next = new Set(state.viewedBattles);
      next.add(battleId);
      return { viewedBattles: next };
    }),
  removeCamp: (campId) =>
    set((state) => ({
      camps: state.camps.filter((c) => c.id !== campId),
      selectedCamp: state.selectedCamp?.id === campId ? null : state.selectedCamp,
    })),
  dismissExpedition: (id) =>
    set((state) => {
      const next = new Set(state.dismissedIds);
      next.add(id);
      return {
        dismissedIds: next,
        expeditions: state.expeditions.filter((e) => e.id !== id),
      };
    }),
  setIncomingAttacks: (attacks) => set({ incomingAttacks: attacks }),
  addIncomingAttack: (attack) =>
    set((state) => ({ incomingAttacks: [...state.incomingAttacks, attack] })),
}));
