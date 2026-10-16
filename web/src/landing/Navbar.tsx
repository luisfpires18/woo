import { useEffect, useId, useRef, useState, type KeyboardEvent } from 'react'
import { AnvilIcon, ChevronDownIcon } from './icons'

export type Account = { nickname: string; avatarUrl?: string }

/**
 * Slim navigation bar: brand at left; Log in (signed out) or the account menu (signed in) at right.
 * No section links and no appearance controls (Settings is its own page).
 */
export function Navbar({ account }: { account: Account | null }) {
  return (
    <header className="navbar">
      <div className="navbar__inner">
        <a className="navbar__brand" href="/" aria-label="Weapons of Order home">
          <span className="navbar__anvil">
            <AnvilIcon />
          </span>
          <span className="navbar__wordmark">WOO</span>
        </a>
        {account ? (
          <AccountMenu account={account} />
        ) : (
          <a className="ui-button ui-button--primary" href="/login">
            Log in
          </a>
        )}
      </div>
    </header>
  )
}

const initials = (name: string) =>
  name
    .split(/\s+/)
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0]!.toUpperCase())
    .join('')

/**
 * Disclosure-style account menu (button + list of links). Keyboard: Enter/Space/ArrowDown opens and focuses the
 * first item, ArrowUp/ArrowDown/Home/End move, Escape closes and returns focus to the button; Tab leaves normally
 * and closes the menu. Log out is wired once real sessions exist.
 */
function AccountMenu({ account }: { account: Account }) {
  const [open, setOpen] = useState(false)
  const menuId = useId()
  const rootRef = useRef<HTMLDivElement>(null)
  const buttonRef = useRef<HTMLButtonElement>(null)
  const itemsRef = useRef<HTMLUListElement>(null)
  const focusTarget = useRef<'first' | 'last' | null>(null)

  const items = () => [...(itemsRef.current?.querySelectorAll<HTMLElement>('a, button') ?? [])]

  useEffect(() => {
    if (!open) return
    if (focusTarget.current) {
      const list = items()
      ;(focusTarget.current === 'first' ? list[0] : list[list.length - 1])?.focus()
      focusTarget.current = null
    }
    const closeOutside = (event: PointerEvent) => {
      if (!rootRef.current?.contains(event.target as Node)) setOpen(false)
    }
    document.addEventListener('pointerdown', closeOutside)
    return () => document.removeEventListener('pointerdown', closeOutside)
  }, [open])

  const openWith = (target: 'first' | 'last') => {
    focusTarget.current = target
    if (open) {
      const list = items()
      ;(target === 'first' ? list[0] : list[list.length - 1])?.focus()
      focusTarget.current = null
    } else setOpen(true)
  }

  const onButtonKeyDown = (event: KeyboardEvent) => {
    if (event.key === 'ArrowDown') {
      event.preventDefault()
      openWith('first')
    } else if (event.key === 'ArrowUp') {
      event.preventDefault()
      openWith('last')
    }
  }

  const onMenuKeyDown = (event: KeyboardEvent) => {
    const list = items()
    const index = list.indexOf(document.activeElement as HTMLElement)
    const move = (to: number) => {
      event.preventDefault()
      list[(to + list.length) % list.length]?.focus()
    }
    if (event.key === 'ArrowDown') move(index + 1)
    else if (event.key === 'ArrowUp') move(index - 1)
    else if (event.key === 'Home') move(0)
    else if (event.key === 'End') move(list.length - 1)
    else if (event.key === 'Escape') {
      event.preventDefault()
      setOpen(false)
      buttonRef.current?.focus()
    }
  }

  return (
    <div
      className="account-menu"
      ref={rootRef}
      onBlur={(event) => {
        if (!rootRef.current?.contains(event.relatedTarget as Node)) setOpen(false)
      }}
    >
      <button
        ref={buttonRef}
        type="button"
        className="account-menu__button"
        aria-label={`${account.nickname}: account menu`}
        aria-expanded={open}
        aria-controls={menuId}
        onClick={() => {
          if (!open) focusTarget.current = 'first'
          setOpen(!open)
        }}
        onKeyDown={onButtonKeyDown}
      >
        <span className="account-menu__avatar" aria-hidden="true">
          {account.avatarUrl ? <img src={account.avatarUrl} alt="" /> : initials(account.nickname)}
        </span>
        <span className="account-menu__name">{account.nickname}</span>
        <ChevronDownIcon />
      </button>
      <ul id={menuId} ref={itemsRef} className="account-menu__list" hidden={!open} onKeyDown={onMenuKeyDown}>
        <li>
          <a href="/profile">Profile</a>
        </li>
        <li>
          <a href="/settings">Settings</a>
        </li>
        <li>
          <button type="button">Log out</button>
        </li>
      </ul>
    </div>
  )
}
