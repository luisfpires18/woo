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
