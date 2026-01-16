# SIPODI API Documentation

Base URL: `https://sipodi.cabdinmalang.go.id/api/v1`

## Overview

API untuk SIPODI - Sistem Informasi Potensi Diri GTK Cabang Dinas Pendidikan Wilayah Malang.

### Authentication

Semua endpoint yang memerlukan autentikasi menggunakan JWT Bearer Token:

```
Authorization: Bearer <access_token>
```

### Response Format

Semua response menggunakan format JSON dengan struktur konsisten:

**Success Response:**
```json
{
  "data": { ... },
  "meta": { ... }
}
```

**Error Response:**
```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable message",
    "details": [ ... ]
  }
}
```

### User Roles

| Role | Deskripsi |
|------|-----------|
| `super_admin` | Cabang Dinas - kelola semua data |
| `admin_sekolah` | Admin Sekolah - kelola GTK di sekolahnya |
| `gtk` | Guru/Tendik - kelola data diri & talenta |

---

## Table of Contents

1. [Authentication](#1-authentication)
2. [Profile (Me)](#2-profile-me)
3. [Sekolah](#3-sekolah)
4. [Users/GTK](#4-usersgtk)
5. [Talenta](#5-talenta)
6. [Verifikasi Talenta](#6-verifikasi-talenta)
7. [Notifikasi](#7-notifikasi)
8. [File Upload (MinIO)](#8-file-upload-minio)
9. [Dashboard & Statistik](#9-dashboard--statistik)


---

## 1. Authentication

### POST /auth/login

Login user dan dapatkan access token.

**Authentication:** None

**Request Body:**
```json
{
  "email": "guru@sekolah.sch.id",
  "password": "securepassword123"
}
```

**Success Response (200):**
```json
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 900,
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "guru@sekolah.sch.id",
      "full_name": "Budi Santoso",
      "role": "gtk",
      "photo_url": "https://cdn.sipodi.go.id/photos/budi.jpg"
    }
  }
}
```

**Response Headers:**
```
Set-Cookie: refresh_token=abc123...; HttpOnly; Secure; SameSite=Strict; Path=/api/v1/auth; Max-Age=604800
```

**Error Responses:**

401 Unauthorized - Kredensial salah:
```json
{
  "error": {
    "code": "INVALID_CREDENTIALS",
    "message": "Email atau password salah"
  }
}
```

403 Forbidden - Akun nonaktif:
```json
{
  "error": {
    "code": "ACCOUNT_DISABLED",
    "message": "Akun Anda telah dinonaktifkan. Hubungi admin."
  }
}
```

422 Unprocessable Entity - Validasi gagal:
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validasi gagal",
    "details": [
      {
        "field": "email",
        "message": "Email wajib diisi"
      },
      {
        "field": "password",
        "message": "Password wajib diisi"
      }
    ]
  }
}
```


---

### POST /auth/refresh

Refresh access token menggunakan refresh token dari cookie.

**Authentication:** None (menggunakan HttpOnly cookie)

**Success Response (200):**
```json
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 900
  }
}
```

**Error Responses:**

401 Unauthorized - Token expired:
```json
{
  "error": {
    "code": "TOKEN_EXPIRED",
    "message": "Refresh token telah expired. Silakan login ulang."
  }
}
```

401 Unauthorized - Token tidak valid:
```json
{
  "error": {
    "code": "INVALID_TOKEN",
    "message": "Refresh token tidak valid"
  }
}
```

---

### POST /auth/logout

Logout dari sesi saat ini.

**Authentication:** Required

**Success Response (200):**
```json
{
  "message": "Berhasil logout"
}
```

**Error Responses:**

401 Unauthorized:
```json
{
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Token tidak valid atau sudah expired"
  }
}
```

---

### POST /auth/logout-all

Logout dari semua perangkat/sesi.

**Authentication:** Required

**Success Response (200):**
```json
{
  "message": "Berhasil logout dari semua perangkat",
  "data": {
    "sessions_terminated": 3
  }
}
```


---

## 2. Profile (Me)

### GET /me

Profil user yang sedang login.

**Authentication:** Required

**Success Response (200):**
```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "budi@sekolah.sch.id",
    "role": "gtk",
    "full_name": "Budi Santoso, S.Pd",
    "photo_url": "https://cdn.sipodi.go.id/photos/budi.jpg",
    "nuptk": "1234567890123456",
    "nip": "198501012010011001",
    "gender": "L",
    "birth_date": "1985-01-01",
    "gtk_type": "guru",
    "position": "Guru Matematika",
    "is_active": true,
    "school": {
      "id": "660e8400-e29b-41d4-a716-446655440000",
      "name": "SMAN 1 Malang",
      "npsn": "20518765"
    },
    "created_at": "2024-01-15T08:00:00Z",
    "updated_at": "2024-12-01T10:00:00Z"
  }
}
```

**Error Responses:**

401 Unauthorized:
```json
{
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Token tidak valid atau sudah expired"
  }
}
```

---

### PATCH /me

Update profil user yang sedang login.

**Authentication:** Required

**Request Body:**
```json
{
  "full_name": "Budi Santoso, S.Pd, M.Pd",
  "position": "Guru Matematika Senior",
  "gender": "L",
  "birth_date": "1985-01-01"
}
```

**Success Response (200):**
```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "full_name": "Budi Santoso, S.Pd, M.Pd",
    "position": "Guru Matematika Senior",
    "gender": "L",
    "birth_date": "1985-01-01",
    "updated_at": "2024-12-10T14:00:00Z"
  },
  "message": "Profil berhasil diperbarui"
}
```

**Error Responses:**

422 Unprocessable Entity:
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validasi gagal",
    "details": [
      {
        "field": "full_name",
        "message": "Nama lengkap wajib diisi"
      },
      {
        "field": "birth_date",
        "message": "Format tanggal tidak valid (gunakan YYYY-MM-DD)"
      }
    ]
  }
}
```


---

### PATCH /me/password

Ubah password.

**Authentication:** Required

**Request Body:**
```json
{
  "current_password": "oldpassword123",
  "new_password": "newpassword456",
  "new_password_confirmation": "newpassword456"
}
```

**Success Response (200):**
```json
{
  "message": "Password berhasil diubah"
}
```

**Error Responses:**

400 Bad Request - Password lama salah:
```json
{
  "error": {
    "code": "INVALID_PASSWORD",
    "message": "Password lama tidak sesuai"
  }
}
```

