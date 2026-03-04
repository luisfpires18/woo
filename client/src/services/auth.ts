// Auth API service — register, login, refresh, logout

import { api } from './api';
import type {
  AuthResponse,
  RegisterRequest,
  LoginRequest,
} from '../types/api';

export async function registerPlayer(
  data: RegisterRequest,
): Promise<AuthResponse> {
  return api.post<AuthResponse>('/auth/register', data, true);
}

export async function loginPlayer(data: LoginRequest): Promise<AuthResponse> {
  return api.post<AuthResponse>('/auth/login', data, true);
}

export async function refreshToken(token: string): Promise<AuthResponse> {
  return api.post<AuthResponse>(
    '/auth/refresh',
    { refresh_token: token },
    true,
  );
}

export async function logoutPlayer(token: string): Promise<void> {
  await api.post('/auth/logout', { refresh_token: token }, true);
}
