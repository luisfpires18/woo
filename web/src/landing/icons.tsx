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