422 Unprocessable Entity:
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validasi gagal",
    "details": [
      {
        "field": "new_password",
        "message": "Password minimal 8 karakter"
      },
      {
        "field": "new_password_confirmation",
        "message": "Konfirmasi password tidak cocok"
      }
    ]
  }
}
```

---

### PATCH /me/photo

Update foto profil (setelah upload via presigned URL).

**Authentication:** Required

**Request Body:**
```json
{
  "upload_id": "cc0e8400-e29b-41d4-a716-446655440000"
}
```

**Success Response (200):**
```json
{
  "data": {
    "photo_url": "https://cdn.sipodi.go.id/photos/550e8400/profile.jpg"
  },
  "message": "Foto profil berhasil diperbarui"
}
```

**Error Responses:**

404 Not Found:
```json
{
  "error": {
    "code": "UPLOAD_NOT_FOUND",
    "message": "Upload tidak ditemukan atau sudah expired"
  }
}
```


---

## 3. Sekolah

### GET /schools

Daftar semua sekolah.

**Authentication:** Required (Super Admin)

**Query Parameters:**

| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| search | string | Cari berdasarkan nama/NPSN | ?search=SMAN |
| status | string | Filter status (negeri/swasta) | ?status=negeri |
| page | integer | Halaman (default: 1) | ?page=2 |
| limit | integer | Jumlah per halaman (default: 20, max: 100) | ?limit=20 |
| sort | string | Sorting field (- untuk DESC) | ?sort=-created_at |

**Success Response (200):**
```json
{
  "data": [
    {
      "id": "660e8400-e29b-41d4-a716-446655440000",
      "name": "SMAN 1 Malang",
      "npsn": "20518765",
      "status": "negeri",
      "address": "Jl. Tugu No. 1, Malang",
      "head_master": {
        "id": "770e8400-e29b-41d4-a716-446655440000",
        "full_name": "Dr. Ahmad Wijaya, M.Pd"
      },
      "gtk_count": 85,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "meta": {
    "current_page": 1,
    "per_page": 20,
    "total_pages": 5,
    "total_count": 100
  }
}
```

**Error Responses:**

403 Forbidden:
```json
{
  "error": {
    "code": "FORBIDDEN",
    "message": "Anda tidak memiliki akses ke resource ini"
  }
}
```

---

### GET /schools/{id}

Detail sekolah berdasarkan ID.

**Authentication:** Required

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| id | UUID | ID sekolah |

**Success Response (200):**
```json
{
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440000",
    "name": "SMAN 1 Malang",
    "npsn": "20518765",
    "status": "negeri",
    "address": "Jl. Tugu No. 1, Malang",
    "head_master": {
      "id": "770e8400-e29b-41d4-a716-446655440000",
      "full_name": "Dr. Ahmad Wijaya, M.Pd",
      "nip": "196501011990011001"
    },
    "gtk_count": 85,
    "guru_count": 65,
    "tendik_count": 18,
    "kepala_sekolah_count": 2,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-12-01T10:00:00Z"
  }
}
```

**Error Responses:**

404 Not Found:
```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Sekolah tidak ditemukan"
  }
}
```


---

### POST /schools

Tambah sekolah baru.

**Authentication:** Required (Super Admin)

**Request Body:**
```json
{
  "name": "SMAN 2 Malang",
  "npsn": "20518766",
  "status": "negeri",
  "address": "Jl. Laksamana Martadinata No. 84, Malang"
}
```

**Success Response (201):**
```json
{
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "name": "SMAN 2 Malang",
    "npsn": "20518766",
    "status": "negeri",
    "address": "Jl. Laksamana Martadinata No. 84, Malang",
    "created_at": "2024-12-10T14:00:00Z"
  },
  "message": "Sekolah berhasil ditambahkan"
}
```

**Error Responses:**

409 Conflict - NPSN sudah ada:
```json
{
  "error": {
    "code": "DUPLICATE_NPSN",
    "message": "NPSN sudah terdaftar"
  }
}
```

422 Unprocessable Entity:
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validasi gagal",
    "details": [
      {
        "field": "name",
        "message": "Nama sekolah wajib diisi"
      },
      {
        "field": "npsn",
        "message": "NPSN wajib diisi"
      },
      {
        "field": "status",
        "message": "Status harus negeri atau swasta"
      }
    ]
  }
}
```

---

### PUT /schools/{id}

Update data sekolah.

**Authentication:** Required (Super Admin)

**Request Body:**
```json
{
  "name": "SMAN 2 Malang",
  "npsn": "20518766",
  "status": "negeri",
  "address": "Jl. Laksamana Martadinata No. 84, Malang (Updated)",
  "head_master_id": "770e8400-e29b-41d4-a716-446655440001"
}
```

**Success Response (200):**
```json
{
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "name": "SMAN 2 Malang",
    "npsn": "20518766",
    "status": "negeri",
    "address": "Jl. Laksamana Martadinata No. 84, Malang (Updated)",
    "head_master": {
      "id": "770e8400-e29b-41d4-a716-446655440001",
      "full_name": "Dr. Siti Rahayu, M.Pd"
    },
    "updated_at": "2024-12-10T14:00:00Z"
  },
  "message": "Sekolah berhasil diperbarui"
}
```

**Error Responses:**

404 Not Found:
```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Sekolah tidak ditemukan"
  }
}
```

400 Bad Request - Kepala sekolah tidak valid:
```json
{
  "error": {
    "code": "INVALID_HEAD_MASTER",
    "message": "User yang dipilih bukan kepala sekolah"
  }
}
```


---

### DELETE /schools/{id}

Hapus sekolah.

**Authentication:** Required (Super Admin)

**Success Response (204):** No Content

**Error Responses:**

404 Not Found:
```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Sekolah tidak ditemukan"
  }
}
```

400 Bad Request - Masih ada GTK:
```json
{
  "error": {
    "code": "HAS_DEPENDENCIES",
    "message": "Tidak dapat menghapus sekolah yang masih memiliki GTK"
  }
}
```

---

### GET /schools/{id}/users

Daftar GTK di sekolah tertentu.

**Authentication:** Required (Super Admin, Admin Sekolah)

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| search | string | Cari nama/NUPTK/NIP |
| gtk_type | string | Filter: guru, tendik, kepala_sekolah |
| is_active | boolean | Filter status aktif |
| page | integer | Halaman |
| limit | integer | Jumlah per halaman |

**Success Response (200):**
```json
{
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "full_name": "Budi Santoso, S.Pd",
      "nuptk": "1234567890123456",
      "nip": "198501012010011001",
      "gtk_type": "guru",
      "position": "Guru Matematika",
      "photo_url": "https://cdn.sipodi.go.id/photos/budi.jpg",
      "is_active": true,
      "talent_summary": {
        "total": 5,
        "approved": 3,
        "pending": 2,
        "rejected": 0
      }
    }
  ],
  "meta": {
    "current_page": 1,
    "per_page": 20,
    "total_pages": 5,
    "total_count": 85
  }
}
```

**Error Responses:**

403 Forbidden - Admin sekolah akses sekolah lain:
```json
{
  "error": {
    "code": "FORBIDDEN",
    "message": "Anda hanya dapat melihat GTK di sekolah Anda"
  }
}
```


---

## 4. Users/GTK

