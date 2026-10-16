import { useEffect, useId, useRef, useState, type ReactNode, type Ref } from 'react'
import { sceneFailureFixture } from '../devFixtures'
import { selectableObjects } from '../scenes/settlement/sceneDefinition'
import type { ArtStatus, SettlementSceneHandle } from '../scenes/settlement/SettlementRenderer'
import { Button, Panel } from '../ui'
import { GameView } from './GameShell'
import { previewBuildings } from './previewData'

type SceneState = 'loading' | 'ready' | 'failed'

const artText: Record<ArtStatus, string> = {
  prototype: 'Prototype geometry (no image slot yet)',
  image: 'Image from its asset slot',
  'image-failed': 'Image slot set but the image is missing; prototype geometry shown',
}

/**
 * Settlement screen (task 013 VILLAGE-UI): PixiJS scene for the picture, React for everything you act on. All state
 * here is local mock state; nothing is saved and no upgrade is ordered.
 */
export function SettlementView({ headingRef, tags }: { headingRef: Ref<HTMLHeadingElement>; tags: ReactNode }) {
  const hostRef = useRef<HTMLDivElement>(null)
  const sceneRef = useRef<SettlementSceneHandle | null>(null)
  const [sceneState, setSceneState] = useState<SceneState>('loading')
  const [selected, setSelected] = useState<string | null>(null)
  const [stages, setStages] = useState<Record<string, number>>({ forge: 1 })
  const [art, setArt] = useState<Record<string, ArtStatus>>({})
  const stagesRef = useRef(stages)
  const descriptionId = useId()

  // Create the scene once per mount; destroy it (canvas, listeners, textures) on unmount.
  useEffect(() => {
    const host = hostRef.current!
    let handle: (SettlementSceneHandle & { ready: Promise<void> }) | undefined
    let cancelled = false
    const callbacks = {
      onSelect: (id: string) => setSelected(id),
      onArtStatus: (id: string, status: ArtStatus) => setArt((current) => (current[id] === status ? current : { ...current, [id]: status })),
    }
    import('../scenes/settlement/SettlementRenderer')
      .then(({ createSettlementScene }) => {
        if (cancelled) return
        if (sceneFailureFixture()) throw new Error('Simulated renderer failure (dev fixture)')
        handle = createSettlementScene(host, callbacks, stagesRef.current)
        sceneRef.current = handle
        return handle.ready
      })
      .then(() => {
        if (!cancelled) setSceneState('ready')
      })
      .catch(() => {
        if (!cancelled) setSceneState('failed')
      })
    return () => {
      cancelled = true
      handle?.destroy()
      sceneRef.current = null
    }
  }, [])

  useEffect(() => {
    if (sceneState === 'ready') sceneRef.current?.select(selected)
  }, [selected, sceneState])

  const changeStage = (id: string, stage: number) => {
    const next = { ...stages, [id]: stage }
    stagesRef.current = next
    setStages(next)
    sceneRef.current?.setStage(id, stage)
  }

  const current = selectableObjects.find((o) => o.id === selected)

  return (
    <GameView
      title="Settlement"
      lead="Select a building in the scene or in the list. This is a layout and interaction prototype."
      headingRef={headingRef}
      tags={tags}
      stageClassName="game-stage--scene"
      stage={
        <div className="settlement-stage">
          <div
            ref={hostRef}
            className="settlement-stage__canvas"
            role="img"
            aria-label="Settlement scene, prototype geometry"
            aria-describedby={descriptionId}
          />
          <p id={descriptionId} className="visually-hidden">
            Placeholder village with a forge, a cottage and a hall inside a wall with a gate and tower, a tree, and a river
            crossed by a bridge. Drag to pan and use the zoom buttons, wheel or pinch to zoom. Every building can also be
            selected from the building list.
          </p>
          {sceneState !== 'failed' && (
            <p className="settlement-stage__note">
              <span className="preview-tag">Prototype geometry · not final art</span>
            </p>
          )}
          {sceneState === 'loading' && <p className="settlement-stage__status">Preparing the scene…</p>}
          {sceneState === 'failed' && (
            <p className="settlement-stage__status settlement-stage__status--failed" role="alert">
              The scene could not be drawn on this device. The building list and inspector still work.
            </p>
          )}
          {sceneState === 'ready' && (
            <div className="settlement-stage__controls" role="group" aria-label="Camera">
              <button type="button" className="scene-control" aria-label="Zoom in" onClick={() => sceneRef.current?.zoomBy(1.25)}>
                +
              </button>
              <button type="button" className="scene-control" aria-label="Zoom out" onClick={() => sceneRef.current?.zoomBy(0.8)}>
                −
              </button>
              <button type="button" className="scene-control scene-control--wide" onClick={() => sceneRef.current?.resetCamera()}>
                Reset view
              </button>
            </div>
          )}
        </div>
      }
      aside={
        <>
          <Panel title="Buildings">
            <ul className="building-list">
              {selectableObjects.map((o) => (
                <li key={o.id}>
                  <button
                    type="button"
                    className="building-list__item"
                    aria-pressed={selected === o.id}
                    onClick={() => setSelected(selected === o.id ? null : o.id)}
                  >
                    <span className="building-list__name">{o.name}</span>
                    <span className="building-list__meta">
                      Mock level {previewBuildings[o.id]?.mockLevel[stages[o.id] ?? 1] ?? '—'}
                    </span>
                  </button>
                </li>
              ))}
            </ul>
            <p role="status" className="visually-hidden">
              {current ? `${current.name} selected` : 'No building selected'}
            </p>
          </Panel>
          <Inspector
            id={current?.id}
            name={current?.name}
            stage={current ? (stages[current.id] ?? 1) : 1}
            stageCount={current ? Object.keys(current.stages).length : 0}
            art={current ? (art[current.id] ?? 'prototype') : 'prototype'}
            sceneReady={sceneState === 'ready'}
            onStage={(stage) => current && changeStage(current.id, stage)}
          />
        </>
      }
    />
  )
}

