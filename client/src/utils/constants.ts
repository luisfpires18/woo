// Client-side constants

import type { Kingdom } from '../types/game';

/** All kingdoms (playable + NPC). Used for theme and validation. */
export const VALID_KINGDOMS: readonly Kingdom[] = [
  'veridor', 'sylvara', 'arkazia', 'draxys',
  'zandres', 'lumus', 'nordalh', 'drakanith',
];

/** Map dimensions: 401x401, coords -200 to +200 */
export const MAP_SIZE = 401;
export const MAP_MIN = -200;
export const MAP_MAX = 200;

/** Starting resources per village */
export const STARTING_RESOURCES = 500;

/** Breakpoints (use raw values in @media queries) */
export const BREAKPOINT_MOBILE = 768;
export const BREAKPOINT_TABLET = 1024;