### GET /users

Daftar semua user.

**Authentication:** Required (Super Admin)

**Query Parameters:**

| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| search | string | Cari nama/email/NUPTK/NIP | ?search=budi |
| role | string | Filter role | ?role=gtk |
| school_id | UUID | Filter berdasarkan sekolah | ?school_id=xxx |
| gtk_type | string | Filter jenis GTK | ?gtk_type=guru |
| is_active | boolean | Filter status aktif | ?is_active=true |
| page | integer | Halaman | ?page=2 |
| limit | integer | Jumlah per halaman | ?limit=20 |
| sort | string | Sorting | ?sort=-created_at |

**Success Response (200):**
```json
{
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "budi@sekolah.sch.id",
      "role": "gtk",
      "full_name": "Budi Santoso, S.Pd",
      "nuptk": "1234567890123456",
      "nip": "198501012010011001",
      "gtk_type": "guru",
      "position": "Guru Matematika",
      "is_active": true,
      "school": {
        "id": "660e8400-e29b-41d4-a716-446655440000",
        "name": "SMAN 1 Malang"
      },
      "created_at": "2024-01-15T08:00:00Z"
    }
  ],
  "meta": {
    "current_page": 1,
    "per_page": 20,
    "total_pages": 50,
    "total_count": 1000
  }
}
```

---

### GET /users/{id}

Detail user berdasarkan ID.

**Authentication:** Required

**Success Response (200):**
```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "budi@sekolah.sch.id",
    "role": "gtk",
    "full_name": "Budi Santoso, S.Pd",
    "photo_url": "https://cdn.sipodi.go.id/photos/budi.jpg",
    "nuptk": "1234567890123456",
    "nip": "198501012010011001",
    "gender": "L",
    "birth_date": "1985-01-01",
    "gtk_type": "guru",
    "position": "Guru Matematika",
    "is_active": true,
    "school": {
      "id": "660e8400-e29b-41d4-a716-446655440000",
      "name": "SMAN 1 Malang",
      "npsn": "20518765"
    },
    "talent_summary": {
      "total": 5,
      "approved": 3,
      "pending": 2,
      "rejected": 0
    },
    "created_at": "2024-01-15T08:00:00Z",
    "updated_at": "2024-12-01T10:00:00Z"
  }
}
```

**Error Responses:**

404 Not Found:
```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "User tidak ditemukan"
  }
}
```

403 Forbidden - Admin sekolah akses user sekolah lain:
```json
{
  "error": {
    "code": "FORBIDDEN",
    "message": "Anda tidak memiliki akses ke user ini"
  }
}
```


---

### POST /users

Tambah user baru.

**Authentication:** Required (Super Admin, Admin Sekolah)

**Request Body:**
```json
{
  "email": "guru.baru@sekolah.sch.id",
  "password": "password123",
  "role": "gtk",
  "full_name": "Siti Aminah, S.Pd",
  "nuptk": "9876543210123456",
  "nip": "199001012015012001",
  "gender": "P",
  "birth_date": "1990-01-01",
  "gtk_type": "guru",
  "position": "Guru Bahasa Indonesia",
  "school_id": "660e8400-e29b-41d4-a716-446655440000"
}
```

**Success Response (201):**
```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "email": "guru.baru@sekolah.sch.id",
    "full_name": "Siti Aminah, S.Pd",
    "role": "gtk",
    "gtk_type": "guru",
    "school": {
      "id": "660e8400-e29b-41d4-a716-446655440000",
      "name": "SMAN 1 Malang"
    },
    "created_at": "2024-12-10T14:00:00Z"
  },
  "message": "User berhasil ditambahkan"
}
```

**Error Responses:**

409 Conflict - Email sudah dipakai:
```json
{
  "error": {
    "code": "EMAIL_TAKEN",
    "message": "Email sudah terdaftar"
  }
}
```

409 Conflict - NUPTK sudah dipakai:
```json
{
  "error": {
    "code": "NUPTK_TAKEN",
    "message": "NUPTK sudah terdaftar"
  }
}
```

409 Conflict - NIP sudah dipakai:
```json
{
  "error": {
    "code": "NIP_TAKEN",
    "message": "NIP sudah terdaftar"
  }
}
```

