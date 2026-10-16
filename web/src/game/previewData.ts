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
