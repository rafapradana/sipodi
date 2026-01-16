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
import type { School } from "@/types";

interface DeleteSchoolDialogProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    school: School | null;
    onSuccess: () => void;
}

export function DeleteSchoolDialog({
    open,
    onOpenChange,
    school,
    onSuccess,
}: DeleteSchoolDialogProps) {
    const [loading, setLoading] = useState(false);

    const handleDelete = async () => {
        if (!school) return;

        setLoading(true);
        try {
            await api.delete(`/schools/${school.id}`);
            toast.success("Sekolah berhasil dihapus");
            onSuccess();
        } catch (error) {
            if (error instanceof ApiException) {
                toast.error(error.message);
            } else {
                toast.error("Gagal menghapus sekolah");
            }
        } finally {
            setLoading(false);
        }
    };

    return (
        <AlertDialog open={open} onOpenChange={onOpenChange}>
            <AlertDialogContent>
                <AlertDialogHeader>
                    <AlertDialogTitle>Hapus Sekolah?</AlertDialogTitle>
                    <AlertDialogDescription>
                        Apakah Anda yakin ingin menghapus <strong>{school?.name}</strong>?
                        Tindakan ini tidak dapat dibatalkan.
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
