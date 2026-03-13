// PixiJS battle scene — lunge attacks, smooth HP bars, speed bars, idle bob, death fade

import {
  Application,
  Container,
  Graphics,
  Sprite,
  Text,
  Texture,
  Assets,
} from 'pixi.js';
import type { BattleReplayResponse, ReplayUnit, ReplayEvent } from '../../../../types/api';
import { animateTo, fadeOut } from './tween';
import { showDamageText, showKillText, type VfxOwner } from './vfx';

// ── Layout ───────────────────────────────────────────────────────────────────

const SCENE_W = 640;
const SCENE_H = 420;
const UNIT_SIZE = 52;
const HP_BAR_W = 42;
const HP_BAR_H = 4;
const SPEED_BAR_H = 3;
const SIDE_PAD = 48;
const ROW_GAP = 72;
const TOP_PAD = 24;

// Colors
const TROOP_COLORS = [0x4a90d9, 0x5ba85b, 0xc9a227, 0xb35c3f, 0x7b5ea7, 0x3db8b8];
const BG_COLOR = 0x1a1a2e;

// ── Per-unit data ────────────────────────────────────────────────────────────

interface UnitData {
  id: number;
  side: 'attacker' | 'defender';
  name: string;
  maxHp: number;
  hp: number;
  dead: boolean;
  attackInterval: number;
  spriteKey: string;
  // Layout
  baseX: number;
  baseY: number;
  idlePhase: number;
  // PIXI objects
  container: Container;
  body: Sprite | Graphics;
  hpBg: Graphics;
  hpFill: Graphics;
  speedBg: Graphics;
  speedFill: Graphics;
  nameLabel: Text;
  // Animation state
  isAttacking: boolean;
  flashTimeout: ReturnType<typeof setTimeout> | null;
  // Speed bar: track last attack tick to interpolate cooldown
  lastAttackTick: number;
  nextAttackTick: number;
}

// ── BattleScene ──────────────────────────────────────────────────────────────

export class BattleScene implements VfxOwner {
  // VfxOwner / TweenOwner fields
  _rafIds: number[] = [];
  battleSpeed = 1;
  stage: Container | null = null;

  private app: Application | null = null;
  private containerEl: HTMLElement;
  private replay: BattleReplayResponse;
  private units: Map<number, UnitData> = new Map();
  private onComplete?: () => void;

  // Playback
  private currentTick = 0;
  private playing = false;
  private idleTime = 0;
  private lastFrameTime = 0;
  private gameLoopId = 0;
  private _timeoutIds: ReturnType<typeof setTimeout>[] = [];

  // Pre-built event index: tick → events
  private eventsByTick: Map<number, ReplayEvent[]> = new Map();

  // Resize
  private resizeHandler: (() => void) | null = null;
  private destroyed = false;

  constructor(
    container: HTMLElement,
    replay: BattleReplayResponse,
    onComplete?: () => void,
  ) {
    this.containerEl = container;
    this.replay = replay;
    this.onComplete = onComplete;

    // Build event index
    for (const e of replay.events) {
      const arr = this.eventsByTick.get(e.tick);
      if (arr) arr.push(e);
      else this.eventsByTick.set(e.tick, [e]);
    }
  }

  // ── Public API ───────────────────────────────────────────────────────────

  async init(): Promise<void> {
    this.app = new Application();
    await this.app.init({
      width: SCENE_W,
      height: SCENE_H,
      background: BG_COLOR,
      antialias: true,
      resolution: window.devicePixelRatio || 1,
      autoDensity: true,
    });

    this.containerEl.appendChild(this.app.canvas as HTMLCanvasElement);
    this.stage = this.app.stage;

    // Resize to fit parent width
    this.fitToParent();
    this.resizeHandler = () => this.fitToParent();
    window.addEventListener('resize', this.resizeHandler);

    this.drawBackground();
    await this.createUnits();
    this.precomputeSpeedBarTicks();
  }

  play(): void {
    if (this.playing) return;
    this.playing = true;
    this.lastFrameTime = performance.now();
    this.gameLoopId = requestAnimationFrame((t) => this.gameLoop(t));
    this._rafIds.push(this.gameLoopId);
  }

