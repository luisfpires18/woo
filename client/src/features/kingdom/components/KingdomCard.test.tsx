import { describe, it, expect } from 'vitest';
import { KINGDOMS } from './KingdomCard';

describe('KINGDOMS constant', () => {
  it('contains all 8 kingdoms', () => {
    expect(KINGDOMS).toHaveLength(8);
  });

  it('has unique IDs for all kingdoms', () => {
    const ids = KINGDOMS.map((k) => k.id);
    expect(new Set(ids).size).toBe(8);
  });

  it('includes the 7 currently playable kingdoms', () => {
    const playable = KINGDOMS.filter((k) => k.playable);
    const playableIds = playable.map((k) => k.id);

    expect(playable).toHaveLength(7);
    expect(playableIds).toContain('veridor');
    expect(playableIds).toContain('sylvara');
    expect(playableIds).toContain('arkazia');
    expect(playableIds).toContain('draxys');
    expect(playableIds).toContain('nordalh');
    expect(playableIds).toContain('zandres');
    expect(playableIds).toContain('lumus');
  });

  it('marks drakanith as not playable', () => {
    const drakanith = KINGDOMS.find((k) => k.id === 'drakanith');
    expect(drakanith).toBeDefined();
    expect(drakanith!.playable).toBe(false);
  });

  it('every kingdom has required display properties', () => {
    for (const k of KINGDOMS) {
      expect(k.name.length).toBeGreaterThan(0);
      expect(k.tagline.length).toBeGreaterThan(0);
      expect(k.description.length).toBeGreaterThan(0);
      expect(k.traits.length).toBeGreaterThanOrEqual(1);
      expect(k.colorVar).toBeTruthy();
      expect(k.glowVar).toBeTruthy();
    }
  });

  it('locked kingdoms have a lockReason', () => {
    const locked = KINGDOMS.filter((k) => !k.playable);
    for (const k of locked) {
      expect(k.lockReason).toBeTruthy();
    }
  });
});
