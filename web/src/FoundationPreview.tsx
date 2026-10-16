// Development preview of the task 009 design foundation on the skeleton page. It shows token and primitive states
// only; it is not a game screen and is replaced when real layouts arrive.
import { useState, type FormEvent } from 'react'
import { kingdoms } from './theme'
import { Button, Cluster, KingdomBadge, Panel, StatusMessage, TextField } from './ui'

export function FoundationPreview() {
  const [sample, setSample] = useState('')
  const [error, setError] = useState<string>()
  const [pressed, setPressed] = useState(false)

  const submit = (event: FormEvent) => {
    event.preventDefault()
    setError(sample.trim() ? undefined : 'Enter some text to see the valid state.')
  }

  return (
    <>
      <Panel title="Design foundation preview">
        <p>Shared controls and states for upcoming screens. Appearance follows your system setting.</p>
        <Cluster>
          <Button variant="primary">Primary action</Button>
          <Button>Secondary action</Button>
          <Button aria-pressed={pressed} onClick={() => setPressed((p) => !p)}>
            Toggle
          </Button>
          <Button variant="primary" disabled>
            Disabled
          </Button>
        </Cluster>
        <form className="ui-form" onSubmit={submit} noValidate>
          <TextField
            label="Sample text"
            hint="Submit empty to see the error state."
            required
            value={sample}
            onChange={(e) => setSample(e.target.value)}
            error={error}
          />
          <Button type="submit" variant="primary">
            Check field
          </Button>
        </form>
        <TextField label="Disabled field" disabled defaultValue="Not editable" />
        <StatusMessage tone="success">Saved changes appear like this.</StatusMessage>
        <StatusMessage tone="warning">Something needs attention.</StatusMessage>
        <StatusMessage tone="danger">Something failed.</StatusMessage>
        <StatusMessage tone="info">Neutral information.</StatusMessage>
      </Panel>

      {kingdoms.map(({ id, name }) => (
        <Panel key={id} title={`${name} accent`} kingdom={id}>
          <Cluster>
            <KingdomBadge kingdom={id} />
            <Button variant="kingdom">Kingdom action</Button>
            <Button variant="primary">Primary action</Button>
          </Cluster>
        </Panel>
      ))}
      <Panel title="No kingdom (neutral accent)">
        <Cluster>
          <Button variant="kingdom">Kingdom action</Button>
          <Button variant="primary">Primary action</Button>
        </Cluster>
      </Panel>
    </>
  )
}
