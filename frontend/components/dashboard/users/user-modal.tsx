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
import {
    Select,
    SelectTrigger,
    SelectValue,
    SelectContent,
    SelectItem,
} from "@/components/ui/select";
ChevronsUpDown
} from "lucide-react";
import { toast } from "sonner";
import { FileUpload } from "@/components/ui/file-upload";
import { useFileUpload } from "@/hooks/use-file-upload";
import type { UserListItem, School, DataResponse, User, CreateUserRequest, UpdateUserRequest, UserRole, Gender, GTKType } from "@/types";

interface UserModalProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    user: UserListItem | null; // or any if type issue persists, but UserListItem is better
    schools: School[];
    onSuccess: () => void;
}

export function UserModal({ open, onOpenChange, user, schools, onSuccess }: UserModalProps) {
    const isEdit = !!user;
    const [loading, setLoading] = useState(false);
    const [errors, setErrors] = useState<Record<string, string>>({});

    // Form States
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [role, setRole] = useState<UserRole>("gtk");
    const [fullName, setFullName] = useState("");
    const [nuptk, setNuptk] = useState("");
    const [nip, setNip] = useState("");
    const [gender, setGender] = useState<Gender | "">("");
    const [birthDate, setBirthDate] = useState("");
    const [gtkType, setGtkType] = useState<GTKType | "">("");
    const [position, setPosition] = useState("");
    const [schoolId, setSchoolId] = useState("");
    const [uploadId, setUploadId] = useState("");

    // File Upload (Manual)
    const { upload, isUploading, progress: uploadProgress, reset: resetUpload } = useFileUpload();
    const [selectedFile, setSelectedFile] = useState<File | null>(null);

    // Populate form when editing
    useEffect(() => {
        if (user) {
            setEmail(user.email);
            setRole(user.role);
            setFullName(user.full_name);
            setNuptk(user.nuptk || "");
            setNip(user.nip || "");
            setGender(user.gender || "");
            setBirthDate(user.birth_date || "");
            setGtkType(user.gtk_type || "");
            setPosition(user.position || "");
            setSchoolId(user.school?.id || "");
            setUploadId("");
        } else {
            setEmail("");
            setPassword("");
            setRole("gtk");
            setFullName("");
            setNuptk("");
            setNip("");
            setGender("");
            setBirthDate("");
            setGtkType("");
            setPosition("");
            setSchoolId("");
            setSchoolId("");
            setUploadId("");
        }
        setErrors({});
        setSelectedFile(null);
        resetUpload();
    }, [user, open]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setErrors({});
        setLoading(true);

        try {
            if (isEdit) {
                const data: UpdateUserRequest = {
                    full_name: fullName,
                };
                if (nuptk) data.nuptk = nuptk;
                if (nip) data.nip = nip;
                if (gender) data.gender = gender as Gender;
                if (birthDate) data.birth_date = birthDate;
                if (gtkType) data.gtk_type = gtkType as GTKType;
                if (position) data.position = position;
                if (schoolId) data.school_id = schoolId;
                if (schoolId) data.school_id = schoolId;

                // Upload file if selected
                if (selectedFile) {
                    try {
                        const newUploadId = await upload(selectedFile, "profile_photo");
                        data.upload_id = newUploadId;
                    } catch (err) {
                        setLoading(false);
                        return;
                    }
                } else if (uploadId) {
                    data.upload_id = uploadId;
                }

                await api.put<DataResponse<User>>(`/users/${user.id}`, data);
                toast.success("User berhasil diperbarui");
            } else {
                const data: CreateUserRequest = {
                    email,
                    password,
                    role,
                    full_name: fullName,
                };
                if (nuptk) data.nuptk = nuptk;
                if (nip) data.nip = nip;
                if (gender) data.gender = gender as Gender;
                if (birthDate) data.birth_date = birthDate;
                if (gtkType) data.gtk_type = gtkType as GTKType;
                if (position) data.position = position;
                if (schoolId) data.school_id = schoolId;
                if (position) data.position = position;
                if (schoolId) data.school_id = schoolId;

                // Upload file if selected
                if (selectedFile) {
                    try {
                        const newUploadId = await upload(selectedFile, "profile_photo");
                        data.upload_id = newUploadId;
                    } catch (err) {
                        setLoading(false);
                        return;
                    }
                } else if (uploadId) {
                    data.upload_id = uploadId;
                }

                await api.post<DataResponse<User>>("/users", data);
                toast.success("User berhasil ditambahkan");
            }
            onSuccess();
        } catch (error) {
            // ... error handling
        } finally {
            setLoading(false);
        }
    };

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="sm:max-w-xl max-h-[90vh] overflow-y-auto">
                <DialogHeader>
                    <DialogTitle>{isEdit ? "Edit User" : "Tambah User"}</DialogTitle>
                    <DialogDescription>
                        {isEdit ? "Perbarui informasi user" : "Tambahkan user baru ke sistem"}
                    </DialogDescription>
                </DialogHeader>
                <form onSubmit={handleSubmit}>
                    <div className="grid gap-4 py-4">
                        {/* Photo Upload */}
                        <div className="grid gap-2">
                            <Label>Foto Profil</Label>
                            <FileUpload
                                uploadType="profile_photo"
                                accept="image/*"
                                value={user?.photo_url || ""} // Show existing URL
                                manualUpload={true}
                                onFileChange={setSelectedFile}
                                onRemove={() => {
                                    setSelectedFile(null);
                                    setUploadId("");
                                }}
                                progress={uploadProgress}
                                isUploading={isUploading}
                                disabled={loading}
                            />
                        </div>

                        {/* Email */}
                        <div className="grid gap-2">
                            <Label htmlFor="email">Email *</Label>
                            <Input
                                id="email"
                                type="email"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                                placeholder="user@sekolah.sch.id"
                                required
                                disabled={loading || isEdit}
                            />
                            {errors.email && <p className="text-sm text-destructive">{errors.email}</p>}
                        </div>

                        {/* Password (only for create) */}
                        {!isEdit && (
                            <div className="grid gap-2">
                                <Label htmlFor="password">Password *</Label>
                                <Input
                                    id="password"
                                    type="password"
                                    value={password}
                                    onChange={(e) => setPassword(e.target.value)}
                                    placeholder="Minimal 8 karakter"
                                    required
                                    disabled={loading}
                                    minLength={8}
                                />
                                {errors.password && <p className="text-sm text-destructive">{errors.password}</p>}
                            </div>
                        )}

                        {/* Role */}
                        <div className="grid gap-2">
                            <Label htmlFor="role">Role *</Label>
                            <Select
                                value={role}
                                onValueChange={(v) => setRole(v as UserRole)}
                                disabled={loading || isEdit}
                            >
                                <SelectTrigger>
                                    <SelectValue placeholder="Pilih role" />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectItem value="super_admin">Super Admin</SelectItem>
                                    <SelectItem value="admin_sekolah">Admin Sekolah</SelectItem>
                                    <SelectItem value="gtk">GTK</SelectItem>
                                </SelectContent>
                            </Select>
                            {errors.role && <p className="text-sm text-destructive">{errors.role}</p>}
                        </div>

                        {/* Full Name */}
                        <div className="grid gap-2">
                            <Label htmlFor="fullName">Nama Lengkap *</Label>
                            <Input
                                id="fullName"
                                value={fullName}
                                onChange={(e) => setFullName(e.target.value)}
                                placeholder="Nama lengkap"
                                required
                                disabled={loading}
                            />
                            {errors.full_name && <p className="text-sm text-destructive">{errors.full_name}</p>}
                        </div>

                        {/* Two column layout */}
                        <div className="grid grid-cols-2 gap-4">
                            {/* NUPTK */}
                            <div className="grid gap-2">
                                <Label htmlFor="nuptk">NUPTK</Label>
                                <Input
                                    id="nuptk"
                                    value={nuptk}
                                    onChange={(e) => setNuptk(e.target.value)}
                                    placeholder="16 digit"
                                    disabled={loading}
                                />
                                {errors.nuptk && <p className="text-sm text-destructive">{errors.nuptk}</p>}
                            </div>

                            {/* NIP */}
                            <div className="grid gap-2">
                                <Label htmlFor="nip">NIP</Label>
                                <Input
                                    id="nip"
                                    value={nip}
                                    onChange={(e) => setNip(e.target.value)}
                                    placeholder="18 digit"
                                    disabled={loading}
                                />
                                {errors.nip && <p className="text-sm text-destructive">{errors.nip}</p>}
                            </div>
                        </div>

                        <div className="grid grid-cols-2 gap-4">
                            {/* Gender */}
                            <div className="grid gap-2">
                                <Label htmlFor="gender">Jenis Kelamin</Label>
                                <Select value={gender} onValueChange={(v) => setGender(v as Gender)} disabled={loading}>
                                    <SelectTrigger>
                                        <SelectValue placeholder="Pilih" />
                                    </SelectTrigger>
                                    <SelectContent>
                                        <SelectItem value="L">Laki-laki</SelectItem>
                                        <SelectItem value="P">Perempuan</SelectItem>
                                    </SelectContent>
                                </Select>
                            </div>

                            {/* Birth Date */}
                            <div className="grid gap-2">
                                <Label htmlFor="birthDate">Tanggal Lahir</Label>
                                <Input
                                    id="birthDate"
                                    type="date"
                                    value={birthDate}
                                    onChange={(e) => setBirthDate(e.target.value)}
                                    disabled={loading}
                                />
                            </div>
                        </div>

                        {/* GTK Type */}
                        {(role === "gtk" || role === "admin_sekolah") && (
                            <div className="grid gap-2">
                                <Label htmlFor="gtkType">Jenis GTK</Label>
                                <Select value={gtkType} onValueChange={(v) => setGtkType(v as GTKType)} disabled={loading}>
                                    <SelectTrigger>
                                        <SelectValue placeholder="Pilih jenis GTK" />
                                    </SelectTrigger>
                                    <SelectContent>
                                        <SelectItem value="guru">Guru</SelectItem>
                                        <SelectItem value="tendik">Tenaga Kependidikan</SelectItem>
                                        <SelectItem value="kepala_sekolah">Kepala Sekolah</SelectItem>
                                    </SelectContent>
                                </Select>
                            </div>
                        )}

                        {/* Position */}
                        <div className="grid gap-2">
                            <Label htmlFor="position">Jabatan</Label>
                            <Input
                                id="position"
                                value={position}
                                onChange={(e) => setPosition(e.target.value)}
                                placeholder="Contoh: Guru Matematika"
                                disabled={loading}
                            />
                        </div>

                        {/* School */}
                        {role !== "super_admin" && (
                            <div className="grid gap-2">
                                <Label htmlFor="school">Sekolah</Label>
                                <Select value={schoolId} onValueChange={setSchoolId} disabled={loading}>
                                    <SelectTrigger>
                                        <SelectValue placeholder="Pilih sekolah" />
                                    </SelectTrigger>
                                    <SelectContent>
                                        {(schools || []).map((school) => (
                                            <SelectItem key={school.id} value={school.id}>
                                                {school.name}
                                            </SelectItem>
                                        ))}
                                    </SelectContent>
                                </Select>
                            </div>
                        )}
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
