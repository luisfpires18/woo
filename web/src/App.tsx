import { AppFooter } from './AppFooter'
import { signedInFixture } from './devFixtures'
import { Landing } from './landing/Landing'
import { Navbar } from './landing/Navbar'
import { NotAvailable } from './landing/NotAvailable'

// No client router yet: "/" is the landing page; every other path (including the future /login, /settings and
// /profile pages) shows an honest "not available yet" page. Links are ordinary page navigations.
export function App() {
  const isLanding = window.location.pathname === '/'
  // No real sessions exist yet; the signed-in navbar appears only through the dev-server fixture.
  const account = signedInFixture()

  return (
    <>
      <a className="skip-link" href="#main">
        Skip to content
      </a>
      <Navbar account={account} />
      <main id="main" tabIndex={-1}>
        {isLanding ? <Landing /> : <NotAvailable />}
      </main>
      <AppFooter />
    </>
  )
}
