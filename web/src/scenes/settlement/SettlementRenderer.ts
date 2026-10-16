// PixiJS renderer for the settlement scene (task 013). Imperative and React-free: React creates it once, calls the
// handle methods on selection/stage changes and receives callbacks only for meaningful events (never per frame).
// Loaded with a dynamic import so PixiJS is fetched only on the Settlement screen.
import { Application, Assets, Container, Graphics, Polygon, Sprite, Text, type Texture } from 'pixi.js'
import { Viewport } from 'pixi-viewport'
import {
  LAYERS,
  OBLIQUE,
  SCENE,
  boxCorners,
  settlementScene,
  type Layer,
  type Point,
  type PrototypeShape,
  type SceneObject,
  type Visual,
} from './sceneDefinition'

export type ArtStatus = 'prototype' | 'image' | 'image-failed'

export type SceneCallbacks = {
  onSelect: (id: string) => void
  onArtStatus: (id: string, status: ArtStatus) => void
}

export type SettlementSceneHandle = {
  select: (id: string | null) => void
  setStage: (id: string, stage: number) => void
  zoomBy: (factor: number) => void
  resetCamera: () => void
  destroy: () => void
}

const MAX_SCALE = 2
const INK = 0x2a2420
const SELECT = 0xffd27a
const prefersReducedMotion = () => window.matchMedia('(prefers-reduced-motion: reduce)').matches

const add = (a: Point, b: Point): Point => ({ x: a.x + b.x, y: a.y + b.y })
const up = (p: Point, h: number): Point => ({ x: p.x, y: p.y - h })
const back = (p: Point, d: number): Point => add(p, { x: OBLIQUE.x * d, y: OBLIQUE.y * d })
const flat = (points: Point[]) => points.flatMap((p) => [p.x, p.y])

/** Fills and outlines one face. */
function face(g: Graphics, points: Point[], color: string | number, outline = 0.55) {
  g.poly(flat(points)).fill(color).stroke({ width: 2, color: INK, alpha: outline, join: 'round' })
}

function drawBuilding(g: Graphics, anchor: Point, s: Extract<PrototypeShape, { kind: 'building' }>) {
  const { fl, fr, br, bl } = boxCorners(anchor, s.width, s.depth)
  const h = s.height
  const ridgeL = up(back(fl, s.depth / 2), h + s.roof)
  const ridgeR = up(back(fr, s.depth / 2), h + s.roof)
  // Shadow falls down-right (light from the upper left).
  g.poly(flat([fl, fr, br, bl].map((p) => add(p, { x: 22, y: 12 })))).fill({ color: 0x000000, alpha: 0.22 })
  face(g, [ridgeL, ridgeR, up(br, h), up(bl, h)], s.palette.roofSide)
  face(g, [fr, br, up(br, h), up(fr, h)], s.palette.wallSide)
  face(g, [up(fr, h), up(br, h), ridgeR], s.palette.wallSide)
  face(g, [fl, fr, up(fr, h), up(fl, h)], s.palette.wall)
  // Door and windows on the front face.
  const door = { x: anchor.x - 18, y: anchor.y }
  g.rect(door.x, door.y - Math.min(56, h * 0.55), 36, Math.min(56, h * 0.55)).fill(0x3b2b1e)
  for (const dx of [-s.width * 0.32, s.width * 0.22]) g.rect(anchor.x + dx, anchor.y - h * 0.72, 26, 22).fill(0x4a3a2a)
  face(g, [up(fl, h), up(fr, h), ridgeR, ridgeL], s.palette.roof)
  if (s.chimney) {
    const base = up(back({ x: fr.x - 46, y: fr.y }, s.depth * 0.35), h + s.roof * 0.55)
    face(g, [base, add(base, { x: 30, y: 0 }), up(add(base, { x: 30, y: 0 }), 62), up(base, 62)], 0x6d6259)
    face(g, [add(base, { x: 30, y: 0 }), add(back(base, 20), { x: 30, y: 0 }), up(add(back(base, 20), { x: 30, y: 0 }), 62), up(add(base, { x: 30, y: 0 }), 62)], 0x544a43)
    g.rect(base.x + 4, base.y - 70, 22, 8).fill(0xff8a3d) // forge glow at the chimney mouth
  }
}

