// Floating damage / kill text VFX adapted from RTUB's shared/vfx.ts

import { Text, type Container } from 'pixi.js';
import { animateTo, type TweenOwner } from './tween';

export interface VfxOwner extends TweenOwner {
  stage: Container | null;
}

const MIN_FLOAT_MS = 350;
const MIN_TEXT_INTERVAL_MS = 50;
let _lastDamageTextTime = 0;
let _lastKillTextTime = 0;

/**
 * Show a damage number floating above (x, y) that drifts up and fades.
 * Throttled at high speeds to prevent flooding.
 */
export function showDamageText(
  owner: VfxOwner,
  x: number,
  y: number,
  damage: number,
  isCrit: boolean,
): void {
  if (!owner.stage) return;

  const now = performance.now();
  if (owner.battleSpeed > 1 && now - _lastDamageTextTime < MIN_TEXT_INTERVAL_MS) return;
  _lastDamageTextTime = now;

  const label = isCrit ? `CRIT! -${damage}` : `-${damage}`;
  const fontSize = isCrit ? 16 : 13;
  const fillColor = isCrit ? 0xffff00 : 0xff4444;
  const floatDist = isCrit ? 50 : 35;
  const duration = isCrit ? 900 : 700;

  const text = new Text({
    text: label,
    style: {
      fontFamily: 'Arial',
      fontSize,
      fontWeight: 'bold',
      fill: fillColor,
      stroke: { color: 0x000000, width: 3 },
    },
  });
  text.anchor.set(0.5);
  text.x = x + (Math.random() - 0.5) * 20;
  text.y = y;
  owner.stage.addChild(text);

  const speed = owner.battleSpeed || 1;
  const compensated = Math.max(duration, MIN_FLOAT_MS * speed);

  animateTo(owner, text, { y: text.y - floatDist, alpha: 0 }, compensated, () => {
    text.destroy();
  });
}

/**
 * Show a "KILLED!" text floating above (x, y).
 * Throttled at high speeds.
 */
export function showKillText(
  owner: VfxOwner,
  x: number,
  y: number,
): void {
  if (!owner.stage) return;

  const now = performance.now();
  if (owner.battleSpeed > 1 && now - _lastKillTextTime < MIN_TEXT_INTERVAL_MS) return;
  _lastKillTextTime = now;

  const text = new Text({
    text: 'KILLED!',
    style: {
      fontFamily: 'Arial',
      fontSize: 14,
      fontWeight: 'bold',
      fill: 0xff6600,
      stroke: { color: 0x000000, width: 3 },
    },
  });
  text.anchor.set(0.5);
  text.x = x;
  text.y = y - 10;
  owner.stage.addChild(text);

  const speed = owner.battleSpeed || 1;
  const compensated = Math.max(1000, MIN_FLOAT_MS * speed);

  animateTo(owner, text, { y: text.y - 45, alpha: 0 }, compensated, () => {
    text.destroy();
  });
}
