"use client";

import { useState, useEffect } from "react";
import { api, ApiException } from "@/lib/api";
import { useAuth } from "@/hooks/use-auth";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Loader2, User as UserIcon, Lock } from "lucide-react";
import { toast } from "sonner";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
    Select,
    SelectTrigger,
    SelectValue,
    SelectContent,
    SelectItem,
} from "@/components/ui/select";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
    DialogFooter,
} from "@/components/ui/dialog";
import { FileUpload } from "@/components/ui/file-upload";
import { useFileUpload } from "@/hooks/use-file-upload";
import { Pencil } from "lucide-react";
import type { DataResponse, User, Gender } from "@/types";

export default function ProfilePage() {
    const { user, refreshUser } = useAuth();
    const [loading, setLoading] = useState(false);
    const [activeTab, setActiveTab] = useState("profile");

    // Profile Form
    const [fullName, setFullName] = useState("");
    const [gender, setGender] = useState<Gender | "">("");
    const [birthDate, setBirthDate] = useState("");
    const [position, setPosition] = useState("");

    // Password Form
    const [currentPassword, setCurrentPassword] = useState("");
    const [newPassword, setNewPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");

    useEffect(() => {
        if (user) {
            setFullName(user.full_name);
            setGender(user.gender || "");
            setBirthDate(user.birth_date || "");
            setPosition(user.position || "");
        }
    }, [user]);

    const handleUpdateProfile = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        try {
            const payload: any = {
                full_name: fullName,
            };
            if (gender) payload.gender = gender;
            if (birthDate) payload.birth_date = birthDate;
            if (position) payload.position = position;

            await api.patch("/me", payload);
            toast.success("Profil berhasil diperbarui");
            refreshUser();
        } catch (error) {
            toast.error("Gagal memperbarui profil");
        } finally {
            setLoading(false);
        }
    };

    const handleChangePassword = async (e: React.FormEvent) => {
        e.preventDefault();
        if (newPassword !== confirmPassword) {
            toast.error("Konfirmasi password tidak cocok");
            return;
        }
        setLoading(true);
        try {
            await api.patch("/me/password", {
                current_password: currentPassword,
                new_password: newPassword,
                new_password_confirmation: confirmPassword,
            });
            toast.success("Password berhasil diubah");
            setCurrentPassword("");
            setNewPassword("");
            setConfirmPassword("");
        } catch (error) {
            if (error instanceof ApiException) {
                toast.error(error.message);
            } else {
                toast.error("Gagal mengubah password");
            }
        } finally {
            setLoading(false);
        }
    };

    // Photo Upload State
    const [photoModalOpen, setPhotoModalOpen] = useState(false);
    const [uploadId, setUploadId] = useState("");

    // Manual Upload
    const { upload, isUploading, progress: uploadProgress, reset: resetUpload } = useFileUpload();
    const [selectedFile, setSelectedFile] = useState<File | null>(null);

    const handleUpdatePhoto = async () => {
        if (!selectedFile) return;
        setLoading(true);
        try {
            const newUploadId = await upload(selectedFile, "profile_photo");
            await api.patch("/me/photo", { upload_id: newUploadId });
            toast.success("Foto profil berhasil diperbarui");
            setPhotoModalOpen(false);
            setUploadId("");
            setSelectedFile(null);
            refreshUser();
            resetUpload();
        } catch (error) {
            if (error instanceof ApiException) {
                toast.error(error.message);
            } else {
                toast.error("Gagal memperbarui foto profil");
            }
        } finally {
            setLoading(false);
        }
    };

    const getInitials = (name: string) => {
        return name
            .split(" ")
            .map((n) => n[0])
            .join("")
            .toUpperCase()
            .slice(0, 2);
    };

    return (
        <div className="max-w-4xl mx-auto space-y-6">
            <div className="flex items-center gap-4">
                <div className="relative group cursor-pointer" onClick={() => setPhotoModalOpen(true)}>
                    <Avatar className="h-20 w-20 group-hover:opacity-75 transition-opacity">
                        <AvatarImage src={user?.photo_url} />
                        <AvatarFallback className="text-2xl">
                            {user?.full_name ? getInitials(user.full_name) : "U"}
                        </AvatarFallback>
                    </Avatar>
                    <div className="absolute inset-0 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity bg-black/20 rounded-full">
                        <Pencil className="h-6 w-6 text-white" />
                    </div>
                </div>
                <div>
                    <h1 className="text-2xl font-bold">{user?.full_name}</h1>
                    <p className="text-muted-foreground">{user?.email}</p>
                    <Button variant="link" className="p-0 h-auto text-sm text-primary" onClick={() => setPhotoModalOpen(true)}>
                        Ganti Foto
                    </Button>
                </div>
            </div>

            <Tabs value={activeTab} onValueChange={setActiveTab}>
                <TabsList>
                    <TabsTrigger value="profile">Profil Saya</TabsTrigger>
                    <TabsTrigger value="password">Ganti Password</TabsTrigger>
                </TabsList>

                <TabsContent value="profile">
                    <div className="rounded-md border p-6">
                        <form onSubmit={handleUpdateProfile} className="space-y-4">
                            <div className="grid gap-2">
                                <Label htmlFor="fullName">Nama Lengkap</Label>
                                <Input
                                    id="fullName"
                                    value={fullName}
                                    onChange={(e) => setFullName(e.target.value)}
                                    disabled={loading}
                                />
                            </div>

                            <div className="grid grid-cols-2 gap-4">
                                <div className="grid gap-2">
                                    <Label htmlFor="gender">Jenis Kelamin</Label>
                                    <Select value={gender} onValueChange={(v) => setGender(v as Gender)} disabled={loading}>
                                        <SelectTrigger><SelectValue placeholder="Pilih" /></SelectTrigger>
                                        <SelectContent>
                                            <SelectItem value="L">Laki-laki</SelectItem>
                                            <SelectItem value="P">Perempuan</SelectItem>
                                        </SelectContent>
                                    </Select>
                                </div>
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

                            <div className="pt-4">
                                <Button type="submit" disabled={loading}>
                                    {loading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                                    Simpan Perubahan
                                </Button>
                            </div>
                        </form>
                    </div>
                </TabsContent>

                <TabsContent value="password">
                    <div className="rounded-md border p-6 max-w-md">
                        <form onSubmit={handleChangePassword} className="space-y-4">
                            <div className="grid gap-2">
                                <Label htmlFor="currentPassword">Password Saat Ini</Label>
                                <Input
                                    id="currentPassword"
                                    type="password"
                                    value={currentPassword}
                                    onChange={(e) => setCurrentPassword(e.target.value)}
                                    required
                                    disabled={loading}
                                />
                            </div>
                            <div className="grid gap-2">
                                <Label htmlFor="newPassword">Password Baru</Label>
                                <Input
                                    id="newPassword"
                                    type="password"
                                    value={newPassword}
                                    onChange={(e) => setNewPassword(e.target.value)}
                                    required
                                    disabled={loading}
                                    minLength={8}
                                />
                            </div>
                            <div className="grid gap-2">
                                <Label htmlFor="confirmPassword">Konfirmasi Password Baru</Label>
                                <Input
                                    id="confirmPassword"
                                    type="password"
                                    value={confirmPassword}
                                    onChange={(e) => setConfirmPassword(e.target.value)}
                                    required
                                    disabled={loading}
                                />
                            </div>
                            <div className="pt-4">
                                <Button type="submit" disabled={loading}>
                                    {loading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                                    Ubah Password
                                </Button>
                            </div>
                        </form>
                    </div>
                </TabsContent>
            </Tabs>

            {/* Photo Upload Modal */}
            <Dialog open={photoModalOpen} onOpenChange={setPhotoModalOpen}>
                <DialogContent className="sm:max-w-md">
                    <DialogHeader>
                        <DialogTitle>Ganti Foto Profil</DialogTitle>
                        <DialogDescription>
                            Upload foto baru untuk profil Anda. Format: JPG, PNG. Max 5MB.
                        </DialogDescription>
                    </DialogHeader>
                    <div className="grid gap-4 py-4">
                        <FileUpload
                            uploadType="profile_photo"
                            accept=".jpg,.jpeg,.png"
                            value={""} // Always empty for new upload
                            manualUpload={true}
                            onFileChange={setSelectedFile}
                            onRemove={() => {
                                setSelectedFile(null);
                                setUploadId("");
                            }}
                            progress={uploadProgress}
                            isUploading={isUploading}
                            disabled={loading || isUploading}
                        />
                    </div>
                    <DialogFooter>
                        <Button variant="outline" onClick={() => setPhotoModalOpen(false)} disabled={loading}>
                            Batal
                        </Button>
                        <Button onClick={handleUpdatePhoto} disabled={!selectedFile || loading || isUploading}>
                            {(loading || isUploading) && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                            Simpan Foto
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
        </div>
    );
}