  pause(): void {
    this.playing = false;
  }

  setSpeed(speed: number): void {
    this.battleSpeed = speed;
  }

  seekToTick(tick: number): void {
    this.currentTick = Math.max(0, Math.min(tick, this.replay.total_ticks));
    this.rebuildStateAtTick(this.currentTick);
  }

  getCurrentTick(): number {
    return Math.floor(this.currentTick);
  }

  isPlaying(): boolean {
    return this.playing;
  }

  getAliveCount(): { attackers: number; defenders: number } {
    let attackers = 0;
    let defenders = 0;
    this.units.forEach((u) => {
      if (!u.dead) {
        if (u.side === 'attacker') attackers++;
        else defenders++;
      }
    });
    return { attackers, defenders };
  }

  destroy(): void {
    if (this.destroyed) return;
    this.destroyed = true;
    this.playing = false;

    // Cancel all RAF
    for (const id of this._rafIds) cancelAnimationFrame(id);
    this._rafIds.length = 0;

    // Cancel all timeouts
    for (const id of this._timeoutIds) clearTimeout(id);
    this._timeoutIds.length = 0;

    // Cancel unit flash timeouts
    this.units.forEach((u) => {
      if (u.flashTimeout) clearTimeout(u.flashTimeout);
    });

    // Remove resize listener
    if (this.resizeHandler) {
      window.removeEventListener('resize', this.resizeHandler);
      this.resizeHandler = null;
    }

    // Destroy PixiJS app
    if (this.app) {
      try {
        const canvas = this.app.canvas as HTMLCanvasElement;
        canvas.parentNode?.removeChild(canvas);
      } catch { /* canvas getter may fail if renderer already torn down */ }
      try {
        this.app.destroy(false, { children: true });
      } catch { /* ignore double-destroy */ }
      this.app = null;
    }
    this.stage = null;
    this.units.clear();
  }

  // ── Internals ────────────────────────────────────────────────────────────

  private fitToParent(): void {
    if (!this.app) return;
    const pw = this.containerEl.clientWidth;
    if (pw <= 0) return;
    const canvas = this.app.canvas as HTMLCanvasElement;
    const scale = pw / SCENE_W;
    canvas.style.width = `${pw}px`;
    canvas.style.height = `${SCENE_H * scale}px`;
  }

  private drawBackground(): void {
    if (!this.stage) return;

    // Centre divider
    const divider = new Graphics();
    divider.setStrokeStyle({ width: 1, color: 0x333333 });
    // Dashed line approximation
    for (let y = 0; y < SCENE_H; y += 10) {
      divider.moveTo(SCENE_W / 2, y);
      divider.lineTo(SCENE_W / 2, y + 6);
    }
    divider.stroke();
    this.stage.addChild(divider);

    // Side labels
    const atkLabel = new Text({
      text: 'ATTACKERS',
      style: { fontFamily: 'Arial', fontSize: 9, fontWeight: 'bold', fill: 0x8cb4e0 },
    });
    atkLabel.x = 4;
    atkLabel.y = 4;
    this.stage.addChild(atkLabel);

    const defLabel = new Text({
      text: 'DEFENDERS',
      style: { fontFamily: 'Arial', fontSize: 9, fontWeight: 'bold', fill: 0xe0a08c },
    });
    defLabel.anchor.set(1, 0);
    defLabel.x = SCENE_W - 4;
    defLabel.y = 4;
    this.stage.addChild(defLabel);
  }

  /**
   * Derive the sprite URL candidates for a unit.
   * Beasts: use sprite_key directly.
   * Troops: derive kingdom from name prefix and build convention URL.
   */
  private getSpriteUrls(unit: ReplayUnit): string[] {
    const urls: string[] = [];
    // Beast sprite_key takes priority
    if (unit.sprite_key) {
      urls.push(`/uploads/sprites/${unit.sprite_key}.png`);
    }
    // Convention URL for troops: /uploads/sprites/{kingdom}/units/{troop_type}.png
    const parts = unit.name.split('_');
    if (parts.length >= 2) {
      const kingdom = parts[0];
      urls.push(`/uploads/sprites/${kingdom}/units/${unit.name}.png`);
    }
    return urls;
  }