function drawTree(g: Graphics, anchor: Point, s: Extract<PrototypeShape, { kind: 'tree' }>) {
  g.ellipse(anchor.x + 40, anchor.y + 6, s.radius * 0.95, s.radius * 0.32).fill({ color: 0x000000, alpha: 0.22 })
  g.rect(anchor.x - 11, anchor.y - s.trunk, 22, s.trunk).fill(0x5a4330).stroke({ width: 2, color: INK, alpha: 0.5 })
  const c = up(anchor, s.trunk + s.radius * 0.55)
  g.circle(c.x + s.radius * 0.35, c.y + 10, s.radius * 0.8).fill(0x2f5a33).stroke({ width: 2, color: INK, alpha: 0.5 })
  g.circle(c.x - s.radius * 0.3, c.y + 4, s.radius * 0.78).fill(0x3c6d3d).stroke({ width: 2, color: INK, alpha: 0.5 })
  g.circle(c.x, c.y - s.radius * 0.35, s.radius * 0.8).fill(0x4f8247).stroke({ width: 2, color: INK, alpha: 0.5 })
}

/** Battlements along the top edge from a to b. */
function merlons(g: Graphics, a: Point, b: Point, color: number) {
  const len = Math.hypot(b.x - a.x, b.y - a.y)
  const n = Math.floor(len / 30)
  for (let i = 0; i < n; i++) {
    const t = (i + 0.25) / n
    const p = { x: a.x + (b.x - a.x) * t, y: a.y + (b.y - a.y) * t }
    g.rect(p.x - 6, p.y - 12, 12, 12).fill(color).stroke({ width: 1.5, color: INK, alpha: 0.5 })
  }
}

function drawWall(g: Graphics, a: Point, s: Extract<PrototypeShape, { kind: 'wall' }>) {
  const b = s.to
  // Walls running across the view are thick towards the back; walls running into the view are thick to the left.
  const across = Math.abs(b.x - a.x) > Math.abs(b.y - a.y)
  const off = across ? { x: OBLIQUE.x * s.thickness, y: OBLIQUE.y * s.thickness } : { x: -s.thickness, y: 0 }
  const top = [up(a, s.height), up(b, s.height), up(add(b, off), s.height), up(add(a, off), s.height)]
  if (across) g.poly(flat([a, b, add(b, { x: 18, y: 12 }), add(a, { x: 18, y: 12 })])).fill({ color: 0x000000, alpha: 0.18 })
  face(g, [a, b, up(b, s.height), up(a, s.height)], across ? 0xa49e92 : 0x7e786d)
  face(g, top, 0xc7c1b4)
  merlons(g, up(a, s.height), up(b, s.height), across ? 0xb4aea2 : 0x8e887d)
}

function drawBlock(g: Graphics, anchor: Point, size: number, height: number) {
  const { fl, fr, br, bl } = boxCorners(anchor, size, size * 0.8)
  face(g, [fr, br, up(br, height), up(fr, height)], 0x7e786d)
  face(g, [fl, fr, up(fr, height), up(fl, height)], 0xa49e92)
  face(g, [up(fl, height), up(fr, height), up(br, height), up(bl, height)], 0xc7c1b4)
  merlons(g, up(fl, height), up(fr, height), 0xb4aea2)
  return { fl, fr, br, bl }
}

function drawTower(g: Graphics, anchor: Point, s: Extract<PrototypeShape, { kind: 'tower' }>) {
  const { fl, fr, br, bl } = drawBlock(g, anchor, s.size, s.height)
  const apex = up(back({ x: anchor.x, y: anchor.y }, s.size * 0.4), s.height + 70)
  face(g, [up(fr, s.height + 12), up(br, s.height + 12), apex], 0x5a2a22)
  face(g, [up(fl, s.height + 12), up(fr, s.height + 12), apex], 0x7c3a2e)
  void bl
  g.rect(anchor.x - 6, anchor.y - s.height * 0.62, 12, 22).fill(0x2e2620)
}

