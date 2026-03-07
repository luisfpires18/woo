// Map API service — fetches map chunks from the server

import { api } from './api';
import type { MapChunkResponse, MapTile } from '../types/map';

/** Fetch a chunk of map tiles centered on (x, y) with the given range. */
export async function fetchMapChunk(
  x: number,
  y: number,
  range: number = 10,
): Promise<MapChunkResponse> {
  return api.get<MapChunkResponse>(`/map?x=${x}&y=${y}&range=${range}`);
}

/** Fetch a single tile's details. */
export async function fetchMapTile(x: number, y: number): Promise<MapTile> {
  return api.get<MapTile>(`/map/tile?x=${x}&y=${y}`);
}
