# SIPODI Backend

Backend API untuk SIPODI (Sistem Informasi Potensi Diri) - Sistem manajemen data GTK Cabang Dinas Pendidikan Wilayah Malang.

## Tech Stack

- **Language:** Go 1.22
- **Framework:** GoFiber v2
- **Database:** PostgreSQL
- **Object Storage:** MinIO
- **Authentication:** JWT (Access Token + Refresh Token)

## Project Structure

```
backend/
├── cmd/
│   └── api/
│       └── main.go          # Application entry point
├── db/
│   └── db.sql               # Database schema
├── internal/
│   ├── config/              # Configuration
│   ├── database/            # Database connection
│   ├── domain/              # Entities and DTOs
│   ├── handler/             # HTTP handlers
│   ├── middleware/          # HTTP middleware
│   ├── repository/          # Data access layer
│   ├── router/              # Route definitions
│   ├── service/             # Business logic
│   └── storage/             # MinIO storage
├── .env.example             # Environment variables example
├── Dockerfile               # Docker build file
├── go.mod                   # Go modules
├── Makefile                 # Build commands
└── README.md                # This file
```

## Getting Started

### Prerequisites

- Go 1.22+
- PostgreSQL 15+
- MinIO

### Installation

1. Clone repository dan masuk ke folder backend:
```bash
cd backend
```

2. Copy environment file:
```bash
cp .env.example .env
```

3. Edit `.env` sesuai konfigurasi lokal

4. Install dependencies:
```bash
go mod download
```

5. Setup database:
```bash
psql -U postgres -c "CREATE DATABASE sipodi"
psql -U postgres -d sipodi -f db/db.sql
```

6. Run application:
```bash
make run
# atau
go run cmd/api/main.go
```

### Development dengan Hot Reload

Install air:
```bash
go install github.com/cosmtrek/air@latest
```

Run:
```bash
make dev
```

## API Documentation

API documentation tersedia di `/docs/api.md`

Base URL: `http://localhost:8080/api/v1`

### Authentication

Semua endpoint (kecuali login) memerlukan JWT Bearer Token:
```
Authorization: Bearer <access_token>
```

### Main Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /auth/login | Login |
| POST | /auth/refresh | Refresh token |
| GET | /me | Get profile |
| GET | /schools | List schools |
| GET | /users | List users |
| GET | /talents | List talents |
| GET | /verifications/talents | List pending verifications |
| GET | /dashboard/summary | Dashboard summary |

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| APP_ENV | Environment (development/production) | development |
| APP_PORT | Server port | 8080 |
| DB_HOST | PostgreSQL host | localhost |
| DB_PORT | PostgreSQL port | 5432 |
| DB_USER | PostgreSQL user | sipodi |
| DB_PASSWORD | PostgreSQL password | sipodi_secret |
| DB_NAME | Database name | sipodi |
| JWT_SECRET | JWT signing secret | - |
| JWT_ACCESS_EXPIRY | Access token expiry | 15m |
| JWT_REFRESH_EXPIRY | Refresh token expiry | 7d |
| MINIO_ENDPOINT | MinIO endpoint | localhost:9000 |
| MINIO_ACCESS_KEY | MinIO access key | minioadmin |
| MINIO_SECRET_KEY | MinIO secret key | minioadmin |
| MINIO_BUCKET | MinIO bucket name | sipodi |

## Docker

Build image:
```bash
make docker-build
```

Run container:
```bash
make docker-run
```

## Testing

Run tests:
```bash
make test
```

Run with coverage:
```bash
make test-coverage
```

## Default Credentials

Super Admin:
- Email: `superadmin@sipodi.go.id`
- Password: `admin123`

⚠️ **Ganti password di production!**

## License

Proprietary - Cabang Dinas Pendidikan Wilayah Malang
