import { useState } from 'react';
import type { BuildingInfo } from '../../../types/api';
import { RESOURCE_BUILDING_GROUPS } from '../../../config/buildings';
import { useResourceBuildingDisplay } from '../../../hooks/useResourceBuildingDisplay';
import styles from './ResourceFieldsGrid.module.css';

interface ResourceFieldsGridProps {
  buildings: BuildingInfo[];
  onBuildingClick: (building: BuildingInfo) => void;
}

export function ResourceFieldsGrid({ buildings, onBuildingClick }: ResourceFieldsGridProps) {
  const { getDisplay } = useResourceBuildingDisplay();
  const [failedSprites, setFailedSprites] = useState<Set<string>>(new Set());
  const buildingMap = new Map<string, BuildingInfo>();
  for (const b of buildings) {
    buildingMap.set(b.building_type, b);
  }

  const handleSpriteError = (key: string) => {
    setFailedSprites((prev) => new Set(prev).add(key));
  };

  const [failedHeaderIcons, setFailedHeaderIcons] = useState<Set<string>>(new Set());

  const handleHeaderIconError = (resource: string) => {
    setFailedHeaderIcons((prev) => new Set(prev).add(resource));
  };

  return (
    <div className={styles.groups}>
      {RESOURCE_BUILDING_GROUPS.map((group) => (
        <div key={group.resource} className={styles.group}>
          <h4 className={styles.groupTitle}>
            {!failedHeaderIcons.has(group.resource) ? (
              <img
                src={`/uploads/sprites/resources/${group.resource}.png`}
                alt={group.label}
                className={styles.groupIcon}
                onError={() => handleHeaderIconError(group.resource)}
                draggable={false}
              />
            ) : (
              <span className={styles.groupEmoji}>{group.emoji}</span>
            )}
            {group.label}
          </h4>
          <div className={styles.grid}>
            {group.types.map((type) => {
              const b = buildingMap.get(type);
              if (!b) return null;
              const isBuilt = b.level > 0;
              const display = getDisplay(type);
              const showSprite = display.spriteUrl && !failedSprites.has(type);

              return (
                <button
                  key={b.id}
                  type="button"
                  className={`${styles.card} ${isBuilt ? '' : styles.unbuilt}`}
                  onClick={() => onBuildingClick(b)}
                >
                  {showSprite ? (
                    <img
                      src={display.spriteUrl!}
                      alt={display.displayName}
                      className={styles.icon}
                      onError={() => handleSpriteError(type)}
                      draggable={false}
                    />
                  ) : (
                    <span
                      className={styles.iconEmoji}
                      role="img"
                      aria-label={display.displayName}
                    >
                      {display.emoji}
                    </span>
                  )}
                  <span className={styles.name}>{display.displayName}</span>
                  <span className={styles.level}>
                    {isBuilt ? `Lv ${b.level}` : 'Not built'}
                  </span>
                </button>
              );
            })}
          </div>
        </div>
      ))}
    </div>
  );
}