function drawGate(g: Graphics, anchor: Point, s: Extract<PrototypeShape, { kind: 'gate' }>) {
  const { fl, fr, br, bl } = boxCorners(anchor, s.width, s.thickness)
  face(g, [fr, br, up(br, s.height), up(fr, s.height)], 0x7e786d)
  face(g, [fl, fr, up(fr, s.height), up(fl, s.height)], 0x9a9488)
  face(g, [up(fl, s.height), up(fr, s.height), up(br, s.height), up(bl, s.height)], 0xc7c1b4)
  merlons(g, up(fl, s.height), up(fr, s.height), 0xb4aea2)
  // Arched opening; the road shows through it.
  const w = s.width * 0.42
  const top = s.height * 0.68
  g.moveTo(anchor.x - w / 2, anchor.y)
    .lineTo(anchor.x - w / 2, anchor.y - top + w / 2)
    .arc(anchor.x, anchor.y - top + w / 2, w / 2, Math.PI, 0)
    .lineTo(anchor.x + w / 2, anchor.y)
    .closePath()
    .fill(0x3a3128)
    .stroke({ width: 2, color: INK, alpha: 0.6 })
}

function drawShape(g: Graphics, anchor: Point, shape: PrototypeShape) {
  switch (shape.kind) {
    case 'building':
      return drawBuilding(g, anchor, shape)
    case 'tree':
      return drawTree(g, anchor, shape)
    case 'wall':
      return drawWall(g, anchor, shape)
    case 'corner':
      return void drawBlock(g, anchor, shape.size, shape.height)
    case 'tower':
      return drawTower(g, anchor, shape)
    case 'gate':
      return drawGate(g, anchor, shape)
  }
}

/** Ground, river, roads, plots and bridge. Identical in light and dark UI appearance. */
function drawTerrain(layers: Record<Layer, Container>) {
  const t = settlementScene.terrain
  const ground = new Graphics().rect(0, 0, SCENE.width, SCENE.height).fill(0x77874d)
  // A few broad tonal patches so the ground reads as terrain, not a flat fill.
  for (const [x, y, rx, ry, c] of [
    [380, 300, 420, 180, 0x6e7e47], [1150, 380, 380, 150, 0x7f8f55], [600, 1420, 520, 130, 0x6c7b45], [1850, 1200, 260, 260, 0x6e7e47],
  ] as const)
    ground.ellipse(x, y, rx, ry).fill({ color: c, alpha: 0.8 })
  layers.terrain.addChild(ground)

  const water = new Graphics()
  const left = t.river.map((p) => ({ x: p.x - t.riverWidth / 2, y: p.y }))
  const right = t.river.map((p) => ({ x: p.x + t.riverWidth / 2, y: p.y })).reverse()
  water.poly(flat([...left, ...right].map((p, i) => add(p, { x: i < left.length ? -16 : 16, y: 0 })))).fill(0x8c7b55) // banks
  water.poly(flat([...left, ...right])).fill(0x4d7d97)
  water.poly(flat([...left.map((p) => add(p, { x: 30, y: 0 })), ...right.map((p) => add(p, { x: -30, y: 0 }))])).fill({ color: 0x6aa0bb, alpha: 0.55 })
  layers.water.addChild(water)

  const roads = new Graphics()
  for (const plot of t.plots) {
    const { fl, fr, br, bl } = boxCorners(plot.anchor, plot.width, plot.depth)
    roads.poly(flat([add(fl, { x: 0, y: 14 }), add(fr, { x: 0, y: 14 }), add(br, { x: 0, y: 14 }), add(bl, { x: 0, y: 14 })])).fill(0x9b8a68).stroke({ width: 3, color: 0x6f6148 })
  }
  for (const path of t.roads) {
    roads.moveTo(path[0]!.x, path[0]!.y)
    for (const p of path.slice(1)) roads.lineTo(p.x, p.y)
    roads.stroke({ width: t.roadWidth, color: 0xb9a37a, cap: 'round', join: 'round' })
  }
  roads.ellipse(t.square.x, t.square.y, t.square.rx, t.square.ry).fill(0xb9a37a)
  layers.roads.addChild(roads)

  // Bridge: deck and rails above the water, split so a back rail and a front rail frame the deck.
  const { from, to, width, railHeight } = t.bridge
  const off = { x: OBLIQUE.x * width, y: OBLIQUE.y * width }
  const bridge = new Graphics()
  bridge.poly(flat([add(from, { x: 0, y: 12 }), add(to, { x: 0, y: 12 }), add(to, { x: 26, y: 28 }), add(from, { x: 26, y: 28 })])).fill({ color: 0x000000, alpha: 0.25 })
  face(bridge, [up(add(from, off), railHeight), up(add(to, off), railHeight), add(to, off), add(from, off)], 0x6b4f35)
  face(bridge, [from, to, add(to, off), add(from, off)], 0x9b7a52)
  face(bridge, [from, to, add(to, { x: 0, y: 14 }), add(from, { x: 0, y: 14 })], 0x5e4530)
  for (let i = 0; i <= 8; i++) {
    const p = { x: from.x + ((to.x - from.x) * i) / 8, y: from.y + ((to.y - from.y) * i) / 8 }
    bridge.moveTo(p.x, p.y).lineTo(p.x + off.x, p.y + off.y).stroke({ width: 1.5, color: INK, alpha: 0.35 })
  }
  face(bridge, [up(from, railHeight), up(to, railHeight), up(to, railHeight - 6), up(from, railHeight - 6)], 0x7d5c3e)
  for (const p of [from, to, { x: (from.x + to.x) / 2, y: (from.y + to.y) / 2 }]) bridge.rect(p.x - 4, p.y - railHeight, 8, railHeight).fill(0x5e4530)
  layers.bridge.addChild(bridge)
}

