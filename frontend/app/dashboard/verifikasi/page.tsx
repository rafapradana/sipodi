"use client";

import { useState, useEffect, useCallback } from "react";
import { api, ApiException } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Checkbox } from "@/components/ui/checkbox";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import {
    Search,
    RefreshCw,
    CheckCircle,
    XCircle,
    Eye,
    FileText,
} from "lucide-react";
import { toast } from "sonner";
import type { Talent, TalentListItem, ListResponse, TalentStatus } from "@/types";
import { TALENT_TYPE_LABELS, TALENT_STATUS_LABELS } from "@/types";
import { TalentDetailModal } from "@/components/dashboard/verifikasi/talent-detail-modal";
import { RejectModal } from "@/components/dashboard/verifikasi/reject-modal";

export default function VerifikasiPage() {
    const [talents, setTalents] = useState<TalentListItem[]>([]);
    const [loading, setLoading] = useState(true);
    const [search, setSearch] = useState("");
    const [statusTab, setStatusTab] = useState<TalentStatus>("pending");
    const [page, setPage] = useState(1);
    const [totalPages, setTotalPages] = useState(1);
    const [totalCount, setTotalCount] = useState(0);

    // Selection for batch operations
    const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set());

    // Modal states
    const [detailModal, setDetailModal] = useState<{ open: boolean; talentId: string | null }>({
        open: false,
        talentId: null,
    });
    const [rejectModal, setRejectModal] = useState<{
        open: boolean;
        talentIds: string[];
        isBatch: boolean;
    }>({
        open: false,
        talentIds: [],
        isBatch: false,
    });

    const fetchTalents = useCallback(async () => {
        setLoading(true);
        try {
            const params: Record<string, string | number> = {
                page,
                limit: 10,
                status: statusTab,
            };
            if (search) params.search = search;

            const response = await api.get<ListResponse<TalentListItem>>("/verifications/talents", params);
            setTalents(response.data || []);
            setTotalPages(response.meta.total_pages);
            setTotalCount(response.meta.total_count);
            setSelectedIds(new Set());
        } catch (error) {
            console.error("Failed to fetch talents:", error);
            toast.error("Gagal memuat data talenta");
        } finally {
            setLoading(false);
        }
    }, [page, statusTab, search]);

    useEffect(() => {
        fetchTalents();
    }, [fetchTalents]);

    useEffect(() => {
        setPage(1);
        setSelectedIds(new Set());
    }, [statusTab]);

    const handleSearch = (e: React.FormEvent) => {
        e.preventDefault();
        setPage(1);
        fetchTalents();
    };

    const handleSelectAll = (checked: boolean) => {
        if (checked) {
            setSelectedIds(new Set(talents.map((t) => t.id)));
        } else {
            setSelectedIds(new Set());
        }
    };

    const handleSelect = (id: string, checked: boolean) => {
        const newSelected = new Set(selectedIds);
        if (checked) {
            newSelected.add(id);
        } else {
            newSelected.delete(id);
        }
        setSelectedIds(newSelected);
    };

    const handleViewDetail = (talent: TalentListItem) => {
        setDetailModal({ open: true, talentId: talent.id });
    };

    const handleApprove = async (id: string) => {
        try {
            await api.post(`/verifications/talents/${id}/approve`);
            toast.success("Talenta berhasil disetujui");
            fetchTalents();
        } catch (error) {
            if (error instanceof ApiException) {
                toast.error(error.message);
            } else {
                toast.error("Gagal menyetujui talenta");
            }
        }
    };

    const handleReject = (id: string) => {
        setRejectModal({ open: true, talentIds: [id], isBatch: false });
    };

    const handleBatchApprove = async () => {
        if (selectedIds.size === 0) return;

        try {
            const response = await api.post<{ data: { approved_count: number; failed_count: number } }>(
                "/verifications/talents/batch/approve",
                { ids: Array.from(selectedIds) }
            );
            toast.success(`${response.data.approved_count} talenta berhasil disetujui`);
            fetchTalents();
        } catch (error) {
            if (error instanceof ApiException) {
                toast.error(error.message);
            } else {
                toast.error("Gagal menyetujui talenta");
            }
        }
    };

    const handleBatchReject = () => {
        if (selectedIds.size === 0) return;
        setRejectModal({ open: true, talentIds: Array.from(selectedIds), isBatch: true });
    };

    const handleRejectSuccess = () => {
        setRejectModal({ open: false, talentIds: [], isBatch: false });
        fetchTalents();
    };

    const getTalentSummary = (talent: TalentListItem): string => {
        const detail = talent.detail;
        if (!detail) return "-";
        if ("activity_name" in detail) return detail.activity_name;
        if ("competition_name" in detail) return detail.competition_name;
        if ("interest_name" in detail) return detail.interest_name;
        return "-";
    };

    return (
        <div className="space-y-6">
            {/* Page Header */}
            <div>
                <h1 className="text-2xl font-bold tracking-tight">Verifikasi Talenta</h1>
                <p className="text-muted-foreground">
                    Tinjau dan verifikasi data talenta yang diajukan GTK
                </p>
            </div>

            {/* Tabs */}
            <Tabs value={statusTab} onValueChange={(v) => setStatusTab(v as TalentStatus)}>
                <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
                    <TabsList>
                        <TabsTrigger value="pending" className="gap-2">
                            <span className="hidden sm:inline">Pending</span>
                            <Badge variant="secondary" className="ml-1">
                                {statusTab === "pending" ? totalCount : "..."}
                            </Badge>
                        </TabsTrigger>
                        <TabsTrigger value="approved">Disetujui</TabsTrigger>
                        <TabsTrigger value="rejected">Ditolak</TabsTrigger>
                    </TabsList>

                    {/* Batch Actions */}
                    {statusTab === "pending" && selectedIds.size > 0 && (
                        <div className="flex gap-2">
                            <Button size="sm" onClick={handleBatchApprove}>
                                <CheckCircle className="mr-2 h-4 w-4" />
                                Setujui ({selectedIds.size})
                            </Button>
                            <Button size="sm" variant="destructive" onClick={handleBatchReject}>
                                <XCircle className="mr-2 h-4 w-4" />
                                Tolak ({selectedIds.size})
                            </Button>
                        </div>
                    )}
                </div>

                {/* Filters */}
                <div className="mt-4 flex flex-col gap-4 md:flex-row md:items-center">
                    <form onSubmit={handleSearch} className="flex flex-1 gap-2">
                        <div className="relative flex-1 max-w-sm">
                            <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                            <Input
                                placeholder="Cari nama GTK..."
                                value={search}
                                onChange={(e) => setSearch(e.target.value)}
                                className="pl-9"
                            />
                        </div>
                        <Button type="submit" variant="secondary">
                            Cari
                        </Button>
                    </form>
                    <Button variant="outline" size="icon" onClick={fetchTalents}>
                        <RefreshCw className="h-4 w-4" />
                    </Button>
                </div>

                {/* Content */}
                <TabsContent value={statusTab} className="mt-4">
                    <div className="rounded-md border">
                        <Table>
                            <TableHeader>
                                <TableRow>
                                    {statusTab === "pending" && (
                                        <TableHead className="w-12">
                                            <Checkbox
                                                checked={talents.length > 0 && selectedIds.size === talents.length}
                                                onCheckedChange={handleSelectAll}
                                            />
                                        </TableHead>
                                    )}
                                    <TableHead>GTK</TableHead>
                                    <TableHead>Sekolah</TableHead>
                                    <TableHead>Jenis Talenta</TableHead>
                                    <TableHead>Detail</TableHead>
                                    <TableHead>Tanggal</TableHead>
                                    <TableHead className="text-right">Aksi</TableHead>
                                </TableRow>
                            </TableHeader>
                            <TableBody>
                                {loading ? (
                                    Array.from({ length: 5 }).map((_, i) => (
                                        <TableRow key={i}>
                                            {statusTab === "pending" && (
                                                <TableCell>
                                                    <Skeleton className="h-4 w-4" />
                                                </TableCell>
                                            )}
                                            <TableCell>
                                                <Skeleton className="h-4 w-32" />
                                            </TableCell>
                                            <TableCell>
                                                <Skeleton className="h-4 w-28" />
                                            </TableCell>
                                            <TableCell>
                                                <Skeleton className="h-4 w-24" />
                                            </TableCell>
                                            <TableCell>
                                                <Skeleton className="h-4 w-40" />
                                            </TableCell>
                                            <TableCell>
                                                <Skeleton className="h-4 w-20" />
                                            </TableCell>
                                            <TableCell>
                                                <Skeleton className="h-4 w-16 ml-auto" />
                                            </TableCell>
                                        </TableRow>
                                    ))
                                ) : talents.length === 0 ? (
                                    <TableRow>
                                        <TableCell
                                            colSpan={statusTab === "pending" ? 7 : 6}
                                            className="h-24 text-center text-muted-foreground"
                                        >
                                            Tidak ada data talenta {TALENT_STATUS_LABELS[statusTab].toLowerCase()}
                                        </TableCell>
                                    </TableRow>
                                ) : (
                                    talents.map((talent) => (
                                        <TableRow key={talent.id}>
                                            {statusTab === "pending" && (
                                                <TableCell>
                                                    <Checkbox
                                                        checked={selectedIds.has(talent.id)}
                                                        onCheckedChange={(checked) => handleSelect(talent.id, !!checked)}
                                                    />
                                                </TableCell>
                                            )}
                                            <TableCell className="font-medium">
                                                {talent.user?.full_name || "-"}
                                            </TableCell>
                                            <TableCell>{talent.user?.school_name || "-"}</TableCell>
                                            <TableCell>
                                                <Badge variant="outline">
                                                    {TALENT_TYPE_LABELS[talent.talent_type]}
                                                </Badge>
                                            </TableCell>
                                            <TableCell className="max-w-xs truncate">
                                                {getTalentSummary(talent)}
                                            </TableCell>
                                            <TableCell>
                                                {new Date(talent.created_at).toLocaleDateString("id-ID")}
                                            </TableCell>
                                            <TableCell className="text-right">
                                                <div className="flex justify-end gap-1">
                                                    <Button
                                                        variant="ghost"
                                                        size="icon"
                                                        onClick={() => handleViewDetail(talent)}
                                                        title="Lihat Detail"
                                                    >
                                                        <Eye className="h-4 w-4" />
                                                    </Button>
                                                    {statusTab === "pending" && (
                                                        <>
                                                            <Button
                                                                variant="ghost"
                                                                size="icon"
                                                                onClick={() => handleApprove(talent.id)}
                                                                title="Setujui"
                                                            >
                                                                <CheckCircle className="h-4 w-4 text-green-600" />
                                                            </Button>
                                                            <Button
                                                                variant="ghost"
                                                                size="icon"
                                                                onClick={() => handleReject(talent.id)}
                                                                title="Tolak"
                                                            >
                                                                <XCircle className="h-4 w-4 text-red-600" />
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
                        <div className="mt-4 flex items-center justify-between">
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
                </TabsContent>
            </Tabs>

            {/* Modals */}
            <TalentDetailModal
                open={detailModal.open}
                onOpenChange={(open) => setDetailModal({ ...detailModal, open })}
                talentId={detailModal.talentId}
                onApprove={() => {
                    setDetailModal({ open: false, talentId: null });
                    fetchTalents();
                }}
                onReject={() => {
                    setDetailModal({ open: false, talentId: null });
                    setRejectModal({ open: true, talentIds: [detailModal.talentId!], isBatch: false });
                }}
            />
            <RejectModal
                open={rejectModal.open}
                onOpenChange={(open) => setRejectModal({ ...rejectModal, open })}
                talentIds={rejectModal.talentIds}
                isBatch={rejectModal.isBatch}
                onSuccess={handleRejectSuccess}
            />
        </div>
    );
}
