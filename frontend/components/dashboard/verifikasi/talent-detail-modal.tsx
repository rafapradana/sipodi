"use client";

import { useState, useEffect } from "react";
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
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { Separator } from "@/components/ui/separator";
import { CheckCircle, XCircle, FileText, ExternalLink } from "lucide-react";
import { toast } from "sonner";
import type { Talent, DataResponse, TrainingDetail, MentorDetail, ParticipantDetail, InterestDetail } from "@/types";
import {
    TALENT_TYPE_LABELS,
    TALENT_STATUS_LABELS,
    COMPETITION_LEVEL_LABELS,
    TALENT_FIELD_LABELS,
} from "@/types";

interface TalentDetailModalProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    talentId: string | null;
    onApprove: () => void;
    onReject: () => void;
}

export function TalentDetailModal({
    open,
    onOpenChange,
    talentId,
    onApprove,
    onReject,
}: TalentDetailModalProps) {
    const [talent, setTalent] = useState<Talent | null>(null);
    const [loading, setLoading] = useState(false);
    const [approving, setApproving] = useState(false);

    useEffect(() => {
        if (open && talentId) {
            fetchTalent();
        } else {
            setTalent(null);
        }
    }, [open, talentId]);

    const fetchTalent = async () => {
        if (!talentId) return;
        setLoading(true);
        try {
            const response = await api.get<DataResponse<Talent>>(`/talents/${talentId}`);
            setTalent(response.data);
        } catch (error) {
            console.error("Failed to fetch talent:", error);
            toast.error("Gagal memuat detail talenta");
            onOpenChange(false);
        } finally {
            setLoading(false);
        }
    };

    const handleApprove = async () => {
        if (!talentId) return;
        setApproving(true);
        try {
            await api.post(`/verifications/talents/${talentId}/approve`);
            toast.success("Talenta berhasil disetujui");
            onApprove();
        } catch (error) {
            if (error instanceof ApiException) {
                toast.error(error.message);
            } else {
                toast.error("Gagal menyetujui talenta");
            }
        } finally {
            setApproving(false);
        }
    };

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

    const renderDetail = () => {
        if (!talent) return null;

        const detail = talent.detail;

        switch (talent.talent_type) {
            case "peserta_pelatihan":
                const training = detail as TrainingDetail;
                return (
                    <div className="grid gap-3">
                        <DetailRow label="Nama Kegiatan" value={training.activity_name} />
                        <DetailRow label="Penyelenggara" value={training.organizer} />
                        <DetailRow
                            label="Tanggal Mulai"
                            value={new Date(training.start_date).toLocaleDateString("id-ID")}
                        />
                        <DetailRow label="Jangka Waktu" value={`${training.duration_days} hari`} />
                    </div>
                );

            case "pembimbing_lomba":
                const mentor = detail as MentorDetail;
                return (
                    <div className="grid gap-3">
                        <DetailRow label="Nama Lomba" value={mentor.competition_name} />
                        <DetailRow label="Jenjang" value={COMPETITION_LEVEL_LABELS[mentor.level]} />
                        <DetailRow label="Penyelenggara" value={mentor.organizer} />
                        <DetailRow label="Bidang" value={TALENT_FIELD_LABELS[mentor.field]} />
                        <DetailRow label="Prestasi" value={mentor.achievement} />
                    </div>
                );

            case "peserta_lomba":
                const participant = detail as ParticipantDetail;
                return (
                    <div className="grid gap-3">
                        <DetailRow label="Nama Lomba" value={participant.competition_name} />
                        <DetailRow label="Jenjang" value={COMPETITION_LEVEL_LABELS[participant.level]} />
                        <DetailRow label="Penyelenggara" value={participant.organizer} />
                        <DetailRow label="Bidang" value={TALENT_FIELD_LABELS[participant.field]} />
                        <DetailRow
                            label="Tanggal Mulai"
                            value={new Date(participant.start_date).toLocaleDateString("id-ID")}
                        />
                        <DetailRow label="Jangka Waktu" value={`${participant.duration_days} hari`} />
                        <DetailRow label="Bidang Lomba" value={participant.competition_field} />
                        <DetailRow label="Prestasi" value={participant.achievement} />
                    </div>
                );

            case "minat_bakat":
                const interest = detail as InterestDetail;
                return (
                    <div className="grid gap-3">
                        <DetailRow label="Nama Minat/Bakat" value={interest.interest_name} />
                        <DetailRow label="Deskripsi" value={interest.description} />
                    </div>
                );

            default:
                return null;
        }
    };

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="sm:max-w-lg max-h-[90vh] overflow-y-auto">
                <DialogHeader>
                    <DialogTitle>Detail Talenta</DialogTitle>
                    <DialogDescription>Informasi lengkap talenta GTK</DialogDescription>
                </DialogHeader>

                {loading ? (
                    <div className="space-y-4 py-4">
                        <Skeleton className="h-4 w-3/4" />
                        <Skeleton className="h-4 w-1/2" />
                        <Skeleton className="h-4 w-2/3" />
                        <Skeleton className="h-4 w-1/2" />
                    </div>
                ) : talent ? (
                    <div className="space-y-4 py-4">
                        {/* User Info */}
                        <div className="rounded-lg bg-muted/50 p-4">
                            <p className="font-medium">{talent.user?.full_name || "Unknown"}</p>
                            <p className="text-sm text-muted-foreground">
                                {talent.user?.school_name || "Tidak ada sekolah"}
                            </p>
                        </div>

                        {/* Status & Type */}
                        <div className="flex items-center gap-2">
                            {getStatusBadge(talent.status)}
                            <Badge variant="outline">{TALENT_TYPE_LABELS[talent.talent_type]}</Badge>
                        </div>

                        <Separator />

                        {/* Detail */}
                        {renderDetail()}

                        {/* Certificate */}
                        {talent.certificate_url && (
                            <>
                                <Separator />
                                <div className="flex items-center justify-between">
                                    <div className="flex items-center gap-2">
                                        <FileText className="h-4 w-4 text-muted-foreground" />
                                        <span className="text-sm">Bukti/Sertifikat</span>
                                    </div>
                                    <Button variant="outline" size="sm" asChild>
                                        <a href={talent.certificate_url} target="_blank" rel="noopener noreferrer">
                                            <ExternalLink className="mr-2 h-4 w-4" />
                                            Lihat
                                        </a>
                                    </Button>
                                </div>
                            </>
                        )}

                        {/* Rejection Reason */}
                        {talent.status === "rejected" && talent.rejection_reason && (
                            <>
                                <Separator />
                                <div className="rounded-lg bg-destructive/10 p-3">
                                    <p className="text-sm font-medium text-destructive">Alasan Penolakan:</p>
                                    <p className="text-sm text-destructive">{talent.rejection_reason}</p>
                                </div>
                            </>
                        )}

                        {/* Timestamp */}
                        <div className="text-xs text-muted-foreground">
                            Dibuat: {new Date(talent.created_at).toLocaleString("id-ID")}
                            {talent.verified_at && (
                                <>
                                    <br />
                                    Diverifikasi: {new Date(talent.verified_at).toLocaleString("id-ID")}
                                    {talent.verified_by && ` oleh ${talent.verified_by.full_name}`}
                                </>
                            )}
                        </div>
                    </div>
                ) : null}

                {talent?.status === "pending" && (
                    <DialogFooter>
                        <Button variant="outline" onClick={onReject}>
                            <XCircle className="mr-2 h-4 w-4" />
                            Tolak
                        </Button>
                        <Button onClick={handleApprove} disabled={approving}>
                            <CheckCircle className="mr-2 h-4 w-4" />
                            Setujui
                        </Button>
                    </DialogFooter>
                )}
            </DialogContent>
        </Dialog>
    );
}

function DetailRow({ label, value }: { label: string; value: string }) {
    return (
        <div className="grid grid-cols-3 gap-2">
            <span className="text-sm text-muted-foreground">{label}</span>
            <span className="col-span-2 text-sm font-medium">{value}</span>
        </div>
    );
}
