// Settlement scene definition (task 013 VILLAGE-UI). Plain data, independent of React and of the renderer: every
// object has a stable ID, scene coordinates, a ground anchor, footprint, hit polygon, label point, layer and an
// optional image slot. PROTOTYPE CONTENT: the geometry below is temporary authored shapes, not Arkazia artwork, and
// the names are visual samples, not a building roster.

export type Point = { x: number; y: number }

/** Canonical scene space for this proof (docs/design/village-asset-specification.md). */
export const SCENE = { width: 2048, height: 1536 } as const

/**
 * One shallow-overhead oblique camera for every object: "back" on the ground runs up and to the right by this
 * vector per unit of depth, verticals stay vertical, and light comes from the upper left (front faces mid,
 * right-hand faces dark, front roof slopes light, shadows fall down-right).
 */
export const OBLIQUE = { x: 0.35, y: -0.6 } as const

/** Draw order, back to front. "objects" is sorted by ground-anchor y; the others keep authored order. */
export const LAYERS = ['terrain', 'water', 'roads', 'bridge', 'shadows', 'objects', 'foreground', 'overlay'] as const
export type Layer = (typeof LAYERS)[number]

/** Optional artwork for a visual. Until an approved image exists (or if it fails to load) the prototype shape is drawn. */
export type ImageSlot = {
  src: string
  /** Pixel in the image that sits on the scene ground anchor (usually bottom centre of the base). */
  anchor: Point
  scale: number
  /** Cache revision; changes when the image is replaced (task 034 IMAGE-CACHE). */
  revision: number
}

/** Temporary geometry used when there is no image. */
export type PrototypeShape =
  | { kind: 'building'; width: number; depth: number; height: number; roof: number; palette: BuildingPalette; chimney?: boolean }
  | { kind: 'tree'; trunk: number; radius: number }
  | { kind: 'wall'; to: Point; height: number; thickness: number }
  | { kind: 'corner'; size: number; height: number }
  | { kind: 'tower'; size: number; height: number }
  | { kind: 'gate'; width: number; height: number; thickness: number }

export type BuildingPalette = { wall: string; wallSide: string; roof: string; roofSide: string }

export type Visual = {
  shape: PrototypeShape
  image?: ImageSlot
  /** Ground polygon (scene coordinates). */
  footprint: Point[]
  /** Clickable polygon (scene coordinates); not the image bounds. Absent = not selectable. */
  hit?: Point[]
  label: Point
}

export type SceneObject = {
  /** Stable logical ID; survives visual-stage changes and image replacement. */
  id: string
  name: string
  layer: Layer
  /** Ground anchor in scene coordinates; stays fixed when the visual stage changes. */
  anchor: Point
  scale: number
  /** Visual stages keyed by stage number; a building may have one or more. */
  stages: Record<number, Visual>
  /** Whether it appears in the building list and can be selected. */
  selectable: boolean
}

const add = (a: Point, b: Point): Point => ({ x: a.x + b.x, y: a.y + b.y })
const back = (p: Point, depth: number): Point => add(p, { x: OBLIQUE.x * depth, y: OBLIQUE.y * depth })
const up = (p: Point, h: number): Point => ({ x: p.x, y: p.y - h })

/** Corners of an oblique box whose front edge is centred on the anchor. */
export function boxCorners(anchor: Point, width: number, depth: number) {
  const fl = { x: anchor.x - width / 2, y: anchor.y }
  const fr = { x: anchor.x + width / 2, y: anchor.y }
  return { fl, fr, br: back(fr, depth), bl: back(fl, depth) }
}

/** Footprint, silhouette hit polygon and label point of a prototype building with a gable roof. */
function buildingVisual(anchor: Point, shape: Extract<PrototypeShape, { kind: 'building' }>): Visual {
  const { fl, fr, br, bl } = boxCorners(anchor, shape.width, shape.depth)
  const ridgeL = up(back(fl, shape.depth / 2), shape.height + shape.roof)
  const ridgeR = up(back(fr, shape.depth / 2), shape.height + shape.roof)
  return {
    shape,
    footprint: [fl, fr, br, bl],
    hit: [fl, fr, br, up(br, shape.height), ridgeR, ridgeL, up(fl, shape.height)],
    label: { x: (ridgeL.x + ridgeR.x) / 2, y: Math.min(ridgeL.y, ridgeR.y) - 34 },
  }
}

function building(id: string, name: string, anchor: Point, stages: Record<number, Extract<PrototypeShape, { kind: 'building' }>>, image?: ImageSlot): SceneObject {
  const visuals: Record<number, Visual> = {}
  for (const [stage, shape] of Object.entries(stages)) visuals[Number(stage)] = { ...buildingVisual(anchor, shape), image }
  return { id, name, layer: 'objects', anchor, scale: 1, stages: visuals, selectable: true }
}

function scenery(id: string, name: string, layer: Layer, anchor: Point, shape: PrototypeShape, footprint: Point[] = [anchor]): SceneObject {
  return { id, name, layer, anchor, scale: 1, selectable: false, stages: { 1: { shape, footprint, label: anchor } } }
}