422 Unprocessable Entity:
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validasi gagal",
    "details": [
      {
        "field": "email",
        "message": "Format email tidak valid"
      },
      {
        "field": "password",
        "message": "Password minimal 8 karakter"
      },
      {
        "field": "gtk_type",
        "message": "Jenis GTK harus guru, tendik, atau kepala_sekolah"
      }
    ]
  }
}
```

403 Forbidden - Admin sekolah buat user di sekolah lain:
```json
{
  "error": {
    "code": "FORBIDDEN",
    "message": "Anda hanya dapat menambah user di sekolah Anda"
  }
}
```


---

### PUT /users/{id}

Update data user.

**Authentication:** Required (Super Admin, Admin Sekolah)

**Request Body:**
```json
{
  "full_name": "Siti Aminah, S.Pd, M.Pd",
  "position": "Guru Bahasa Indonesia Senior",
  "gtk_type": "guru",
  "gender": "P",
  "birth_date": "1990-01-01",
  "school_id": "660e8400-e29b-41d4-a716-446655440000"
}
```

**Success Response (200):**
```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "full_name": "Siti Aminah, S.Pd, M.Pd",
    "position": "Guru Bahasa Indonesia Senior",
    "updated_at": "2024-12-10T14:00:00Z"
  },
  "message": "User berhasil diperbarui"
}
```

**Error Responses:**

404 Not Found:
```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "User tidak ditemukan"
  }
}
```

---

### DELETE /users/{id}

Hapus user (soft delete - set is_active = false).

**Authentication:** Required (Super Admin)

**Success Response (204):** No Content

**Error Responses:**

404 Not Found:
```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "User tidak ditemukan"
  }
}
```

400 Bad Request - Tidak bisa hapus diri sendiri:
```json
{
  "error": {
    "code": "CANNOT_DELETE_SELF",
    "message": "Tidak dapat menghapus akun sendiri"
  }
}
```

---

### PATCH /users/{id}/activate

Aktifkan user.

**Authentication:** Required (Super Admin, Admin Sekolah)

**Success Response (200):**
```json
{
  "message": "User berhasil diaktifkan"
}
```

---

### PATCH /users/{id}/deactivate

Nonaktifkan user.

**Authentication:** Required (Super Admin, Admin Sekolah)

**Success Response (200):**
```json
{
  "message": "User berhasil dinonaktifkan"
}
```


---

## 5. Talenta

### GET /talents

Daftar semua talenta (untuk admin).

**Authentication:** Required (Super Admin, Admin Sekolah)

**Query Parameters:**

| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| user_id | UUID | Filter berdasarkan user | ?user_id=xxx |
| school_id | UUID | Filter berdasarkan sekolah | ?school_id=xxx |
| talent_type | string | Filter jenis talenta | ?talent_type=peserta_pelatihan |
| status | string | Filter status verifikasi | ?status=pending |
| page | integer | Halaman | ?page=2 |
| limit | integer | Jumlah per halaman | ?limit=20 |
| sort | string | Sorting | ?sort=-created_at |

**Jenis Talenta:**
- `peserta_pelatihan`
- `pembimbing_lomba`
- `peserta_lomba`
- `minat_bakat`

**Status Verifikasi:**
- `pending`
- `approved`
- `rejected`

**Success Response (200):**
```json
{
  "data": [
    {
      "id": "880e8400-e29b-41d4-a716-446655440000",
      "user": {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "full_name": "Budi Santoso, S.Pd",
        "school_name": "SMAN 1 Malang"
      },
      "talent_type": "peserta_pelatihan",
      "status": "pending",
      "detail": {
        "activity_name": "Pelatihan Kurikulum Merdeka",
        "organizer": "Kemendikbud",
        "start_date": "2024-06-01",
        "duration_days": 5
      },
      "created_at": "2024-12-01T10:00:00Z",
      "updated_at": "2024-12-01T10:00:00Z"
    }
  ],
  "meta": {
    "current_page": 1,
    "per_page": 20,
    "total_pages": 10,
    "total_count": 200
  }
}
```


---

### GET /me/talents

Daftar talenta milik user yang login.

**Authentication:** Required (GTK)

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| talent_type | string | Filter jenis talenta |
| status | string | Filter status verifikasi |
| page | integer | Halaman |
| limit | integer | Jumlah per halaman |

**Success Response (200):**
```json
{
  "data": [
    {
      "id": "880e8400-e29b-41d4-a716-446655440000",
      "talent_type": "peserta_pelatihan",
      "status": "approved",
      "detail": {
        "activity_name": "Pelatihan Kurikulum Merdeka",
        "organizer": "Kemendikbud",
        "start_date": "2024-06-01",
        "duration_days": 5
      },
      "certificate_url": "https://cdn.sipodi.go.id/talents/sertifikat.pdf",
      "verified_at": "2024-12-05T14:00:00Z",
      "created_at": "2024-12-01T10:00:00Z"
    },
    {
      "id": "880e8400-e29b-41d4-a716-446655440001",
      "talent_type": "pembimbing_lomba",
      "status": "pending",
      "detail": {
        "competition_name": "Olimpiade Sains Nasional",
        "level": "nasional",
        "organizer": "Kemendikbud",
        "field": "akademik",
        "achievement": "Juara 1 Tingkat Nasional"
      },
      "certificate_url": "https://cdn.sipodi.go.id/talents/osn.pdf",
      "created_at": "2024-12-08T10:00:00Z"
    }
  ],
  "meta": {
    "current_page": 1,
    "per_page": 20,
    "total_pages": 1,
    "total_count": 5
  }
}
```

---

### GET /talents/{id}

Detail talenta berdasarkan ID.

**Authentication:** Required

**Success Response (200) - Peserta Pelatihan:**
```json
{
  "data": {
    "id": "880e8400-e29b-41d4-a716-446655440000",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "full_name": "Budi Santoso, S.Pd",
      "school_name": "SMAN 1 Malang"
    },
    "talent_type": "peserta_pelatihan",
    "status": "approved",
    "detail": {
      "activity_name": "Pelatihan Kurikulum Merdeka",
      "organizer": "Kemendikbud",
      "start_date": "2024-06-01",
      "duration_days": 5
    },
    "certificate_url": "https://cdn.sipodi.go.id/talents/sertifikat.pdf",
    "verified_by": {
      "id": "aa0e8400-e29b-41d4-a716-446655440000",
      "full_name": "Admin Sekolah"
    },
    "verified_at": "2024-12-05T14:00:00Z",
    "rejection_reason": null,
    "created_at": "2024-12-01T10:00:00Z",
    "updated_at": "2024-12-05T14:00:00Z"
  }
}
```

**Success Response (200) - Pembimbing Lomba:**
```json
{
  "data": {
    "id": "880e8400-e29b-41d4-a716-446655440001",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "full_name": "Budi Santoso, S.Pd",
      "school_name": "SMAN 1 Malang"
    },
    "talent_type": "pembimbing_lomba",
    "status": "approved",
    "detail": {
      "competition_name": "Olimpiade Sains Nasional",
      "level": "nasional",
      "organizer": "Kemendikbud",
      "field": "akademik",
      "achievement": "Juara 1 Tingkat Nasional"
    },
    "certificate_url": "https://cdn.sipodi.go.id/talents/osn.pdf",
    "verified_by": {
      "id": "aa0e8400-e29b-41d4-a716-446655440000",
      "full_name": "Admin Sekolah"
    },
    "verified_at": "2024-12-05T14:00:00Z",
    "rejection_reason": null,
    "created_at": "2024-12-01T10:00:00Z",
    "updated_at": "2024-12-05T14:00:00Z"
  }
}
```


**Success Response (200) - Peserta Lomba:**
```json
{
  "data": {
    "id": "880e8400-e29b-41d4-a716-446655440002",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "full_name": "Budi Santoso, S.Pd",
      "school_name": "SMAN 1 Malang"
    },
    "talent_type": "peserta_lomba",
    "status": "pending",
    "detail": {
      "competition_name": "Lomba Guru Berprestasi",
      "level": "provinsi",
      "organizer": "Dinas Pendidikan Jatim",
      "field": "akademik",
      "start_date": "2024-08-15",
      "duration_days": 3,
      "competition_field": "Guru SMA",
      "achievement": "Juara 2 Tingkat Provinsi"
    },
    "certificate_url": "https://cdn.sipodi.go.id/talents/lgb.pdf",
    "verified_by": null,
    "verified_at": null,
    "rejection_reason": null,
    "created_at": "2024-12-01T10:00:00Z",
    "updated_at": "2024-12-01T10:00:00Z"
  }
}
```

**Success Response (200) - Minat/Bakat:**
```json
{
  "data": {
    "id": "880e8400-e29b-41d4-a716-446655440003",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "full_name": "Budi Santoso, S.Pd",
      "school_name": "SMAN 1 Malang"
    },
    "talent_type": "minat_bakat",
    "status": "approved",
    "detail": {
      "interest_name": "Menulis Buku",
      "description": "Penulis buku pelajaran matematika untuk SMA, telah menerbitkan 3 buku"
    },
    "certificate_url": "https://cdn.sipodi.go.id/talents/buku.pdf",
    "verified_by": {
      "id": "aa0e8400-e29b-41d4-a716-446655440000",
      "full_name": "Admin Sekolah"
    },
    "verified_at": "2024-12-05T14:00:00Z",
    "rejection_reason": null,
    "created_at": "2024-12-01T10:00:00Z",
    "updated_at": "2024-12-05T14:00:00Z"
  }
}
```

**Error Responses:**

404 Not Found:
```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Talenta tidak ditemukan"
  }
}
```

403 Forbidden:
```json
{
  "error": {
    "code": "FORBIDDEN",
    "message": "Anda tidak memiliki akses ke talenta ini"
  }
}
```


---

### POST /me/talents

Tambah talenta baru (oleh GTK).

**Authentication:** Required (GTK)

**Request Body - Peserta Pelatihan:**
```json
{
  "talent_type": "peserta_pelatihan",
  "detail": {
    "activity_name": "Pelatihan Kurikulum Merdeka",
    "organizer": "Kemendikbud",
    "start_date": "2024-06-01",
    "duration_days": 5
  }
}
```

**Request Body - Pembimbing Lomba:**
```json
{
  "talent_type": "pembimbing_lomba",
  "detail": {
    "competition_name": "Olimpiade Sains Nasional",
    "level": "nasional",
    "organizer": "Kemendikbud",
    "field": "akademik",
    "achievement": "Juara 1 Tingkat Nasional"
  },
  "upload_id": "cc0e8400-e29b-41d4-a716-446655440000"
}
```

**Request Body - Peserta Lomba:**
```json
{
  "talent_type": "peserta_lomba",
  "detail": {
    "competition_name": "Lomba Guru Berprestasi",
    "level": "provinsi",
    "organizer": "Dinas Pendidikan Jatim",
    "field": "akademik",
    "start_date": "2024-08-15",
    "duration_days": 3,
    "competition_field": "Guru SMA",
    "achievement": "Juara 2 Tingkat Provinsi"
  },
  "upload_id": "cc0e8400-e29b-41d4-a716-446655440001"
}
```

**Request Body - Minat/Bakat:**
```json
{
  "talent_type": "minat_bakat",
  "detail": {
    "interest_name": "Menulis Buku",
    "description": "Penulis buku pelajaran matematika untuk SMA"
  },
  "upload_id": "cc0e8400-e29b-41d4-a716-446655440002"
}
```

**Jenjang Lomba (level):**
- `kota`
- `provinsi`
- `nasional`
- `internasional`

**Bidang (field):**
- `akademik`
- `inovasi`
- `teknologi`
- `sosial`
- `olahraga`
- `seni`
- `kepemimpinan`

**Success Response (201):**
```json
{
  "data": {
    "id": "880e8400-e29b-41d4-a716-446655440001",
    "talent_type": "peserta_pelatihan",
    "status": "pending",
    "detail": {
      "activity_name": "Pelatihan Kurikulum Merdeka",
      "organizer": "Kemendikbud",
      "start_date": "2024-06-01",
      "duration_days": 5
    },
    "created_at": "2024-12-10T14:00:00Z"
  },
  "message": "Talenta berhasil ditambahkan dan menunggu verifikasi"
}
```


**Error Responses:**

422 Unprocessable Entity - Peserta Pelatihan:
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validasi gagal",
    "details": [
      {
        "field": "detail.activity_name",
        "message": "Nama kegiatan wajib diisi"
      },
      {
        "field": "detail.organizer",
        "message": "Penyelenggara wajib diisi"
      },
      {
        "field": "detail.start_date",
        "message": "Tanggal mulai wajib diisi"
      },
      {
        "field": "detail.duration_days",
        "message": "Jangka waktu harus lebih dari 0"
      }
    ]
  }
}
```

