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
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
    Select,
    SelectTrigger,
    SelectValue,
    SelectContent,
    SelectItem,
} from "@/components/ui/select";
import { Loader2 } from "lucide-react";
import { toast } from "sonner";
import { FileUpload } from "@/components/ui/file-upload";
import type {
    Talent,
    TalentType,
    DataResponse,
    CompetitionLevel,
    TalentField,
    TrainingDetail,
    MentorDetail,
    ParticipantDetail,
    InterestDetail,
} from "@/types";
import {
    TALENT_TYPE_LABELS,
    COMPETITION_LEVEL_LABELS,
    TALENT_FIELD_LABELS,
} from "@/types";

interface TalentModalProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    talent: Talent | null; // If null, create mode
    onSuccess: () => void;
}

export function TalentModal({ open, onOpenChange, talent, onSuccess }: TalentModalProps) {
    const isEdit = !!talent;
    const [loading, setLoading] = useState(false);
    const [talentType, setTalentType] = useState<TalentType>("peserta_pelatihan");
    const [uploadId, setUploadId] = useState<string>("");
    const [fileUrl, setFileUrl] = useState<string>("");

    // Common Form Fields
    const [activityName, setActivityName] = useState("");
    const [organizer, setOrganizer] = useState("");
    const [startDate, setStartDate] = useState("");
    const [durationDays, setDurationDays] = useState(1);
    const [competitionName, setCompetitionName] = useState("");
    const [level, setLevel] = useState<CompetitionLevel>("kota");
    const [field, setField] = useState<TalentField>("akademik");
    const [achievement, setAchievement] = useState("");
    const [competitionField, setCompetitionField] = useState("");
    const [interestName, setInterestName] = useState("");
    const [description, setDescription] = useState("");

    useEffect(() => {
        if (talent) {
            setTalentType(talent.talent_type);
            setFileUrl(talent.certificate_url || "");
            // Reset uploadId as we have url, unless user changes it
            setUploadId("");

            const detail = talent.detail;
            if (talent.talent_type === "peserta_pelatihan") {
                const d = detail as TrainingDetail;
                setActivityName(d.activity_name);
                setOrganizer(d.organizer);
                setStartDate(d.start_date);
                setDurationDays(d.duration_days);
            } else if (talent.talent_type === "pembimbing_lomba") {
                const d = detail as MentorDetail;
                setCompetitionName(d.competition_name);
                setLevel(d.level);
                setOrganizer(d.organizer);
                setField(d.field);
                setAchievement(d.achievement);
            } else if (talent.talent_type === "peserta_lomba") {
                const d = detail as ParticipantDetail;
                setCompetitionName(d.competition_name);
                setLevel(d.level);
                setOrganizer(d.organizer);
                setField(d.field);
                setStartDate(d.start_date);
                setDurationDays(d.duration_days);
                setCompetitionField(d.competition_field);
                setAchievement(d.achievement);
            } else if (talent.talent_type === "minat_bakat") {
                const d = detail as InterestDetail;
                setInterestName(d.interest_name);
                setDescription(d.description);
            }
        } else {
            // Reset defaults
            setTalentType("peserta_pelatihan");
            setFileUrl("");
            setUploadId("");
            setActivityName("");
            setOrganizer("");
            setStartDate("");
            setDurationDays(1);
            setCompetitionName("");
            setLevel("kota");
            setField("akademik");
            setAchievement("");
            setCompetitionField("");
            setInterestName("");
            setDescription("");
        }
    }, [talent, open]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);

        try {
            let detail: any = {};
            if (talentType === "peserta_pelatihan") {
                detail = { activity_name: activityName, organizer, start_date: startDate, duration_days: Number(durationDays) };
            } else if (talentType === "pembimbing_lomba") {
                detail = { competition_name: competitionName, level, organizer, field, achievement };
            } else if (talentType === "peserta_lomba") {
                detail = { competition_name: competitionName, level, organizer, field, start_date: startDate, duration_days: Number(durationDays), competition_field: competitionField, achievement };
            } else if (talentType === "minat_bakat") {
                detail = { interest_name: interestName, description };
            }

            const payload: any = {
                talent_type: talentType,
                detail,
            };

            if (uploadId) {
                payload.upload_id = uploadId;
            }

            if (isEdit) {
                await api.put(`/me/talents/${talent.id}`, payload);
                toast.success("Talenta berhasil diperbarui");
            } else {
                await api.post("/me/talents", payload);
                toast.success("Talenta berhasil ditambahkan");
            }
            onSuccess();
        } catch (error) {
            if (error instanceof ApiException) {
                toast.error(error.message);
            } else {
                toast.error("Gagal menyimpan talenta");
            }
        } finally {
            setLoading(false);
        }
    };

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="sm:max-w-xl max-h-[90vh] overflow-y-auto">
                <DialogHeader>
                    <DialogTitle>{isEdit ? "Edit Talenta" : "Tambah Talenta Baru"}</DialogTitle>
                    <DialogDescription>
                        Isi formulir berikut untuk mengajukan data talenta baru.
                    </DialogDescription>
                </DialogHeader>

                <form onSubmit={handleSubmit} className="space-y-4 py-4">
                    <div className="grid gap-2">
                        <Label htmlFor="type">Jenis Talenta</Label>
                        <Select
                            value={talentType}
                            onValueChange={(v) => !isEdit && setTalentType(v as TalentType)}
                            disabled={loading || isEdit}
                        >
                            <SelectTrigger>
                                <SelectValue placeholder="Pilih jenis talenta" />
                            </SelectTrigger>
                            <SelectContent>
                                {Object.entries(TALENT_TYPE_LABELS).map(([key, label]) => (
                                    <SelectItem key={key} value={key}>{label}</SelectItem>
                                ))}
                            </SelectContent>
                        </Select>
                    </div>

                    {/* Type Specific Fields */}
                    {talentType === "peserta_pelatihan" && (
                        <>
                            <div className="grid gap-2">
                                <Label htmlFor="activityName">Nama Kegiatan</Label>
                                <Input id="activityName" value={activityName} onChange={(e) => setActivityName(e.target.value)} required disabled={loading} />
                            </div>
                            <div className="grid gap-2">
                                <Label htmlFor="organizer">Penyelenggara</Label>
                                <Input id="organizer" value={organizer} onChange={(e) => setOrganizer(e.target.value)} required disabled={loading} />
                            </div>
                            <div className="grid grid-cols-2 gap-4">
                                <div className="grid gap-2">
                                    <Label htmlFor="startDate">Tanggal Mulai</Label>
                                    <Input id="startDate" type="date" value={startDate} onChange={(e) => setStartDate(e.target.value)} required disabled={loading} />
                                </div>
                                <div className="grid gap-2">
                                    <Label htmlFor="durationDays">Durasi (Hari)</Label>
                                    <Input id="durationDays" type="number" min="1" value={durationDays} onChange={(e) => setDurationDays(Number(e.target.value))} required disabled={loading} />
                                </div>
                            </div>
                        </>
                    )}

                    {talentType === "pembimbing_lomba" && (
                        <>
                            <div className="grid gap-2">
                                <Label htmlFor="competitionName">Nama Lomba</Label>
                                <Input id="competitionName" value={competitionName} onChange={(e) => setCompetitionName(e.target.value)} required disabled={loading} />
                            </div>
                            <div className="grid grid-cols-2 gap-4">
                                <div className="grid gap-2">
                                    <Label htmlFor="level">Jenjang</Label>
                                    <Select value={level} onValueChange={(v) => setLevel(v as CompetitionLevel)} disabled={loading}>
                                        <SelectTrigger><SelectValue /></SelectTrigger>
                                        <SelectContent>
                                            {Object.entries(COMPETITION_LEVEL_LABELS).map(([key, label]) => (
                                                <SelectItem key={key} value={key}>{label}</SelectItem>
                                            ))}
                                        </SelectContent>
                                    </Select>
                                </div>
                                <div className="grid gap-2">
                                    <Label htmlFor="field">Bidang</Label>
                                    <Select value={field} onValueChange={(v) => setField(v as TalentField)} disabled={loading}>
                                        <SelectTrigger><SelectValue /></SelectTrigger>
                                        <SelectContent>
                                            {Object.entries(TALENT_FIELD_LABELS).map(([key, label]) => (
                                                <SelectItem key={key} value={key}>{label}</SelectItem>
                                            ))}
                                        </SelectContent>
                                    </Select>
                                </div>
                            </div>
                            <div className="grid gap-2">
                                <Label htmlFor="organizer">Penyelenggara</Label>
                                <Input id="organizer" value={organizer} onChange={(e) => setOrganizer(e.target.value)} required disabled={loading} />
                            </div>
                            <div className="grid gap-2">
                                <Label htmlFor="achievement">Prestasi</Label>
                                <Input id="achievement" value={achievement} onChange={(e) => setAchievement(e.target.value)} required disabled={loading} placeholder="Contoh: Juara 1" />
                            </div>
                        </>
                    )}

                    {talentType === "peserta_lomba" && (
                        <>
                            <div className="grid gap-2">
                                <Label htmlFor="competitionName">Nama Lomba</Label>
                                <Input id="competitionName" value={competitionName} onChange={(e) => setCompetitionName(e.target.value)} required disabled={loading} />
                            </div>
                            <div className="grid grid-cols-2 gap-4">
                                <div className="grid gap-2">
                                    <Label htmlFor="level">Jenjang</Label>
                                    <Select value={level} onValueChange={(v) => setLevel(v as CompetitionLevel)} disabled={loading}>
                                        <SelectTrigger><SelectValue /></SelectTrigger>
                                        <SelectContent>
                                            {Object.entries(COMPETITION_LEVEL_LABELS).map(([key, label]) => (
                                                <SelectItem key={key} value={key}>{label}</SelectItem>
                                            ))}
                                        </SelectContent>
                                    </Select>
                                </div>
                                <div className="grid gap-2">
                                    <Label htmlFor="field">Bidang</Label>
                                    <Select value={field} onValueChange={(v) => setField(v as TalentField)} disabled={loading}>
                                        <SelectTrigger><SelectValue /></SelectTrigger>
                                        <SelectContent>
                                            {Object.entries(TALENT_FIELD_LABELS).map(([key, label]) => (
                                                <SelectItem key={key} value={key}>{label}</SelectItem>
                                            ))}
                                        </SelectContent>
                                    </Select>
                                </div>
                            </div>
                            <div className="grid gap-2">
                                <Label htmlFor="organizer">Penyelenggara</Label>
                                <Input id="organizer" value={organizer} onChange={(e) => setOrganizer(e.target.value)} required disabled={loading} />
                            </div>
                            <div className="grid grid-cols-2 gap-4">
                                <div className="grid gap-2">
                                    <Label htmlFor="startDate">Tanggal Mulai</Label>
                                    <Input id="startDate" type="date" value={startDate} onChange={(e) => setStartDate(e.target.value)} required disabled={loading} />
                                </div>
                                <div className="grid gap-2">
                                    <Label htmlFor="durationDays">Durasi (Hari)</Label>
                                    <Input id="durationDays" type="number" min="1" value={durationDays} onChange={(e) => setDurationDays(Number(e.target.value))} required disabled={loading} />
                                </div>
                            </div>
                            <div className="grid gap-2">
                                <Label htmlFor="competitionField">Bidang Lomba</Label>
                                <Input id="competitionField" value={competitionField} onChange={(e) => setCompetitionField(e.target.value)} required disabled={loading} placeholder="Contoh: Guru Matematika" />
                            </div>
                            <div className="grid gap-2">
                                <Label htmlFor="achievement">Prestasi</Label>
                                <Input id="achievement" value={achievement} onChange={(e) => setAchievement(e.target.value)} required disabled={loading} placeholder="Contoh: Juara 2" />
                            </div>
                        </>
                    )}

                    {talentType === "minat_bakat" && (
                        <>
                            <div className="grid gap-2">
                                <Label htmlFor="interestName">Nama Minat/Bakat</Label>
                                <Input id="interestName" value={interestName} onChange={(e) => setInterestName(e.target.value)} required disabled={loading} />
                            </div>
                            <div className="grid gap-2">
                                <Label htmlFor="description">Deskripsi</Label>
                                <Textarea id="description" value={description} onChange={(e) => setDescription(e.target.value)} required disabled={loading} rows={3} placeholder="Jelaskan detail minat/bakat" />
                            </div>
                        </>
                    )}

                    {/* Upload Certificate */}
                    <div className="grid gap-2">
                        <Label>Bukti / Sertifikat (PDF/Image, Max 5MB)</Label>
                        <FileUpload
                            accept=".pdf,.jpg,.jpeg,.png"
                            value={fileUrl}
                            onChange={setUploadId}
                            disabled={loading}
                            onRemove={() => {
                                setFileUrl("");
                                setUploadId("");
                            }}
                        />
                    </div>

                    <DialogFooter>
                        <Button type="button" variant="outline" onClick={() => onOpenChange(false)} disabled={loading}>
                            Batal
                        </Button>
                        <Button type="submit" disabled={loading}>
                            {loading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                            {isEdit ? "Simpan Perubahan" : "Ajukan Talenta"}
                        </Button>
                    </DialogFooter>
                </form>
            </DialogContent>
        </Dialog>
    );
}
