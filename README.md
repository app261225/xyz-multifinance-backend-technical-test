# API Backend PT XYZ Multifinance

API Backend yang siap produksi untuk PT XYZ Multifinance dibangun dengan Go. Fitur meliputi Clean Architecture, transaksi yang sesuai dengan ACID dengan penanganan konkurensi, dan standar keamanan OWASP.

## Gambaran Umum

PT XYZ Multifinance adalah salah satu perusahaan pembiayaan terbesar di Indonesia di sektor White Goods, Motor, dan Otomotif. API Backend ini menyediakan solusi yang dapat diskalakan, mudah dipelihara, andal, dan aman untuk mengelola:

- **Pendaftaran & Manajemen Konsumen**: Menangani data KTP pelanggan, informasi pribadi, dan foto
- **Batas Kredit**: Mengelola batas kredit berbasis tenor (1, 2, 3, 6 bulan) per konsumen
- **Transaksi Pembiayaan**: Merekam dan melacak semua transaksi pembiayaan dengan perhitungan angsuran
- **Penanganan Transaksi Konkuren**: Pemrosesan transaksi konkuren yang aman dengan kunci berbasis mutex

## Struktur Proyek

```
xyz-multifinance-backend-technical-test/
├── config/              # Konfigurasi & setup Database
│   └── database.go      # Koneksi database dengan SSL/TLS
├── internal/
│   ├── handler/         # HTTP Handlers (endpoint API)
│   │   ├── consumer_handler.go
│   │   └── transaction_handler.go
│   ├── middleware/      # Security middleware (OWASP)
│   │   └── security.go
│   ├── model/           # Model data (entities)
│   │   └── entity.go
│   ├── repository/      # Lapisan akses data
│   │   └── consumer_repository.go
│   └── usecase/         # Lapisan logika bisnis
│       ├── consumer_usecase.go
│       └── consumer_usecase_test.go
├── main.go              # Titik masuk aplikasi
├── go.mod              # Dependensi Go
├── Dockerfile          # Container image
├── docker-compose.yml  # Environment lokal
├── database_schema.sql # Schema database dengan ACID
└── README.md           # File ini
```

## Arsitektur

### Implementasi Clean Architecture

Aplikasi mengikuti prinsip Clean Architecture dengan pemisahan tanggung jawab yang jelas:

```
┌─────────────────────────────────────────────────────┐
│         HTTP Handlers (Presentasi)                  │
├─────────────────────────────────────────────────────┤
│         Middleware (Keamanan)                       │
├─────────────────────────────────────────────────────┤
│         Usecase (Logika Bisnis)                     │
├─────────────────────────────────────────────────────┤
│         Repository (Akses Data)                     │
├─────────────────────────────────────────────────────┤
│         Database (MySQL dengan SSL/TLS)             │
└─────────────────────────────────────────────────────┘
```

### Komponen Utama

1. **Handlers**: Endpoint REST API untuk konsumen, batas, dan transaksi
2. **Middleware**: Header keamanan, validasi input, CORS, dan pembatasan laju
3. **Usecase**: Logika bisnis dengan validasi dan kepatuhan ACID
4. **Repository**: Abstraksi akses data dengan GORM ORM
5. **Models**: Struktur data yang diketik dengan kuat dengan tag validasi

## Fitur

### 1. Keamanan (Perlindungan OWASP Top 10)

API menerapkan minimal 3 ukuran keamanan OWASP:

#### A. Header Keamanan
- **X-Content-Type-Options: nosniff** - Mencegah MIME type sniffing (proteksi XSS)
- **X-Frame-Options: DENY** - Mencegah serangan Clickjacking
- **X-XSS-Protection: 1; mode=block** - Perlindungan XSS warisan
- **Content-Security-Policy** - Mencegah eksekusi skrip inline
- **Strict-Transport-Security** - Memaksa HTTPS (OWASP A02:2021)
- **Referrer-Policy** - Mengontrol informasi referrer

#### B. Validasi Input
- Pencegahan SQL injection melalui query terparameterisasi dan sanitasi input
- Validasi permintaan di tingkat middleware
- Validasi Content-Type (application/json dipaksakan)
- Validasi format NIK (16 digit ID Indonesia)

#### C. Perlindungan CORS
- Origin yang diizinkan dikonfigurasi
- Pembatasan metode (GET, POST, PUT, DELETE)
- Validasi header

### 2. Kepatuhan ACID

Semua transaksi keuangan sesuai dengan ACID:

- **Atomicity**: Eksekusi transaksi semua-atau-tidak-sama-sekali dengan prosedur tersimpan
- **Consistency**: Validasi data di tingkat aplikasi dan database
- **Isolation**: Tingkat isolasi transaksi diatur ke READ_COMMITTED
- **Durability**: Penyimpanan persisten dengan mesin MySQL InnoDB

