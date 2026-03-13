// WebSocket message types

export interface WSMessage {
  type: string;
  data?: unknown;
}

export interface WSSubscribe extends WSMessage {
  type: 'subscribe';
  data: { topics: string[] };
}

export interface WSUnsubscribe extends WSMessage {
  type: 'unsubscribe';
  data: { topics: string[] };
}

export interface WSResourceUpdate extends WSMessage {
  type: 'resource_update';
  data: {
    village_id: number;
    food: number;
    water: number;
    lumber: number;
    stone: number;
  };
}

export interface WSTrainComplete extends WSMessage {
  type: 'train_complete';
  data: {
    village_id: number;
    troop_type: string;
    new_total: number;
  };
}

export interface WSGoldUpdate extends WSMessage {
  type: 'gold_update';
  data: {
    player_id: number;
    gold: number;
  };
}

export interface WSExpeditionComplete extends WSMessage {
  type: 'expedition_complete';
  data: {
    village_id: number;
    expedition_id: number;
    camp_id: number;
    result: 'attacker_won' | 'defender_won' | 'draw';
  };
}

export interface WSExpeditionReturn extends WSMessage {
  type: 'expedition_return';
  data: {
    village_id: number;
    expedition_id: number;
  };
}
