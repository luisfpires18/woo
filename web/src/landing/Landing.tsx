import { useId } from 'react'
import { kingdoms } from '../theme'
import { Button, KingdomBadge } from '../ui'
import { Banner } from './Banner'
import { CrownIcon, SettlementIcon, SwordsIcon } from './icons'
import { previewMetrics, previewWorlds, type PreviewWorld } from './previewData'

const features = [
  { label: 'Settlement building', Icon: SettlementIcon },
  { label: 'Weapons & armour', Icon: SwordsIcon },
  { label: 'Kingdom warfare', Icon: CrownIcon },
] as const

export function Landing() {
  return (
    <>
      <Banner />
      <div className="landing">
        <section className="landing__intro" aria-labelledby="intro-title">
          <h2 id="intro-title" className="landing__headline">
            Forge your kingdom’s future
          </h2>
          <p className="landing__tagline">Build your settlement. Equip your troops. Defend your kingdom.</p>
          <ul className="features" aria-label="Game features">
            {features.map(({ label, Icon }) => (
              <li key={label} className="features__item">
                <span className="features__icon">
                  <Icon />
                </span>
                {label}
              </li>
            ))}
          </ul>
        </section>

        <section className="metrics" aria-labelledby="metrics-title">
          <h2 id="metrics-title" className="visually-hidden">
            World activity
          </h2>
          <dl className="metrics__list">
            {previewMetrics.map(({ label, value }) => (
              <div key={label} className="metrics__item">
                <dt>{label}</dt>
                <dd>
                  <span aria-hidden="true">{value}</span>
                  <span className="visually-hidden">Not available yet</span>
                </dd>
              </div>
            ))}
          </dl>
          <p className="preview-tag">Preview · live figures appear when worlds open</p>
        </section>

        <section className="worlds" aria-labelledby="worlds-title">
          <div className="worlds__header">
            <h2 id="worlds-title">Choose a world</h2>
            <p className="worlds__lead">Each world is a separate game.</p>
          </div>
          <p className="worlds__note" id="worlds-note">
            <span className="preview-tag">Preview</span> World selection opens once worlds are available. These cards
            show the layout only.
          </p>
          <ul className="worlds__list">
            {previewWorlds.map((world) => (
              <li key={world.id}>
                <WorldCard world={world} />
              </li>
            ))}
          </ul>
        </section>
      </div>
    </>
  )
}

function WorldCard({ world }: { world: PreviewWorld }) {
  const titleId = useId()
  return (
    <article className="world-card" aria-labelledby={titleId}>
      <div className="world-card__main">
        <div className="world-card__heading">
          <h3 id={titleId}>{world.name}</h3>
          <span className="world-card__status">Preview</span>
        </div>
        <p className="world-card__meta">
          {world.era} · {world.pace}
        </p>
        <p className="world-card__meta">
          <span aria-hidden="true">—</span> players · <span aria-hidden="true">—</span> online
          <span className="visually-hidden">Player counts not available yet</span>
        </p>
      </div>
      <div className="world-card__kingdoms">
        <h4 className="world-card__kingdoms-title">Kingdoms</h4>
        <ul>
          {kingdoms
            .filter((k) => world.kingdoms.includes(k.id))
            .map((k) => (
              <li key={k.id}>
                <KingdomBadge kingdom={k.id} />
              </li>
            ))}
        </ul>
      </div>
      <Button variant="primary" className="world-card__action" disabled aria-describedby="worlds-note">
        Select
      </Button>
    </article>
  )
}
