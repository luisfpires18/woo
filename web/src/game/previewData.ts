// ILLUSTRATIVE PREVIEW DATA ONLY (task 012). The final resource types are not decided (docs/open-questions.md), so
// these are numbered placeholder slots, not a resource roster: no names, icons or amounts. Real definitions arrive
// with 025 RESOURCES-ADMIN and real amounts with 039 RESOURCE-STATE; they plug into the same ResourceSlot shape.

/** Resource-header contract: one entry per resource the server reports for the current settlement. */
export type ResourceSlot = {
  id: string
  /** Player-facing name from the admin resource definition. */
  label: string
  /** Current amount; null while not available (shown as a dash and announced as "Not available yet"). */
  amount: number | null
  /** Storage limit when the resource has one. */
  capacity?: number | null
  /** Admin-uploaded image; the shell shows a neutral fallback symbol without one. */
  iconUrl?: string
}

export const previewResourceSlots: readonly ResourceSlot[] = [1, 2, 3, 4].map((n) => ({
  id: `preview-${n}`,
  label: `Resource ${n}`,
  amount: null,
}))

/**
 * Settlement inspector fixtures (task 013). ILLUSTRATIVE ONLY: mock levels and example upgrade figures that shape the
 * inspector layout. They are not costs, timings or mechanics; resource slots reuse the numbered preview names above.
 */
export type PreviewBuilding = {
  summary: string
  /** Mock level shown for each prototype visual stage. */
  mockLevel: Record<number, number>
  illustrativeUpgrade?: { cost: { label: string; amount: number }[]; duration: string }
}

export const previewBuildings: Record<string, PreviewBuilding> = {
  forge: {
    summary: 'Where weapons and armour will be improved. Two visual stages prove that a building can be replaced at the same plot.',
    mockLevel: { 1: 1, 2: 5 },
    illustrativeUpgrade: { cost: [{ label: 'Resource 1', amount: 120 }, { label: 'Resource 2', amount: 80 }], duration: '12 min' },
  },
  cottage: {
    summary: 'An ordinary dwelling placeholder. Its image slot points at artwork that does not exist yet, so the fallback shape is shown.',
    mockLevel: { 1: 1 },
  },
  hall: {
    summary: 'A larger placeholder building whose base sits behind the south wall to prove foreground overlap.',
    mockLevel: { 1: 2 },
  },
}
