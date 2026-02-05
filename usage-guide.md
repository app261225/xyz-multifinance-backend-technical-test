# Panduan Penggunaan API Backend PT XYZ Multifinance

Panduan lengkap untuk menggunakan API Backend PT XYZ Multifinance dengan contoh-contoh praktis.

## Daftar Isi
1. [Persiapan Awal](#persiapan-awal)
2. [Menjalankan Aplikasi](#menjalankan-aplikasi)
3. [Endpoint Manajemen Konsumen](#endpoint-manajemen-konsumen)
4. [Endpoint Manajemen Batas Kredit](#endpoint-manajemen-batas-kredit)
5. [Endpoint Manajemen Transaksi](#endpoint-manajemen-transaksi)
6. [Testing & Debugging](#testing--debugging)
7. [Troubleshooting](#troubleshooting)
8. [Best Practices](#best-practices)

---

## Persiapan Awal

### 1. Instalasi Dependensi

**Menggunakan Docker (Recommended):**
```bash
cd xyz-multifinance-backend-technical-test
docker-compose up -d
```

**Atau Instalasi Manual:**

```bash
# 1. Clone repository
git clone https://github.com/yourusername/xyz-multifinance-backend-technical-test.git
cd xyz-multifinance-backend-technical-test

# 2. Install Go dependencies
go mod download
go mod tidy

# 3. Setup database
mysql -u root -p < database_schema.sql

# 4. Konfigurasi environment
cp .env.example .env
# Edit .env sesuai konfigurasi lokal Anda
```

### 2. Verifikasi Instalasi

```bash
# Cek koneksi database
mysql -h localhost -u xyz_user -p xyz_multifinance

# Cek dependensi Go
go mod verify

# List dependensi
go list -m all
```

---

## Menjalankan Aplikasi

### Metode 1: Langsung dengan Go

```bash
go run main.go
```

**Output yang diharapkan:**
```
✓ Database terhubung dan dimigrasikan dengan sukses
✓ Database terhubung dengan sukses
✓ Memulai server API di port 8080
✓ Health check: http://localhost:8080/health
✓ OWASP Security Headers: ENABLED
✓ Input Validation: ENABLED
✓ CORS Protection: ENABLED
```

### Metode 2: Build & Run

```bash
# Build aplikasi
go build -o xyz-multifinance .

# Jalankan binary
./xyz-multifinance
```

### Metode 3: Docker

```bash
# Build image
docker build -t xyz-multifinance:1.0.0 .

# Run container
docker run -p 8080:8080 \
  -e DB_HOST=mysql \
  -e DB_USER=xyz_user \
  -e DB_PASS=xyz_password \
  xyz-multifinance:1.0.0
```

### Metode 4: Docker Compose

```bash
# Start semua services
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop services
docker-compose down
```

### Verifikasi Server Berjalan

```bash
# Test health endpoint
curl http://localhost:8080/health

# Expected response:
# {"status":"healthy","version":"1.0.0"}
```

---

## Endpoint Manajemen Konsumen

### 1. Daftarkan Konsumen Baru

**URL:** `POST /api/consumers`

**Content-Type:** `application/json`

**Request Body:**
```json
{
  "nik": "3173011234567890",
  "full_name": "John Doe",
  "legal_name": "John Doe",
  "place_of_birth": "Jakarta",
  "date_of_birth": "1990-01-15T00:00:00Z",
  "salary": 5000000
}
```

**Contoh dengan cURL:**
```bash
curl -X POST http://localhost:8080/api/consumers \
  -H "Content-Type: application/json" \
  -d '{
    "nik": "3173011234567890",
    "full_name": "John Doe",
    "legal_name": "John Doe",
    "place_of_birth": "Jakarta",
    "date_of_birth": "1990-01-15T00:00:00Z",
    "salary": 5000000
  }'
```

**Contoh dengan Postman:**
1. Buat request baru: POST
2. URL: http://localhost:8080/api/consumers
3. Headers:
   - Content-Type: application/json
4. Body (raw JSON):
```json
{
  "nik": "3173011234567890",
  "full_name": "John Doe",
  "legal_name": "John Doe",
  "place_of_birth": "Jakarta",
  "date_of_birth": "1990-01-15T00:00:00Z",
  "salary": 5000000
}
```

**Response (201 Created):**
```json
{
  "message": "Consumer registered successfully",
  "data": {
    "id": 1,
    "nik": "3173011234567890",
    "full_name": "John Doe",
    "legal_name": "John Doe",
    "place_of_birth": "Jakarta",
    "date_of_birth": "1990-01-15T00:00:00Z",
    "salary": 5000000,
    "created_at": "2025-02-05T10:30:00Z",
    "updated_at": "2025-02-05T10:30:00Z"
  }
}
```

**Validasi yang Dilakukan:**
- ✓ NIK harus 16 digit
- ✓ Nama lengkap tidak boleh kosong
- ✓ Nama sah tidak boleh kosong
- ✓ Gaji minimum 1.000.000 rupiah
- ✓ Gaji tidak boleh negatif

**Error Response:**
```json
{
  "error": "gaji minimum 1 juta rupiah"
}
```

### 2. Ambil Detail Konsumen

**URL:** `GET /api/consumers/get?id=1`

**Contoh dengan cURL:**
```bash
curl http://localhost:8080/api/consumers/get?id=1
```

**Response (200 OK):**
```json
{
  "id": 1,
  "nik": "3173011234567890",
  "full_name": "John Doe",
  "legal_name": "John Doe",
  "place_of_birth": "Jakarta",
  "date_of_birth": "1990-01-15T00:00:00Z",
  "salary": 5000000,
  "ktp_photo": null,
  "selfie_photo": null,
  "created_at": "2025-02-05T10:30:00Z",
  "updated_at": "2025-02-05T10:30:00Z"
}
```

---

## Endpoint Manajemen Batas Kredit

### 1. Tetapkan Batas Kredit

**URL:** `POST /api/consumers/limits`

**Request Body:**
```json
{
  "consumer_id": 1,
  "tenor": 6,
  "limit_amount": 2000000
}
```

**Contoh dengan cURL:**
```bash
curl -X POST http://localhost:8080/api/consumers/limits \
  -H "Content-Type: application/json" \
  -d '{
    "consumer_id": 1,
    "tenor": 6,
    "limit_amount": 2000000
  }'
```

**Response (201 Created):**
```json
{
  "message": "Limit assigned successfully",
  "data": {
    "id": 1,
    "consumer_id": 1,
    "tenor": 6,
    "limit_amount": 2000000,
    "used_amount": 0,
    "created_at": "2025-02-05T10:35:00Z",
    "updated_at": "2025-02-05T10:35:00Z"
  }
}
```

**Tenor yang Diizinkan:** 1, 2, 3, 6 bulan

**Validasi yang Dilakukan:**
- ✓ Tenor harus 1, 2, 3, atau 6 bulan
- ✓ Jumlah limit harus > 0
- ✓ Consumer ID harus valid (harus terdaftar)

### 2. Ambil Batas Kredit Konsumen

**URL:** `GET /api/consumers/limits/get?id=1`

**Contoh dengan cURL:**
```bash
curl http://localhost:8080/api/consumers/limits/get?id=1
```

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "consumer_id": 1,
    "tenor": 1,
    "limit_amount": 1000000,
    "used_amount": 500000,
    "created_at": "2025-02-05T10:35:00Z",
    "updated_at": "2025-02-05T10:35:00Z"
  },
  {
    "id": 2,
    "consumer_id": 1,
    "tenor": 6,
    "limit_amount": 2000000,
    "used_amount": 0,
    "created_at": "2025-02-05T10:35:00Z",
    "updated_at": "2025-02-05T10:35:00Z"
  }
]
```

### 3. Menghitung Batas Tersedia

**Formula:**
```
Batas Tersedia = limit_amount - used_amount
```

**Contoh:**
```
Limit Total      : 2.000.000
Telah Digunakan  : 500.000
Batas Tersedia   : 1.500.000
```

---

## Endpoint Manajemen Transaksi

### 1. Buat Transaksi Baru

**URL:** `POST /api/transactions`

**Request Body:**
```json
{
  "consumer_id": 1,
  "contract_number": "CONT-001-2025",
  "tenor": 6,
  "otr": 1500000,
  "admin_fee": 50000,
  "installment_amount": 250000,
  "interest_amount": 100000,
  "asset_name": "TV Samsung 55 Inch"
}
```

**Contoh dengan cURL:**
```bash
curl -X POST http://localhost:8080/api/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "consumer_id": 1,
    "contract_number": "CONT-001-2025",
    "tenor": 6,
    "otr": 1500000,
    "admin_fee": 50000,
    "installment_amount": 250000,
    "interest_amount": 100000,
    "asset_name": "TV Samsung 55 Inch"
  }'
```

**Response (201 Created):**
```json
{
  "message": "Transaction created successfully",
  "data": {
    "id": 1,
    "consumer_id": 1,
    "contract_number": "CONT-001-2025",
    "tenor": 6,
    "otr": 1500000,
    "admin_fee": 50000,
    "installment_amount": 250000,
    "interest_amount": 100000,
    "asset_name": "TV Samsung 55 Inch",
    "status": "ACTIVE",
    "created_at": "2025-02-05T10:40:00Z",
    "updated_at": "2025-02-05T10:40:00Z"
  }
}
```

**Validasi yang Dilakukan:**
- ✓ Consumer harus terdaftar
- ✓ Nomor kontrak harus unik
- ✓ Tenor harus 1, 2, 3, atau 6
- ✓ OTR harus > 0
- ✓ Batas kredit harus cukup
- ✓ **Penanganan Konkurensi**: Jika dua transaksi datang bersamaan, keduanya tidak akan melebihi batas

**Error Response - Batas Tidak Cukup:**
```json
{
  "error": "limit tidak cukup untuk transaksi ini"
}
```

### 2. Ambil Detail Transaksi

**URL:** `GET /api/transactions/get?id=1`

**Contoh dengan cURL:**
```bash
curl http://localhost:8080/api/transactions/get?id=1
```

**Response (200 OK):**
```json
{
  "id": 1,
  "consumer_id": 1,
  "contract_number": "CONT-001-2025",
  "tenor": 6,
  "otr": 1500000,
  "admin_fee": 50000,
  "installment_amount": 250000,
  "interest_amount": 100000,
  "asset_name": "TV Samsung 55 Inch",
  "status": "ACTIVE",
  "created_at": "2025-02-05T10:40:00Z",
  "updated_at": "2025-02-05T10:40:00Z"
}
```

### 3. Ambil Semua Transaksi Konsumen

**URL:** `GET /api/transactions/consumer?id=1`

**Contoh dengan cURL:**
```bash
curl http://localhost:8080/api/transactions/consumer?id=1
```

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "consumer_id": 1,
    "contract_number": "CONT-001-2025",
    "tenor": 6,
    "otr": 1500000,
    "status": "ACTIVE",
    "created_at": "2025-02-05T10:40:00Z"
  },
  {
    "id": 2,
    "consumer_id": 1,
    "contract_number": "CONT-002-2025",
    "tenor": 3,
    "otr": 800000,
    "status": "COMPLETED",
    "created_at": "2025-02-04T15:20:00Z"
  }
]
```

### 4. Perbarui Status Transaksi

**URL:** `PUT /api/transactions/status?id=1`

**Request Body:**
```json
{
  "status": "COMPLETED"
}
```

**Status yang Diizinkan:**
- `ACTIVE` - Transaksi sedang berlangsung
- `COMPLETED` - Transaksi selesai
- `DEFAULTED` - Transaksi macet/gagal

**Contoh dengan cURL:**
```bash
curl -X PUT http://localhost:8080/api/transactions/status?id=1 \
  -H "Content-Type: application/json" \
  -d '{
    "status": "COMPLETED"
  }'
```

**Response (200 OK):**
```json
{
  "message": "Transaction status updated successfully"
}
```

---

## Testing & Debugging

### 1. Menjalankan Unit Test

**Jalankan Semua Test:**
```bash
go test ./...
```

**Jalankan Test dengan Verbose:**
```bash
go test -v ./...
```

**Jalankan Test Spesifik:**
```bash
go test -v -run TestRegisterConsumer_Valid ./...
```

**Lihat Coverage:**
```bash
go test -cover ./...
```

**Detail Coverage per File:**
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 2. Debugging dengan Logs

**Lihat Log Output:**
```bash
# Jika menggunakan Docker Compose
docker-compose logs -f app

# Jika running lokal, output langsung di terminal
go run main.go
```

**Contoh Log yang Dilihat:**
```
2025/02/05 10:30:00 ✓ Database terhubung dengan sukses
2025/02/05 10:30:05 ✓ Logika Bisnis OK. Menyimpan konsumen...
2025/02/05 10:30:10 Upaya SQL injection mencurigakan dalam parameter: id='1 OR 1=1'
```

### 3. Testing dengan Postman

**Import Collection:**
1. Buat folder baru: "XYZ Multifinance"
2. Buat requests untuk setiap endpoint

**Saran Request Order:**
1. Health Check
2. Register Consumer
3. Assign Limit
4. Create Transaction
5. Update Transaction Status

### 4. Testing dengan Thunder Client (VS Code)

**Instalasi Extension:**
1. Buka VS Code
2. Go to Extensions
3. Cari "Thunder Client"
4. Click Install

**Membuat Request:**
1. Alt+Cmd+R (macOS) atau Ctrl+Shift+R (Windows/Linux)
2. Masukkan URL dan method
3. Tambah headers dan body
4. Click "Send"

---

## Troubleshooting

### 1. Koneksi Database Gagal

**Error:**
```
Gagal koneksi database: [SSL: certificate verify failed]
```

**Solusi:**
```bash
# Verifikasi file ca.pem ada
ls -la ca.pem

# Pastikan path benar di config/database.go
# Check .env settings
cat .env
```

### 2. Port 8080 Sudah Digunakan

**Error:**
```
listen tcp :8080: bind: address already in use
```

**Solusi:**
```bash
# Cek proses yang menggunakan port 8080
netstat -tlnp | grep 8080
lsof -i :8080

# Kill process (ganti PID sesuai hasil di atas)
kill -9 <PID>

# Atau ubah port di .env
API_PORT=8081
```

### 3. Migration Database Gagal

**Error:**
```
Gagal melakukan migration: [SQL syntax error]
```

**Solusi:**
```bash
# Jalankan schema secara manual
mysql -u xyz_user -p xyz_multifinance < database_schema.sql

# Atau gunakan GORM dengan mode safe
# Edit config/database.go dan set DryRun: true
```

### 4. NIK Format Invalid

**Error Response:**
```json
{
  "error": "NIK must contain only numbers"
}
```

**Solusi:**
- NIK harus 16 digit angka
- Contoh valid: `3173011234567890`
- Contoh invalid: `3173-011-234-567-890` (ada karakter khusus)

### 5. Gaji di Bawah Minimum

**Error Response:**
```json
{
  "error": "gaji minimum 1 juta rupiah"
}
```

**Solusi:**
- Masukkan gaji minimal 1.000.000
- Contoh: `"salary": 1000000`

### 6. Tenor Tidak Valid

**Error Response:**
```json
{
  "error": "tenor harus 1, 2, 3, atau 6 bulan"
}
```

**Solusi:**
- Gunakan hanya tenor: 1, 2, 3, atau 6
- Contoh valid: `"tenor": 6`
- Contoh invalid: `"tenor": 12` atau `"tenor": 4`

### 7. Docker Container Tidak Jalan

**Error:**
```
docker: Error response from daemon: ...
```

**Solusi:**
```bash
# Rebuild tanpa cache
docker-compose build --no-cache

# Cek logs
docker-compose logs app

# Hapus volume dan rebuild
docker-compose down -v
docker-compose up -d
```

---

## Best Practices

### 1. Keamanan

- ✅ Selalu gunakan HTTPS di produksi (bukan HTTP)
- ✅ Jangan hardcode kredensial database di code
- ✅ Gunakan .env untuk konfigurasi sensitif
- ✅ Jangan expose error detail ke client
- ✅ Validasi semua input dari user

### 2. Performance

- ✅ Gunakan pagination untuk list besar
- ✅ Cache hasil query yang sering diakses
- ✅ Gunakan index database untuk query frequent
- ✅ Monitor slow query dengan logging
- ✅ Limit request rate dengan rate limiting

### 3. Data Management

- ✅ Gunakan soft delete untuk data penting
- ✅ Backup database secara berkala
- ✅ Implementasi audit trail untuk perubahan
- ✅ Archive old data ke cold storage
- ✅ Encrypt sensitive data (NIK, foto)

### 4. Testing

- ✅ Tulis unit test untuk business logic
- ✅ Gunakan mock untuk external services
- ✅ Test edge cases dan error scenarios
- ✅ Target coverage minimal 80%
- ✅ Integrasi test dalam CI/CD pipeline

### 5. Monitoring & Logging

- ✅ Log semua request dan response
- ✅ Monitor resource usage (CPU, Memory)
- ✅ Set up alerts untuk error rate tinggi
- ✅ Centralize logs (ELK, Datadog, dll)
- ✅ Monitor database performance

### 6. Dokumentasi

- ✅ Update README setiap ada perubahan
- ✅ Document API dengan OpenAPI/Swagger
- ✅ Maintein changelog yang detail
- ✅ Dokumentasi konfigurasi production
- ✅ Create runbooks untuk operations

---

## Contoh Skenario Lengkap

### Skenario: Pelanggan Baru Membeli TV dengan Cicilan

**Step 1: Daftarkan Pelanggan Baru**
```bash
curl -X POST http://localhost:8080/api/consumers \
  -H "Content-Type: application/json" \
  -d '{
    "nik": "3173019876543210",
    "full_name": "Siti Nurhaliza",
    "legal_name": "Siti Nurhaliza",
    "place_of_birth": "Bandung",
    "date_of_birth": "1992-06-10T00:00:00Z",
    "salary": 8000000
  }'
```

**Response:**
```json
{
  "message": "Consumer registered successfully",
  "data": { "id": 2, ... }
}
```

**Step 2: Tetapkan Batas Kredit untuk 6 Bulan**
```bash
curl -X POST http://localhost:8080/api/consumers/limits \
  -H "Content-Type: application/json" \
  -d '{
    "consumer_id": 2,
    "tenor": 6,
    "limit_amount": 3000000
  }'
```

**Step 3: Buat Transaksi Pembelian TV**
```bash
curl -X POST http://localhost:8080/api/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "consumer_id": 2,
    "contract_number": "CONT-TV-001-2025",
    "tenor": 6,
    "otr": 2500000,
    "admin_fee": 100000,
    "installment_amount": 433334,
    "interest_amount": 100000,
    "asset_name": "Smart TV Samsung 65 Inch 4K UHD"
  }'
```

**Step 4: Cek Status Transaksi**
```bash
curl http://localhost:8080/api/transactions/get?id=1
```

**Step 5: Setelah 6 Bulan, Update Status ke COMPLETED**
```bash
curl -X PUT http://localhost:8080/api/transactions/status?id=1 \
  -H "Content-Type: application/json" \
  -d '{
    "status": "COMPLETED"
  }'
```

---

## Support & Help

Jika mengalami masalah:

1. **Cek dokumentasi:** Lihat README.md
2. **Review changelog:** Lihat full-changelog.md
3. **Run tests:** `go test -v ./...`
4. **Check logs:** Lihat output application
5. **Contact team:** Email dev@xyz-multifinance.com

---

**Versi**: 1.0.0  
**Pembaruan Terakhir**: 5 Februari 2025  
**Status**: Siap Produksi