422 Unprocessable Entity - Pembimbing/Peserta Lomba:
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validasi gagal",
    "details": [
      {
        "field": "detail.competition_name",
        "message": "Nama lomba wajib diisi"
      },
      {
        "field": "detail.level",
        "message": "Jenjang harus kota, provinsi, nasional, atau internasional"
      },
      {
        "field": "detail.field",
        "message": "Bidang tidak valid"
      },
      {
        "field": "detail.achievement",
        "message": "Prestasi wajib diisi"
      }
    ]
  }
}
```

422 Unprocessable Entity - Minat/Bakat:
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validasi gagal",
    "details": [
      {
        "field": "detail.interest_name",
        "message": "Nama minat/bakat wajib diisi"
      },
      {
        "field": "detail.description",
        "message": "Deskripsi wajib diisi"
      }
    ]
  }
}
```

---

### PUT /me/talents/{id}

Update talenta milik sendiri.

**Authentication:** Required (GTK)

**Note:** Setelah update, status akan kembali menjadi `pending` dan perlu verifikasi ulang.

**Request Body:**
```json
{
  "detail": {
    "activity_name": "Pelatihan Kurikulum Merdeka (Updated)",
    "organizer": "Kemendikbud RI",
    "start_date": "2024-06-01",
    "duration_days": 7
  },
  "upload_id": "cc0e8400-e29b-41d4-a716-446655440003"
}
```

**Success Response (200):**
```json
{
  "data": {
    "id": "880e8400-e29b-41d4-a716-446655440001",
    "talent_type": "peserta_pelatihan",
    "status": "pending",
    "detail": {
      "activity_name": "Pelatihan Kurikulum Merdeka (Updated)",
      "organizer": "Kemendikbud RI",
      "start_date": "2024-06-01",
      "duration_days": 7
    },
    "updated_at": "2024-12-10T14:00:00Z"
  },
  "message": "Talenta berhasil diperbarui dan menunggu verifikasi ulang"
}
```

---

### DELETE /me/talents/{id}

Hapus talenta milik sendiri.

**Authentication:** Required (GTK)

**Success Response (204):** No Content

**Error Responses:**

404 Not Found:
```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Talenta tidak ditemukan"
  }
}
```

403 Forbidden - Bukan milik sendiri:
```json
{
  "error": {
    "code": "FORBIDDEN",
    "message": "Anda hanya dapat menghapus talenta milik sendiri"
  }
}
```


---

## 6. Verifikasi Talenta

### GET /verifications/talents

Daftar talenta yang perlu diverifikasi.

**Authentication:** Required (Super Admin, Admin Sekolah)

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| status | string | Filter status (default: pending) |
| school_id | UUID | Filter berdasarkan sekolah |
| talent_type | string | Filter jenis talenta |
| page | integer | Halaman |
| limit | integer | Jumlah per halaman |

**Success Response (200):**
```json
{
  "data": [
    {
      "id": "880e8400-e29b-41d4-a716-446655440000",
      "user": {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "full_name": "Budi Santoso, S.Pd",
        "photo_url": "https://cdn.sipodi.go.id/photos/budi.jpg"
      },
      "school": {
        "id": "660e8400-e29b-41d4-a716-446655440000",
        "name": "SMAN 1 Malang"
      },
      "talent_type": "peserta_pelatihan",
      "status": "pending",
      "detail": {
        "activity_name": "Pelatihan Kurikulum Merdeka",
        "organizer": "Kemendikbud"
      },
      "has_certificate": true,
      "created_at": "2024-12-01T10:00:00Z"
    }
  ],
  "meta": {
    "current_page": 1,
    "per_page": 20,
    "total_pages": 5,
    "total_count": 100
  }
}
```

