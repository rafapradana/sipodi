"use client";

import { useState } from "react";
import { api, ApiException } from "@/lib/api";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
    DialogFooter,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { Loader2 } from "lucide-react";
import { toast } from "sonner";

interface RejectModalProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    talentIds: string[];
    isBatch: boolean;
    onSuccess: () => void;
}

export function RejectModal({
    open,
    onOpenChange,
    talentIds,
    isBatch,
    onSuccess,
}: RejectModalProps) {
    const [reason, setReason] = useState("");
    const [loading, setLoading] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!reason.trim()) {
            toast.error("Alasan penolakan wajib diisi");
            return;
        }

        setLoading(true);
        try {
            if (isBatch) {
                const response = await api.post<{ data: { rejected_count: number } }>(
                    "/verifications/talents/batch/reject",
                    { ids: talentIds, rejection_reason: reason }
                );
                toast.success(`${response.data.rejected_count} talenta berhasil ditolak`);
            } else {
                await api.post(`/verifications/talents/${talentIds[0]}/reject`, {
                    rejection_reason: reason,
                });
                toast.success("Talenta berhasil ditolak");
            }
            setReason("");
            onSuccess();
        } catch (error) {
            if (error instanceof ApiException) {
                toast.error(error.message);
            } else {
                toast.error("Gagal menolak talenta");
            }
        } finally {
            setLoading(false);
        }
    };

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="sm:max-w-md">
                <DialogHeader>
                    <DialogTitle>Tolak Talenta</DialogTitle>
                    <DialogDescription>
                        {isBatch
                            ? `Anda akan menolak ${talentIds.length} talenta. Berikan alasan penolakan.`
                            : "Berikan alasan penolakan untuk talenta ini."}
                    </DialogDescription>
                </DialogHeader>
                <form onSubmit={handleSubmit}>
                    <div className="grid gap-4 py-4">
                        <div className="grid gap-2">
                            <Label htmlFor="reason">Alasan Penolakan *</Label>
                            <Textarea
                                id="reason"
                                value={reason}
                                onChange={(e) => setReason(e.target.value)}
                                placeholder="Contoh: Dokumen bukti tidak valid atau tidak terbaca"
                                required
                                disabled={loading}
                                rows={4}
                            />
                        </div>
                    </div>
                    <DialogFooter>
                        <Button
                            type="button"
                            variant="outline"
                            onClick={() => onOpenChange(false)}
                            disabled={loading}
                        >
                            Batal
                        </Button>
                        <Button type="submit" variant="destructive" disabled={loading}>
                            {loading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                            Tolak
                        </Button>
                    </DialogFooter>
                </form>
            </DialogContent>
        </Dialog>
    );
}