function Inspector(props: {
  id?: string
  name?: string
  stage: number
  stageCount: number
  art: ArtStatus
  sceneReady: boolean
  onStage: (stage: number) => void
}) {
  if (!props.id) {
    return (
      <Panel title="Inspector">
        <p className="inspector__empty">Select a building in the scene or in the list above.</p>
      </Panel>
    )
  }
  const info = previewBuildings[props.id]
  return (
    <Panel title={props.name}>
      <p className="inspector__mock">
        <span className="preview-tag">Mock values · nothing is saved</span>
      </p>
      {info && <p>{info.summary}</p>}
      <dl className="inspector__facts">
        <div>
          <dt>Mock level</dt>
          <dd>{info?.mockLevel[props.stage] ?? '—'}</dd>
        </div>
        <div>
          <dt>Visual stage</dt>
          <dd>
            {props.stage} of {props.stageCount}
          </dd>
        </div>
        <div>
          <dt>Art</dt>
          <dd>{props.sceneReady ? artText[props.art] : 'Scene not drawn'}</dd>
        </div>
      </dl>
      {info?.illustrativeUpgrade && (
        <section className="inspector__upgrade" aria-labelledby={`${props.id}-upgrade`}>
          <h3 id={`${props.id}-upgrade`}>Illustrative upgrade (mock)</h3>
          <ul>
            {info.illustrativeUpgrade.cost.map((c) => (
              <li key={c.label}>
                {c.label}: {c.amount}
              </li>
            ))}
            <li>Time: {info.illustrativeUpgrade.duration}</li>
          </ul>
          <p className="inspector__hint">Example figures only; resources, costs and timings are not decided.</p>
        </section>
      )}
      {props.stageCount > 1 && (
        <fieldset className="inspector__stages">
          <legend>Prototype visual stage (local only, no upgrade is ordered)</legend>
          <div className="ui-cluster">
            {Object.keys(previewBuildings[props.id]?.mockLevel ?? { 1: 1 }).map((s) => (
              <Button key={s} aria-pressed={props.stage === Number(s)} onClick={() => props.onStage(Number(s))}>
                Stage {s}
              </Button>
            ))}
          </div>
        </fieldset>
      )}
    </Panel>
  )
}
