import { Page, PageHeader } from '../ui'

/** Shown for every path except "/": Log in, Profile and Settings pages belong to later tasks. */
export function NotAvailable() {
  return (
    <Page narrow>
      <PageHeader
        title="Not available yet"
        lead="This page is part of a later development step. Nothing has been changed."
      />
      <p>
        <a className="ui-button ui-button--secondary" href="/">
          Back to the start page
        </a>
      </p>
    </Page>
  )
}