// Scene-owned image textures are reference-counted per URL (across live scenes): an owner is taken when a load
// starts and given back when the load is abandoned or fails, when its view is replaced, or when the scene is
// destroyed. The texture is unloaded from PixiJS Assets only when the last owner gives it back, so a texture shared by
// two visual stages is never pulled from under a live view, and nothing outside this module's own loads is unloaded.
const textureOwners = new Map<string, number>()

function acquireTexture(url: string) {
  textureOwners.set(url, (textureOwners.get(url) ?? 0) + 1)
}

function releaseTexture(url: string) {
  const remaining = (textureOwners.get(url) ?? 1) - 1
  if (remaining > 0) return void textureOwners.set(url, remaining)
  textureOwners.delete(url)
  if (Assets.cache.has(url)) void Assets.unload(url)
}

/** Development-only hooks for lifecycle tests (absent from production builds): delay adoption after a load resolves. */
export const devSceneHooks = import.meta.env.DEV
  ? { afterLoad: undefined as undefined | ((url: string) => Promise<void>), log: [] as string[] }
  : undefined
if (import.meta.env.DEV) Object.assign(window, { __wooSceneTest: { hooks: devSceneHooks, Assets, textureOwners } })

/**
 * Creates the scene inside `host`. The returned handle can be destroyed at any time, including before PixiJS has
 * finished initialising (React StrictMode mounts effects twice in development).
 */
