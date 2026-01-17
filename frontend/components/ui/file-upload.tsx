"use client";

import * as React from "react";
import { useCallback, useState } from "react";
import { UploadCloud, File, X, Check, Loader2, AlertCircle } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Progress } from "@/components/ui/progress";
import { cn } from "@/lib/utils";
import { api, ApiException } from "@/lib/api";
import type { DataResponse, PresignResponse } from "@/types";
import { toast } from "sonner";

interface FileUploadProps {
    accept?: string;
    maxSize?: number; // in bytes
    value?: string; // URL of existing file
    onChange?: (uploadId: string) => void;
    onFileChange?: (file: File | null) => void; // New: used for manual upload
    onRemove?: () => void;
    className?: string;
    disabled?: boolean;
    uploadType: "profile_photo" | "talent_certificate";
    manualUpload?: boolean;
    progress?: number;
    isUploading?: boolean;
}

export function FileUpload({
    accept = "*",
    maxSize = 5 * 1024 * 1024, // 5MB default
    value,
    onChange,
    onFileChange,
    onRemove,
    className,
    disabled = false,
    uploadType,
    manualUpload = false,
    progress: externalProgress,
    isUploading: externalIsUploading,
}: FileUploadProps) {
    const [internalFile, setInternalFile] = useState<File | null>(null);
    const [internalUploading, setInternalUploading] = useState(false);
    const [internalProgress, setInternalProgress] = useState(0);
    const [uploadId, setUploadId] = useState<string | null>(null);
    const [internalError, setInternalError] = useState<string | null>(null);

    // Derived states
    const uploading = manualUpload ? externalIsUploading : internalUploading;
    const progress = manualUpload ? (externalProgress || 0) : internalProgress;
    const error = internalError;

    // Preview logic
    const [previewUrl, setPreviewUrl] = useState<string>("");

    // Effect to handle preview URL
    React.useEffect(() => {
        if (value) {
            setPreviewUrl(value);
            // If value exists (edit mode), we might want to clear internal file if it was set
        } else if (internalFile) {
            const url = URL.createObjectURL(internalFile);
            setPreviewUrl(url);
            return () => URL.revokeObjectURL(url);
        } else {
            setPreviewUrl("");
        }
    }, [value, internalFile]);


    const formatFileSize = (bytes: number) => {
        if (bytes === 0) return "0 Bytes";
        const k = 1024;
        const sizes = ["Bytes", "KB", "MB", "GB", "TB"];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
    };

    const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (e.target.files && e.target.files.length > 0) {
            const selectedFile = e.target.files[0];

            // Validate size
            if (selectedFile.size > maxSize) {
                setInternalError(`File terlalu besar. Maksimal ${formatFileSize(maxSize)}`);
                return;
            }

            setInternalError(null);
            setInternalFile(selectedFile);

            if (manualUpload) {
                if (onFileChange) onFileChange(selectedFile);
            } else {
                handleUpload(selectedFile);
            }
        }
    };

    const handleUpload = async (fileToUpload: File) => {
        setInternalUploading(true);
        setInternalProgress(0);
        setInternalError(null);

        try {
            // 1. Get presigned URL
            const presignRes = await api.post<DataResponse<PresignResponse>>("/uploads/presign", {
                filename: fileToUpload.name,
                size: fileToUpload.size,
                content_type: fileToUpload.type,
                upload_type: uploadType,
            });

            const { upload_id, presigned_url } = presignRes.data;

            // 2. Upload to MinIO
            const xhr = new XMLHttpRequest();
            xhr.open("PUT", presigned_url, true);
            xhr.setRequestHeader("Content-Type", fileToUpload.type);

            xhr.upload.onprogress = (e) => {
                if (e.lengthComputable) {
                    const percentComplete = (e.loaded / e.total) * 100;
                    setInternalProgress(Math.round(percentComplete)); // Use Math.round
                }
            };

            xhr.onload = async () => {
                if (xhr.status >= 200 && xhr.status < 300) {
                    // 3. Confirm upload
                    try {
                        await api.post(`/uploads/${upload_id}/confirm`);
                        setUploadId(upload_id);
                        if (onChange) onChange(upload_id);
                        toast.success("File berhasil diupload");
                    } catch (err) {
                        console.error("Confirmation error:", err);
                        setInternalError("Gagal konfirmasi upload");
                        toast.error("Gagal konfirmasi upload");
                    } finally {
                        setInternalUploading(false);
                    }
                } else {
                    setInternalError("Gagal upload file");
                    setInternalUploading(false);
                    toast.error("Gagal upload file ke server storage");
                }
            };

            xhr.onerror = () => {
                setInternalError("Network error saat upload");
                setInternalUploading(false);
                toast.error("Network error saat upload");
            };

            xhr.send(fileToUpload);

        } catch (err) {
            console.error("Upload error:", err);
            if (err instanceof ApiException) {
                setInternalError(err.message);
            } else {
                setInternalError("Gagal memulai upload");
            }
            setInternalUploading(false);
            toast.error("Gagal memulai upload");
        }
    };

    const handleRemove = () => {
        setInternalFile(null);
        setUploadId(null);
        setInternalProgress(0);
        setInternalError(null);
        if (manualUpload && onFileChange) onFileChange(null);
        if (onRemove) onRemove();
        if (onChange) onChange("");
    };

    return (
        <div className={cn("space-y-3", className)}>
            {!previewUrl ? (
                <div className={cn(
                    "relative flex flex-col items-center justify-center rounded-lg border border-dashed border-muted-foreground/25 px-6 py-10 text-center transition hover:bg-muted/50",
                    disabled && "cursor-not-allowed opacity-60 hover:bg-transparent"
                )}>
                    <div className="flex flex-col items-center gap-2">
                        <div className="rounded-full bg-muted/50 p-2">
                            <UploadCloud className="h-6 w-6 text-muted-foreground" />
                        </div>
                        <div className="text-sm">
                            <label
                                htmlFor="file-upload"
                                className={cn(
                                    "relative cursor-pointer font-semibold text-primary focus-within:outline-none focus-within:ring-2 focus-within:ring-primary focus-within:ring-offset-2 hover:text-primary/80",
                                    disabled && "pointer-events-none"
                                )}
                            >
                                <span>Pilih file</span>
                                <input
                                    id="file-upload"
                                    name="file-upload"
                                    type="file"
                                    className="sr-only"
                                    accept={accept}
                                    onChange={handleFileSelect}
                                    disabled={disabled}
                                />
                            </label>
                            <span className="text-muted-foreground"> atau drag and drop</span>
                        </div>
                        <p className="text-xs text-muted-foreground">
                            Maksimal {formatFileSize(maxSize)}
                        </p>
                    </div>
                </div>
            ) : (
                <div className="rounded-md border p-3">
                    <div className="flex items-center justify-between">
                        <div className="flex items-center gap-3 overflow-hidden">
                            {/* Preview Thumbnail */}
                            {accept.includes("image") || previewUrl.match(/\.(jpg|jpeg|png)$/i) ? (
                                <img
                                    src={previewUrl}
                                    alt="Preview"
                                    className="h-10 w-10 rounded object-cover border"
                                />
                            ) : (
                                <div className="rounded-full bg-muted p-2">
                                    <File className="h-4 w-4" />
                                    end
                            )}

                                    <div className="grid gap-0.5">
                                        <p className="text-sm font-medium truncate max-w-[200px]">{internalFile?.name || "File Tersimpan"}</p>
                                        <p className="text-xs text-muted-foreground">{internalFile ? formatFileSize(internalFile.size) : "Siap"}</p>
                                    </div>
                                </div>
                        
                        {!uploading && !disabled && (
                                <Button
                                    type="button"
                                    variant="ghost"
                                    size="icon"
                                    className="text-muted-foreground hover:text-destructive"
                                    onClick={handleRemove}
                                >
                                    <X className="h-4 w-4" />
                                </Button>
                            )}
                        </div>

                        {(uploading || progress > 0) && (
                            <div className="mt-3 space-y-1">
                                <div className="flex justify-between text-xs">
                                    <span>{uploading ? "Mengupload..." : (manualUpload && !uploadId ? "Menunggu Submit" : "Selesai")}</span>
                                    <span>{Math.round(progress)}%</span>
                                </div>
                                <Progress value={progress} className="h-1" />
                            </div>
                        )}

                        {error && (
                            <div className="mt-2 text-xs text-destructive flex items-center gap-1">
                                <AlertCircle className="h-3 w-3" />
                                {error}
                            </div>
                        )}
                    </div>
            )}
                </div>
            );
}