  private async createUnits(): Promise<void> {
    const allUnits = [
      ...this.replay.attackers.map((u) => ({ ...u, side: 'attacker' as const })),
      ...this.replay.defenders.map((u) => ({ ...u, side: 'defender' as const })),
    ];

    // Preload sprites — try each candidate URL, keep the first that loads
    const spriteTextures = new Map<number, Texture>();
    const loadPromises: Promise<void>[] = [];

    for (const u of allUnits) {
      const urls = this.getSpriteUrls(u);
      if (urls.length === 0) continue;
      loadPromises.push(
        (async () => {
          for (const url of urls) {
            try {
              const tex = await Assets.load(url) as Texture;
              if (tex) { spriteTextures.set(u.id, tex); return; }
            } catch { /* try next URL */ }
          }
        })(),
      );
    }
    await Promise.all(loadPromises);

    // Layout units
    const maxPerCol = Math.max(1, Math.floor((SCENE_H - TOP_PAD - 10) / ROW_GAP));

    const attackers = allUnits.filter((u) => u.side === 'attacker');
    const defenders = allUnits.filter((u) => u.side === 'defender');

    const createSide = (units: (ReplayUnit & { side: 'attacker' | 'defender' })[], side: 'attacker' | 'defender') => {
      const baseX = side === 'attacker' ? SIDE_PAD : SCENE_W - SIDE_PAD - UNIT_SIZE;

      units.forEach((u, i) => {
        const col = Math.floor(i / maxPerCol);
        const row = i % maxPerCol;
        const xOff = side === 'attacker' ? col * (UNIT_SIZE + 10) : -col * (UNIT_SIZE + 10);
        const x = baseX + xOff;
        const y = TOP_PAD + row * ROW_GAP;

        this.createUnitVisuals(u, side, x, y, i, spriteTextures);
      });
    };

    createSide(attackers, 'attacker');
    createSide(defenders, 'defender');
  }

