// Entry title screen, ported from the discontinued weapons-of-order app (issue #65; provenance in
// docs/tasks/011-legacy-title-screen.md). Only the presentation was carried over: no session, world guard or settings.
import forgeArt from '../assets/title/forge-16x9.png'

/** Path of the main page (the task 010 landing) that ENTER WORLD opens. */
export const MAIN_PATH = '/worlds'

export function TitleScreen() {
  return (
    <main className="title-screen">
      <img className="title-screen__art" src={forgeArt} alt="" />
      <div className="title-screen__scrim" aria-hidden="true" />
      <div className="title-screen__content">
        <h1 className="wordmark">
          <span className="visually-hidden">Weapons of Order</span>
          <span className="wordmark__line" aria-hidden="true">
            WEAPONS
          </span>
          <span className="wordmark__line wordmark__line--of" aria-hidden="true">
            <span className="wordmark__rule" />
            <span className="wordmark__of">OF</span>
            <span className="wordmark__rule" />
          </span>
          <span className="wordmark__line" aria-hidden="true">
            ORDER
          </span>
        </h1>
        <nav aria-label="Main menu">
          <ul className="title-menu">
            <li>
              <a className="title-menu__item" href={MAIN_PATH}>
                ENTER WORLD
              </a>
            </li>
          </ul>
        </nav>
      </div>
    </main>
  )
}