const palettes = {
  forge: { wall: '#cbbd9f', wallSide: '#9f9177', roof: '#8e3b2e', roofSide: '#6b2a21' },
  cottage: { wall: '#d6c9ab', wallSide: '#a99b7e', roof: '#7a5b3a', roofSide: '#5a4129' },
  hall: { wall: '#c8bfae', wallSide: '#9a9180', roof: '#55524f', roofSide: '#3d3b39' },
} satisfies Record<string, BuildingPalette>

// ---- Fixture composition (scene coordinates) ----

const WALL_Y = 1300
const CORNER: Point = { x: 300, y: WALL_Y }
const WEST_END = back(CORNER, 1000)
const GATE_X = 840
const TOWER_X = 1320

export const settlementScene = {
  version: 1,
  /** Region framed when the scene first opens. */
  initialFrame: { x: 230, y: 420, width: 1560, height: 1000 },
  /** Tighter framing for narrow stages (phones), so buildings stay readable at the first view. */
  compactFrame: { x: 470, y: 430, width: 980, height: 960 },
  objects: [
    building('forge', 'Forge', { x: 780, y: 900 }, {
      1: { kind: 'building', width: 200, depth: 150, height: 105, roof: 70, palette: palettes.forge },
      2: { kind: 'building', width: 240, depth: 170, height: 150, roof: 85, palette: palettes.forge, chimney: true },
    }),
    // The cottage points at an image that does not exist yet: it proves the missing-art fallback.
    building('cottage', 'Cottage', { x: 1080, y: 760 }, {
      1: { kind: 'building', width: 150, depth: 120, height: 90, roof: 60, palette: palettes.cottage },
    }, { src: '/scene-assets/settlement/cottage-stage-1.png', anchor: { x: 128, y: 230 }, scale: 1, revision: 1 }),
    building('hall', 'Hall', { x: 1060, y: 1258 }, {
      1: { kind: 'building', width: 260, depth: 160, height: 120, roof: 80, palette: palettes.hall },
    }),
    scenery('oak', 'Tree', 'foreground', { x: 905, y: 1000 }, { kind: 'tree', trunk: 95, radius: 78 }),
    scenery('wall-west', 'Wall', 'objects', CORNER, { kind: 'wall', to: WEST_END, height: 70, thickness: 26 }, [CORNER, WEST_END]),
    scenery('wall-corner', 'Wall corner', 'foreground', CORNER, { kind: 'corner', size: 46, height: 82 }),
    scenery('wall-south-west', 'Wall', 'foreground', { x: CORNER.x + 23, y: WALL_Y }, { kind: 'wall', to: { x: GATE_X - 75, y: WALL_Y }, height: 70, thickness: 26 }),
    scenery('gate', 'Gate', 'foreground', { x: GATE_X, y: WALL_Y }, { kind: 'gate', width: 150, height: 112, thickness: 34 }),
    scenery('wall-south-east', 'Wall', 'foreground', { x: GATE_X + 75, y: WALL_Y }, { kind: 'wall', to: { x: TOWER_X - 36, y: WALL_Y }, height: 70, thickness: 26 }),
    scenery('tower', 'Tower', 'foreground', { x: TOWER_X, y: WALL_Y + 8 }, { kind: 'tower', size: 72, height: 165 }),
  ] satisfies SceneObject[],
  terrain: {
    river: [
      { x: 1505, y: -20 }, { x: 1470, y: 300 }, { x: 1525, y: 620 }, { x: 1560, y: 940 }, { x: 1525, y: 1240 }, { x: 1590, y: 1560 },
    ],
    riverWidth: 165,
    roads: [
      [{ x: GATE_X, y: 1420 }, { x: GATE_X, y: 1010 }, { x: 1440, y: 940 }],
      [{ x: 1680, y: 940 }, { x: 2060, y: 900 }],
      [{ x: GATE_X, y: 1010 }, { x: 780, y: 910 }],
      [{ x: GATE_X, y: 1010 }, { x: 1080, y: 770 }],
      [{ x: GATE_X, y: 1120 }, { x: 560, y: 1150 }],
    ],
    roadWidth: 54,
    square: { x: 840, y: 1030, rx: 120, ry: 62 },
    /** Building plots (foundations); the last one is empty. */
    plots: [
      { anchor: { x: 780, y: 900 }, width: 270, depth: 190 },
      { anchor: { x: 1080, y: 760 }, width: 180, depth: 140 },
      { anchor: { x: 1060, y: 1258 }, width: 290, depth: 180 },
      { anchor: { x: 560, y: 1165 }, width: 170, depth: 130 },
    ],
    bridge: { from: { x: 1440, y: 940 }, to: { x: 1680, y: 932 }, width: 62, railHeight: 16 },
  },
}

export type SettlementScene = typeof settlementScene
export const selectableObjects = settlementScene.objects.filter((o) => o.selectable)
