// ILLUSTRATIVE PREVIEW DATA ONLY. No worlds, players or activity exist yet. These values shape the landing layout
// until world APIs exist (tasks 022 WORLDS and later); they are deliberately non-authoritative (dashes, generic
// names) and are labelled as previews wherever they are shown.
import type { Kingdom } from '../theme'

export type PreviewWorld = {
  id: string
  name: string
  era: string
  pace: string
  kingdoms: Kingdom[]
}

export const previewMetrics = [
  { label: 'Open worlds', value: '—' },
  { label: 'Registered players', value: '—' },
  { label: 'Online now', value: '—' },
] as const

export const previewWorlds: PreviewWorld[] = [
  { id: 'example-1', name: 'Example world 1', era: 'Medieval', pace: 'Standard pace', kingdoms: ['arkazia', 'veridor', 'sylvara'] },
  { id: 'example-2', name: 'Example world 2', era: 'Medieval', pace: 'Standard pace', kingdoms: ['arkazia', 'veridor', 'sylvara'] },
]
