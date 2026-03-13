// RAF-based tween engine adapted from RTUB's shared/tween.ts

import type { Container } from 'pixi.js';

export interface TweenOwner {
  _rafIds: number[];
  battleSpeed: number;
}

/**
 * Animate properties of a PIXI DisplayObject from current to target over `duration` ms.
 * Duration is scaled by owner.battleSpeed (floored at 20ms).
 */
export function animateTo(
  owner: TweenOwner,
  target: Container,
  properties: Record<string, number>,
  duration: number,
  onComplete?: () => void,
): void {
  if (!target || (target as unknown as { destroyed?: boolean }).destroyed) {
    onComplete?.();
    return;
  }

  const speed = owner.battleSpeed || 1;
  const adjustedDuration = Math.max(20, speed > 0 ? duration / speed : duration);

  const startTime = Date.now();
  const container = target as unknown as Record<string, unknown>;

  // Snapshot start values
  const startValues: Record<string, number> = {};
  for (const key of Object.keys(properties)) {
    if (key === 'scale') {
      startValues[key] = target.scale?.x ?? 1;
    } else {
      startValues[key] = (container[key] as number) ?? 0;
    }
  }

  const animate = (): void => {
    const elapsed = Date.now() - startTime;
    const t = Math.min(elapsed / adjustedDuration, 1);

    try {
      for (const key of Object.keys(properties)) {
        const from = startValues[key] ?? 0;
        const to = properties[key] ?? 0;
        const value = from + (to - from) * t;
        if (key === 'scale') {
          target.scale?.set(value);
        } else {
          (container[key] as number) = value;
        }
      }
    } catch {
      onComplete?.();
      return;
    }

    if (t < 1) {
      const id = requestAnimationFrame(animate);
      owner._rafIds.push(id);
    } else {
      onComplete?.();
    }
  };

  const id = requestAnimationFrame(animate);
  owner._rafIds.push(id);
}

/**
 * Fade a target to alpha=0.
 */
export function fadeOut(
  owner: TweenOwner,
  target: Container,
  duration: number,
  onComplete?: () => void,
): void {
  animateTo(owner, target, { alpha: 0 }, duration, onComplete);
}
