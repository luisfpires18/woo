import { AppFooter } from './AppFooter'
import { signedInFixture } from './devFixtures'
import { Landing } from './landing/Landing'
import { Navbar } from './landing/Navbar'
import { NotAvailable } from './landing/NotAvailable'
import { MAIN_PATH, TitleScreen } from './title/TitleScreen'

// No client router yet: "/" is the full-screen title (issue #65), /worlds is the main landing page and every other
// path (including the future /login, /settings and /profile pages) shows an honest "not available yet" page.
// Links are ordinary page navigations.
export function App() {
  const pathname = window.location.pathname.replace(/(.)\/+$/, '$1')
  // The title screen has no navbar or footer; its only action leads to the main page.
  if (pathname === '/') return <TitleScreen />

  const isMain = pathname === MAIN_PATH
  // No real sessions exist yet; the signed-in navbar appears only through the dev-server fixture.
  const account = signedInFixture()

  return (
    <>
      <a className="skip-link" href="#main">
        Skip to content
      </a>
      <Navbar account={account} />
      <main id="main" tabIndex={-1}>
        {isMain ? <Landing /> : <NotAvailable />}
      </main>
      <AppFooter />
    </>
  )
}
