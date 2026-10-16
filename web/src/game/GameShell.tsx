import { useEffect, useId, useRef, useState, type ReactNode, type Ref } from 'react'
import { AppFooter } from '../AppFooter'
import { kingdomFixture, signedInFixture } from '../devFixtures'
import { Navbar } from '../landing/Navbar'
import { followLink } from '../routing'
import { kingdoms } from '../theme'
import { Panel } from '../ui'
import { previewResourceSlots, type ResourceSlot } from './previewData'
import { SettlementView } from './SettlementView'
import { findGameView, gameViews, defaultGamePath, type GameView as View } from './views'

/**
 * Game application shell (task 012): charcoal command bar with the resource header, a navigation rail (bottom tab
 * bar on narrow screens), the page content and the shared footer. Navigation stays inside the app without reloads;
 * every /game/... path also loads directly. No session, membership or settlement exists yet, so nothing is saved.
 */
export function GameShell({ pathname }: { pathname: string }) {
  const view = findGameView(pathname)
  // Dev-server visual fixtures only (null in production); read once so they survive in-app navigation.
  const [account] = useState(signedInFixture)
  const [kingdom] = useState(kingdomFixture)
  const headingRef = useRef<HTMLHeadingElement>(null)
  // The view the heading last showed; focus moves only when it changes (not on load, not on StrictMode re-runs).
  const shownView = useRef(view)

  // Name the page, and after in-app navigation move focus to the new heading so screen readers announce it.
  useEffect(() => {
    document.title = `${view?.label ?? 'Not available yet'} · Weapons of Order`
    if (shownView.current === view) return
    shownView.current = view
    headingRef.current?.focus()
  }, [view])

  const tags = (
    <>
      <span className="preview-tag">Layout prototype · nothing is saved</span>
      {kingdom && <span className="preview-tag">Accent preview: {kingdoms.find((k) => k.id === kingdom)?.name}</span>}
    </>
  )

  return (
    <div className="game" data-kingdom={kingdom ?? undefined}>
      <a className="skip-link" href="#main">
        Skip to content
      </a>
      <Navbar account={account} wide>
        <ResourceBar slots={previewResourceSlots} preview />
      </Navbar>
      <GameNav current={view} />
      <div className="game__body">
        <main id="main" tabIndex={-1} className="game__main">
          {view?.id === 'settlement' ? (
            <SettlementView headingRef={headingRef} tags={tags} />
          ) : view ? (
            <GameView
              title={view.label}
              lead={view.lead}
              headingRef={headingRef}
              tags={tags}
              stage={<StagePlaceholder view={view} />}
              aside={<PlannedPanel view={view} />}
            />
          ) : (
            <GameView
              title="Not available yet"
              lead="This game screen does not exist. Nothing has been changed."
              headingRef={headingRef}
              stage={
                <a className="ui-button ui-button--secondary" href={defaultGamePath} onClick={followLink}>
                  Go to Settlement
                </a>
              }
            />
          )}
        </main>
        <AppFooter />
      </div>
    </div>
  )
}

/** Main game navigation. The current view has aria-current="page", a raised plate, a marker bar and bold text. */
function GameNav({ current }: { current: View | undefined }) {
  return (
    <nav className="game-nav" aria-label="Game">
      <ul className="game-nav__list">
        {gameViews.map((view) => (
          <li key={view.id}>
            <a
              className="game-nav__link"
              href={view.path}
              aria-label={view.label}
              aria-current={view === current ? 'page' : undefined}
              onClick={followLink}
            >
              <span className="game-nav__icon">
                <view.Icon />
              </span>
              {/* One label is displayed at a time (the short one in the bottom bar); the full label is always the
                  accessible name, and each short label is contained in it. */}
              <span className="game-nav__label">{view.label}</span>
              <span className="game-nav__short">{view.shortLabel}</span>
            </a>
          </li>
        ))}
      </ul>
    </nav>
  )
}

const formatAmount = (value: number) => value.toLocaleString()

/** Resource header. Renders whatever slots it is given; unknown amounts show a dash and read "Not available yet". */
export function ResourceBar({ slots, preview = false }: { slots: readonly ResourceSlot[]; preview?: boolean }) {
  return (
    <section className="resources" aria-label={preview ? 'Resources (preview)' : 'Resources'}>
      <dl className="resources__list">
        {slots.map((slot) => (
          <div key={slot.id} className="resources__slot">
            <dt>
              <span className="resources__icon" aria-hidden="true">
                {slot.iconUrl && <img src={slot.iconUrl} alt="" />}
              </span>
              {slot.label}
            </dt>
            <dd>
              {slot.amount === null ? (
                <>
                  <span aria-hidden="true">—</span>
                  <span className="visually-hidden">Not available yet</span>
                </>
              ) : (
                <>
                  {formatAmount(slot.amount)}
                  {slot.capacity != null && <span className="resources__capacity"> / {formatAmount(slot.capacity)}</span>}
                </>
              )}
            </dd>
          </div>
        ))}
      </dl>
      {preview && <p className="preview-tag">Preview · resource types not decided</p>}
    </section>
  )
}

/**
 * Reusable game page structure: heading row, a framed stage for the screen's main scene (later the PixiJS settlement
 * or map) and an optional side column for inspectors and lists. The side column stacks under the stage when narrow.
 */
export function GameView(props: {
  title: string
  lead: string
  headingRef?: Ref<HTMLHeadingElement>
  tags?: ReactNode
  stage: ReactNode
  /** Extra class for the stage frame (e.g. a full-bleed scene). */
  stageClassName?: string
  aside?: ReactNode
}) {
  const id = useId()
  return (
    <div className="game-view">
      <header className="game-view__header">
        <div>
          <h1 ref={props.headingRef} tabIndex={-1}>
            {props.title}
          </h1>
          <p className="game-view__lead">{props.lead}</p>
        </div>
        {props.tags && <p className="game-view__tags">{props.tags}</p>}
      </header>
      <div className={props.aside ? 'game-view__body game-view__body--aside' : 'game-view__body'}>
        <section className={props.stageClassName ? `game-stage ${props.stageClassName}` : 'game-stage'} aria-labelledby={`${id}-stage`}>
          <h2 id={`${id}-stage`} className="visually-hidden">
            {props.title} main area
          </h2>
          {props.stage}
        </section>
        {props.aside && <div className="game-view__aside">{props.aside}</div>}
      </div>
    </div>
  )
}

function StagePlaceholder({ view }: { view: View }) {
  return (
    <div className="game-stage__placeholder">
      <span className="game-stage__icon">
        <view.Icon />
      </span>
      <p className="game-stage__title">{view.label} screen</p>
      <p>Placeholder area. The real screen is built in a later task.</p>
    </div>
  )
}

function PlannedPanel({ view }: { view: View }) {
  return (
    <Panel title="Planned for this screen">
      <ul className="planned">
        {view.planned.map(({ task, text }) => (
          <li key={task}>
            <span className="planned__task">{task}</span>
            {text}
          </li>
        ))}
      </ul>
    </Panel>
  )
}