  private createUnitVisuals(
    u: ReplayUnit & { side: 'attacker' | 'defender' },
    side: 'attacker' | 'defender',
    x: number,
    y: number,
    index: number,
    spriteTextures: Map<number, Texture>,
  ): void {
    if (!this.stage) return;

    const unitContainer = new Container();
    unitContainer.x = x;
    unitContainer.y = y;
    this.stage.addChild(unitContainer);

    // Body: sprite or shield placeholder
    let body: Sprite | Graphics;
    const tex = spriteTextures.get(u.id);

    if (tex) {
      const sprite = new Sprite(tex);
      sprite.width = UNIT_SIZE;
      sprite.height = UNIT_SIZE;
      sprite.anchor.set(0, 0);
      body = sprite;
    } else {
      // Colored shield placeholder
      const color = TROOP_COLORS[index % TROOP_COLORS.length];
      const shield = new Graphics();
      shield.roundRect(0, 0, UNIT_SIZE, UNIT_SIZE * 0.85, 6);
      shield.fill(color);
      // Triangle bottom
      shield.moveTo(0, UNIT_SIZE * 0.7);
      shield.lineTo(UNIT_SIZE / 2, UNIT_SIZE);
      shield.lineTo(UNIT_SIZE, UNIT_SIZE * 0.7);
      shield.fill(color);

      body = shield;
    }
    unitContainer.addChild(body);

    // If placeholder, add letter to unitContainer (Graphics can't have children in PixiJS 8)
    if (!tex) {
      const letter = new Text({
        text: u.name.charAt(0).toUpperCase(),
        style: { fontFamily: 'Arial', fontSize: 16, fontWeight: 'bold', fill: 0xffffff },
      });
      letter.anchor.set(0.5);
      letter.x = UNIT_SIZE / 2;
      letter.y = UNIT_SIZE / 2 - 2;
      unitContainer.addChild(letter);
    }

    // HP bar
    const barY = UNIT_SIZE + 3;
    const barX = (UNIT_SIZE - HP_BAR_W) / 2;

    const hpBg = new Graphics();
    hpBg.roundRect(barX, barY, HP_BAR_W, HP_BAR_H, HP_BAR_H / 2);
    hpBg.fill({ color: 0x1a1a1a, alpha: 0.85 });
    hpBg.stroke({ color: 0x333333, width: 0.5 });
    unitContainer.addChild(hpBg);

    const hpFill = new Graphics();
    hpFill.roundRect(0, 0, HP_BAR_W, HP_BAR_H, HP_BAR_H / 2);
    hpFill.fill(0x4caf50);
    hpFill.x = barX;
    hpFill.y = barY;
    unitContainer.addChild(hpFill);

    // Speed bar (below HP)
    const spBarY = barY + HP_BAR_H + 2;
    const speedBg = new Graphics();
    speedBg.roundRect(barX, spBarY, HP_BAR_W, SPEED_BAR_H, SPEED_BAR_H / 2);
    speedBg.fill({ color: 0x1a1a1a, alpha: 0.7 });
    unitContainer.addChild(speedBg);

    const speedFill = new Graphics();
    speedFill.roundRect(0, 0, HP_BAR_W, SPEED_BAR_H, SPEED_BAR_H / 2);
    speedFill.fill(0x00bcd4);
    speedFill.x = barX;
    speedFill.y = spBarY;
    speedFill.width = 0;
    unitContainer.addChild(speedFill);

    // Name label
    const nameLabel = new Text({
      text: u.name.length > 8 ? u.name.slice(0, 7) + '…' : u.name,
      style: {
        fontFamily: 'Arial',
        fontSize: 8,
        fill: side === 'attacker' ? 0x8cb4e0 : 0xe0a08c,
      },
    });
    nameLabel.anchor.set(0.5, 0);
    nameLabel.x = UNIT_SIZE / 2;
    nameLabel.y = spBarY + SPEED_BAR_H + 1;
    unitContainer.addChild(nameLabel);

    const data: UnitData = {
      id: u.id,
      side,
      name: u.name,
      maxHp: u.max_hp,
      hp: u.max_hp,
      dead: false,
      attackInterval: u.attack_interval || 5,
      spriteKey: u.sprite_key,
      baseX: x,
      baseY: y,
      idlePhase: Math.random() * Math.PI * 2,
      container: unitContainer,
      body,
      hpBg,
      hpFill,
      speedBg,
      speedFill,
      nameLabel,
      isAttacking: false,
      flashTimeout: null,
      lastAttackTick: 0,
      nextAttackTick: u.attack_interval || 5,
    };

    this.units.set(u.id, data);
  }

  /**
   * Pre-scan events to set nextAttackTick for speed bar interpolation.
   */
  private precomputeSpeedBarTicks(): void {
    // For each unit, build a sorted list of ticks where it attacks
    const attackTicks = new Map<number, number[]>();
    for (const e of this.replay.events) {
      if (e.type === 'attack') {
        let arr = attackTicks.get(e.source_id);
        if (!arr) {
          arr = [];
          attackTicks.set(e.source_id, arr);
        }
        arr.push(e.tick);
      }
    }
    // Store on units — we'll use it during playback
    this.units.forEach((u) => {
      const ticks = attackTicks.get(u.id) || [];
      // Stash for lookup later
      (u as UnitData & { _attackTicks?: number[] })._attackTicks = ticks;
    });
  }

  // ── Game loop ────────────────────────────────────────────────────────────

  private gameLoop(now: number): void {
    if (!this.playing) return;

    const delta = now - this.lastFrameTime;
    this.lastFrameTime = now;

    // Advance idle time (always runs at real-time speed for smooth bob)
    this.idleTime += delta / 1000;

    // Advance tick
    const tickRateMs = this.replay.tick_rate_ms || 100;
    const tickAdvance = (delta * this.battleSpeed) / tickRateMs;
    const prevTick = Math.floor(this.currentTick);
    this.currentTick += tickAdvance;
    const newTick = Math.floor(this.currentTick);

    // Process events for any ticks we passed
    for (let t = prevTick + 1; t <= newTick; t++) {
      this.processEventsAtTick(t);
    }

    // Update idle animation
    this.updateIdle();

    // Update speed bars
    this.updateSpeedBars(Math.floor(this.currentTick));

    // Check end
    if (this.currentTick >= this.replay.total_ticks) {
      this.currentTick = this.replay.total_ticks;
      this.playing = false;
      this.onComplete?.();
    }

    if (this.playing) {
      this.gameLoopId = requestAnimationFrame((t) => this.gameLoop(t));
      this._rafIds.push(this.gameLoopId);
    }
  }

