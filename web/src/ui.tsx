// Shared UI primitives (task 009). They read design tokens only; kingdom colours come from the nearest data-kingdom.
import { useId, type ButtonHTMLAttributes, type InputHTMLAttributes, type ReactNode } from 'react'
import { kingdoms, type Kingdom } from './theme'

const cx = (...names: (string | false | undefined)[]) => names.filter(Boolean).join(' ')

/** Page body: centred content column with responsive gutters. Use once per layout inside <main>. */
export function Page({ children, narrow = false }: { children: ReactNode; narrow?: boolean }) {
  return <div className={cx('ui-page', narrow && 'ui-page--narrow')}>{children}</div>
}

export function PageHeader({ title, lead }: { title: string; lead?: ReactNode }) {
  return (
    <header className="ui-page-header">
      <h1>{title}</h1>
      {lead && <p className="ui-lead">{lead}</p>}
    </header>
  )
}

/** Bordered surface. With a title it becomes a labelled region. */
export function Panel({ title, children, kingdom }: { title?: string; children: ReactNode; kingdom?: Kingdom }) {
  const id = useId()
  return (
    <section
      className={cx('ui-panel', kingdom && 'ui-panel--kingdom')}
      data-kingdom={kingdom}
      aria-labelledby={title ? id : undefined}
    >
      {title && (
        <h2 id={id} className="ui-panel__title">
          {title}
        </h2>
      )}
      {children}
    </section>
  )
}

/** Horizontal group that wraps on narrow screens (button rows, badge lists). */
export function Cluster({ children }: { children: ReactNode }) {
  return <div className="ui-cluster">{children}</div>
}

type ButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  /** primary: brand crimson action; secondary: neutral outline; kingdom: nearest kingdom accent. */
  variant?: 'primary' | 'secondary' | 'kingdom'
}

export function Button({ variant = 'secondary', className, type = 'button', ...props }: ButtonProps) {
  return <button type={type} className={cx('ui-button', `ui-button--${variant}`, className)} {...props} />
}

type TextFieldProps = Omit<InputHTMLAttributes<HTMLInputElement>, 'id'> & {
  label: string
  hint?: string
  /** Shown below the input; also sets aria-invalid and is announced via aria-describedby. */
  error?: string
}

export function TextField({ label, hint, error, className, ...input }: TextFieldProps) {
  const id = useId()
  const hintId = hint ? `${id}-hint` : undefined
  const errorId = error ? `${id}-error` : undefined
  return (
    <div className={cx('ui-field', error && 'ui-field--invalid', className)}>
      <label htmlFor={id} className="ui-field__label">
        {label}
        {input.required && <span className="ui-field__required"> (required)</span>}
      </label>
      {hint && (
        <p id={hintId} className="ui-field__hint">
          {hint}
        </p>
      )}
      <input
        id={id}
        className="ui-input"
        aria-invalid={error ? true : undefined}
        aria-describedby={[hintId, errorId].filter(Boolean).join(' ') || undefined}
        {...input}
      />
      {error && (
        <p id={errorId} className="ui-field__error">
          <span aria-hidden="true">✕ </span>
          {error}
        </p>
      )}
    </div>
  )
}

const statusIcons = { success: '✓', warning: '!', danger: '✕', info: 'i' } as const
const statusLabels = { success: 'Success', warning: 'Warning', danger: 'Error', info: 'Note' } as const

/** Semantic message: colour plus icon plus a text label, never colour alone. Not a live region by itself. */
export function StatusMessage({ tone, children }: { tone: keyof typeof statusIcons; children: ReactNode }) {
  return (
    <div className={`ui-status ui-status--${tone}`}>
      <span className="ui-status__icon" aria-hidden="true">
        {statusIcons[tone]}
      </span>
      <p>
        <strong>{statusLabels[tone]}:</strong> {children}
      </p>
    </div>
  )
}

/** Kingdom ownership marker: always shows the kingdom name next to the colour swatch. */
export function KingdomBadge({ kingdom }: { kingdom: Kingdom }) {
  const name = kingdoms.find((k) => k.id === kingdom)?.name ?? kingdom
  return (
    <span className="ui-kingdom-badge" data-kingdom={kingdom}>
      <span className="ui-kingdom-badge__swatch" aria-hidden="true" />
      {name}
    </span>
  )
}
