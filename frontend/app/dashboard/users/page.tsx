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
import { Plus, Search, Pencil, Trash2, RefreshCw, UserCheck, UserX, Download } from "lucide-react";
import { toast } from "sonner";
import type { UserListItem, ListResponse, School } from "@/types";
import { USER_ROLE_LABELS, GTK_TYPE_LABELS } from "@/types";
import { UserModal } from "@/components/dashboard/users/user-modal";
import { DeleteUserDialog } from "@/components/dashboard/users/delete-user-dialog";

export default function UsersPage() {
    const [users, setUsers] = useState<UserListItem[]>([]);
    const [schools, setSchools] = useState<School[]>([]);
    const [loading, setLoading] = useState(true);
    const [search, setSearch] = useState("");
    const [roleFilter, setRoleFilter] = useState<string>("all");
    const [schoolFilter, setSchoolFilter] = useState<string>("all");
    const [statusFilter, setStatusFilter] = useState<string>("all");
    const [page, setPage] = useState(1);
    const [totalPages, setTotalPages] = useState(1);
    const [totalCount, setTotalCount] = useState(0);

    // Modal states
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [editingUser, setEditingUser] = useState<any | null>(null); // Use any temporarily to avoid type issues during switch, or better import User
    const [deleteDialog, setDeleteDialog] = useState<{ open: boolean; user: UserListItem | null }>({
        open: false,
        user: null,
    });

    // Fetch schools for filter
    useEffect(() => {
        api.get<ListResponse<School>>("/schools", { limit: 200 })
            .then((res) => setSchools(res.data || []))
            .catch(() => { });
    }, []);

    const fetchUsers = useCallback(async () => {
        setLoading(true);
        try {
            const params: Record<string, string | number> = {
                page,
                limit: 10,
            };
            if (search) params.search = search;
            if (roleFilter !== "all") params.role = roleFilter;
            if (schoolFilter !== "all") params.school_id = schoolFilter;
            if (statusFilter !== "all") params.is_active = statusFilter === "active" ? "true" : "false";

            const response = await api.get<ListResponse<UserListItem>>("/users", params);
            setUsers(response.data || []);
            setTotalPages(response.meta.total_pages);
            setTotalCount(response.meta.total_count);
        } catch (error) {
            console.error("Failed to fetch users:", error);
            toast.error("Gagal memuat data users");
        } finally {
            setLoading(false);
        }
    }, [page, search, roleFilter, schoolFilter, statusFilter]);

    useEffect(() => {
        fetchUsers();
    }, [fetchUsers]);

    const handleSearch = (e: React.FormEvent) => {
        e.preventDefault();
        setPage(1);
        fetchUsers();
    };

    const handleCreate = () => {
        setEditingUser(null);
        setIsModalOpen(true);
    };

    const handleEdit = async (user: UserListItem) => {
        try {
            const toastId = toast.loading("Memuat detail user...");
            // Use DataResponse<User> but for now using any to bypass strict type check if import User is not set
            // Wait, I should import User.
            const res = await api.get<any>(`/users/${user.id}`);
            toast.dismiss(toastId);
            setEditingUser(res.data);
            setIsModalOpen(true);
        } catch (error) {
            toast.error("Gagal memuat details user");
        }
    };

    const handleDelete = (user: UserListItem) => {
        setDeleteDialog({ open: true, user });
    };

    const handleToggleActive = async (user: UserListItem) => {
        try {
            if (user.is_active) {
                await api.patch(`/users/${user.id}/deactivate`);
                toast.success("User berhasil dinonaktifkan");
            } else {
                await api.patch(`/users/${user.id}/activate`);
                toast.success("User berhasil diaktifkan");
            }
            fetchUsers();
        } catch (error) {
            if (error instanceof ApiException) {
                toast.error(error.message);
            } else {
                toast.error("Gagal mengubah status user");
            }
        }
    };

    const handleModalSuccess = () => {
        setIsModalOpen(false);
        setEditingUser(null);
        fetchUsers();
    };

    const handleDeleteSuccess = () => {
        setDeleteDialog({ open: false, user: null });
        fetchUsers();
    };

    return (
        <div className="space-y-6">
            {/* Page Header */}
            <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
                <div>
                    <h1 className="text-2xl font-bold tracking-tight">Users</h1>
                    <p className="text-muted-foreground">
                        Kelola data pengguna sistem SIPODI
                    </p>
                </div>
                <div className="flex gap-2">
                    <Button variant="outline" onClick={() => api.download('/exports/gtk', 'data-gtk.xlsx')}>
                        <Download className="mr-2 h-4 w-4" />
                        Export GTK
                    </Button>
                    <Button onClick={handleCreate}>
                        <Plus className="mr-2 h-4 w-4" />
                        Tambah User
                    </Button>
                </div>
            </div>

            {/* Filters */}
            <div className="flex flex-col gap-4 lg:flex-row lg:items-center">
                <form onSubmit={handleSearch} className="flex flex-1 gap-2">
                    <div className="relative flex-1 max-w-sm">
                        <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                        <Input
                            placeholder="Cari nama, email, NUPTK..."
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
                    <Select value={roleFilter} onValueChange={(value) => { setRoleFilter(value); setPage(1); }}>
                        <SelectTrigger className="w-36">
                            <SelectValue placeholder="Role" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectItem value="all">Semua Role</SelectItem>
                            <SelectItem value="super_admin">Super Admin</SelectItem>
                            <SelectItem value="admin_sekolah">Admin Sekolah</SelectItem>
                            <SelectItem value="gtk">GTK</SelectItem>
                        </SelectContent>
                    </Select>

                    <Select value={schoolFilter} onValueChange={(value) => { setSchoolFilter(value); setPage(1); }}>
                        <SelectTrigger className="w-40">
                            <SelectValue placeholder="Sekolah" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectItem value="all">Semua Sekolah</SelectItem>
                            {schools.map((school) => (
                                <SelectItem key={school.id} value={school.id}>
                                    {school.name}
                                </SelectItem>
                            ))}
                        </SelectContent>
                    </Select>

                    <Select value={statusFilter} onValueChange={(value) => { setStatusFilter(value); setPage(1); }}>
                        <SelectTrigger className="w-32">
                            <SelectValue placeholder="Status" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectItem value="all">Semua</SelectItem>
                            <SelectItem value="active">Aktif</SelectItem>
                            <SelectItem value="inactive">Nonaktif</SelectItem>
                        </SelectContent>
                    </Select>

                    <Button variant="outline" size="icon" onClick={fetchUsers}>
                        <RefreshCw className="h-4 w-4" />
                    </Button>
                </div>
            </div>

            {/* Data Table */}
            <div className="rounded-md border">
                <Table>
                    <TableHeader>
                        <TableRow>
                            <TableHead>Nama</TableHead>
                            <TableHead>Email</TableHead>
                            <TableHead>Role</TableHead>
                            <TableHead>Sekolah</TableHead>
                            <TableHead>Jenis GTK</TableHead>
                            <TableHead>Status</TableHead>
                            <TableHead className="text-right">Aksi</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {loading ? (
                            Array.from({ length: 5 }).map((_, i) => (
                                <TableRow key={i}>
                                    <TableCell><Skeleton className="h-4 w-32" /></TableCell>
                                    <TableCell><Skeleton className="h-4 w-40" /></TableCell>
                                    <TableCell><Skeleton className="h-4 w-24" /></TableCell>
                                    <TableCell><Skeleton className="h-4 w-32" /></TableCell>
                                    <TableCell><Skeleton className="h-4 w-20" /></TableCell>
                                    <TableCell><Skeleton className="h-4 w-16" /></TableCell>
                                    <TableCell><Skeleton className="h-4 w-24 ml-auto" /></TableCell>
                                </TableRow>
                            ))
                        ) : users.length === 0 ? (
                            <TableRow>
                                <TableCell colSpan={7} className="h-24 text-center text-muted-foreground">
                                    Tidak ada data user
                                </TableCell>
                            </TableRow>
                        ) : (
                            users.map((user) => (
                                <TableRow key={user.id}>
                                    <TableCell className="font-medium">{user.full_name}</TableCell>
                                    <TableCell>{user.email}</TableCell>
                                    <TableCell>
                                        <Badge variant="outline">{USER_ROLE_LABELS[user.role]}</Badge>
                                    </TableCell>
                                    <TableCell>{user.school?.name || "-"}</TableCell>
                                    <TableCell>
                                        {user.gtk_type ? GTK_TYPE_LABELS[user.gtk_type] : "-"}
                                    </TableCell>
                                    <TableCell>
                                        <Badge variant={user.is_active ? "default" : "secondary"}>
                                            {user.is_active ? "Aktif" : "Nonaktif"}
                                        </Badge>
                                    </TableCell>
                                    <TableCell className="text-right">
                                        <div className="flex justify-end gap-1">
                                            <Button
                                                variant="ghost"
                                                size="icon"
                                                onClick={() => handleToggleActive(user)}
                                                title={user.is_active ? "Nonaktifkan" : "Aktifkan"}
                                            >
                                                {user.is_active ? (
                                                    <UserX className="h-4 w-4 text-orange-500" />
                                                ) : (
                                                    <UserCheck className="h-4 w-4 text-green-500" />
                                                )}
                                            </Button>
                                            <Button variant="ghost" size="icon" onClick={() => handleEdit(user)}>
                                                <Pencil className="h-4 w-4" />
                                            </Button>
                                            <Button variant="ghost" size="icon" onClick={() => handleDelete(user)}>
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
                        Menampilkan {users.length} dari {totalCount} users
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
            <UserModal
                open={isModalOpen}
                onOpenChange={setIsModalOpen}
                user={editingUser}
                schools={schools}
                onSuccess={handleModalSuccess}
            />
            <DeleteUserDialog
                open={deleteDialog.open}
                onOpenChange={(open) => setDeleteDialog({ ...deleteDialog, open })}
                user={deleteDialog.user}
                onSuccess={handleDeleteSuccess}
            />
        </div>
    );
}
