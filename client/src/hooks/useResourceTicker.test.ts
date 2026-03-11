import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useResourceTicker } from './useResourceTicker';
import type { ResourcesResponse } from '../types/api';

function makeSnapshot(overrides: Partial<ResourcesResponse> = {}): ResourcesResponse {
  return {
    food: 500,
    water: 500,
    lumber: 500,
    stone: 500,
    food_rate: 10,
    water_rate: 10,
    lumber_rate: 10,
    stone_rate: 10,
    food_consumption: 2,
    max_food: 5000,
    max_water: 5000,
    max_lumber: 5000,
    max_stone: 5000,
    ...overrides,
  };
}

describe('useResourceTicker', () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('returns initial snapshot values immediately', () => {
    const snapshot = makeSnapshot();
    const { result } = renderHook(() => useResourceTicker(snapshot));

    expect(result.current.food).toBe(500);
    expect(result.current.water).toBe(500);
    expect(result.current.lumber).toBe(500);
    expect(result.current.stone).toBe(500);
  });

  it('interpolates resources over time using rates', () => {
    const snapshot = makeSnapshot({ food: 100, food_rate: 10, food_consumption: 0 });
    const { result } = renderHook(() => useResourceTicker(snapshot));

    // Advance 5 seconds
    act(() => {
      vi.advanceTimersByTime(5000);
    });

    // food should be ~100 + 10*5 = 150 (approximately)
    expect(result.current.food).toBeGreaterThanOrEqual(140);
    expect(result.current.food).toBeLessThanOrEqual(160);
  });

  it('subtracts food_consumption from food rate', () => {
    const snapshot = makeSnapshot({ food: 100, food_rate: 10, food_consumption: 8 });
    const { result } = renderHook(() => useResourceTicker(snapshot));

    // Advance 5 seconds: net food rate = 10-8 = 2/s → +10 food
    act(() => {
      vi.advanceTimersByTime(5000);
    });

    // food should be ~110 (100 + 2*5)
    expect(result.current.food).toBeGreaterThanOrEqual(105);
    expect(result.current.food).toBeLessThanOrEqual(115);
  });

  it('clamps resources at their per-resource max', () => {
    const snapshot = makeSnapshot({
      food: 990,
      food_rate: 100,
      food_consumption: 0,
      max_food: 1000,
    });
    const { result } = renderHook(() => useResourceTicker(snapshot));

    act(() => {
      vi.advanceTimersByTime(5000);
    });

    expect(result.current.food).toBe(1000); // capped at max_food
  });

  it('uses different caps for different resources', () => {
    const snapshot = makeSnapshot({
      food: 900, water: 900, lumber: 900, stone: 900,
      food_rate: 100, water_rate: 100, lumber_rate: 100, stone_rate: 100,
      food_consumption: 0,
      max_food: 1000,
      max_water: 2000,
      max_lumber: 1500,
      max_stone: 1000,
    });
    const { result } = renderHook(() => useResourceTicker(snapshot));

    act(() => {
      vi.advanceTimersByTime(15000); // 15 seconds — enough for all to max
    });

    expect(result.current.food).toBe(1000);   // capped at max_food
    expect(result.current.water).toBe(2000);   // capped at max_water
    expect(result.current.lumber).toBe(1500);  // capped at max_lumber
    expect(result.current.stone).toBe(1000);   // capped at max_stone
  });

  it('floors food at zero when consumption exceeds rate', () => {
    const snapshot = makeSnapshot({
      food: 10,
      food_rate: 1,
      food_consumption: 100,
    });
    const { result } = renderHook(() => useResourceTicker(snapshot));

    act(() => {
      vi.advanceTimersByTime(2000);
    });

    expect(result.current.food).toBe(0);
  });

  it('resets baseline when snapshot reference changes', () => {
    const snapshot1 = makeSnapshot({ food: 100, food_rate: 10, food_consumption: 0 });
    const { result, rerender } = renderHook(
      ({ s }) => useResourceTicker(s),
      { initialProps: { s: snapshot1 } },
    );

    // Advance 3 seconds: food ≈ 100 + 10*3 = 130
    act(() => {
      vi.advanceTimersByTime(3000);
    });
    expect(result.current.food).toBeGreaterThan(120);

    // New snapshot from server with higher food
    const snapshot2 = makeSnapshot({ food: 500, food_rate: 10, food_consumption: 0 });
    rerender({ s: snapshot2 });

    // Should reset to new baseline (~500, not ~630)
    expect(result.current.food).toBeGreaterThanOrEqual(499);
    expect(result.current.food).toBeLessThan(510);
  });
});
