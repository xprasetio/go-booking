# Go booking

booking untuk aplikasi Go dengan fitur-fitur modern dan best practices.

## Fitur

- Echo sebagai web framework
- GORM untuk ORM database
- Redis untuk caching
- JWT untuk autentikasi
- Dependency Injection menggunakan sarulabs/di
- Konfigurasi menggunakan Viper
- Logging menggunakan Logrus
- Docker dan Docker Compose untuk containerization

## Persyaratan

- Go 1.23.4 atau lebih baru
- Docker dan Docker Compose
- MySQL 8.0
- Redis 7.0

## Struktur Proyek Menggunakan DDD serta Rich Domain

```
.
├── cmd/            # Entry point aplikasi
├── config/         # Konfigurasi aplikasi
├── container/      # File-file terkait Docker
├── internal/       # Kode internal aplikasi
├── pkg/           # Package yang dapat digunakan ulang
├── routes/        # Definisi routing
└── shared/        # Kode yang dibagi antar package
```

## Cara Menjalankan

### Menggunakan Docker

1. Clone repository:

```bash
git clone <repository-url>
cd booking
```

2. Jalankan dengan Docker Compose:

```bash
docker-compose up -d
```

Aplikasi akan berjalan di `http://localhost:8080`

### Menjalankan Secara Lokal

1. Install dependensi:

```bash
go mod download
```

2. Jalankan aplikasi:

```bash
go run cmd/main.go
```

## Konfigurasi

Konfigurasi aplikasi dapat diatur melalui environment variables atau file konfigurasi di `config/`. Beberapa konfigurasi penting:

- `DB_HOST`: Host database MySQL
- `DB_PORT`: Port database MySQL
- `DB_USER`: Username database
- `DB_PASSWORD`: Password database
- `DB_NAME`: Nama database
- `REDIS_HOST`: Host Redis
- `REDIS_PORT`: Port Redis

## Pengembangan

### Menjalankan Tests

```bash
go test ./...
```

### Menjalankan Linter

```bash
golangci-lint run
```

## Lisensi

MIT License
