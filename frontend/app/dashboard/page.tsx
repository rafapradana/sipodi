"use client";

import { useEffect, useState } from "react";
import { useAuth } from "@/hooks/use-auth";
import { api } from "@/lib/api";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Badge } from "@/components/ui/badge";
import {
    School,
    Users,
    Star,
    Clock,
    CheckCircle,
    XCircle,
    GraduationCap,
    Briefcase,
    UserCheck,
} from "lucide-react";
import type { DashboardSummary, DataResponse, TalentListItem } from "@/types";
import { TALENT_TYPE_LABELS, TALENT_STATUS_LABELS } from "@/types";

export default function DashboardPage() {
    const { user } = useAuth();
    const [summary, setSummary] = useState<DashboardSummary | null>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchSummary = async () => {
            try {
                const response = await api.get<DataResponse<DashboardSummary>>("/dashboard/summary");
                setSummary(response.data);
            } catch (error) {
                console.error("Failed to fetch dashboard summary:", error);
            } finally {
                setLoading(false);
            }
        };

        fetchSummary();
    }, []);

    if (loading) {
        return <DashboardSkeleton />;
    }

    return (
        <div className="space-y-6">
            {/* Page Title */}
            <div>
                <h1 className="text-2xl font-bold tracking-tight">Dashboard</h1>
                <p className="text-muted-foreground">
                    Selamat datang kembali, {user?.full_name}!
                </p>
            </div>

            {/* Summary Cards */}
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                {user?.role === "super_admin" && (
                    <>
                        <SummaryCard
                            title="Total Sekolah"
                            value={summary?.total_schools || 0}
                            icon={School}
                            description="Sekolah terdaftar"
                        />
                        <SummaryCard
                            title="Total GTK"
                            value={summary?.total_gtk || 0}
                            icon={Users}
                            description="Guru & Tenaga Kependidikan"
                        />
                        <SummaryCard
                            title="Total Talenta"
                            value={summary?.total_talents || 0}
                            icon={Star}
                            description="Data talenta GTK"
                        />
                        <SummaryCard
                            title="Menunggu Verifikasi"
                            value={summary?.talents_by_status?.pending || 0}
                            icon={Clock}
                            description="Perlu ditinjau"
                            variant="warning"
                        />
                    </>
                )}
            </div>

            {/* GTK by Type */}
            {user?.role === "super_admin" && summary?.gtk_by_type && (
                <Card>
                    <CardHeader>
                        <CardTitle>GTK Berdasarkan Jenis</CardTitle>
                        <CardDescription>Distribusi GTK berdasarkan jenis</CardDescription>
                    </CardHeader>
                    <CardContent>
                        <div className="grid gap-4 md:grid-cols-3">
                            <div className="flex items-center gap-4 rounded-lg border p-4">
                                <div className="flex h-12 w-12 items-center justify-center rounded-full bg-blue-100 dark:bg-blue-900/30">
                                    <GraduationCap className="h-6 w-6 text-blue-600" />
                                </div>
                                <div>
                                    <p className="text-2xl font-bold">{summary.gtk_by_type.guru || 0}</p>
                                    <p className="text-sm text-muted-foreground">Guru</p>
                                </div>
                            </div>
                            <div className="flex items-center gap-4 rounded-lg border p-4">
                                <div className="flex h-12 w-12 items-center justify-center rounded-full bg-green-100 dark:bg-green-900/30">
                                    <Briefcase className="h-6 w-6 text-green-600" />
                                </div>
                                <div>
                                    <p className="text-2xl font-bold">{summary.gtk_by_type.tendik || 0}</p>
                                    <p className="text-sm text-muted-foreground">Tenaga Kependidikan</p>
                                </div>
                            </div>
                            <div className="flex items-center gap-4 rounded-lg border p-4">
                                <div className="flex h-12 w-12 items-center justify-center rounded-full bg-purple-100 dark:bg-purple-900/30">
                                    <UserCheck className="h-6 w-6 text-purple-600" />
                                </div>
                                <div>
                                    <p className="text-2xl font-bold">{summary.gtk_by_type.kepala_sekolah || 0}</p>
                                    <p className="text-sm text-muted-foreground">Kepala Sekolah</p>
                                </div>
                            </div>
                        </div>
                    </CardContent>
                </Card>
            )}

            {/* Talents by Status */}
            {user?.role === "super_admin" && summary?.talents_by_status && (
                <Card>
                    <CardHeader>
                        <CardTitle>Talenta Berdasarkan Status</CardTitle>
                        <CardDescription>Distribusi talenta berdasarkan status verifikasi</CardDescription>
                    </CardHeader>
                    <CardContent>
                        <div className="grid gap-4 md:grid-cols-3">
                            <div className="flex items-center gap-4 rounded-lg border p-4">
                                <div className="flex h-12 w-12 items-center justify-center rounded-full bg-yellow-100 dark:bg-yellow-900/30">
                                    <Clock className="h-6 w-6 text-yellow-600" />
                                </div>
                                <div>
                                    <p className="text-2xl font-bold">{summary.talents_by_status.pending || 0}</p>
                                    <p className="text-sm text-muted-foreground">Pending</p>
                                </div>
                            </div>
                            <div className="flex items-center gap-4 rounded-lg border p-4">
                                <div className="flex h-12 w-12 items-center justify-center rounded-full bg-green-100 dark:bg-green-900/30">
                                    <CheckCircle className="h-6 w-6 text-green-600" />
                                </div>
                                <div>
                                    <p className="text-2xl font-bold">{summary.talents_by_status.approved || 0}</p>
                                    <p className="text-sm text-muted-foreground">Disetujui</p>
                                </div>
                            </div>
                            <div className="flex items-center gap-4 rounded-lg border p-4">
                                <div className="flex h-12 w-12 items-center justify-center rounded-full bg-red-100 dark:bg-red-900/30">
                                    <XCircle className="h-6 w-6 text-red-600" />
                                </div>
                                <div>
                                    <p className="text-2xl font-bold">{summary.talents_by_status.rejected || 0}</p>
                                    <p className="text-sm text-muted-foreground">Ditolak</p>
                                </div>
                            </div>
                        </div>
                    </CardContent>
                </Card>
            )}

            {/* Recent Talents */}
            {summary?.recent_talents && summary.recent_talents.length > 0 && (
                <Card>
                    <CardHeader>
                        <CardTitle>Talenta Terbaru</CardTitle>
                        <CardDescription>Talenta yang baru ditambahkan</CardDescription>
                    </CardHeader>
                    <CardContent>
                        <div className="space-y-4">
                            {summary.recent_talents.slice(0, 5).map((talent) => (
                                <RecentTalentItem key={talent.id} talent={talent} />
                            ))}
                        </div>
                    </CardContent>
                </Card>
            )}
        </div>
    );
}