  private processEventsAtTick(tick: number): void {
    const events = this.eventsByTick.get(tick);
    if (!events) return;

    for (const e of events) {
      const source = this.units.get(e.source_id);
      const target = this.units.get(e.target_id);
      if (!source || !target) continue;

      // Update HP
      target.hp = e.target_hp_after;

      // Update speed bar tracking for source
      source.lastAttackTick = tick;
      const nextTicks = (source as UnitData & { _attackTicks?: number[] })._attackTicks;
      if (nextTicks) {
        const idx = nextTicks.indexOf(tick);
        const nextVal = idx >= 0 && idx < nextTicks.length - 1
          ? nextTicks[idx + 1]
          : tick + source.attackInterval;
        source.nextAttackTick = nextVal ?? tick + source.attackInterval;
      }

      // Lunge attack animation
      this.animateAttack(source, target, e.is_crit);

      // Hit flash on target
      this.flashUnit(target, e.is_crit);

      // Smooth HP bar
      this.animateHpBar(target);

      // Damage text
      if (e.damage > 0) {
        showDamageText(
          this,
          target.container.x + UNIT_SIZE / 2,
          target.container.y - 5,
          e.damage,
          e.is_crit,
        );
      }

      // Kill
      if (e.is_kill) {
        target.dead = true;
        showKillText(
          this,
          target.container.x + UNIT_SIZE / 2,
          target.container.y - 20,
        );
        this.animateDeath(target);
      }
    }
  }

  // ── Attack lunge ─────────────────────────────────────────────────────────

  private animateAttack(source: UnitData, target: UnitData, isCrit: boolean): void {
    if (source.isAttacking || source.dead) return;
    source.isAttacking = true;

    const lungeDistance = isCrit ? 80 : 60;
    const lungeDuration = isCrit ? 120 : 150;

    // Direction: toward opponent
    const dx = target.baseX - source.baseX;
    const direction = dx > 0 ? 1 : -1;
    const lungeX = source.baseX + direction * lungeDistance;

    // Lunge toward target
    animateTo(this, source.container, { x: lungeX }, lungeDuration, () => {
      // Return to home
      animateTo(this, source.container, { x: source.baseX }, 240, () => {
        source.isAttacking = false;
      });
    });
  }

  // ── Hit flash ────────────────────────────────────────────────────────────

  private flashUnit(unit: UnitData, isCrit: boolean): void {
    // Cancel pending flash timeout
    if (unit.flashTimeout) {
      clearTimeout(unit.flashTimeout);
      unit.flashTimeout = null;
    }

    if (unit.body instanceof Sprite) {
      unit.body.tint = isCrit ? 0xcc0000 : 0xff0000;
    }

    const flashDuration = isCrit ? 180 : 100;
    const scaled = Math.max(50, flashDuration / this.battleSpeed);

    unit.flashTimeout = setTimeout(() => {
      if (unit.body instanceof Sprite) {
        unit.body.tint = 0xffffff;
      }
      unit.flashTimeout = null;
    }, scaled);
    this._timeoutIds.push(unit.flashTimeout);
  }

  // ── HP bar animation ─────────────────────────────────────────────────────

  private animateHpBar(unit: UnitData): void {
    const ratio = Math.max(0, Math.min(1, unit.hp / unit.maxHp));
    const color = ratio > 0.5 ? 0x4caf50 : ratio > 0.25 ? 0xff9800 : 0xf44336;

    // Redraw with new color
    unit.hpFill.clear();
    unit.hpFill.roundRect(0, 0, HP_BAR_W, HP_BAR_H, HP_BAR_H / 2);
    unit.hpFill.fill(color);

    // Animate width
    animateTo(this, unit.hpFill, { width: HP_BAR_W * ratio }, 200);
  }

  // ── Speed bar ────────────────────────────────────────────────────────────