### 3. Penanganan Transaksi Konkuren

`TransactionUsecase` menerapkan kontrol konkurensi berbasis mutex:

```go
func (u *transactionUsecase) CreateTransaction(transaction *model.Transaction) error {
    u.mu.Lock()
    defer u.mu.Unlock()
    // Verifikasi batas dan pembaruan atomik
    // Mencegah kondisi balapan dalam permintaan konkuren
}
```

### 4. Fitur Database

- **Constraint Foreign Key**: Integritas referensial
- **Constraint Unik**: Cegah kontrak duplikat
- **Check Constraint**: Validasi domain di tingkat database
- **Index**: Optimasi performa kueri
- **View**: Kueri business intelligence
- **Stored Procedure**: Operasi transaksi atomik

## Endpoint API

### Manajemen Konsumen

```bash
# Daftarkan konsumen baru
POST /api/consumers
Content-Type: application/json

{
  "nik": "3173011234567890",
  "full_name": "John Doe",
  "legal_name": "John Doe",
  "place_of_birth": "Jakarta",
  "date_of_birth": "1990-01-15T00:00:00Z",
  "salary": 5000000
}

# Dapatkan detail konsumen
GET /api/consumers/get?id=1

# Tetapkan batas kredit
POST /api/consumers/limits
{
  "consumer_id": 1,
  "tenor": 6,
  "limit_amount": 2000000
}

# Dapatkan batas konsumen
GET /api/consumers/limits/get?id=1
```

### Manajemen Transaksi

```bash
# Buat transaksi (dengan pemeriksaan batas konkuren)
POST /api/transactions
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

# Dapatkan transaksi
GET /api/transactions/get?id=1

# Dapatkan transaksi konsumen
GET /api/transactions/consumer?id=1

# Perbarui status transaksi
PUT /api/transactions/status?id=1
{
  "status": "COMPLETED"
}

# Health check
GET /health
```

## Instalasi & Setup

### Prasyarat

- Go 1.25+
- MySQL 8.0+ (untuk produksi)
- Docker & Docker Compose (opsional, untuk setup containerized)
- Sertifikat ca.pem (untuk koneksi database SSL/TLS)

### Pengembangan Lokal (tanpa Docker)

```bash
# 1. Clone repository
git clone https://github.com/yourusername/xyz-multifinance-backend-technical-test.git
cd xyz-multifinance-backend-technical-test

# 2. Install dependensi
go mod download

# 3. Setup environment
cp .env.example .env
# Edit .env dengan kredensial database Anda

# 4. Buat schema database
mysql -u root -p < database_schema.sql

# 5. Jalankan aplikasi
go run main.go

# 6. Aplikasi akan dimulai di http://localhost:8080
curl http://localhost:8080/health
```

### Setup Docker (Direkomendasikan)

```bash
# Build dan jalankan dengan Docker Compose
docker-compose up -d

# Lihat log
docker-compose logs -f app

# Hentikan layanan
docker-compose down
```

## Testing

### Menjalankan Unit Test

```bash
# Jalankan semua test
go test ./...

# Jalankan dengan output verbose
go test -v ./...

# Jalankan file test spesifik
go test -v ./internal/usecase/consumer_usecase_test.go

# Jalankan dengan coverage
go test -cover ./...
```

### Cakupan Test

Proyek mencakup test komprehensif untuk:

- **Pendaftaran Konsumen**: Data valid/tidak valid, validasi format NIK, pemeriksaan gaji
- **Penugasan Batas**: Validasi tenor, validasi jumlah, verifikasi konsumen
- **Pembuatan Transaksi**: Keamanan konkuren, penegakan batas, integritas data

### Contoh Unit Test

```bash
TestRegisterConsumer_Valid
TestRegisterConsumer_InvalidNIK
TestRegisterConsumer_MissingFields
TestRegisterConsumer_LowSalary
TestGetConsumer
TestAssignLimit_Valid
TestAssignLimit_InvalidTenor
TestAssignLimit_InvalidAmount
```

## Schema Database

### Tabel

#### consumers
- Tabel utama menyimpan informasi pribadi konsumen
- Field: NIK, Nama Lengkap, Nama Sah, TTL, Gaji, data Foto
- Index: NIK (unik), created_at, deleted_at

#### consumer_limits
- Batas kredit per tenor (1, 2, 3, 6 bulan)
- Melacak limit_amount dan used_amount untuk setiap tenor
- Memaksa validasi tenor dan kombinasi unik consumer-tenor

