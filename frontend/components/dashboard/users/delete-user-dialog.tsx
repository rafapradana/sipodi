"use client";

import { useState } from "react";
import { api, ApiException } from "@/lib/api";
import {
    AlertDialog,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
import { Loader2 } from "lucide-react";
import { toast } from "sonner";
import type { UserListItem } from "@/types";

interface DeleteUserDialogProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    user: UserListItem | null;
    onSuccess: () => void;
}

export function DeleteUserDialog({
    open,
    onOpenChange,
    user,
    onSuccess,
}: DeleteUserDialogProps) {
    const [loading, setLoading] = useState(false);

    const handleDelete = async () => {
        if (!user) return;

        setLoading(true);
        try {
            await api.delete(`/users/${user.id}`);
            toast.success("User berhasil dihapus");
            onSuccess();
        } catch (error) {
            if (error instanceof ApiException) {
                toast.error(error.message);
            } else {
                toast.error("Gagal menghapus user");
            }
        } finally {
            setLoading(false);
        }
    };

    return (
        <AlertDialog open={open} onOpenChange={onOpenChange}>
            <AlertDialogContent>
                <AlertDialogHeader>
                    <AlertDialogTitle>Hapus User?</AlertDialogTitle>
                    <AlertDialogDescription>
                        Apakah Anda yakin ingin menghapus <strong>{user?.full_name}</strong>?
                        Tindakan ini akan menonaktifkan akun user.
                    </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                    <Button variant="outline" onClick={() => onOpenChange(false)} disabled={loading}>
                        Batal
                    </Button>
                    <Button variant="destructive" onClick={handleDelete} disabled={loading}>
                        {loading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                        Hapus
                    </Button>
                </AlertDialogFooter>
            </AlertDialogContent>
        </AlertDialog>
    );
}