  private updateSpeedBars(currentTick: number): void {
    this.units.forEach((u) => {
      if (u.dead) {
        u.speedFill.width = 0;
        return;
      }

      const interval = u.nextAttackTick - u.lastAttackTick;
      if (interval <= 0) {
        u.speedFill.width = HP_BAR_W;
        return;
      }

      const elapsed = currentTick - u.lastAttackTick;
      const ratio = Math.max(0, Math.min(1, elapsed / interval));
      u.speedFill.width = HP_BAR_W * ratio;
    });
  }

  // ── Death animation ──────────────────────────────────────────────────────

  private animateDeath(unit: UnitData): void {
    // Fade out + float up
    animateTo(this, unit.container, { alpha: 0, y: unit.baseY - 40 }, 500);

    // Fade bars slightly faster
    fadeOut(this, unit.hpBg, 300);
    fadeOut(this, unit.hpFill, 300);
    fadeOut(this, unit.speedBg, 300);
    fadeOut(this, unit.speedFill, 300);
    fadeOut(this, unit.nameLabel, 300);
  }

  // ── Idle bob ─────────────────────────────────────────────────────────────

  private updateIdle(): void {
    const t = this.idleTime;
    this.units.forEach((u) => {
      if (u.dead || u.isAttacking) return;
      u.container.y = u.baseY + Math.sin(t * 1.2 + u.idlePhase) * 3;
      u.container.x = u.baseX + Math.sin(t * 0.8 + u.idlePhase + 1) * 2;
    });
  }

  // ── Seek / scrub ─────────────────────────────────────────────────────────

  private rebuildStateAtTick(tick: number): void {
    // Reset all units to initial state
    this.units.forEach((u) => {
      u.hp = u.maxHp;
      u.dead = false;
      u.isAttacking = false;
      u.lastAttackTick = 0;
      u.nextAttackTick = u.attackInterval;
      u.container.alpha = 1;
      u.container.x = u.baseX;
      u.container.y = u.baseY;
      if (u.body instanceof Sprite) u.body.tint = 0xffffff;

      // Reset HP bar
      u.hpFill.clear();
      u.hpFill.roundRect(0, 0, HP_BAR_W, HP_BAR_H, HP_BAR_H / 2);
      u.hpFill.fill(0x4caf50);
      u.hpFill.width = HP_BAR_W;
      u.hpBg.alpha = 1;
      u.speedBg.alpha = 1;
      u.nameLabel.alpha = 1;
      u.speedFill.width = 0;

      // Cancel flash
      if (u.flashTimeout) {
        clearTimeout(u.flashTimeout);
        u.flashTimeout = null;
      }
    });

    // Replay events up to tick (instant, no animation)
    for (const e of this.replay.events) {
      if (e.tick > tick) break;

      const target = this.units.get(e.target_id);
      if (target) {
        target.hp = e.target_hp_after;
        if (e.is_kill) target.dead = true;
      }

      const source = this.units.get(e.source_id);
      if (source && e.type === 'attack') {
        source.lastAttackTick = e.tick;
        const nextTicks = (source as UnitData & { _attackTicks?: number[] })._attackTicks;
        if (nextTicks) {
          const idx = nextTicks.indexOf(e.tick);
          const nextVal = idx >= 0 && idx < nextTicks.length - 1
            ? nextTicks[idx + 1]
            : e.tick + source.attackInterval;
          source.nextAttackTick = nextVal ?? e.tick + source.attackInterval;
        }
      }
    }

    // Apply visual state (instant)
    this.units.forEach((u) => {
      if (u.dead) {
        u.container.alpha = 0;
      }

      const ratio = Math.max(0, Math.min(1, u.hp / u.maxHp));
      const color = ratio > 0.5 ? 0x4caf50 : ratio > 0.25 ? 0xff9800 : 0xf44336;
      u.hpFill.clear();
      u.hpFill.roundRect(0, 0, HP_BAR_W, HP_BAR_H, HP_BAR_H / 2);
      u.hpFill.fill(color);
      u.hpFill.width = HP_BAR_W * ratio;
    });

    this.updateSpeedBars(tick);
  }
}
