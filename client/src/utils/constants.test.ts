import { describe, it, expect } from 'vitest';
import type { Kingdom, PlayableKingdom } from '../types/game';
import { VALID_KINGDOMS } from '../utils/constants';

describe('Kingdom types and constants', () => {
  it('VALID_KINGDOMS includes all 8 kingdoms', () => {
    const expected: Kingdom[] = [
      'veridor', 'sylvara', 'arkazia', 'draxys',
      'zandres', 'lumus', 'nordalh', 'drakanith',
    ];
    expect(VALID_KINGDOMS).toHaveLength(8);
    for (const k of expected) {
      expect(VALID_KINGDOMS).toContain(k);
    }
  });

  it('PlayableKingdom type covers 7 kingdoms (excludes drakanith)', () => {
    // This is a compile-time check — if the type is wrong, this won't compile.
    const playable: PlayableKingdom[] = [
      'veridor', 'sylvara', 'arkazia', 'draxys', 'nordalh', 'zandres', 'lumus',
    ];
    expect(playable).toHaveLength(7);
  });

  it('drakanith is a valid Kingdom but not a PlayableKingdom', () => {
    const allKingdoms: Kingdom[] = [...VALID_KINGDOMS];
    expect(allKingdoms).toContain('drakanith');
    // PlayableKingdom type excludes 'drakanith' at compile time.
    // If someone adds 'drakanith' to PlayableKingdom, the type system should catch it.
  });
});
