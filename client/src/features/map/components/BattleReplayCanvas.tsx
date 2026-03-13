// PixiJS-powered battle replay viewer — auto-plays, skip-only controls

import { useEffect, useRef, useState } from 'react';
import type { BattleReplayResponse } from '../../../types/api';
import { BattleScene } from './battle/BattleScene';
import styles from './BattleReplayCanvas.module.css';

interface BattleReplayCanvasProps {
  replay: BattleReplayResponse;
  onClose: () => void;
}

export function BattleReplayCanvas({ replay, onClose }: BattleReplayCanvasProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const sceneRef = useRef<BattleScene | null>(null);
  const [finished, setFinished] = useState(false);

  // Init scene on mount, auto-play immediately
  useEffect(() => {
    const el = containerRef.current;
    if (!el) return;

    let autoCloseTimer: ReturnType<typeof setTimeout> | undefined;

    const scene = new BattleScene(el, replay, () => {
      // onComplete — battle ended, auto-transition after 1.5s
      setFinished(true);
      autoCloseTimer = setTimeout(onClose, 1500);
    });

    sceneRef.current = scene;

    scene.init().then(() => {
      scene.play();
    });

    return () => {
      if (autoCloseTimer) clearTimeout(autoCloseTimer);
      scene.destroy();
      sceneRef.current = null;
    };
  }, [replay, onClose]);

  return (
    <div className={styles.container}>
      <div ref={containerRef} className={styles.sceneContainer} />
      {finished && <div className={styles.finishedLabel}>Battle Complete</div>}
      <button className={styles.skipBtn} onClick={onClose}>
        Skip
      </button>
    </div>
  );
}