function SummaryCard({
    title,
    value,
    icon: Icon,
    description,
    variant = "default",
}: {
    title: string;
    value: number;
    icon: React.ComponentType<{ className?: string }>;
    description: string;
    variant?: "default" | "warning";
}) {
    return (
        <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">{title}</CardTitle>
                <Icon
                    className={`h-4 w-4 ${variant === "warning" ? "text-yellow-600" : "text-muted-foreground"
                        }`}
                />
            </CardHeader>
            <CardContent>
                <div className="text-2xl font-bold">{value.toLocaleString("id-ID")}</div>
                <p className="text-xs text-muted-foreground">{description}</p>
            </CardContent>
        </Card>
    );
}

function RecentTalentItem({ talent }: { talent: TalentListItem }) {
    const getStatusBadge = (status: string) => {
        switch (status) {
            case "pending":
                return <Badge variant="secondary">Pending</Badge>;
            case "approved":
                return <Badge className="bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400">Disetujui</Badge>;
            case "rejected":
                return <Badge variant="destructive">Ditolak</Badge>;
            default:
                return null;
        }
    };

    return (
        <div className="flex items-center justify-between rounded-lg border p-3">
            <div className="flex items-center gap-3">
                <div className="flex h-10 w-10 items-center justify-center rounded-full bg-muted">
                    <Star className="h-5 w-5 text-muted-foreground" />
                </div>
                <div>
                    <p className="font-medium">{talent.user?.full_name || "Unknown"}</p>
                    <p className="text-sm text-muted-foreground">
                        {TALENT_TYPE_LABELS[talent.talent_type]} • {talent.user?.school_name || ""}
                    </p>
                </div>
            </div>
            <div className="flex items-center gap-2">
                {getStatusBadge(talent.status)}
                <span className="text-xs text-muted-foreground">
                    {new Date(talent.created_at).toLocaleDateString("id-ID")}
                </span>
            </div>
        </div>
    );
}

function DashboardSkeleton() {
    return (
        <div className="space-y-6">
            <div>
                <Skeleton className="h-8 w-48" />
                <Skeleton className="mt-2 h-4 w-64" />
            </div>
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                {[1, 2, 3, 4].map((i) => (
                    <Card key={i}>
                        <CardHeader className="pb-2">
                            <Skeleton className="h-4 w-24" />
                        </CardHeader>
                        <CardContent>
                            <Skeleton className="h-8 w-16" />
                            <Skeleton className="mt-1 h-3 w-32" />
                        </CardContent>
                    </Card>
                ))}
            </div>
        </div>
    );
}
