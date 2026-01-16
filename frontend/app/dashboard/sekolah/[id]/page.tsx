"use client";

import { useState, useEffect, useCallback } from "react";
import { useParams, useRouter } from "next/navigation";
import { api } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import {
    Card,
    CardContent,
    CardHeader,
    CardTitle,
} from "@/components/ui/card";
import { ArrowLeft, Users, School, MapPin, Hash } from "lucide-react";
import { toast } from "sonner";
import type { SchoolDetail, UserListItem, ListResponse, DataResponse } from "@/types";
import { USER_ROLE_LABELS, GTK_TYPE_LABELS, SCHOOL_STATUS_LABELS } from "@/types";

export default function SchoolDetailPage() {
    const params = useParams();
    const router = useRouter();
    const id = params.id as string;

    const [school, setSchool] = useState<SchoolDetail | null>(null);
    const [users, setUsers] = useState<UserListItem[]>([]);
    const [loading, setLoading] = useState(true);

    const fetchData = useCallback(async () => {
        setLoading(true);
        try {
            // Fetch school detail
            const schoolRes = await api.get<DataResponse<SchoolDetail>>(`/schools/${id}`);
            setSchool(schoolRes.data);

            // Fetch users in school
            const usersRes = await api.get<ListResponse<UserListItem>>(`/schools/${id}/users`);
            setUsers(usersRes.data);
        } catch (error) {
            console.error("Failed to fetch data:", error);
            toast.error("Gagal memuat detail sekolah");
            router.push("/dashboard/sekolah");
        } finally {
            setLoading(false);
        }
    }, [id, router]);

    useEffect(() => {
        if (id) fetchData();
    }, [fetchData, id]);

    if (loading) {
        return (
            <div className="space-y-6">
                <div className="flex items-center gap-4">
                    <Skeleton className="h-8 w-8" />
                    <div className="space-y-2">
                        <Skeleton className="h-8 w-64" />
                        <Skeleton className="h-4 w-32" />
                    </div>
                </div>
                <div className="grid gap-4 md:grid-cols-3">
                    <Skeleton className="h-32" />
                    <Skeleton className="h-32" />
                    <Skeleton className="h-32" />
                </div>
                <Skeleton className="h-64" />
            </div>
        );
    }

    if (!school) return null;

    return (
        <div className="space-y-6">
            {/* Header */}
            <div className="flex items-center gap-4">
                <Button variant="ghost" size="icon" onClick={() => router.back()}>
                    <ArrowLeft className="h-4 w-4" />
                </Button>
                <div>
                    <h1 className="text-2xl font-bold tracking-tight">{school.name}</h1>
                    <div className="flex items-center gap-2 text-muted-foreground">
                        <Badge variant="outline">{SCHOOL_STATUS_LABELS[school.status]}</Badge>
                        <span>•</span>
                        <span className="flex items-center gap-1">
                            <MapPin className="h-3 w-3" /> {school.address}
                        </span>
                    </div>
                </div>
            </div>

            {/* Stats Cards */}
            <div className="grid gap-4 md:grid-cols-3">
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">NPSN</CardTitle>
                        <Hash className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">{school.npsn}</div>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Total GTK</CardTitle>
                        <Users className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">{school.gtk_count}</div>
                        <p className="text-xs text-muted-foreground">
                            {school.guru_count} Guru, {school.tendik_count} Tendik
                        </p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Kepala Sekolah</CardTitle>
                        <School className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-lg font-bold truncate">
                            {school.head_master?.full_name || "-"}
                        </div>
                        <p className="text-xs text-muted-foreground">
                            {school.head_master?.nip ? `NIP: ${school.head_master.nip}` : "Belum diatur"}
                        </p>
                    </CardContent>
                </Card>
            </div>

            {/* GTK List */}
            <div className="space-y-4">
                <h3 className="text-lg font-semibold">Daftar GTK</h3>
                <div className="rounded-md border">
                    <Table>
                        <TableHeader>
                            <TableRow>
                                <TableHead>Nama</TableHead>
                                <TableHead>NIP/NUPTK</TableHead>
                                <TableHead>Jenis GTK</TableHead>
                                <TableHead>Jabatan</TableHead>
                                <TableHead>Status</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {users.length === 0 ? (
                                <TableRow>
                                    <TableCell colSpan={5} className="h-24 text-center text-muted-foreground">
                                        Tidak ada data GTK di sekolah ini
                                    </TableCell>
                                </TableRow>
                            ) : (
                                users.map((user) => (
                                    <TableRow key={user.id}>
                                        <TableCell className="font-medium">{user.full_name}</TableCell>
                                        <TableCell>{user.nip || user.nuptk || "-"}</TableCell>
                                        <TableCell>
                                            {user.gtk_type ? GTK_TYPE_LABELS[user.gtk_type] : "-"}
                                        </TableCell>
                                        <TableCell>{user.position || "-"}</TableCell>
                                        <TableCell>
                                            <Badge variant={user.is_active ? "secondary" : "destructive"}>
                                                {user.is_active ? "Aktif" : "Nonaktif"}
                                            </Badge>
                                            <div className="mt-1">
                                                <Badge variant="outline" className="text-xs">
                                                    {USER_ROLE_LABELS[user.role]}
                                                </Badge>
                                            </div>
                                        </TableCell>
                                    </TableRow>
                                ))
                            )}
                        </TableBody>
                    </Table>
                </div>
            </div>
        </div>
    );
}
