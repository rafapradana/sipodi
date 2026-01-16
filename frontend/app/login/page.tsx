'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/hooks/use-auth';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { ThemeToggle } from '@/components/theme-toggle';
import { ApiException } from '@/lib/api';
import { Loader2, GraduationCap } from 'lucide-react';

export default function LoginPage() {
    const router = useRouter();
    const { login, isAuthenticated, isLoading: authLoading } = useAuth();
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const [isLoading, setIsLoading] = useState(false);

    // Redirect if already authenticated
    if (authLoading) {
        return (
            <div className="flex min-h-screen items-center justify-center">
                <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
            </div>
        );
    }

    if (isAuthenticated) {
        router.replace('/dashboard');
        return null;
    }

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');
        setIsLoading(true);

        try {
            await login({ email, password });
            router.push('/dashboard');
        } catch (err) {
            if (err instanceof ApiException) {
                setError(err.message);
            } else {
                setError('Terjadi kesalahan. Silakan coba lagi.');
            }
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className="relative flex min-h-screen items-center justify-center bg-gradient-to-br from-background via-background to-muted/30 p-4">
            {/* Theme Toggle */}
            <div className="absolute right-4 top-4">
                <ThemeToggle />
            </div>

            {/* Login Card */}
            <Card className="w-full max-w-md shadow-lg">
                <CardHeader className="space-y-4 text-center">
                    {/* Logo */}
                    <div className="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-primary">
                        <GraduationCap className="h-8 w-8 text-primary-foreground" />
                    </div>
                    <div>
                        <CardTitle className="text-2xl font-bold">SIPODI</CardTitle>
                        <CardDescription className="text-sm">
                            Sistem Informasi Potensi Diri GTK
                            <br />
                            Cabang Dinas Pendidikan Wilayah Malang
                        </CardDescription>
                    </div>
                </CardHeader>
                <CardContent>
                    <form onSubmit={handleSubmit} className="space-y-4">
                        {/* Error Message */}
                        {error && (
                            <div className="rounded-md bg-destructive/10 p-3 text-sm text-destructive">
                                {error}
                            </div>
                        )}

                        {/* Email Field */}
                        <div className="space-y-2">
                            <Label htmlFor="email">Email</Label>
                            <Input
                                id="email"
                                type="email"
                                placeholder="contoh@sekolah.sch.id"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                                required
                                disabled={isLoading}
                                autoComplete="email"
                            />
                        </div>

                        {/* Password Field */}
                        <div className="space-y-2">
                            <Label htmlFor="password">Password</Label>
                            <Input
                                id="password"
                                type="password"
                                placeholder="Masukkan password"
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                required
                                disabled={isLoading}
                                autoComplete="current-password"
                            />
                        </div>

                        {/* Submit Button */}
                        <Button type="submit" className="w-full" disabled={isLoading}>
                            {isLoading ? (
                                <>
                                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                    Masuk...
                                </>
                            ) : (
                                'Masuk'
                            )}
                        </Button>
                    </form>
                </CardContent>
            </Card>

            {/* Footer */}
            <div className="absolute bottom-4 text-center text-xs text-muted-foreground">
                © 2026 Cabang Dinas Pendidikan Wilayah Malang
            </div>
        </div>
    );
}
