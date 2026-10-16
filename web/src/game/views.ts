// Game-shell route/view model (task 012). One entry per navigation destination; the shell, its navigation and the
// document title all read from this list. Add a later screen by adding an entry and giving it real content.
// "Settlement" is the player-facing term; the legacy VILLAGE task keys remain internal.
import type { ComponentType } from 'react'
import { HammerIcon, MapIcon, ScrollIcon, SettlementIcon, ShieldIcon } from '../landing/icons'

export type GameView = {
  id: string
  path: `/game/${string}`
  label: string
  /** Bottom-bar label on narrow screens. */
  shortLabel: string
  Icon: ComponentType
  lead: string
  /** Later roadmap tasks that fill this view; shown so the placeholder never looks like a working screen. */
  planned: readonly { task: string; text: string }[]
}

export const gameViews: readonly GameView[] = [
  {
    id: 'settlement',
    path: '/game/settlement',
    label: 'Settlement',
    shortLabel: 'Settlement',
    Icon: SettlementIcon,
    lead: 'Your settlement’s buildings will be shown and managed here.',
    planned: [
      { task: '013 VILLAGE-UI', text: 'Illustrated settlement scene and building selection' },
      { task: '041 CONSTRUCTION', text: 'Building construction and upgrades' },
    ],
  },
  {
    id: 'map',
    path: '/game/map',
    label: 'Kingdom Map',
    shortLabel: 'Map',
    Icon: MapIcon,
    lead: 'The world map with districts, routes and settlements will be shown here.',
    planned: [
      { task: '014 MAP-UI', text: 'World map rendering' },
      { task: '045–048', text: 'Districts, routes, scouting and troop movement' },
    ],
  },
  {
    id: 'forge',
    path: '/game/forge',
    label: 'Forge',
    shortLabel: 'Forge',
    Icon: HammerIcon,
    lead: 'Weapon and armour work and blacksmith progression will be shown here.',
    planned: [
      { task: '016 FORGE-UI', text: 'Forge, equipment and talent screens' },
      { task: '043 FORGING', text: 'Weaponsmith and armorsmith upgrades' },
    ],
  },
  {
    id: 'army',
    path: '/game/army',
    label: 'Army',
    shortLabel: 'Army',
    Icon: ShieldIcon,
    lead: 'Your troops and their equipment will be shown here.',
    planned: [
      { task: '042 RECRUIT', text: 'Recruitment and troop inventory' },
      { task: '015 COMBAT-UI', text: 'Battalion board and battle results' },
    ],
  },
  {
    id: 'duties',
    path: '/game/duties',
    label: 'Kingdom Duties',
    shortLabel: 'Duties',
    Icon: ScrollIcon,
    lead: 'Duties you can take on for your kingdom will be listed here.',
    planned: [{ task: '049 DUTIES', text: 'Useful kingdom duties' }],
  },
]

export const defaultGamePath = gameViews[0]!.path

/** The view for a /game/... path, or undefined for an unknown one. A trailing slash is ignored. */
export const findGameView = (pathname: string) => gameViews.find((view) => view.path === pathname.replace(/\/+$/, ''))
