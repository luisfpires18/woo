import { useState } from 'react';
import { useAssetStore } from '../../stores/assetStore';
import styles from './GameIcon.module.css';

interface GameIconProps {
  /** The asset id, e.g. "iron_mine", "iron", "swordsman". */
  assetId: string;
  /** Emoji fallback if no sprite or asset not found. */
  fallback?: string;
  /** Size override in px. Defaults to the asset's native sprite width. */
  size?: number;
  /** Additional CSS class name. */
  className?: string;
}

/**
 * Renders a sprite `<img>` if the asset has an uploaded sprite,
 * otherwise falls back to the emoji from the asset table (or the
 * explicit `fallback` prop).
 */
export function GameIcon({ assetId, fallback, size, className }: GameIconProps) {
  const asset = useAssetStore((s) => s.getById(assetId));
  const [imgError, setImgError] = useState(false);

  const emoji = asset?.default_icon ?? fallback ?? '❓';
  const spriteUrl = asset?.sprite_url;
  const displaySize = size ?? asset?.sprite_width ?? 32;

  const wrapClass = [styles.icon, className].filter(Boolean).join(' ');

  if (spriteUrl && !imgError) {
    return (
      <img
        src={spriteUrl}
        alt={asset?.display_name ?? assetId}
        width={displaySize}
        height={displaySize}
        className={wrapClass}
        onError={() => setImgError(true)}
        draggable={false}
      />
    );
  }

  return (
    <span
      className={wrapClass}
      style={{ fontSize: `${displaySize}px`, lineHeight: 1 }}
      role="img"
      aria-label={asset?.display_name ?? assetId}
    >
      {emoji}
    </span>
  );
}
