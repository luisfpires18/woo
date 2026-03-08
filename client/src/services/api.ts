// API client — wraps fetch with auth headers, response envelope, and error handling

import type { ApiResponse } from '../types/api';

const API_BASE = '/api';

/** Get the current access token from localStorage */
function getAccessToken(): string | null {
  return localStorage.getItem('access_token');
}

/** Get the current refresh token from localStorage */
function getRefreshToken(): string | null {
  return localStorage.getItem('refresh_token');
}

/** Attempt to refresh the access token using the stored refresh token */
async function attemptTokenRefresh(): Promise<string | null> {
  const refreshToken = getRefreshToken();
  if (!refreshToken) return null;

  try {
    const response = await fetch(`${API_BASE}/auth/refresh`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: refreshToken }),
    });

    if (!response.ok) return null;

    const result: ApiResponse<{ access_token: string; refresh_token: string }> =
      await response.json();

    localStorage.setItem('access_token', result.data.access_token);
    localStorage.setItem('refresh_token', result.data.refresh_token);
    return result.data.access_token;
  } catch {
    return null;
  }
}

async function request<T>(
  path: string,
  options?: RequestInit,
  skipAuth = false,
): Promise<T> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options?.headers as Record<string, string>),
  };

  if (!skipAuth) {
    const token = getAccessToken();
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }
  }

  let response = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers,
  });

  // On 401, try to refresh the token and retry once
  if (response.status === 401 && !skipAuth) {
    const newToken = await attemptTokenRefresh();
    if (newToken) {
      headers['Authorization'] = `Bearer ${newToken}`;
      response = await fetch(`${API_BASE}${path}`, {
        ...options,
        headers,
      });
    }
  }

  // Guard against non-JSON responses (e.g. HTML fallback when backend is down)
  const contentType = response.headers.get('content-type') ?? '';
  const isJson = contentType.includes('application/json');

  if (!response.ok) {
    if (isJson) {
      const body = await response.json().catch(() => ({ error: 'Unknown error' }));
      throw new ApiRequestError(
        body.error || `HTTP ${response.status}`,
        response.status,
      );
    }
    throw new ApiRequestError(
      response.status === 404
        ? 'Server not reachable — please try again in a moment'
        : `HTTP ${response.status}`,
      response.status,
    );
  }

  // Read the body once as text so we can inspect it on failure
  const bodyText = await response.text();

  try {
    const result: ApiResponse<T> = JSON.parse(bodyText);
    return result.data;
  } catch {
    // Log the actual response for debugging
    console.error(
      `[api] Non-JSON response for ${path}`,
      `status=${response.status}`,
      `content-type=${contentType}`,
      `body=${bodyText.substring(0, 300)}`,
    );
    throw new ApiRequestError('Server returned an unexpected response', 502);
  }
}

/** Typed API error with status code */
export class ApiRequestError extends Error {
  constructor(
    message: string,
    public status: number,
  ) {
    super(message);
    this.name = 'ApiRequestError';
  }
}

export const api = {
  get: <T>(path: string, skipAuth = false) => request<T>(path, undefined, skipAuth),

  post: <T>(path: string, body: unknown, skipAuth = false) =>
    request<T>(path, { method: 'POST', body: JSON.stringify(body) }, skipAuth),

  put: <T>(path: string, body: unknown) =>
    request<T>(path, { method: 'PUT', body: JSON.stringify(body) }),

  patch: <T>(path: string, body: unknown) =>
    request<T>(path, { method: 'PATCH', body: JSON.stringify(body) }),

  delete: <T>(path: string) => request<T>(path, { method: 'DELETE' }),
};
