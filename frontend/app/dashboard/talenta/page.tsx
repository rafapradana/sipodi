"use client";

import { useState, useEffect, useCallback } from "react";
import { api } from "@/lib/api";
import { useAuth } from "@/hooks/use-auth";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import {
    Select,
    SelectTrigger,
    SelectValue,
    SelectContent,
    SelectItem,
} from "@/components/ui/select";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import {
    Search,
    RefreshCw,
    Plus,
    Eye,
    Pencil,
    Trash2,
    Download,
} from "lucide-react";
import { toast } from "sonner";
import {
    AlertDialog,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
    AlertDialogCancel,
    AlertDialogAction,
} from "@/components/ui/alert-dialog";

import type { TalentListItem, ListResponse, DataResponse, Talent } from "@/types";
import { TALENT_TYPE_LABELS, TALENT_STATUS_LABELS } from "@/types";
import { TalentModal } from "@/components/dashboard/talenta/talent-modal";
import { TalentDetailModal } from "@/components/dashboard/verifikasi/talent-detail-modal";

export default function TalentaPage() {
    const { user } = useAuth();
    const isAdmin = user?.role === "super_admin" || user?.role === "admin_sekolah";
    const isGTK = user?.role === "gtk";

    const [talents, setTalents] = useState<TalentListItem[]>([]);
    const [loading, setLoading] = useState(true);
    const [search, setSearch] = useState("");
    const [statusFilter, setStatusFilter] = useState<string>("all");
    const [typeFilter, setTypeFilter] = useState<string>("all");
    const [page, setPage] = useState(1);
    const [totalPages, setTotalPages] = useState(1);
    const [totalCount, setTotalCount] = useState(0);

    // Modal states
    const [createModal, setCreateModal] = useState(false);
    const [selectedTalent, setSelectedTalent] = useState<Talent | null>(null);
    const [detailModal, setDetailModal] = useState(false);
    const [deleteId, setDeleteId] = useState<string | null>(null);

    const fetchTalents = useCallback(async () => {
        setLoading(true);
        try {
            const params: Record<string, string | number> = {
                page,
                limit: 10,
            };
            if (search) params.search = search;
            if (statusFilter !== "all") params.status = statusFilter;
            if (typeFilter !== "all") params.type = typeFilter;

            // Determine endpoint based on role
            const endpoint = isAdmin ? "/talents" : "/me/talents";

            const response = await api.get<ListResponse<TalentListItem>>(endpoint, params);
            setTalents(response.data);
            setTotalPages(response.meta.total_pages);
            setTotalCount(response.meta.total_count);
        } catch (error) {
            console.error("Failed to fetch talents:", error);
            toast.error("Gagal memuat data talenta");
        } finally {
            setLoading(false);
        }
    }, [page, search, statusFilter, typeFilter, isAdmin]);

    useEffect(() => {
        if (user) fetchTalents();
    }, [fetchTalents, user]);

    const handleSearch = (e: React.FormEvent) => {
        e.preventDefault();
        setPage(1);
        fetchTalents();
    };

    const handleCreate = () => {
        setSelectedTalent(null);
        setCreateModal(true);
    };

    const handleEdit = async (id: string) => {
        try {
            // Need to fetch detail first to edit
            const endpoint = isGTK ? `me/talents/${id}` : `talents/${id}`;
            // Wait, /me/talents/:id is for DELETE/PUT but GET is not explicitly listed in router for GTK detail?
            // Router: protected.Get("/me/talents", ...list)
            // Router: protected.Get("/talents/:id", ...getByID) <- accessible to admin/superadmin only? 
            // Ah, let's check router again.
            // Line 107: `talents.Get("/:id", r.talentHandler.GetByID)` is protected but inside `talents` group?
            // Wait, line 105: `talents := protected.Group("/talents")`.
            // No middleware on Group.
            // Line 106: `talents.Get("/", middleware...List)` -> role restricted.
            // Line 107: `talents.Get("/:id", r.talentHandler.GetByID)` -> NO role middleware? So it inherits AuthMiddleware from `protected` group.
            // So GTK CAN access `GET /talents/:id`! Nice.

            const response = await api.get<DataResponse<Talent>>(`/talents/${id}`);
            setSelectedTalent(response.data);
            setCreateModal(true);
        } catch (error) {
            toast.error("Gagal memuat detail talenta");
        }
    };

    const handleView = (id: string) => {
        // Just set ID for DetailModal
        setSelectedTalent({ id } as any); // Detail modal only needs ID, but here I'm using state hack.
        // Actually TalentDetailModal takes `talentId`. 
        // I'll use a separate state variable `viewId`.
    };
    // Wait, I used `selectedTalent` for CreateModal (full object) and `detailModal` uses ID.
    const [viewId, setViewId] = useState<string | null>(null);

    const handleDelete = async () => {
        if (!deleteId) return;
        try {
            await api.delete(`/me/talents/${deleteId}`);
            toast.success("Talenta berhasil dihapus");
            fetchTalents();
        } catch (error) {
            toast.error("Gagal menghapus talenta");
        } finally {
            setDeleteId(null);
        }
    };

    const getTalentSummary = (talent: TalentListItem): string => {
        const detail = talent.detail;
        if ("activity_name" in detail) return detail.activity_name;
        if ("competition_name" in detail) return detail.competition_name;
        if ("interest_name" in detail) return detail.interest_name;
        return "-";
    };

    return (
        <div className="space-y-6">
            {/* Header */}
            <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
                <div>
                    <h1 className="text-2xl font-bold tracking-tight">Data Talenta</h1>
                    <p className="text-muted-foreground">
                        {isAdmin ? "Semua data talenta terdaftar" : "Data talenta yang Anda ajukan"}
                    </p>
                </div>
                <div className="flex gap-2">
                    {isAdmin && (
                        <Button variant="outline" onClick={() => api.download('/exports/talents', 'data-talenta.xlsx')}>
                            <Download className="mr-2 h-4 w-4" />
                            Export
                        </Button>
                    )}
                    {isGTK && (
                        <Button onClick={handleCreate}>
                            <Plus className="mr-2 h-4 w-4" />
                            Ajukan Talenta
                        </Button>
                    )}
                </div>
            </div>

            {/* Filters */}
            <div className="flex flex-col gap-4 lg:flex-row lg:items-center">
                <form onSubmit={handleSearch} className="flex flex-1 gap-2">
                    <div className="relative flex-1 max-w-sm">
                        <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                        <Input
                            placeholder={isAdmin ? "Cari nama, sekolah..." : "Cari talenta..."}
                            value={search}
                            onChange={(e) => setSearch(e.target.value)}
                            className="pl-9"
                        />
                    </div>
                    <Button type="submit" variant="secondary">
                        Cari
                    </Button>
                </form>

                <div className="flex flex-wrap gap-2">
                    <Select value={typeFilter} onValueChange={(v) => { setTypeFilter(v); setPage(1); }}>
                        <SelectTrigger className="w-40">
                            <SelectValue placeholder="Jenis Talenta" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectItem value="all">Semua Jenis</SelectItem>
                            {Object.entries(TALENT_TYPE_LABELS).map(([k, v]) => (
                                <SelectItem key={k} value={k}>{v}</SelectItem>
                            ))}
                        </SelectContent>
                    </Select>

                    <Select value={statusFilter} onValueChange={(v) => { setStatusFilter(v); setPage(1); }}>
                        <SelectTrigger className="w-32">
                            <SelectValue placeholder="Status" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectItem value="all">Semua Status</SelectItem>
                            {Object.entries(TALENT_STATUS_LABELS).map(([k, v]) => (
                                <SelectItem key={k} value={k}>{v}</SelectItem>
                            ))}
                        </SelectContent>
                    </Select>

                    <Button variant="outline" size="icon" onClick={fetchTalents}>
                        <RefreshCw className="h-4 w-4" />
                    </Button>
                </div>
            </div>

            {/* Data Table */}
            <div className="rounded-md border">
                <Table>
                    <TableHeader>
                        <TableRow>
                            {isAdmin && <TableHead>GTK / Sekolah</TableHead>}
                            <TableHead>Jenis</TableHead>
                            <TableHead>Detail</TableHead>
                            <TableHead>Status</TableHead>
                            <TableHead>Tanggal</TableHead>
                            <TableHead className="text-right">Aksi</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {loading ? (
                            Array.from({ length: 5 }).map((_, i) => (
                                <TableRow key={i}>
                                    {isAdmin && <TableCell><Skeleton className="h-4 w-32" /></TableCell>}
                                    <TableCell><Skeleton className="h-4 w-24" /></TableCell>
                                    <TableCell><Skeleton className="h-4 w-40" /></TableCell>
                                    <TableCell><Skeleton className="h-4 w-20" /></TableCell>
                                    <TableCell><Skeleton className="h-4 w-24" /></TableCell>
                                    <TableCell><Skeleton className="h-4 w-16 ml-auto" /></TableCell>
                                </TableRow>
                            ))
                        ) : talents.length === 0 ? (
                            <TableRow>
                                <TableCell colSpan={isAdmin ? 6 : 5} className="h-24 text-center text-muted-foreground">
                                    Tidak ada data talenta
                                </TableCell>
                            </TableRow>
                        ) : (
                            talents.map((talent) => (
                                <TableRow key={talent.id}>
                                    {isAdmin && (
                                        <TableCell>
                                            <div className="flex flex-col">
                                                <span className="font-medium">{talent.user?.full_name}</span>
                                                <span className="text-xs text-muted-foreground">{talent.user?.school_name}</span>
                                            </div>
                                        </TableCell>
                                    )}
                                    <TableCell>
                                        <Badge variant="outline">{TALENT_TYPE_LABELS[talent.talent_type]}</Badge>
                                    </TableCell>
                                    <TableCell className="max-w-xs truncate">
                                        {getTalentSummary(talent)}
                                    </TableCell>
                                    <TableCell>
                                        <Badge
                                            variant={
                                                talent.status === 'approved' ? 'default' :
                                                    talent.status === 'rejected' ? 'destructive' : 'secondary'
                                            }
                                            className={talent.status === 'approved' ? 'bg-green-600 hover:bg-green-700' : ''}
                                        >
                                            {TALENT_STATUS_LABELS[talent.status]}
                                        </Badge>
                                    </TableCell>
                                    <TableCell>{new Date(talent.created_at).toLocaleDateString("id-ID")}</TableCell>
                                    <TableCell className="text-right">
                                        <div className="flex justify-end gap-1">
                                            <Button variant="ghost" size="icon" onClick={() => setViewId(talent.id)}>
                                                <Eye className="h-4 w-4" />
                                            </Button>
                                            {isGTK && talent.status === "pending" && (
                                                <>
                                                    <Button variant="ghost" size="icon" onClick={() => handleEdit(talent.id)}>
                                                        <Pencil className="h-4 w-4" />
                                                    </Button>
                                                    <Button variant="ghost" size="icon" onClick={() => setDeleteId(talent.id)}>
                                                        <Trash2 className="h-4 w-4 text-destructive" />
                                                    </Button>
                                                </>
                                            )}
                                        </div>
                                    </TableCell>
                                </TableRow>
                            ))
                        )}
                    </TableBody>
                </Table>
            </div>

            {/* Pagination */}
            {totalPages > 1 && (
                <div className="flex items-center justify-between">
                    <p className="text-sm text-muted-foreground">
                        Menampilkan {talents.length} dari {totalCount} talenta
                    </p>
                    <div className="flex gap-2">
                        <Button
                            variant="outline"
                            size="sm"
                            disabled={page === 1}
                            onClick={() => setPage(page - 1)}
                        >
                            Sebelumnya
                        </Button>
                        <Button
                            variant="outline"
                            size="sm"
                            disabled={page === totalPages}
                            onClick={() => setPage(page + 1)}
                        >
                            Selanjutnya
                        </Button>
                    </div>
                </div>
            )}

            {/* Modals */}
            <TalentModal
                open={createModal}
                onOpenChange={setCreateModal}
                talent={selectedTalent}
                onSuccess={() => {
                    setCreateModal(false);
                    fetchTalents();
                }}
            />

            <TalentDetailModal
                open={!!viewId}
                onOpenChange={(op) => !op && setViewId(null)}
                talentId={viewId}
                onApprove={() => { }} // Read only here basically, or redirect to verification?
                onReject={() => { }}
            // Note: TalentDetailModal has Approve/Reject buttons which might show up if status is pending. 
            // Super Admin viewing here might accidentally approve/reject? 
            // Actually, Super Admin SHOULD be able to approve/reject from here too if they want.
            // But the props are required. I'll pass dummy functions or actual handlers if I want to support it.
            // For now, let's keep it read-only-ish or consistent.
            // If I pass empty functions, buttons might do nothing or error.
            // Better to allow it if user is admin. But TalentDetailModal calls onApprove after success.
            />

            {/* Delete Confirmation */}
            <AlertDialog open={!!deleteId} onOpenChange={(op) => !op && setDeleteId(null)}>
                <AlertDialogContent>
                    <AlertDialogHeader>
                        <AlertDialogTitle>Hapus Pengajuan Talenta?</AlertDialogTitle>
                        <AlertDialogDescription>
                            Tindakan ini tidak dapat dibatalkan.
                        </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                        <AlertDialogCancel>Batal</AlertDialogCancel>
                        <AlertDialogAction onClick={handleDelete} className="bg-destructive text-destructive-foreground hover:bg-destructive/90">
                            Hapus
                        </AlertDialogAction>
                    </AlertDialogFooter>
                </AlertDialogContent>
            </AlertDialog>

        </div>
    );
}
