# Code Review

## Lokasi File
```
cmd/pastebin/main.go
```

## Penjelasn
Pejelasan ada di komentar dalam file



# Sistem Tiket Event 🎫

Aplikasi manajemen tiket event berbasis Go dengan Fiber framework. Sistem ini memungkinkan organizer untuk membuat event, dan pengguna untuk membeli tiket event.

## Fitur Aplikasi 🚀

- Autentikasi user (register, login, verifikasi email)
- Manajemen event (create, read, update, delete)
- Transaksi pembelian tiket
- Verifikasi pembayaran
- Dan masih banyak lagi...

## Tech Stack 💻

- Go (Golang) sebagai bahasa pemrograman
- Fiber sebagai web framework
- PostgreSQL sebagai database
- JWT untuk autentikasi
- Argon2 untuk hashing password

## Cara Menjalankan Aplikasi 🏃‍♂️

### Prasyarat

- Go 1.16+
- PostgreSQL
- SMTP Server (untuk fitur email)

### Setup Database

1. Buat database PostgreSQL baru
   ```bash
   createdb ticket_system
   ```

2. Setup environment variables di file `.env` (buat sendiri filenya)
   ```
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=password_kamu
   DB_NAME=ticket_system
   DB_SSLMODE=disable
   SERVER_PORT=8080
   JWT_SECRET=rahasia_aku_kamu_dan_jwt
   TOKEN_EXPIRY=24
   
   # SMTP Settings
   SMTP_HOST=smtp.example.com
   SMTP_PORT=587
   SMTP_USERNAME=email@example.com
   SMTP_PASSWORD=password_email
   SMTP_FROM_NAME=Sistem Tiket Event
   ```

3. Jalankan migrasi database
   ```bash
   go run cmd/migrate/main.go
   ```

### Memulai Aplikasi

1. Clone repo ini
   ```bash
   git clone https://repository-url/ticket-system.git
   cd ticket-system
   ```

2. Install dependencies
   ```bash
   go mod download
   ```

3. Jalankan aplikasi
   ```bash
   go run cmd/api/main.go
   ```

4. Aplikasi akan berjalan di `http://localhost:8080`

## Menjalankan Tests 🧪

### Menjalankan Seluruh Test

```bash
go test -v ./test/...
```

### Menjalankan Test Spesifik

Untuk menjalankan test pada usecase saja:
```bash
go test -v ./test/usecase/...
```

Untuk menjalankan test pada handler saja:
```bash
go test -v ./test/handler/...
```

Untuk menjalankan satu file test saja:
```bash
go test -v ./test/usecase/event_usecase_test.go
```

Untuk menjalankan satu test function saja:
```bash
go test -v ./test/usecase -run TestCreateEvent
```

## API Endpoints 🌐

### Authentication

- `POST /api/register` - Register user baru
- `POST /api/login` - Login user
- `GET /api/verify-email` - Verifikasi email
- `POST /api/resend-verification` - Kirim ulang email verifikasi

### User Profile

- `PUT /api/profile` - Update profil user

### Events

- `GET /api/events` - List semua event
- `GET /api/events/:id` - Detail event
- `POST /api/organizer/events` - Buat event baru (organizer only)
- `PUT /api/organizer/events/:id` - Update event (organizer only)
- `DELETE /api/organizer/events/:id` - Hapus event (organizer only)
- `GET /api/organizer/events` - List event by organizer
- `GET /api/organizer/events/:id/sales` - Data penjualan event

### Transactions

- `POST /api/transactions` - Buat transaksi baru
- `GET /api/transactions` - List transaksi user
- `GET /api/transactions/:id` - Detail transaksi
- `GET /api/transactions/code` - Cari transaksi by code
- `POST /api/transactions/proof` - Upload bukti pembayaran
- `PUT /api/transactions/:id/cancel` - Batalkan transaksi
- `PUT /api/organizer/transactions/:id/verify` - Verifikasi pembayaran (organizer only)


## Saran Pengembangan 💡

- Implementasi caching (Redis)
