"use client";

import { useState, useEffect, useCallback } from "react";
import { api } from "@/lib/api";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Check, Bell, CheckCheck } from "lucide-react";
import { toast } from "sonner";
import { formatDistanceToNow } from "date-fns";
import { id as idLocale } from "date-fns/locale";

import type { Notification, ListResponse } from "@/types";

export default function NotificationsPage() {
    const [notifications, setNotifications] = useState<Notification[]>([]);
    const [loading, setLoading] = useState(true);

    const fetchNotifications = useCallback(async () => {
        setLoading(true);
        try {
            const response = await api.get<ListResponse<Notification>>("/me/notifications");
            setNotifications(response.data);
        } catch (error) {
            console.error("Failed to list notifications:", error);
            toast.error("Gagal memuat notifikasi");
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        fetchNotifications();
    }, [fetchNotifications]);

    const handleMarkAsRead = async (id: string) => {
        try {
            await api.patch(`/me/notifications/${id}/read`);
            setNotifications((prev) =>
                prev.map((n) => (n.id === id ? { ...n, is_read: true } : n))
            );
        } catch (error) {
            toast.error("Gagal update status notifikasi");
        }
    };

    const handleMarkAllRead = async () => {
        try {
            await api.patch("/me/notifications/read-all");
            setNotifications((prev) => prev.map((n) => ({ ...n, is_read: true })));
            toast.success("Semua notifikasi ditandai sudah dibaca");
        } catch (error) {
            toast.error("Gagal update status notifikasi");
        }
    };

    const getNotificationIcon = (type: string) => {
        // Add logic if different types have different icons
        return <Bell className="h-5 w-5 text-primary" />;
    };

    return (
        <div className="space-y-6">
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-2xl font-bold tracking-tight">Notifikasi</h1>
                    <p className="text-muted-foreground">Pemberitahuan aktivitas terbaru</p>
                </div>
                {notifications.some((n) => !n.is_read) && (
                    <Button variant="outline" size="sm" onClick={handleMarkAllRead} disabled={loading}>
                        <CheckCheck className="mr-2 h-4 w-4" />
                        Tandai Semua Dibaca
                    </Button>
                )}
            </div>

            <div className="rounded-md border p-4">
                {loading ? (
                    <div className="space-y-4">
                        <Skeleton className="h-20 w-full" />
                        <Skeleton className="h-20 w-full" />
                        <Skeleton className="h-20 w-full" />
                    </div>
                ) : notifications.length === 0 ? (
                    <div className="text-center py-10 text-muted-foreground">
                        <Bell className="mx-auto h-8 w-8 mb-3 opacity-50" />
                        Tidak ada notifikasi
                    </div>
                ) : (
                    <div className="space-y-4">
                        {notifications.map((notif) => (
                            <div
                                key={notif.id}
                                className={`flex items-start gap-4 p-4 rounded-lg border transition-colors ${notif.is_read ? 'bg-background' : 'bg-muted/30 border-primary/20'}`}
                            >
                                <div className={`mt-1 rounded-full p-2 ${notif.is_read ? 'bg-muted' : 'bg-primary/10'}`}>
                                    {getNotificationIcon(notif.type)}
                                </div>
                                <div className="flex-1 space-y-1">
                                    <p className="text-sm font-medium leading-none">{notif.message}</p>
                                    <p className="text-xs text-muted-foreground">
                                        {formatDistanceToNow(new Date(notif.created_at), { addSuffix: true, locale: idLocale })}
                                    </p>
                                </div>
                                {!notif.is_read && (
                                    <Button
                                        variant="ghost"
                                        size="icon"
                                        className="h-8 w-8 text-muted-foreground hover:text-primary"
                                        onClick={() => handleMarkAsRead(notif.id)}
                                    >
                                        <Check className="h-4 w-4" />
                                    </Button>
                                )}
                            </div>
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
}
