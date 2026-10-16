// Small inline icons (no icon library). Decorative: always paired with visible text and hidden from assistive tech.
const svg = { width: 24, height: 24, viewBox: '0 0 24 24', 'aria-hidden': true, focusable: false } as const

export const AnvilIcon = () => (
  <svg {...svg} viewBox="0 0 32 24">
    <path
      fill="currentColor"
      d="M1 4h22c4 0 7 2 8 5-3 0-6 0-8 1-2 1-3 3-3 5h-3l1 3h4v3H6v-3h4l1-3H9C9 10 6 8 1 8z"
    />
  </svg>
)

export const SettlementIcon = () => (
  <svg {...svg}>
    <path fill="currentColor" d="M12 3 2 11h3v9h5v-6h4v6h5v-9h3z" />
  </svg>
)

export const SwordsIcon = () => (
  <svg {...svg} fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round">
    <path d="M4 4l10 10M20 4 10 14M6 16l2 2M18 16l-2 2M4 20l3-3M20 20l-3-3" />
  </svg>
)

export const CrownIcon = () => (
  <svg {...svg}>
    <path fill="currentColor" d="M3 7l4.5 4L12 4l4.5 7L21 7l-2 11H5zM5 19.5h14V21H5z" />
  </svg>
)

export const ChevronDownIcon = () => (
  <svg {...svg} width={16} height={16} fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round">
    <path d="M6 9l6 6 6-6" />
  </svg>
)

// Game navigation (task 012).
export const MapIcon = () => (
  <svg {...svg} fill="none" stroke="currentColor" strokeWidth="2" strokeLinejoin="round">
    <path d="M3 6l6-2 6 2 6-2v14l-6 2-6-2-6 2zM9 4v14M15 6v14" />
  </svg>
)

export const HammerIcon = () => (
  <svg {...svg}>
    <path fill="currentColor" d="M20.5 12 17 15.5 8.5 7 12 3.5z" />
    <path stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" d="M4 20l9.5-9.5" />
  </svg>
)

export const ShieldIcon = () => (
  <svg {...svg}>
    <path fill="currentColor" d="M12 2 4 5v6c0 5 3.4 9.3 8 11 4.6-1.7 8-6 8-11V5zm0 3v14.8C8.6 18.2 6.5 15 6.5 11V6.7z" />
  </svg>
)

export const ScrollIcon = () => (
  <svg {...svg} fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <path d="M7 3h11a2 2 0 0 1 2 2v1h-4M7 3a2 2 0 0 0-2 2v13a3 3 0 0 0 3 3h10a2 2 0 0 0 2-2v-1H10v1a2 2 0 0 1-4 0M9 8h6M9 12h6" />
  </svg>
)
