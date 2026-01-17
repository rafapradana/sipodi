import { useState } from "react";
import { api, ApiException } from "@/lib/api";
import type { DataResponse, PresignResponse } from "@/types";

export type UploadType = "profile_photo" | "talent_certificate";

interface UseFileUploadReturn {
    upload: (file: File, type: UploadType) => Promise<string>;
    isUploading: boolean;
    progress: number;
    error: string | null;
    reset: () => void;
}

export function useFileUpload(): UseFileUploadReturn {
    const [isUploading, setIsUploading] = useState(false);
    const [progress, setProgress] = useState(0);
    const [error, setError] = useState<string | null>(null);

    const reset = () => {
        setIsUploading(false);
        setProgress(0);
        setError(null);
    };

    const upload = async (file: File, type: UploadType): Promise<string> => {
        setIsUploading(true);
        setProgress(0);
        setError(null);

        try {
            // 1. Get presigned URL
            const presignRes = await api.post<DataResponse<PresignResponse>>("/uploads/presign", {
                filename: file.name,
                size: file.size,
                content_type: file.type,
                upload_type: type,
            });

            const { upload_id, presigned_url } = presignRes.data;

            // 2. Upload to MinIO
            await new Promise<void>((resolve, reject) => {
                const xhr = new XMLHttpRequest();
                xhr.open("PUT", presigned_url, true);
                xhr.setRequestHeader("Content-Type", file.type);

                xhr.upload.onprogress = (e) => {
                    if (e.lengthComputable) {
                        const percentComplete = (e.loaded / e.total) * 100;
                        setProgress(Math.round(percentComplete));
                    }
                };

                xhr.onload = () => {
                    if (xhr.status >= 200 && xhr.status < 300) {
                        resolve();
                    } else {
                        reject(new Error("Upload failed"));
                    }
                };

                xhr.onerror = () => reject(new Error("Network Error"));
                xhr.send(file);
            });

            // 3. Confirm upload
            await api.post(`/uploads/${upload_id}/confirm`);

            return upload_id;
        } catch (err: any) {
            const errorMessage = err instanceof ApiException ? err.message : (err.message || "Gagal mengupload file");
            setError(errorMessage);
            throw new Error(errorMessage);
        } finally {
            setIsUploading(false);
        }
    };

    return { upload, isUploading, progress, error, reset };
}
