# MinIO Storage Microservice (Golang)

API Gateway ringan berbasis Go untuk mengelola unggah, lihat, daftar, dan hapus file ke Object Storage MinIO/S3 dengan autentikasi `SERVER_SECRET_KEY`.

---

## 🛠️ Persyaratan Sistem

- **Golang** 1.22+ (untuk jalankan lokal)
- **Docker** & **Docker Compose** (untuk deployment)
- **MinIO / S3 Object Storage Server**

---

## 🚀 Panduan Memulai & Deployment

### 1. Konfigurasi Environment

Buat file `.env` berdasarkan template berikut:

```env
PORT=8080
SERVER_SECRET_KEY=rahasia_server_anda_123

# Konfigurasi MinIO (Isi domain saja tanpa https://)
NEO_ENDPOINT=s3.domainanda.com
NEO_ACCESS_KEY=your_access_key
NEO_SECRET_KEY=your_secret_key
NEO_USE_SSL=true
NEO_BUCKET=my-bucket
```

### 2. Jalankan di Lokal (Tanpa Docker)


#### Download dependencies
```bash
go mod tidy
```

#### Jalankan server
```bash
go run cmd/server/main.go
```

3. Deployment Menggunakan Docker
Build Docker Image:
```bash
docker build -t minio-app:latest .
```
Run Docker Container dengan File `.env`:

```bash
docker run -d \
  --name minio-storage-app \
  -p 8080:8080 \
  --env-file .env \
  minio-app:latest
```

## 🧪 Panduan Uji Coba API (Testing)
Anda dapat menggunakan cURL, Postman, atau Insomnia untuk menguji endpoint.

### 1. Upload File (POST /upload)
```
Content-Type: multipart/form-data
```
```bash
curl -X POST http://localhost:8080/upload \
  -F "secret=rahasia_server_anda_123" \
  -F "client=clientA" \
  -F "folder=invoices" \
  -F "file=@/path/ke/file_gambar.jpg"
```

#### Response (202 Accepted):

```JSON

{
  "message": "Sukses Upload (Queued)",
  "success": true,
  "url": "[https://s3.domainanda.com/my-bucket/clientA/uploads/invoices/550e8400-e29b-41d4-a716-446655440000.jpg](https://s3.domainanda.com/my-bucket/clientA/uploads/invoices/550e8400-e29b-41d4-a716-446655440000.jpg)"
}
```
### 2. List File (GET /list)

Query Params: secret, client, folder (optional)

```
curl -X GET "http://localhost:8080/list?secret=rahasia_server_anda_123&client=clientA&folder=invoices"

```

#### Response (200 OK):

```JSON
{
  "files": [
    "[https://s3.domainanda.com/my-bucket/clientA/uploads/invoices/550e8400-e29b-41d4-a716-446655440000.jpg](https://s3.domainanda.com/my-bucket/clientA/uploads/invoices/550e8400-e29b-41d4-a716-446655440000.jpg)"
  ],
  "success": true
}
```

### 3. View / Stat Metadata File (POST /view)
```
Content-Type: application/x-www-form-urlencoded
```

```Bash
curl -X POST http://localhost:8080/view \
  -d "secret=rahasia_server_anda_123" \
  -d "url=[https://s3.domainanda.com/my-bucket/clientA/uploads/invoices/550e8400-e29b-41d4-a716-446655440000.jpg](https://s3.domainanda.com/my-bucket/clientA/uploads/invoices/550e8400-e29b-41d4-a716-446655440000.jpg)"
```

#### Response (200 OK):

```JSON
{
  "key": "clientA/uploads/invoices/550e8400-e29b-41d4-a716-446655440000.jpg",
  "last_modified": "2026-07-21 13:40:00",
  "size": 1048576,
  "type": "image/jpeg",
  "url": "[https://s3.domainanda.com/my-bucket/clientA/uploads/invoices/550e8400-e29b-41d4-a716-446655440000.jpg](https://s3.domainanda.com/my-bucket/clientA/uploads/invoices/550e8400-e29b-41d4-a716-446655440000.jpg)"
}
```

### 4. Hapus File (POST /delete)

```
Content-Type: application/x-www-form-urlencoded
```

```Bash
curl -X POST http://localhost:8080/delete \
  -d "secret=rahasia_server_anda_123" \
  -d "file=[https://s3.domainanda.com/my-bucket/clientA/uploads/invoices/550e8400-e29b-41d4-a716-446655440000.jpg](https://s3.domainanda.com/my-bucket/clientA/uploads/invoices/550e8400-e29b-41d4-a716-446655440000.jpg)"
```
#### Response (200 OK):

```JSON


{
  "message": "File deleted successfully",
  "success": true
}
```