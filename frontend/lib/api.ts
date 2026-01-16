// =================================
// SIPODI API Client
// =================================

import type { ApiError, RefreshResponse } from '@/types';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

// Token storage (in-memory for security)
let accessToken: string | null = null;

export function setAccessToken(token: string | null) {
    accessToken = token;
}

export function getAccessToken(): string | null {
    return accessToken;
}

// Error class for API errors
export class ApiException extends Error {
    code: string;
    status: number;
    details?: { field: string; message: string }[];

    constructor(status: number, error: ApiError['error']) {
        super(error.message);
        this.name = 'ApiException';
        this.code = error.code;
        this.status = status;
        this.details = error.details;
    }
}

// Request options type
interface RequestOptions extends RequestInit {
    params?: Record<string, string | number | boolean | undefined>;
}

// Build URL with query params
function buildUrl(endpoint: string, params?: Record<string, string | number | boolean | undefined>): string {
    const url = new URL(`${API_BASE_URL}${endpoint}`);
    if (params) {
        Object.entries(params).forEach(([key, value]) => {
            if (value !== undefined && value !== '') {
                url.searchParams.append(key, String(value));
            }
        });
    }
    return url.toString();
}

// Refresh token
async function refreshAccessToken(): Promise<string | null> {
    try {
        const response = await fetch(`${API_BASE_URL}/auth/refresh`, {
            method: 'POST',
            credentials: 'include', // Include cookies
        });

        if (!response.ok) {
            return null;
        }

        const data: { data: RefreshResponse } = await response.json();
        setAccessToken(data.data.access_token);
        return data.data.access_token;
    } catch {
        return null;
    }
}

// Core fetch function with auth
async function apiFetch<T>(
    endpoint: string,
    options: RequestOptions = {}
): Promise<T> {
    const { params, ...fetchOptions } = options;
    const url = buildUrl(endpoint, params);

    const headers: HeadersInit = {
        'Content-Type': 'application/json',
        ...fetchOptions.headers,
    };

    if (accessToken) {
        (headers as Record<string, string>)['Authorization'] = `Bearer ${accessToken}`;
    }

    let response = await fetch(url, {
        ...fetchOptions,
        headers,
        credentials: 'include', // Include cookies for refresh token
    });

    // Handle 401 - try to refresh token
    if (response.status === 401 && accessToken) {
        const newToken = await refreshAccessToken();
        if (newToken) {
            (headers as Record<string, string>)['Authorization'] = `Bearer ${newToken}`;
            response = await fetch(url, {
                ...fetchOptions,
                headers,
                credentials: 'include',
            });
        }
    }

    // Handle non-OK responses
    if (!response.ok) {
        let errorData: ApiError;
        try {
            errorData = await response.json();
        } catch {
            errorData = {
                error: {
                    code: 'UNKNOWN_ERROR',
                    message: 'Terjadi kesalahan pada server',
                },
            };
        }
        throw new ApiException(response.status, errorData.error);
    }

    // Handle 204 No Content
    if (response.status === 204) {
        return {} as T;
    }

    return response.json();
}

// API methods
export const api = {
    get<T>(endpoint: string, params?: Record<string, string | number | boolean | undefined>): Promise<T> {
        return apiFetch<T>(endpoint, { method: 'GET', params });
    },

    post<T>(endpoint: string, body?: unknown): Promise<T> {
        return apiFetch<T>(endpoint, {
            method: 'POST',
            body: body ? JSON.stringify(body) : undefined,
        });
    },

    put<T>(endpoint: string, body?: unknown): Promise<T> {
        return apiFetch<T>(endpoint, {
            method: 'PUT',
            body: body ? JSON.stringify(body) : undefined,
        });
    },

    patch<T>(endpoint: string, body?: unknown): Promise<T> {
        return apiFetch<T>(endpoint, {
            method: 'PATCH',
            body: body ? JSON.stringify(body) : undefined,
        });
    },

    delete<T>(endpoint: string): Promise<T> {
        return apiFetch<T>(endpoint, { method: 'DELETE' });
    },

    async download(endpoint: string, filename: string): Promise<void> {
        const url = buildUrl(endpoint);
        const headers: HeadersInit = {};
        if (accessToken) {
            (headers as Record<string, string>)['Authorization'] = `Bearer ${accessToken}`;
        }

        const response = await fetch(url, {
            method: 'GET',
            headers,
            credentials: 'include',
        });

        if (!response.ok) {
            throw new Error("Gagal mengunduh file");
        }

        const blob = await response.blob();
        const downloadUrl = window.URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = downloadUrl;
        link.download = filename;
        document.body.appendChild(link);
        link.click();
        link.remove();
        window.URL.revokeObjectURL(downloadUrl);
    },
};

export default api;
