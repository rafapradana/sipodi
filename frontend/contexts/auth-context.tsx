'use client';

import React, { createContext, useContext, useState, useCallback, useEffect } from 'react';
import { api, setAccessToken, getAccessToken } from '@/lib/api';
import type { User, LoginRequest, LoginResponse, DataResponse } from '@/types';

interface AuthContextType {
    user: User | null;
    isAuthenticated: boolean;
    isLoading: boolean;
    login: (credentials: LoginRequest) => Promise<void>;
    logout: () => Promise<void>;
    refreshUser: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
    const [user, setUser] = useState<User | null>(null);
    const [isLoading, setIsLoading] = useState(true);

    // Fetch current user profile
    const refreshUser = useCallback(async () => {
        try {
            const response = await api.get<DataResponse<User>>('/me');
            setUser(response.data);
        } catch {
            setUser(null);
            setAccessToken(null);
        }
    }, []);

    // Try to restore session on mount
    useEffect(() => {
        const initAuth = async () => {
            // Try to refresh token (will use HttpOnly cookie)
            try {
                const response = await api.post<DataResponse<{ access_token: string; user: User }>>('/auth/refresh');
                if (response.data.access_token) {
                    setAccessToken(response.data.access_token);
                    await refreshUser();
                }
            } catch {
                // No valid session
                setUser(null);
            } finally {
                setIsLoading(false);
            }
        };

        initAuth();
    }, [refreshUser]);

    // Login
    const login = useCallback(async (credentials: LoginRequest) => {
        const response = await api.post<DataResponse<LoginResponse>>('/auth/login', credentials);
        setAccessToken(response.data.access_token);
        setUser(response.data.user);
    }, []);

    // Logout
    const logout = useCallback(async () => {
        try {
            await api.post('/auth/logout');
        } catch {
            // Ignore errors on logout
        } finally {
            setAccessToken(null);
            setUser(null);
        }
    }, []);

    const value: AuthContextType = {
        user,
        isAuthenticated: !!user && !!getAccessToken(),
        isLoading,
        login,
        logout,
        refreshUser,
    };

    return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
}

export default AuthContext;
