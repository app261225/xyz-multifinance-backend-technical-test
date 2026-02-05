# Panduan Testing API dengan cURL

Panduan lengkap untuk melakukan testing terhadap API Backend PT XYZ Multifinance menggunakan cURL command line.

## Daftar Isi
1. [Persiapan Testing](#persiapan-testing)
2. [Health Check](#health-check)
3. [Testing Manajemen Konsumen](#testing-manajemen-konsumen)
4. [Testing Manajemen Batas Kredit](#testing-manajemen-batas-kredit)
5. [Testing Manajemen Transaksi](#testing-manajemen-transaksi)
6. [Testing Error Cases](#testing-error-cases)
7. [Script Testing Otomatis](#script-testing-otomatis)
8. [Validasi Hasil Testing](#validasi-hasil-testing)

---

## Persiapan Testing

### Pastikan Server Berjalan

```bash
# Terminal 1: Jalankan server
cd d:\SIGMA-TECH\git\xyz-multifinance-backend-technical-test
go run main.go

# Terminal 2: Lakukan testing
# (Gunakan terminal terpisah untuk testing)
```

### Verifikasi cURL Tersedia

```bash
# Windows
where curl

# macOS/Linux
which curl
```

### Optional: Simpan Config Variable

```bash
# Simpan dalam variable (opsional, untuk testing interaktif)
API_URL="http://localhost:8080"
CONTENT_TYPE="Content-Type: application/json"
```

---

## Health Check

### 1. Test Endpoint Kesehatan Server

```bash
curl -i http://localhost:8080/health
```

**Expected Response (200 OK):**
```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Wed, 05 Feb 2026 10:00:00 GMT
Content-Length: 42

{"status":"healthy","version":"1.0.0"}
```

**Penjelasan:**
- `Status 200 OK` = Server berjalan dengan baik
- `"status":"healthy"` = Aplikasi siap melayani
- `"version":"1.0.0"` = Versi API saat ini

---

## Testing Manajemen Konsumen

### 1. Register Konsumen Pertama

```bash
curl -X POST http://localhost:8080/api/consumers \
  -H "Content-Type: application/json" \
  -d '{
    "nik": "3173011234567890",
    "full_name": "Budi Santoso",
    "legal_name": "Budi Santoso",
    "place_of_birth": "Jakarta",
    "date_of_birth": "1990-05-15T00:00:00Z",
    "salary": 5000000
  }'
```

**Penjelasan Parameter:**
- `nik` - 16 digit nomor induk kependudukan (wajib unik)
- `full_name` - Nama lengkap konsumen
- `legal_name` - Nama sah di dokumen resmi
- `place_of_birth` - Tempat lahir
- `date_of_birth` - Tanggal lahir (format ISO 8601)
- `salary` - Gaji per bulan (minimal 1.000.000)

**Expected Response (201 Created):**
```json
{
  "message": "Consumer registered successfully",
  "data": {
    "id": 1,
    "nik": "3173011234567890",
    "full_name": "Budi Santoso",
    "legal_name": "Budi Santoso",
    "place_of_birth": "Jakarta",
    "date_of_birth": "1990-05-15T00:00:00Z",
    "salary": 5000000,
    "created_at": "2026-02-05T10:00:00Z",
    "updated_at": "2026-02-05T10:00:00Z"
  }
}
```

**✓ Verifikasi:**
- HTTP Status 201 (bukan 200, 400, atau 500)
- Response memiliki `message` field
- Response memiliki `data` dengan field lengkap
- `id` terbuat otomatis (biasanya 1 untuk konsumen pertama)
- Timestamp `created_at` dan `updated_at` terisi

### 2. Register Konsumen Kedua

```bash
curl -X POST http://localhost:8080/api/consumers \
  -H "Content-Type: application/json" \
  -d '{
    "nik": "3173019876543210",
    "full_name": "Siti Nurhaliza",
    "legal_name": "Siti Nurhaliza",
    "place_of_birth": "Bandung",
    "date_of_birth": "1992-03-22T00:00:00Z",
    "salary": 8000000
  }'
```

**Expected Response:**
```json
{
  "message": "Consumer registered successfully",
  "data": {
    "id": 2,
    ...
  }
}
```

**✓ Verifikasi:**
- `id` adalah 2 (increment dari sebelumnya)

### 3. Ambil Detail Konsumen

```bash
# Ambil konsumen pertama
curl http://localhost:8080/api/consumers/get?id=1
```

**Expected Response (200 OK):**
```json
{
  "id": 1,
  "nik": "3173011234567890",
  "full_name": "Budi Santoso",
  "legal_name": "Budi Santoso",
  "place_of_birth": "Jakarta",
  "date_of_birth": "1990-05-15T00:00:00Z",
  "salary": 5000000,
  "created_at": "2026-02-05T10:00:00Z",
  "updated_at": "2026-02-05T10:00:00Z"
}
```

**✓ Verifikasi:**
- HTTP Status 200
- Data matching dengan data yang disimpan sebelumnya
- Semua field terisi lengkap

```bash
# Ambil konsumen kedua
curl http://localhost:8080/api/consumers/get?id=2
```

**✓ Verifikasi:**
- Response data untuk konsumen kedua (Siti Nurhaliza)
- `id` adalah 2

---

## Testing Manajemen Batas Kredit

### 1. Tetapkan Batas Kredit Tenor 1 Bulan untuk Konsumen 1

```bash
curl -X POST http://localhost:8080/api/consumers/limits \
  -H "Content-Type: application/json" \
  -d '{
    "consumer_id": 1,
    "tenor": 1,
    "limit_amount": 1000000
  }'
```

**Expected Response (201 Created):**
```json
{
  "message": "Limit assigned successfully",
  "data": {
    "id": 1,
    "consumer_id": 1,
    "tenor": 1,
    "limit_amount": 1000000,
    "used_amount": 0,
    "created_at": "2026-02-05T10:05:00Z",
    "updated_at": "2026-02-05T10:05:00Z"
  }
}
```

**✓ Verifikasi:**
- HTTP Status 201
- `tenor` adalah 1
- `limit_amount` adalah 1.000.000
- `used_amount` dimulai dari 0

### 2. Tetapkan Batas Kredit Tenor 6 Bulan untuk Konsumen 1

```bash
curl -X POST http://localhost:8080/api/consumers/limits \
  -H "Content-Type: application/json" \
  -d '{
    "consumer_id": 1,
    "tenor": 6,
    "limit_amount": 3000000
  }'
```

**Expected Response:**
```json
{
  "message": "Limit assigned successfully",
  "data": {
    "id": 2,
    "consumer_id": 1,
    "tenor": 6,
    "limit_amount": 3000000,
    "used_amount": 0,
    ...
  }
}
```

**✓ Verifikasi:**
- `id` adalah 2 (bukan duplikat id 1)
- `tenor` adalah 6 bulan

### 3. Tetapkan Batas Kredit untuk Konsumen 2

```bash
curl -X POST http://localhost:8080/api/consumers/limits \
  -H "Content-Type: application/json" \
  -d '{
    "consumer_id": 2,
    "tenor": 3,
    "limit_amount": 5000000
  }'
```

**Expected Response:**
```json
{
  "message": "Limit assigned successfully",
  "data": {
    "id": 3,
    "consumer_id": 2,
    "tenor": 3,
    "limit_amount": 5000000,
    ...
  }
}
```

### 4. Ambil Semua Batas Kredit Konsumen 1

```bash
curl http://localhost:8080/api/consumers/limits/get?id=1
```

**Expected Response (200 OK):**
```json
[
  {
    "id": 1,
    "consumer_id": 1,
    "tenor": 1,
    "limit_amount": 1000000,
    "used_amount": 0,
    "created_at": "2026-02-05T10:05:00Z",
    "updated_at": "2026-02-05T10:05:00Z"
  },
  {
    "id": 2,
    "consumer_id": 1,
    "tenor": 6,
    "limit_amount": 3000000,
    "used_amount": 0,
    "created_at": "2026-02-05T10:06:00Z",
    "updated_at": "2026-02-05T10:06:00Z"
  }
]
```

**✓ Verifikasi:**
- HTTP Status 200
- Response adalah array dengan 2 item
- Keduanya punya `consumer_id: 1`
- Tenor berbeda (1 dan 6)
- Limit amount sesuai input

### 5. Ambil Batas Kredit Konsumen 2

```bash
curl http://localhost:8080/api/consumers/limits/get?id=2
```

**Expected Response:**
```json
[
  {
    "id": 3,
    "consumer_id": 2,
    "tenor": 3,
    "limit_amount": 5000000,
    ...
  }
]
```

**✓ Verifikasi:**
- Array berisi 1 item (hanya satu limit untuk konsumen 2)
- `tenor` adalah 3

---

## Testing Manajemen Transaksi

### 1. Buat Transaksi Pertama (Tenor 1 Bulan)

```bash
curl -X POST http://localhost:8080/api/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "consumer_id": 1,
    "contract_number": "CONT-001-2026",
    "tenor": 1,
    "otr": 800000,
    "admin_fee": 30000,
    "installment_amount": 830000,
    "interest_amount": 0,
    "asset_name": "Laptop HP Pavilion 15"
  }'
```

**Expected Response (201 Created):**
```json
{
  "message": "Transaction created successfully",
  "data": {
    "id": 1,
    "consumer_id": 1,
    "contract_number": "CONT-001-2026",
    "tenor": 1,
    "otr": 800000,
    "admin_fee": 30000,
    "installment_amount": 830000,
    "interest_amount": 0,
    "asset_name": "Laptop HP Pavilion 15",
    "status": "ACTIVE",
    "created_at": "2026-02-05T10:10:00Z",
    "updated_at": "2026-02-05T10:10:00Z"
  }
}
```

**✓ Verifikasi:**
- HTTP Status 201
- `status` adalah "ACTIVE"
- `consumer_id` adalah 1
- Semua monetary field sesuai input
- `contract_number` unik dan sesuai input

**Penjelasan Transaksi:**
- OTR (On The Road) = harga barang = 800.000
- Admin Fee = 30.000
- Installment per bulan = 830.000 (OTR + Admin Fee)
- Limit yang digunakan = 800.000 (hanya OTR, bukan admin fee)
- Batas tersisa konsumen 1 tenor 1 bulan = 1.000.000 - 800.000 = 200.000

### 2. Buat Transaksi Kedua (Tenor 6 Bulan, OTR Besar)

```bash
curl -X POST http://localhost:8080/api/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "consumer_id": 1,
    "contract_number": "CONT-002-2026",
    "tenor": 6,
    "otr": 2500000,
    "admin_fee": 100000,
    "installment_amount": 433334,
    "interest_amount": 100000,
    "asset_name": "Smart TV Samsung 65 Inch 4K UHD"
  }'
```

**Expected Response:**
```json
{
  "message": "Transaction created successfully",
  "data": {
    "id": 2,
    "consumer_id": 1,
    "contract_number": "CONT-002-2026",
    "tenor": 6,
    "otr": 2500000,
    ...
    "status": "ACTIVE",
    ...
  }
}
```

**✓ Verifikasi:**
- HTTP Status 201
- `id` adalah 2
- `tenor` adalah 6
- Limit digunakan = 2.500.000 (dari limit 3.000.000 tenor 6 bulan)

### 3. Buat Transaksi untuk Konsumen 2

```bash
curl -X POST http://localhost:8080/api/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "consumer_id": 2,
    "contract_number": "CONT-003-2026",
    "tenor": 3,
    "otr": 4500000,
    "admin_fee": 150000,
    "installment_amount": 1550000,
    "interest_amount": 150000,
    "asset_name": "Microwave Panasonic 25L"
  }'
```

**Expected Response:**
```json
{
  "message": "Transaction created successfully",
  "data": {
    "id": 3,
    "consumer_id": 2,
    "contract_number": "CONT-003-2026",
    ...
  }
}
```

### 4. Ambil Detail Transaksi

```bash
# Ambil transaksi pertama
curl http://localhost:8080/api/transactions/get?id=1
```

**Expected Response:**
```json
{
  "id": 1,
  "consumer_id": 1,
  "contract_number": "CONT-001-2026",
  "tenor": 1,
  "otr": 800000,
  "admin_fee": 30000,
  "installment_amount": 830000,
  "interest_amount": 0,
  "asset_name": "Laptop HP Pavilion 15",
  "status": "ACTIVE",
  "created_at": "2026-02-05T10:10:00Z",
  "updated_at": "2026-02-05T10:10:00Z"
}
```

### 5. Ambil Semua Transaksi Konsumen 1

```bash
curl http://localhost:8080/api/transactions/consumer?id=1
```

**Expected Response:**
```json
[
  {
    "id": 1,
    "consumer_id": 1,
    "contract_number": "CONT-001-2026",
    "tenor": 1,
    "otr": 800000,
    "status": "ACTIVE",
    ...
  },
  {
    "id": 2,
    "consumer_id": 1,
    "contract_number": "CONT-002-2026",
    "tenor": 6,
    "otr": 2500000,
    "status": "ACTIVE",
    ...
  }
]
```

**✓ Verifikasi:**
- Array berisi 2 transaksi untuk konsumen 1
- Keduanya punya `consumer_id: 1`

### 6. Update Status Transaksi

```bash
# Update transaksi 1 ke COMPLETED
curl -X PUT http://localhost:8080/api/transactions/status?id=1 \
  -H "Content-Type: application/json" \
  -d '{
    "status": "COMPLETED"
  }'
```

**Expected Response (200 OK):**
```json
{
  "message": "Transaction status updated successfully"
}
```

**Verifikasi Update:**
```bash
# Ambil transaksi 1 lagi
curl http://localhost:8080/api/transactions/get?id=1
```

**Expected Response:**
```json
{
  "id": 1,
  "status": "COMPLETED",
  "updated_at": "2026-02-05T10:15:00Z",
  ...
}
```

**✓ Verifikasi:**
- `status` berubah menjadi "COMPLETED"
- `updated_at` timestamp berubah

### 7. Update Transaksi ke DEFAULTED

```bash
curl -X PUT http://localhost:8080/api/transactions/status?id=2 \
  -H "Content-Type: application/json" \
  -d '{
    "status": "DEFAULTED"
  }'
```

**Expected Response:**
```json
{
  "message": "Transaction status updated successfully"
}
```

---

## Testing Error Cases

### 1. Register Konsumen dengan Gaji di Bawah Minimum

```bash
curl -X POST http://localhost:8080/api/consumers \
  -H "Content-Type: application/json" \
  -d '{
    "nik": "3173011111111111",
    "full_name": "Test User",
    "legal_name": "Test User",
    "place_of_birth": "Jakarta",
    "date_of_birth": "1990-01-01T00:00:00Z",
    "salary": 500000
  }'
```

**Expected Error Response (400 Bad Request):**
```json
{
  "error": "gaji minimum 1 juta rupiah"
}
```

**✓ Verifikasi:**
- HTTP Status 400 (bukan 201)
- Error message jelas dan dalam bahasa Indonesia

### 2. Register dengan NIK Duplikat

```bash
curl -X POST http://localhost:8080/api/consumers \
  -H "Content-Type: application/json" \
  -d '{
    "nik": "3173011234567890",
    "full_name": "Duplicate NIK",
    "legal_name": "Duplicate NIK",
    "place_of_birth": "Jakarta",
    "date_of_birth": "1990-01-01T00:00:00Z",
    "salary": 5000000
  }'
```

**Expected Error Response (400 Bad Request):**
```json
{
  "error": "NIK sudah terdaftar"
}
```

### 3. Tetapkan Limit dengan Tenor Invalid

```bash
curl -X POST http://localhost:8080/api/consumers/limits \
  -H "Content-Type: application/json" \
  -d '{
    "consumer_id": 1,
    "tenor": 12,
    "limit_amount": 5000000
  }'
```

**Expected Error Response (400 Bad Request):**
```json
{
  "error": "tenor harus 1, 2, 3, atau 6 bulan"
}
```

### 4. Buat Transaksi dengan Limit Tidak Cukup

**Situasi:** Konsumen 1 tenor 1 bulan hanya punya limit 1.000.000, sudah pakai 800.000, sisa 200.000.
Coba transaksi OTR 300.000.

```bash
curl -X POST http://localhost:8080/api/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "consumer_id": 1,
    "contract_number": "CONT-OVER-001-2026",
    "tenor": 1,
    "otr": 300000,
    "admin_fee": 15000,
    "installment_amount": 315000,
    "interest_amount": 0,
    "asset_name": "Mouse Wireless"
  }'
```

**Expected Error Response (400 Bad Request):**
```json
{
  "error": "limit tidak cukup untuk transaksi ini"
}
```

**✓ Verifikasi:**
- HTTP Status 400
- Error message menunjukkan limit tidak cukup
- Transaksi tidak terbuat (cek dengan GET /api/transactions/get?id=X)

### 5. Ambil Konsumen yang Tidak Ada

```bash
curl http://localhost:8080/api/consumers/get?id=999
```

**Expected Error Response (404 Not Found atau 400):**
```json
{
  "error": "data not found"
}
```

### 6. Buat Transaksi dengan Consumer ID Invalid

```bash
curl -X POST http://localhost:8080/api/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "consumer_id": 999,
    "contract_number": "CONT-INVALID-2026",
    "tenor": 1,
    "otr": 500000,
    "admin_fee": 20000,
    "installment_amount": 520000,
    "interest_amount": 0,
    "asset_name": "Test Item"
  }'
```

**Expected Error Response (400 Bad Request):**
```json
{
  "error": "consumer tidak ditemukan"
}
```

### 7. Test Input Validation - SQL Injection Attempt

```bash
curl -X POST http://localhost:8080/api/consumers \
  -H "Content-Type: application/json" \
  -d '{
    "nik": "3173011234567890' OR '1'='1",
    "full_name": "Test",
    "legal_name": "Test",
    "place_of_birth": "Jakarta",
    "date_of_birth": "1990-01-01T00:00:00Z",
    "salary": 5000000
  }'
```

**Expected:**
- Request ditolak atau di-sanitize
- NIK tidak valid (tidak 16 digit angka)
- Error response: "NIK must contain only numbers"

**✓ Verifikasi:**
- Injeksi SQL tidak berhasil
- Data tidak corrupt di database

---

## Script Testing Otomatis

### Script 1: Quick Test (Windows PowerShell)

Simpan sebagai file `test.ps1`:

```powershell
# Quick Test Script untuk PT XYZ Multifinance API
$BASE_URL = "http://localhost:8080"

Write-Host "=== PT XYZ Multifinance API Testing ===" -ForegroundColor Green
Write-Host ""

# Test 1: Health Check
Write-Host "1. Testing Health Endpoint..." -ForegroundColor Cyan
$health = Invoke-WebRequest -Uri "$BASE_URL/health" | ConvertFrom-Json
Write-Host "✓ Status: $($health.status)" -ForegroundColor Green
Write-Host ""

# Test 2: Register Consumer
Write-Host "2. Registering Consumer..." -ForegroundColor Cyan
$consumer = @{
    nik = "3173011234567890"
    full_name = "Test User"
    legal_name = "Test User"
    place_of_birth = "Jakarta"
    date_of_birth = "1990-01-01T00:00:00Z"
    salary = 5000000
} | ConvertTo-Json

$resp = Invoke-WebRequest -Uri "$BASE_URL/api/consumers" `
    -Method POST `
    -Headers @{"Content-Type"="application/json"} `
    -Body $consumer
$respData = $resp.Content | ConvertFrom-Json
$CONSUMER_ID = $respData.data.id
Write-Host "✓ Consumer registered with ID: $CONSUMER_ID" -ForegroundColor Green
Write-Host ""

# Test 3: Assign Limit
Write-Host "3. Assigning Credit Limit..." -ForegroundColor Cyan
$limit = @{
    consumer_id = $CONSUMER_ID
    tenor = 6
    limit_amount = 3000000
} | ConvertTo-Json

$resp = Invoke-WebRequest -Uri "$BASE_URL/api/consumers/limits" `
    -Method POST `
    -Headers @{"Content-Type"="application/json"} `
    -Body $limit
Write-Host "✓ Credit limit assigned" -ForegroundColor Green
Write-Host ""

# Test 4: Create Transaction
Write-Host "4. Creating Transaction..." -ForegroundColor Cyan
$transaction = @{
    consumer_id = $CONSUMER_ID
    contract_number = "CONT-TEST-001-2026"
    tenor = 6
    otr = 2000000
    admin_fee = 100000
    installment_amount = 350000
    interest_amount = 100000
    asset_name = "Test Asset"
} | ConvertTo-Json

$resp = Invoke-WebRequest -Uri "$BASE_URL/api/transactions" `
    -Method POST `
    -Headers @{"Content-Type"="application/json"} `
    -Body $transaction
$transData = $resp.Content | ConvertFrom-Json
$TRX_ID = $transData.data.id
Write-Host "✓ Transaction created with ID: $TRX_ID" -ForegroundColor Green
Write-Host ""

Write-Host "=== All Tests Passed! ===" -ForegroundColor Green
```

Jalankan:
```powershell
powershell -ExecutionPolicy Bypass -File test.ps1
```

### Script 2: Comprehensive Test (Bash/Linux/macOS)

Simpan sebagai file `test.sh`:

```bash
#!/bin/bash

BASE_URL="http://localhost:8080"
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== PT XYZ Multifinance API Testing ===${NC}"
echo ""

# Test 1: Health Check
echo -e "${YELLOW}1. Testing Health Endpoint...${NC}"
HEALTH=$(curl -s $BASE_URL/health)
if echo "$HEALTH" | grep -q "healthy"; then
    echo -e "${GREEN}✓ Server is healthy${NC}"
else
    echo -e "${RED}✗ Server health check failed${NC}"
    exit 1
fi
echo ""

# Test 2: Register Consumer
echo -e "${YELLOW}2. Registering Consumer...${NC}"
CONSUMER_RESPONSE=$(curl -s -X POST $BASE_URL/api/consumers \
  -H "Content-Type: application/json" \
  -d '{
    "nik": "3173011234567890",
    "full_name": "Test User",
    "legal_name": "Test User",
    "place_of_birth": "Jakarta",
    "date_of_birth": "1990-01-01T00:00:00Z",
    "salary": 5000000
  }')

CONSUMER_ID=$(echo $CONSUMER_RESPONSE | grep -o '"id":[0-9]*' | head -1 | grep -o '[0-9]*')
if [ ! -z "$CONSUMER_ID" ]; then
    echo -e "${GREEN}✓ Consumer registered with ID: $CONSUMER_ID${NC}"
else
    echo -e "${RED}✗ Consumer registration failed${NC}"
    echo "Response: $CONSUMER_RESPONSE"
    exit 1
fi
echo ""

# Test 3: Assign Limit
echo -e "${YELLOW}3. Assigning Credit Limit...${NC}"
LIMIT_RESPONSE=$(curl -s -X POST $BASE_URL/api/consumers/limits \
  -H "Content-Type: application/json" \
  -d "{
    \"consumer_id\": $CONSUMER_ID,
    \"tenor\": 6,
    \"limit_amount\": 3000000
  }")

if echo "$LIMIT_RESPONSE" | grep -q "Limit assigned"; then
    echo -e "${GREEN}✓ Credit limit assigned${NC}"
else
    echo -e "${RED}✗ Limit assignment failed${NC}"
    exit 1
fi
echo ""

# Test 4: Create Transaction
echo -e "${YELLOW}4. Creating Transaction...${NC}"
TRX_RESPONSE=$(curl -s -X POST $BASE_URL/api/transactions \
  -H "Content-Type: application/json" \
  -d "{
    \"consumer_id\": $CONSUMER_ID,
    \"contract_number\": \"CONT-TEST-001-2026\",
    \"tenor\": 6,
    \"otr\": 2000000,
    \"admin_fee\": 100000,
    \"installment_amount\": 350000,
    \"interest_amount\": 100000,
    \"asset_name\": \"Test Asset\"
  }")

TRX_ID=$(echo $TRX_RESPONSE | grep -o '"id":[0-9]*' | head -1 | grep -o '[0-9]*')
if [ ! -z "$TRX_ID" ]; then
    echo -e "${GREEN}✓ Transaction created with ID: $TRX_ID${NC}"
else
    echo -e "${RED}✗ Transaction creation failed${NC}"
    exit 1
fi
echo ""

# Test 5: Get Transaction
echo -e "${YELLOW}5. Retrieving Transaction...${NC}"
GET_TRX=$(curl -s $BASE_URL/api/transactions/get?id=$TRX_ID)
if echo "$GET_TRX" | grep -q "ACTIVE"; then
    echo -e "${GREEN}✓ Transaction retrieved and status is ACTIVE${NC}"
else
    echo -e "${RED}✗ Transaction retrieval failed${NC}"
    exit 1
fi
echo ""

echo -e "${GREEN}=== All Tests Passed! ===${NC}"
```

Jalankan:
```bash
chmod +x test.sh
./test.sh
```

---

## Validasi Hasil Testing

### Checklist Sukses Testing

Jika semua hal di bawah tercapai, project **BERHASIL**:

#### ✓ Konektivitas
- [ ] Health endpoint respond dengan status 200
- [ ] Server berjalan tanpa error
- [ ] Database terhubung

#### ✓ Manajemen Konsumen
- [ ] Bisa register konsumen baru (status 201)
- [ ] Bisa retrieve detail konsumen (status 200)
- [ ] NIK validation berfungsi
- [ ] Gaji minimum validation berfungsi
- [ ] Duplikat NIK ditolak

#### ✓ Manajemen Batas Kredit
- [ ] Bisa assign limit (status 201)
- [ ] Bisa retrieve limits konsumen (status 200)
- [ ] Tenor validation (1,2,3,6 bulan)
- [ ] Limit amount validation

#### ✓ Manajemen Transaksi
- [ ] Bisa create transaksi (status 201)
- [ ] Bisa retrieve transaksi (status 200)
- [ ] Bisa retrieve transaksi per konsumen
- [ ] Bisa update status transaksi
- [ ] Limit checking berfungsi (reject jika limit tidak cukup)
- [ ] Concurrent transaction handling aman

#### ✓ Error Handling
- [ ] Invalid input return 400
- [ ] Not found return 404 atau 400
- [ ] SQL injection protection
- [ ] XSS protection (header validation)

#### ✓ Database
- [ ] Data tersimpan di database
- [ ] Relationship konsumen-limit-transaksi valid
- [ ] Timestamps (created_at, updated_at) terisi

### Metrics Sukses

```
Total Test Cases      : 20+ scenarios
Test Pass Rate        : 100%
Database Records      : >3 consumers, >3 limits, >3 transactions
Response Time         : <200ms per request
Error Handling        : All error cases handled correctly
```

---

## Troubleshooting Testing

### Error: "Connection refused"

```bash
# Pastikan server running
netstat -an | grep 8080
# atau
lsof -i :8080
```

### Error: "Invalid JSON"

```bash
# Validate JSON sebelum send
echo '{"nik": "3173011234567890"}' | jq .
```

### Error: "Cannot parse date"

Format yang benar: `"2026-02-05T10:00:00Z"` (ISO 8601)

---

**Versi**: 1.0.0  
**Updated**: 5 Februari 2026  
**Status**: Siap Digunakan
