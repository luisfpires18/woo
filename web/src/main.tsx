import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { App } from './App'
import { applyAppearance } from './theme'
import './tokens.css'
import './styles.css'
import './landing.css'
import './title.css'

// Follow the operating system until a Settings page stores an explicit light/dark choice.
applyAppearance('system')

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
