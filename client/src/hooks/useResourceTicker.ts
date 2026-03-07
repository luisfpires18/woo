import { useEffect, useRef, useState } from 'react';
import type { ResourcesResponse } from '../types/api';

/**
 * Client-side resource interpolation.
 *
 * Takes the latest server snapshot and ticks every second using the rates,
 * clamped at max_storage.  Resets its baseline whenever the snapshot reference
 * changes (i.e. after an API refetch).
 */
export function useResourceTicker(snapshot: ResourcesResponse): ResourcesResponse {
  // Track when we received this snapshot so we can compute elapsed seconds.
  const baselineRef = useRef<{ snapshot: ResourcesResponse; receivedAt: number }>({
    snapshot,
    receivedAt: Date.now(),
  });

  // Reset baseline whenever the snapshot object identity changes (new fetch)
  const prevSnapshot = useRef(snapshot);
  if (snapshot !== prevSnapshot.current) {
    prevSnapshot.current = snapshot;
    baselineRef.current = { snapshot, receivedAt: Date.now() };
  }

  const [tick, setTick] = useState(0);

  useEffect(() => {
    const id = setInterval(() => setTick((t) => t + 1), 1000);
    return () => clearInterval(id);
  }, []);

  // Suppress unused-var lint — tick drives re-renders
  void tick;

  const { snapshot: base, receivedAt } = baselineRef.current;
  const elapsed = (Date.now() - receivedAt) / 1000; // seconds since snapshot
  const cap = base.max_storage;

  return {
    ...base,
    food: Math.max(0, Math.min(base.food + (base.food_rate - base.food_consumption) * elapsed, cap)),
    water: Math.min(base.water + base.water_rate * elapsed, cap),
    lumber: Math.min(base.lumber + base.lumber_rate * elapsed, cap),
    stone: Math.min(base.stone + base.stone_rate * elapsed, cap),
  };
}