#### transactions
- Transaksi keuangan dengan detail angsuran
- Link ke konsumen via consumer_id
- Lacak status: ACTIVE, COMPLETED, DEFAULTED
- Index untuk konsumen, nomor kontrak, dan status

### Prosedur ACID

#### sp_create_transaction
Prosedur atomik untuk membuat transaksi dengan penegakan batas:
- Validasi keberadaan konsumen
- Periksa keunikan nomor kontrak
- Verifikasi batas kredit yang cukup
- Perbarui penggunaan batas secara atomik
- Buat catatan transaksi

## Implementasi Git Flow

Proyek mengadopsi Git Flow untuk manajemen versi:

```
main (rilis produksi)
  └── release/* (calon rilis)
develop (cabang integrasi)
  ├── feature/* (fitur baru)
  ├── bugfix/* (perbaikan bug)
  └── hotfix/* (perbaikan urgen produksi)
```

## Deployment

### Build Docker Image

```bash
docker build -t xyz-multifinance:1.0.0 .
```

### Push ke Registry

```bash
docker tag xyz-multifinance:1.0.0 your-registry/xyz-multifinance:1.0.0
docker push your-registry/xyz-multifinance:1.0.0
```

### Variabel Environment

```
DB_USER=xyz_user
DB_PASS=xyz_password
DB_HOST=localhost
DB_PORT=3306
DB_NAME=xyz_multifinance
API_PORT=8080
```

## Pertimbangan Performa

1. **Connection Pooling**: GORM menangani connection pooling database
2. **Optimasi Kueri**: Field yang diindex untuk pencarian cepat
3. **Permintaan Konkuren**: Sinkronisasi berbasis mutex untuk keamanan transaksi
4. **Transaksi Database**: Operasi atomik via prosedur tersimpan
5. **Caching**: Pertimbangkan integrasi Redis untuk versi mendatang

## Checklist Keamanan

- [x] Koneksi database SSL/TLS
- [x] Validasi dan sanitasi input
- [x] Pencegahan SQL injection
- [x] Proteksi XSS via CSP
- [x] Proteksi CSRF (CORS)
- [x] Header aman
- [x] Kepatuhan transaksi ACID
- [x] Keamanan transaksi konkuren
- [x] Validasi format NIK
- [x] Penegakan gaji minimum

## Monitoring & Logging

Aplikasi mencakup logging di titik-titik kritis:

```go
log.Println("✓ Database terhubung dengan sukses")
log.Println("✓ Logika Bisnis OK. Menyimpan konsumen...")
log.Printf("Upaya SQL injection mencurigakan dalam parameter: %s=%s\n", key, value)
```

Untuk produksi, integrasikan dengan:
- ELK Stack (Elasticsearch, Logstash, Kibana)
- Prometheus + Grafana
- Datadog atau New Relic

## Troubleshooting

### Masalah Koneksi Database
```bash
# Periksa layanan MySQL
systemctl status mysql

# Verifikasi sertifikat ca.pem
ls -la ca.pem

# Test koneksi database
mysql -h localhost -u xyz_user -p
```

### Port Sudah Digunakan
```bash
# Periksa penggunaan port
netstat -tlnp | grep 8080

# Hentikan proses
kill -9 <PID>
```

### Masalah Docker
```bash
# Periksa log container
docker logs xyz-multifinance-app

# Rebuild tanpa cache
docker-compose build --no-cache
```

## Peningkatan Mendatang

1. **Autentikasi**: Autentikasi API berbasis JWT
2. **Otorisasi**: Kontrol akses berbasis peran (RBAC)
3. **Caching**: Redis untuk data yang sering diakses
4. **Message Queue**: Kafka untuk pemrosesan transaksi asinkron
5. **Metrics**: Metrics Prometheus untuk monitoring
6. **Dokumentasi**: Dokumentasi OpenAPI/Swagger
7. **Versioning API**: Dukungan untuk multiple versi API
8. **Replika Database**: Setup master-slave untuk high availability

## Kontribusi

1. Clone repository
2. Buat cabang fitur: `git checkout -b feature/fitur-Anda`
3. Commit perubahan: `git commit -am 'Tambah fitur baru'`
4. Push ke cabang: `git push origin feature/fitur-Anda`
5. Submit pull request

## Lisensi

Proyek ini dilisensikan di bawah MIT License - lihat file LICENSE untuk detail.

## Kontak

- Email: salambudiarto@gmail.com

## Changelog

Lihat [full-changelog.md](full-changelog.md) untuk sejarah versi terperinci dan catatan rilis.

---

**Versi**: 1.0.0  
**Pembaruan Terakhir**: Februari 2025  
**Status**: Siap Produksi