export function createSettlementScene(host: HTMLElement, callbacks: SceneCallbacks, initialStages: Record<string, number>): SettlementSceneHandle & { ready: Promise<void> } {
  let destroyed = false
  let app: Application | undefined
  let viewport: Viewport | undefined
  let observer: ResizeObserver | undefined
  let selected: string | null = null
  const stages = { ...initialStages }
  const layers = Object.fromEntries(LAYERS.map((l) => [l, new Container({ label: l })])) as Record<Layer, Container>
  const views = new Map<string, Container>()
  const labels = new Map<string, Container>()
  /** Image URL adopted by each live view (that view owns one reference). */
  const adopted = new Map<Container, string>()
  const selection = new Graphics()
  const objects = settlementScene.objects as SceneObject[]
  const byId = new Map(objects.map((o) => [o.id, o]))

  const visualOf = (o: SceneObject): Visual => o.stages[stages[o.id] ?? 1] ?? o.stages[1]!
  const minScale = () => (viewport ? Math.min(viewport.screenWidth / SCENE.width, viewport.screenHeight / SCENE.height) : 0.1)

  /** (Re)builds one object's display at its fixed anchor; used for first draw and visual-stage swaps. */
  function renderObject(o: SceneObject) {
    const visual = visualOf(o)
    const old = views.get(o.id)
    const view = new Container({ label: o.id })
    view.zIndex = o.anchor.y
    const g = new Graphics()
    drawShape(g, o.anchor, visual.shape)
    view.addChild(g)
    if (old) {
      const index = old.parent!.getChildIndex(old)
      old.parent!.addChildAt(view, index)
      old.destroy({ children: true })
    } else layers[o.layer].addChild(view)
    views.set(o.id, view)
    // Start the new load (which takes its texture reference synchronously) before giving back the old view's, so an
    // image shared by both stages is not unloaded and fetched again.
    if (visual.image) void loadImage(o, visual, view, g)
    else if (o.selectable) callbacks.onArtStatus(o.id, 'prototype')
    const oldUrl = old && adopted.get(old)
    if (oldUrl) {
      adopted.delete(old)
      releaseTexture(oldUrl)
    }
  }

  /** Image slot: swaps the prototype shape for the image when it loads; keeps the shape if it is missing or broken. */
  async function loadImage(o: SceneObject, visual: Visual, view: Container, fallback: Graphics) {
    const { src, revision } = visual.image!
    // data:/blob: sources (test fixtures) cannot take a query string; served images carry their revision.
    const url = /^(data|blob):/.test(src) ? src : `${src}?v=${revision}`
    // Abandoned: the scene was destroyed or this view was replaced (e.g. a visual-stage change) while loading.
    const abandoned = () => destroyed || view.destroyed || views.get(o.id) !== view
    acquireTexture(url)
    let texture: Texture
    try {
      texture = await Assets.load<Texture>({ src: url, parser: 'texture' })
      if (import.meta.env.DEV) await devSceneHooks?.afterLoad?.(url)
    } catch {
      releaseTexture(url)
      if (!abandoned()) callbacks.onArtStatus(o.id, 'image-failed')
      return
    }
    if (abandoned()) {
      if (import.meta.env.DEV) devSceneHooks?.log.push(`abandoned ${url.slice(-16)}`)
      return releaseTexture(url)
    }
    if (import.meta.env.DEV) devSceneHooks?.log.push(`adopted ${url.slice(-16)}`)
    adopted.set(view, url)
    const sprite = new Sprite(texture)
    sprite.anchor.set(visual.image!.anchor.x / texture.width, visual.image!.anchor.y / texture.height)
    sprite.position.set(o.anchor.x, o.anchor.y)
    sprite.scale.set(visual.image!.scale * o.scale)
    fallback.visible = false
    view.addChild(sprite)
    callbacks.onArtStatus(o.id, 'image')
  }

  function renderLabel(o: SceneObject) {
    labels.get(o.id)?.destroy({ children: true })
    const visual = visualOf(o)
    const isSelected = selected === o.id
    const label = new Container({ label: `label-${o.id}` })
    const text = new Text({
      text: o.name,
      resolution: 2,
      style: { fontFamily: 'system-ui, Segoe UI, sans-serif', fontSize: 15, fontWeight: isSelected ? '700' : '600', fill: isSelected ? 0x1f1b16 : 0xf6f2ea },
    })
    text.anchor.set(0.5)
    const pill = new Graphics()
      .roundRect(-text.width / 2 - 10, -text.height / 2 - 5, text.width + 20, text.height + 10, 12)
      .fill({ color: isSelected ? SELECT : 0x1e2024, alpha: isSelected ? 1 : 0.85 })
      .stroke({ width: 1.5, color: isSelected ? 0x1f1b16 : 0xf6f2ea, alpha: 0.7 })
    label.addChild(pill, text)
    label.position.set(visual.label.x, visual.label.y)
    label.zIndex = isSelected ? 1 : 0
    layers.overlay.addChild(label)
    labels.set(o.id, label)
  }

  /** Selection outline, drawn above every layer so occluders (tree, walls) never hide it. Constant screen width. */
  function renderSelection() {
    selection.clear()
    const o = selected ? byId.get(selected) : undefined
    if (!o || !viewport) return
    const visual = visualOf(o)
    const w = 1 / viewport.scale.x
    selection.poly(flat(visual.footprint)).fill({ color: SELECT, alpha: 0.18 })
    selection.poly(flat(visual.hit ?? visual.footprint)).stroke({ width: 7 * w, color: INK, alpha: 0.65, join: 'round' })
    selection.poly(flat(visual.hit ?? visual.footprint)).stroke({ width: 3.5 * w, color: SELECT, join: 'round' })
  }

  /** Keeps labels at a constant on-screen size and the selection stroke crisp at any zoom. */
  function onZoom() {
    if (!viewport) return
    const inverse = 1 / viewport.scale.x
    for (const label of labels.values()) label.scale.set(inverse)
    renderSelection()
  }

  /** Topmost selectable object whose hit polygon contains the world point (draw order, back to front). */
  function hitTest(p: Point): string | null {
    const order = objects
      .filter((o) => o.selectable)
      .map((o) => ({ o, rank: LAYERS.indexOf(o.layer) * 1e6 + o.anchor.y }))
      .sort((a, b) => b.rank - a.rank)
    for (const { o } of order) {
      const hit = visualOf(o).hit
      if (hit && new Polygon(flat(hit)).contains(p.x, p.y)) return o.id
    }
    return null
  }

  function frameInitial() {
    if (!viewport) return
    const f = viewport.screenWidth < 600 ? settlementScene.compactFrame : settlementScene.initialFrame
    const scale = Math.max(minScale(), Math.min(MAX_SCALE, viewport.screenWidth / f.width, viewport.screenHeight / f.height))
    viewport.setZoom(scale, true)
    viewport.moveCenter(f.x + f.width / 2, f.y + f.height / 2)
    onZoom()
  }

  function resize() {
    if (!app || !viewport) return
    const width = Math.max(1, host.clientWidth)
    const height = Math.max(1, host.clientHeight)
    const center = viewport.center
    app.renderer.resize(width, height)
    viewport.resize(width, height, SCENE.width, SCENE.height)
    viewport.clampZoom({ minScale: minScale(), maxScale: MAX_SCALE })
    viewport.setZoom(Math.max(minScale(), Math.min(MAX_SCALE, viewport.scale.x)), true)
    viewport.moveCenter(center)
    onZoom()
  }

  const ready = (async () => {
    const instance = new Application()
    await instance.init({
      width: Math.max(1, host.clientWidth),
      height: Math.max(1, host.clientHeight),
      background: 0x2b3320,
      antialias: true,
      resolution: Math.min(window.devicePixelRatio || 1, 2),
      autoDensity: true,
      preference: 'webgl',
    })
    if (destroyed) {
      instance.destroy({ removeView: true }, { children: true })
      return
    }
    app = instance
    app.canvas.setAttribute('aria-hidden', 'true')
    host.appendChild(app.canvas)

    viewport = new Viewport({
      screenWidth: app.screen.width,
      screenHeight: app.screen.height,
      worldWidth: SCENE.width,
      worldHeight: SCENE.height,
      events: app.renderer.events,
      passiveWheel: false,
      threshold: 8,
    })
    viewport.drag().pinch().wheel({ smooth: prefersReducedMotion() ? undefined : 4 }).clamp({ direction: 'all', underflow: 'center' }).clampZoom({ minScale: minScale(), maxScale: MAX_SCALE })
    if (!prefersReducedMotion()) viewport.decelerate({ friction: 0.92 })
    app.stage.addChild(viewport)
    for (const name of LAYERS) viewport.addChild(layers[name])
    layers.objects.sortableChildren = true
    layers.overlay.sortableChildren = true
    layers.overlay.addChild(selection)

    drawTerrain(layers)
    for (const o of objects) renderObject(o)
    for (const o of objects) if (o.selectable) renderLabel(o)

    // "clicked" fires only when the pointer did not move past the drag threshold, with world coordinates mapped
    // through the camera; a pan or pinch never selects.
    viewport.on('clicked', (e) => {
      const id = hitTest(e.world)
      if (id) callbacks.onSelect(id)
    })
    viewport.on('zoomed', onZoom)
    frameInitial()
    observer = new ResizeObserver(resize)
    observer.observe(host)
    if (import.meta.env.DEV) (window as unknown as { __wooSettlement?: unknown }).__wooSettlement = { app, viewport, hitTest }
  })()

  return {
    ready,
    select(id) {
      selected = id
      for (const o of objects) if (o.selectable && app) renderLabel(o)
      onZoom()
    },
    setStage(id, stage) {
      stages[id] = stage
      const o = byId.get(id)
      if (!o || !app) return
      renderObject(o)
      renderLabel(o)
      onZoom()
    },
    zoomBy(factor) {
      if (!viewport) return
      const center = viewport.center
      viewport.setZoom(Math.max(minScale(), Math.min(MAX_SCALE, viewport.scale.x * factor)), true)
      viewport.moveCenter(center)
      onZoom()
    },
    resetCamera: frameInitial,
    destroy() {
      destroyed = true
      observer?.disconnect()
      if (import.meta.env.DEV) delete (window as unknown as { __wooSettlement?: unknown }).__wooSettlement
      if (app) {
        app.destroy({ removeView: true }, { children: true })
        app = undefined
        viewport = undefined
      }
      // Views were destroyed with the stage; give back their textures. In-flight loads release themselves on arrival.
      for (const url of adopted.values()) releaseTexture(url)
      adopted.clear()
    },
  }
}