---

### POST /verifications/talents/{id}/approve

Approve talenta.

**Authentication:** Required (Super Admin, Admin Sekolah)

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| id | UUID | ID talenta |

**Success Response (200):**
```json
{
  "data": {
    "id": "880e8400-e29b-41d4-a716-446655440000",
    "status": "approved",
    "verified_at": "2024-12-10T14:00:00Z"
  },
  "message": "Talenta berhasil disetujui"
}
```

**Note:** Notifikasi akan otomatis dikirim ke GTK yang bersangkutan.

**Error Responses:**

404 Not Found:
```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Talenta tidak ditemukan"
  }
}
```

400 Bad Request - Sudah diverifikasi:
```json
{
  "error": {
    "code": "ALREADY_VERIFIED",
    "message": "Talenta sudah diverifikasi sebelumnya"
  }
}
```

403 Forbidden - Admin sekolah verifikasi talenta sekolah lain:
```json
{
  "error": {
    "code": "FORBIDDEN",
    "message": "Anda hanya dapat memverifikasi talenta GTK di sekolah Anda"
  }
}
```


---

### POST /verifications/talents/{id}/reject

Reject talenta.

**Authentication:** Required (Super Admin, Admin Sekolah)

**Request Body:**
```json
{
  "rejection_reason": "Dokumen bukti tidak valid atau tidak terbaca"
}
```

**Success Response (200):**
```json
{
  "data": {
    "id": "880e8400-e29b-41d4-a716-446655440000",
    "status": "rejected",
    "rejection_reason": "Dokumen bukti tidak valid atau tidak terbaca",
    "verified_at": "2024-12-10T14:00:00Z"
  },
  "message": "Talenta ditolak"
}
```

**Note:** Notifikasi akan otomatis dikirim ke GTK yang bersangkutan.

**Error Responses:**

422 Unprocessable Entity:
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validasi gagal",
    "details": [
      {
        "field": "rejection_reason",
        "message": "Alasan penolakan wajib diisi"
      }
    ]
  }
}
```

---

### POST /verifications/talents/batch/approve

Approve multiple talenta sekaligus.

**Authentication:** Required (Super Admin, Admin Sekolah)

**Request Body:**
```json
{
  "ids": [
    "880e8400-e29b-41d4-a716-446655440000",
    "880e8400-e29b-41d4-a716-446655440001",
    "880e8400-e29b-41d4-a716-446655440002"
  ]
}
```

**Success Response (200):**
```json
{
  "data": {
    "approved_count": 3,
    "failed_count": 0,
    "failed_ids": []
  },
  "message": "3 talenta berhasil disetujui"
}
```

**Partial Success Response (200):**
```json
{
  "data": {
    "approved_count": 2,
    "failed_count": 1,
    "failed_ids": [
      {
        "id": "880e8400-e29b-41d4-a716-446655440002",
        "reason": "Talenta tidak ditemukan"
      }
    ]
  },
  "message": "2 talenta berhasil disetujui, 1 gagal"
}
```

---

### POST /verifications/talents/batch/reject

Reject multiple talenta sekaligus.

**Authentication:** Required (Super Admin, Admin Sekolah)

**Request Body:**
```json
{
  "ids": [
    "880e8400-e29b-41d4-a716-446655440000",
    "880e8400-e29b-41d4-a716-446655440001"
  ],
  "rejection_reason": "Dokumen tidak lengkap"
}
```

**Success Response (200):**
```json
{
  "data": {
    "rejected_count": 2,
    "failed_count": 0,
    "failed_ids": []
  },
  "message": "2 talenta ditolak"
}
```


---

## 7. Notifikasi

### GET /me/notifications

Daftar notifikasi user yang login.

**Authentication:** Required

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| is_read | boolean | Filter berdasarkan status baca |
| page | integer | Halaman |
| limit | integer | Jumlah per halaman |

**Success Response (200):**
```json
{
  "data": [
    {
      "id": "990e8400-e29b-41d4-a716-446655440000",
      "type": "talent_approved",
      "message": "Talenta Anda telah disetujui",
      "talent_id": "880e8400-e29b-41d4-a716-446655440000",
      "is_read": false,
      "created_at": "2024-12-10T14:00:00Z"
    },
    {
      "id": "990e8400-e29b-41d4-a716-446655440001",
      "type": "talent_rejected",
      "message": "Talenta Anda ditolak. Alasan: Dokumen tidak valid",
      "talent_id": "880e8400-e29b-41d4-a716-446655440001",
      "is_read": true,
      "created_at": "2024-12-09T10:00:00Z"
    }
  ],
  "meta": {
    "current_page": 1,
    "per_page": 20,
    "total_pages": 1,
    "total_count": 2,
    "unread_count": 1
  }
}
```

---

### GET /me/notifications/unread-count

Jumlah notifikasi yang belum dibaca.

**Authentication:** Required

**Success Response (200):**
```json
{
  "data": {
    "unread_count": 5
  }
}
```

---

### PATCH /me/notifications/{id}/read

Tandai notifikasi sebagai sudah dibaca.

**Authentication:** Required

**Success Response (200):**
```json
{
  "message": "Notifikasi ditandai sudah dibaca"
}
```

**Error Responses:**

404 Not Found:
```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Notifikasi tidak ditemukan"
  }
}
```

---

### PATCH /me/notifications/read-all

Tandai semua notifikasi sebagai sudah dibaca.

**Authentication:** Required

**Success Response (200):**
```json
{
  "data": {
    "marked_count": 5
  },
  "message": "Semua notifikasi ditandai sudah dibaca"
}
```


---

## 8. File Upload (MinIO)

SIPODI menggunakan MinIO sebagai object storage dengan arsitektur presigned URL untuk upload file. Flow ini memastikan file diupload langsung ke MinIO tanpa melalui backend, mengurangi beban server.

### Upload Flow

```
1. Client request presigned URL dari backend
2. Backend generate presigned URL dari MinIO
3. Client upload file langsung ke MinIO menggunakan presigned URL
4. Client confirm upload ke backend
5. Backend update database dengan file URL
```

### POST /uploads/presign

Request presigned URL untuk upload file.

**Authentication:** Required

**Request Body:**
```json
{
  "filename": "sertifikat_osn.pdf",
  "content_type": "application/pdf",
  "upload_type": "talent_certificate"
}
```

**Upload Types:**
- `profile_photo` - Foto profil (max 2MB, image/*)
- `talent_certificate` - Sertifikat/bukti talenta (max 10MB, application/pdf, image/*)

**Success Response (200):**
```json
{
  "data": {
    "upload_id": "cc0e8400-e29b-41d4-a716-446655440000",
    "presigned_url": "https://minio.sipodi.go.id/sipodi-bucket/uploads/2024/12/cc0e8400.pdf?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=...",
    "method": "PUT",
    "expires_in": 3600,
    "max_size": 10485760,
    "allowed_types": ["application/pdf", "image/jpeg", "image/png"]
  }
}
```

**Error Responses:**

400 Bad Request - Tipe file tidak diizinkan:
```json
{
  "error": {
    "code": "INVALID_FILE_TYPE",
    "message": "Tipe file tidak diizinkan. Gunakan PDF atau gambar (JPG, PNG)"
  }
}
```

422 Unprocessable Entity:
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validasi gagal",
    "details": [
      {
        "field": "filename",
        "message": "Nama file wajib diisi"
      },
      {
        "field": "content_type",
        "message": "Content type wajib diisi"
      },
      {
        "field": "upload_type",
        "message": "Upload type tidak valid"
      }
    ]
  }
}
```


