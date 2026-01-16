"use client";

import { useState, useEffect, useCallback } from "react";
import { api, ApiException } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
    Table,
    TableHeader,
    TableBody,
    TableRow,
    TableHead,
    TableCell,
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
import { Plus, Search, Pencil, Trash2, RefreshCw, Download } from "lucide-react";
import { toast } from "sonner";
import type { School, SchoolStatus, ListResponse, DataResponse } from "@/types";
import { SCHOOL_STATUS_LABELS } from "@/types";
import { SchoolModal } from "@/components/dashboard/sekolah/school-modal";
import { DeleteSchoolDialog } from "@/components/dashboard/sekolah/delete-school-dialog";

export default function SekolahPage() {
    const [schools, setSchools] = useState<School[]>([]);
    const [loading, setLoading] = useState(true);
    const [search, setSearch] = useState("");
    const [statusFilter, setStatusFilter] = useState<string>("all");
    const [page, setPage] = useState(1);
    const [totalPages, setTotalPages] = useState(1);
    const [totalCount, setTotalCount] = useState(0);

    // Modal states
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [editingSchool, setEditingSchool] = useState<School | null>(null);
    const [deleteDialog, setDeleteDialog] = useState<{ open: boolean; school: School | null }>({
        open: false,
        school: null,
    });

    const fetchSchools = useCallback(async () => {
        setLoading(true);
        try {
            const params: Record<string, string | number> = {
                page,
                limit: 10,
            };
            if (search) params.search = search;
            if (statusFilter !== "all") params.status = statusFilter;

            const response = await api.get<ListResponse<School>>("/schools", params);
            setSchools(response.data);
            setTotalPages(response.meta.total_pages);
            setTotalCount(response.meta.total_count);
        } catch (error) {
            console.error("Failed to fetch schools:", error);
            toast.error("Gagal memuat data sekolah");
        } finally {
            setLoading(false);
        }
    }, [page, search, statusFilter]);

    useEffect(() => {
        fetchSchools();
    }, [fetchSchools]);

    const handleSearch = (e: React.FormEvent) => {
        e.preventDefault();
        setPage(1);
        fetchSchools();
    };

    const handleCreate = () => {
        setEditingSchool(null);
        setIsModalOpen(true);
    };

    const handleEdit = (school: School) => {
        setEditingSchool(school);
        setIsModalOpen(true);
    };

    const handleDelete = (school: School) => {
        setDeleteDialog({ open: true, school });
    };

    const handleModalSuccess = () => {
        setIsModalOpen(false);
        setEditingSchool(null);
        fetchSchools();
    };

    const handleDeleteSuccess = () => {
        setDeleteDialog({ open: false, school: null });
        fetchSchools();
    };

    return (
        <div className="space-y-6">
            {/* Page Header */}
            <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
                <div>
                    <h1 className="text-2xl font-bold tracking-tight">Sekolah</h1>
                    <p className="text-muted-foreground">
                        Kelola data sekolah di wilayah Cabang Dinas Malang
                    </p>
                </div>
                <div className="flex gap-2">
                    <Button variant="outline" onClick={() => api.download('/exports/schools', 'data-sekolah.xlsx')}>
                        <Download className="mr-2 h-4 w-4" />
                        Export
                    </Button>
                    <Button onClick={handleCreate}>
                        <Plus className="mr-2 h-4 w-4" />
                        Tambah Sekolah
                    </Button>
                </div>
            </div>

            {/* Filters */}
            <div className="flex flex-col gap-4 md:flex-row md:items-center">
                <form onSubmit={handleSearch} className="flex flex-1 gap-2">
                    <div className="relative flex-1 max-w-sm">
                        <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                        <Input
                            placeholder="Cari nama atau NPSN..."
                            value={search}
                            onChange={(e) => setSearch(e.target.value)}
                            className="pl-9"
                        />
                    </div>
                    <Button type="submit" variant="secondary">
                        Cari
                    </Button>
                </form>
                <Select value={statusFilter} onValueChange={(value) => { setStatusFilter(value); setPage(1); }}>
                    <SelectTrigger className="w-40">
                        <SelectValue placeholder="Status" />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectItem value="all">Semua Status</SelectItem>
                        <SelectItem value="negeri">Negeri</SelectItem>
                        <SelectItem value="swasta">Swasta</SelectItem>
                    </SelectContent>
                </Select>
                <Button variant="outline" size="icon" onClick={fetchSchools}>
                    <RefreshCw className="h-4 w-4" />
                </Button>
            </div>

            {/* Data Table */}
            <div className="rounded-md border">
                <Table>
                    <TableHeader>
                        <TableRow>
                            <TableHead>Nama Sekolah</TableHead>
                            <TableHead>NPSN</TableHead>
                            <TableHead>Status</TableHead>
                            <TableHead>Alamat</TableHead>
                            <TableHead>Kepala Sekolah</TableHead>
                            <TableHead className="text-center">GTK</TableHead>
                            <TableHead className="text-right">Aksi</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {loading ? (
                            Array.from({ length: 5 }).map((_, i) => (
                                <TableRow key={i}>
                                    <TableCell><Skeleton className="h-4 w-40" /></TableCell>
                                    <TableCell><Skeleton className="h-4 w-20" /></TableCell>
                                    <TableCell><Skeleton className="h-4 w-16" /></TableCell>
                                    <TableCell><Skeleton className="h-4 w-48" /></TableCell>
                                    <TableCell><Skeleton className="h-4 w-32" /></TableCell>
                                    <TableCell><Skeleton className="h-4 w-8 mx-auto" /></TableCell>
                                    <TableCell><Skeleton className="h-4 w-20 ml-auto" /></TableCell>
                                </TableRow>
                            ))
                        ) : schools.length === 0 ? (
                            <TableRow>
                                <TableCell colSpan={7} className="h-24 text-center text-muted-foreground">
                                    Tidak ada data sekolah
                                </TableCell>
                            </TableRow>
                        ) : (
                            schools.map((school) => (
                                <TableRow key={school.id}>
                                    <TableCell className="font-medium">{school.name}</TableCell>
                                    <TableCell>{school.npsn}</TableCell>
                                    <TableCell>
                                        <Badge variant={school.status === "negeri" ? "default" : "secondary"}>
                                            {SCHOOL_STATUS_LABELS[school.status]}
                                        </Badge>
                                    </TableCell>
                                    <TableCell className="max-w-xs truncate">{school.address}</TableCell>
                                    <TableCell>{school.head_master?.full_name || "-"}</TableCell>
                                    <TableCell className="text-center">{school.gtk_count}</TableCell>
                                    <TableCell className="text-right">
                                        <div className="flex justify-end gap-2">
                                            <Button variant="ghost" size="icon" onClick={() => handleEdit(school)}>
                                                <Pencil className="h-4 w-4" />
                                            </Button>
                                            <Button variant="ghost" size="icon" onClick={() => handleDelete(school)}>
                                                <Trash2 className="h-4 w-4 text-destructive" />
                                            </Button>
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
                        Menampilkan {schools.length} dari {totalCount} sekolah
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
            <SchoolModal
                open={isModalOpen}
                onOpenChange={setIsModalOpen}
                school={editingSchool}
                onSuccess={handleModalSuccess}
            />
            <DeleteSchoolDialog
                open={deleteDialog.open}
                onOpenChange={(open) => setDeleteDialog({ ...deleteDialog, open })}
                school={deleteDialog.school}
                onSuccess={handleDeleteSuccess}
            />
        </div>
    );
}
