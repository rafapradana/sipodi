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
import type { School, CreateSchoolRequest, UpdateSchoolRequest, SchoolStatus, DataResponse } from "@/types";

interface SchoolModalProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    school: School | null;
    onSuccess: () => void;
}

export function SchoolModal({ open, onOpenChange, school, onSuccess }: SchoolModalProps) {
    const isEdit = !!school;
    const [loading, setLoading] = useState(false);
    const [errors, setErrors] = useState<Record<string, string>>({});

    // Form state
    const [name, setName] = useState("");
    const [npsn, setNpsn] = useState("");
    const [status, setStatus] = useState<SchoolStatus>("negeri");
    const [address, setAddress] = useState("");

    // Populate form when editing
    useEffect(() => {
        if (school) {
            setName(school.name);
            setNpsn(school.npsn);
            setStatus(school.status);
            setAddress(school.address);
        } else {
            setName("");
            setNpsn("");
            setStatus("negeri");
            setAddress("");
        }
        setErrors({});
    }, [school, open]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setErrors({});
        setLoading(true);

        try {
            if (isEdit) {
                const data: UpdateSchoolRequest = { name, npsn, status, address };
                await api.put<DataResponse<School>>(`/schools/${school.id}`, data);
                toast.success("Sekolah berhasil diperbarui");
            } else {
                const data: CreateSchoolRequest = { name, npsn, status, address };
                await api.post<DataResponse<School>>("/schools", data);
                toast.success("Sekolah berhasil ditambahkan");
            }
            onSuccess();
        } catch (error) {
            if (error instanceof ApiException) {
                if (error.details) {
                    const fieldErrors: Record<string, string> = {};
                    error.details.forEach((d) => {
                        fieldErrors[d.field] = d.message;
                    });
                    setErrors(fieldErrors);
                } else {
                    toast.error(error.message);
                }
            } else {
                toast.error("Terjadi kesalahan");
            }
        } finally {
            setLoading(false);
        }
    };

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="sm:max-w-lg">
                <DialogHeader>
                    <DialogTitle>{isEdit ? "Edit Sekolah" : "Tambah Sekolah"}</DialogTitle>
                    <DialogDescription>
                        {isEdit
                            ? "Perbarui informasi sekolah"
                            : "Tambahkan sekolah baru ke sistem"}
                    </DialogDescription>
                </DialogHeader>
                <form onSubmit={handleSubmit}>
                    <div className="grid gap-4 py-4">
                        {/* Name */}
                        <div className="grid gap-2">
                            <Label htmlFor="name">Nama Sekolah *</Label>
                            <Input
                                id="name"
                                value={name}
                                onChange={(e) => setName(e.target.value)}
                                placeholder="SMAN 1 Malang"
                                required
                                disabled={loading}
                            />
                            {errors.name && <p className="text-sm text-destructive">{errors.name}</p>}
                        </div>

                        {/* NPSN */}
                        <div className="grid gap-2">
                            <Label htmlFor="npsn">NPSN *</Label>
                            <Input
                                id="npsn"
                                value={npsn}
                                onChange={(e) => setNpsn(e.target.value)}
                                placeholder="20518765"
                                required
                                disabled={loading}
                            />
                            {errors.npsn && <p className="text-sm text-destructive">{errors.npsn}</p>}
                        </div>

                        {/* Status */}
                        <div className="grid gap-2">
                            <Label htmlFor="status">Status *</Label>
                            <Select value={status} onValueChange={(v) => setStatus(v as SchoolStatus)} disabled={loading}>
                                <SelectTrigger>
                                    <SelectValue placeholder="Pilih status" />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectItem value="negeri">Negeri</SelectItem>
                                    <SelectItem value="swasta">Swasta</SelectItem>
                                </SelectContent>
                            </Select>
                            {errors.status && <p className="text-sm text-destructive">{errors.status}</p>}
                        </div>

                        {/* Address */}
                        <div className="grid gap-2">
                            <Label htmlFor="address">Alamat *</Label>
                            <Textarea
                                id="address"
                                value={address}
                                onChange={(e) => setAddress(e.target.value)}
                                placeholder="Jl. Tugu No. 1, Malang"
                                required
                                disabled={loading}
                                rows={3}
                            />
                            {errors.address && <p className="text-sm text-destructive">{errors.address}</p>}
                        </div>
                    </div>

                    <DialogFooter>
                        <Button type="button" variant="outline" onClick={() => onOpenChange(false)} disabled={loading}>
                            Batal
                        </Button>
                        <Button type="submit" disabled={loading}>
                            {loading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                            {isEdit ? "Simpan" : "Tambah"}
                        </Button>
                    </DialogFooter>
                </form>
            </DialogContent>
        </Dialog>
    );
}