---

### Upload ke MinIO (Client-side)

Setelah mendapat presigned URL, client upload file langsung ke MinIO.

**Request:**
```http
PUT {presigned_url}
Content-Type: application/pdf
Content-Length: 1024000

[binary file data]
```

**Success Response (200):** Empty body dari MinIO

**Error Responses:**

403 Forbidden - URL expired:
```
AccessDenied: Request has expired
```

400 Bad Request - File terlalu besar:
```
EntityTooLarge: Your proposed upload exceeds the maximum allowed size
```

---

### POST /uploads/{upload_id}/confirm

Konfirmasi upload berhasil dan dapatkan final URL.

**Authentication:** Required

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| upload_id | UUID | ID upload dari presign response |

**Success Response (200):**
```json
{
  "data": {
    "upload_id": "cc0e8400-e29b-41d4-a716-446655440000",
    "file_url": "https://cdn.sipodi.go.id/talents/2024/12/cc0e8400.pdf",
    "filename": "sertifikat_osn.pdf",
    "file_size": 1024000,
    "content_type": "application/pdf"
  },
  "message": "Upload berhasil dikonfirmasi"
}
```

**Error Responses:**

404 Not Found - Upload tidak ditemukan:
```json
{
  "error": {
    "code": "UPLOAD_NOT_FOUND",
    "message": "Upload tidak ditemukan atau sudah expired"
  }
}
```

400 Bad Request - File belum diupload:
```json
{
  "error": {
    "code": "FILE_NOT_UPLOADED",
    "message": "File belum diupload ke storage"
  }
}
```

---

### DELETE /uploads/{upload_id}

Batalkan upload (hapus file dari MinIO jika sudah diupload).

**Authentication:** Required

**Success Response (204):** No Content

**Error Responses:**

404 Not Found:
```json
{
  "error": {
    "code": "UPLOAD_NOT_FOUND",
    "message": "Upload tidak ditemukan"
  }
}
```


---

### Complete Upload Example

**Step 1: Request Presigned URL**
```http
POST /api/v1/uploads/presign
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "filename": "sertifikat_osn.pdf",
  "content_type": "application/pdf",
  "upload_type": "talent_certificate"
}
```

**Response:**
```json
{
  "data": {
    "upload_id": "cc0e8400-e29b-41d4-a716-446655440000",
    "presigned_url": "https://minio.sipodi.go.id/sipodi-bucket/uploads/2024/12/cc0e8400.pdf?X-Amz-Algorithm=...",
    "method": "PUT",
    "expires_in": 3600
  }
}
```

**Step 2: Upload to MinIO**
```http
PUT https://minio.sipodi.go.id/sipodi-bucket/uploads/2024/12/cc0e8400.pdf?X-Amz-Algorithm=...
Content-Type: application/pdf

[binary file data]
```

