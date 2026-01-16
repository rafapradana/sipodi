"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/hooks/use-auth";
import { AppSidebar } from "@/components/dashboard/app-sidebar";
import { DashboardHeader } from "@/components/dashboard/header";
import {
    SidebarInset,
    SidebarProvider,
} from "@/components/ui/sidebar";
import { Toaster } from "@/components/ui/sonner";
import { Loader2 } from "lucide-react";
import { api } from "@/lib/api";
import type { DashboardSummary, DataResponse } from "@/types";

export default function DashboardLayout({
    children,
}: {
    children: React.ReactNode;
}) {
    const router = useRouter();
    const { user, isAuthenticated, isLoading } = useAuth();
    const [pendingCount, setPendingCount] = useState(0);

    // Fetch pending verification count for Super Admin
    useEffect(() => {
        if (user?.role === "super_admin") {
            api
                .get<DataResponse<DashboardSummary>>("/dashboard/summary")
                .then((res) => {
                    if (res.data.talents_by_status?.pending) {
                        setPendingCount(res.data.talents_by_status.pending);
                    }
                })
                .catch(() => { });
        }
    }, [user?.role]);

    // Show loading state
    if (isLoading) {
        return (
            <div className="flex min-h-screen items-center justify-center">
                <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
            </div>
        );
    }

    // Redirect to login if not authenticated
    if (!isAuthenticated) {
        router.replace("/login");
        return null;
    }

    return (
        <SidebarProvider>
            <AppSidebar pendingCount={pendingCount} />
            <SidebarInset>
                <DashboardHeader />
                <main className="flex-1 overflow-auto p-4 md:p-6">
                    {children}
                </main>
            </SidebarInset>
            <Toaster richColors position="top-right" />
        </SidebarProvider>
    );
}