**Step 3: Confirm Upload**
```http
POST /api/v1/uploads/cc0e8400-e29b-41d4-a716-446655440000/confirm
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response:**
```json
{
  "data": {
    "upload_id": "cc0e8400-e29b-41d4-a716-446655440000",
    "file_url": "https://cdn.sipodi.go.id/talents/2024/12/cc0e8400.pdf"
  }
}
```

**Step 4: Use upload_id when creating talent**
```http
POST /api/v1/me/talents
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "talent_type": "pembimbing_lomba",
  "detail": {
    "competition_name": "Olimpiade Sains Nasional",
    "level": "nasional",
    "organizer": "Kemendikbud",
    "field": "akademik",
    "achievement": "Juara 1"
  },
  "upload_id": "cc0e8400-e29b-41d4-a716-446655440000"
}
```


---

## 9. Dashboard & Statistik

### GET /dashboard/summary

Ringkasan statistik untuk dashboard.

**Authentication:** Required

**Response berbeda berdasarkan role:**

**Super Admin Response (200):**
```json
{
  "data": {
    "total_schools": 150,
    "total_users": 5000,
    "total_gtk": 4800,
    "total_admin_sekolah": 150,
    "gtk_by_type": {
      "guru": 3500,
      "tendik": 1100,
      "kepala_sekolah": 200
    },
    "total_talents": 12000,
    "talents_by_status": {
      "pending": 500,
      "approved": 11000,
      "rejected": 500
    },
    "talents_by_type": {
      "peserta_pelatihan": 5000,
      "pembimbing_lomba": 2000,
      "peserta_lomba": 3000,
      "minat_bakat": 2000
    },
    "recent_talents": [
      {
        "id": "880e8400-e29b-41d4-a716-446655440000",
        "user_name": "Budi Santoso",
        "school_name": "SMAN 1 Malang",
        "talent_type": "peserta_pelatihan",
        "status": "pending",
        "created_at": "2024-12-10T14:00:00Z"
      }
    ]
  }
}
```

**Admin Sekolah Response (200):**
```json
{
  "data": {
    "school": {
      "id": "660e8400-e29b-41d4-a716-446655440000",
      "name": "SMAN 1 Malang"
    },
    "total_gtk": 85,
    "gtk_by_type": {
      "guru": 65,
      "tendik": 18,
      "kepala_sekolah": 2
    },
    "total_talents": 250,
    "talents_by_status": {
      "pending": 20,
      "approved": 220,
      "rejected": 10
    },
    "talents_by_type": {
      "peserta_pelatihan": 100,
      "pembimbing_lomba": 50,
      "peserta_lomba": 60,
      "minat_bakat": 40
    },
    "pending_verifications": 20,
    "recent_talents": [
      {
        "id": "880e8400-e29b-41d4-a716-446655440000",
        "user_name": "Budi Santoso",
        "talent_type": "peserta_pelatihan",
        "status": "pending",
        "created_at": "2024-12-10T14:00:00Z"
      }
    ]
  }
}
```

**GTK Response (200):**
```json
{
  "data": {
    "my_talents": {
      "total": 5,
      "pending": 2,
      "approved": 3,
      "rejected": 0
    },
    "talents_by_type": {
      "peserta_pelatihan": 2,
      "pembimbing_lomba": 1,
      "peserta_lomba": 1,
      "minat_bakat": 1
    },
    "unread_notifications": 2,
    "recent_talents": [
      {
        "id": "880e8400-e29b-41d4-a716-446655440000",
        "talent_type": "peserta_pelatihan",
        "status": "approved",
        "created_at": "2024-12-01T10:00:00Z"
      }
    ]
  }
}
```


---

### GET /dashboard/schools/statistics

Statistik per sekolah (untuk Super Admin).

**Authentication:** Required (Super Admin)

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| status | string | Filter status sekolah (negeri/swasta) |
| sort | string | Sorting: gtk_count, talent_count |
| page | integer | Halaman |
| limit | integer | Jumlah per halaman |

**Success Response (200):**
```json
{
  "data": [
    {
      "school": {
        "id": "660e8400-e29b-41d4-a716-446655440000",
        "name": "SMAN 1 Malang",
        "npsn": "20518765",
        "status": "negeri"
      },
      "gtk_count": 85,
      "talent_count": 250,
      "talents_approved": 220,
      "talents_pending": 20,
      "talents_rejected": 10
    }
  ],
  "meta": {
    "current_page": 1,
    "per_page": 20,
    "total_pages": 8,
    "total_count": 150
  }
}
```

---

### GET /dashboard/talents/statistics

Statistik talenta berdasarkan berbagai dimensi.

**Authentication:** Required (Super Admin, Admin Sekolah)

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| school_id | UUID | Filter berdasarkan sekolah |
| group_by | string | Grouping: type, status, level, field |
| date_from | date | Filter dari tanggal |
| date_to | date | Filter sampai tanggal |

**Success Response (200) - Group by Type:**
```json
{
  "data": {
    "group_by": "type",
    "statistics": [
      {
        "talent_type": "peserta_pelatihan",
        "total": 5000,
        "approved": 4500,
        "pending": 300,
        "rejected": 200
      },
      {
        "talent_type": "pembimbing_lomba",
        "total": 2000,
        "approved": 1800,
        "pending": 150,
        "rejected": 50
      },
      {
        "talent_type": "peserta_lomba",
        "total": 3000,
        "approved": 2700,
        "pending": 200,
        "rejected": 100
      },
      {
        "talent_type": "minat_bakat",
        "total": 2000,
        "approved": 1800,
        "pending": 150,
        "rejected": 50
      }
    ]
  }
}
```

**Success Response (200) - Group by Level (untuk lomba):**
```json
{
  "data": {
    "group_by": "level",
    "statistics": [
      {
        "level": "kota",
        "total": 2000
      },
      {
        "level": "provinsi",
        "total": 1500
      },
      {
        "level": "nasional",
        "total": 1000
      },
      {
        "level": "internasional",
        "total": 500
      }
    ]
  }
}
```

**Success Response (200) - Group by Field:**
```json
{
  "data": {
    "group_by": "field",
    "statistics": [
      {
        "field": "akademik",
        "total": 2000
      },
      {
        "field": "teknologi",
        "total": 1000
      },
      {
        "field": "seni",
        "total": 800
      },
      {
        "field": "olahraga",
        "total": 700
      },
      {
        "field": "inovasi",
        "total": 500
      },
      {
        "field": "sosial",
        "total": 300
      },
      {
        "field": "kepemimpinan",
        "total": 200
      }
    ]
  }
}
```


---

## 10. Export Laporan

### GET /exports/gtk

Export data GTK ke Excel/PDF.

**Authentication:** Required (Super Admin, Admin Sekolah)

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| format | string | Format export: excel, pdf |
| school_id | UUID | Filter berdasarkan sekolah |
| gtk_type | string | Filter jenis GTK |

**Success Response (200):**
```json
{
  "data": {
    "download_url": "https://cdn.sipodi.go.id/exports/gtk_2024-12-10_abc123.xlsx",
    "filename": "data_gtk_2024-12-10.xlsx",
    "expires_in": 3600
  }
}
```

---

### GET /exports/talents

Export data talenta ke Excel/PDF.

**Authentication:** Required (Super Admin, Admin Sekolah)

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| format | string | Format export: excel, pdf |
| school_id | UUID | Filter berdasarkan sekolah |
| talent_type | string | Filter jenis talenta |
| status | string | Filter status |
| date_from | date | Filter dari tanggal |
| date_to | date | Filter sampai tanggal |

**Success Response (200):**
```json
{
  "data": {
    "download_url": "https://cdn.sipodi.go.id/exports/talents_2024-12-10_abc123.xlsx",
    "filename": "data_talenta_2024-12-10.xlsx",
    "expires_in": 3600
  }
}
```

---

### GET /exports/schools

Export data sekolah ke Excel/PDF.

**Authentication:** Required (Super Admin)

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| format | string | Format export: excel, pdf |
| status | string | Filter status sekolah |

**Success Response (200):**
```json
{
  "data": {
    "download_url": "https://cdn.sipodi.go.id/exports/schools_2024-12-10_abc123.xlsx",
    "filename": "data_sekolah_2024-12-10.xlsx",
    "expires_in": 3600
  }
}
```


---

## Common Error Responses

### 401 Unauthorized

Token tidak valid atau sudah expired:
```json
{
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Token tidak valid atau sudah expired"
  }
}
```

### 403 Forbidden

Tidak memiliki akses:
```json
{
  "error": {
    "code": "FORBIDDEN",
    "message": "Anda tidak memiliki akses ke resource ini"
  }
}
```

### 404 Not Found

Resource tidak ditemukan:
```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Resource tidak ditemukan"
  }
}
```

### 422 Unprocessable Entity

Validasi gagal:
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validasi gagal",
    "details": [
      {
        "field": "field_name",
        "message": "Error message"
      }
    ]
  }
}
```

### 429 Too Many Requests

Rate limit exceeded:
```json
{
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Terlalu banyak request. Coba lagi nanti.",
    "retry_after": 60
  }
}
```

### 500 Internal Server Error

Server error:
```json
{
  "error": {
    "code": "INTERNAL_ERROR",
    "message": "Terjadi kesalahan pada server. Silakan coba lagi."
  }
}
```

---

## Rate Limiting

API menerapkan rate limiting untuk mencegah abuse:

- **Authenticated requests:** 1000 requests per 15 menit
- **Login endpoint:** 10 requests per menit per IP

Response headers:
```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1702036800
```

---

## Changelog

### v1.0.0 (2024-12-10)
- Initial API release
- Authentication (login, refresh, logout)
- User management (CRUD)
- School management (CRUD)
- Talent management (CRUD)
- Talent verification (approve/reject)
- Notifications
- File upload with MinIO presigned URL
- Dashboard & statistics
- Export reports (Excel/PDF)
