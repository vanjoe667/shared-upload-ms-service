# Shared Upload Service

A lightweight, provider‑agnostic microservice for generating presigned upload/download URLs and handling multipart file uploads. It currently supports **AWS S3**, with the architecture designed so additional providers (Cloudflare R2, Cloudinary, etc.) can be plugged in without changing business logic.

---

## 🚀 Features

- **Provider‑agnostic architecture** using clean architecture patterns
- **S3 presigned upload URLs** (single + batch)
- **S3 presigned multipart uploads** for large files
- **Presigned downloads**
- **Delete uploaded objects**
- **Strict content‑type validation**
- **Configurable via environment variables**
- **Extensible provider factory** (plug in any provider later)

---

## 📦 Architecture Overview

```
├── cmd/server           # Entry point (minimal bootstrapping)
├── internal
│   ├── adapter
|   |   |── http        # HTTP handlers & routes
│   │   └── storage
│   │       ├── s3      # S3 implementation
│   │       └── factory # Provider registry + factory
│   ├── port            # Interfaces (StorageProvider, DTOs)
│   ├── usecase         # Business logic
│   └── transport
│       
└── pkg                 # Shared helpers
```

The **main.go** never knows which provider is used — it only reads:
```
UPLOAD_PROVIDER=s3
```
The factory handles the rest.

---

## ⚙️ Environment Variables

Create a `.env` file:

```env
# Active provider (only this matters for activation)
UPLOAD_PROVIDER=s3

# Server
SERVER_PORT=8080

# AWS S3 Config
AWS_REGION=us-east-1
S3_BUCKET=your-bucket-name
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
```

AWS credentials work automatically via `LoadDefaultConfig` inside the AWS SDK.

---

## 🏁 Running the Service

### Install dependencies
```bash
go mod tidy
```

### Start in development mode
```bash
go run ./cmd/server
```

### Auto‑reload on file changes
Install `air`:
```bash
brew install air
```
Run:
```bash
air
```

---

## 🔌 Endpoints

### **1. Generate Presigned Upload URL**
`POST /upload/single`

**Request:**
```json
{
  "key": "uploads/profile.png",
  "content_type": "image/png"
}
```

**Response:**
```json
{
  "url": "https://s3-presigned-url...","url_expires_in": 900,
  "key": "uploads/profile.png"
}
```

---

### **2. Presigned Download URL**
`GET /upload/download?key=uploads/profile.png`

**Response:**
```json
{
  "url": "https://presigned-download...",
  "expires": 900
}
```

---

### **3. Start Multipart Upload**
`POST /upload/multipart/start`

```json
{
  "key": "videos/bigfile.mp4",
  "content_type": "video/mp4",
  "max_size": 104857600
}
```

**Response:**
```json
{
  "upload_id": "XYZ123",
  "key": "videos/bigfile.mp4",
  "expires": "2025-01-01T12:00:00Z"
}
```

---

### **4. Generate Upload Part URL**
`POST /upload/multipart/{uploadID}/part`

```json
{
  "key": "videos/bigfile.mp4",
  "part_number": 1
}
```

---

### **5. Complete Multipart Upload**
`POST /upload/multipart/complete`

```json
{
  "key": "videos/bigfile.mp4",
  "upload_id": "XYZ123",
  "parts": [
    { "etag": "abc123", "part_number": 1 }
  ]
}
```

---

## 🔧 Provider Extensibility

Add a new provider by implementing:
```go
type StorageProvider interface {
    PresignUpload(...)
    PresignDownload(...)
    DeleteObject(...)
    PresignMultiUpload(...)
    InitiateMultipart(...)
    PresignMultipartPart(...)
    CompleteMultipart(...)
    AbortMultipart(...)
}
```

Then register it in the provider registry. No changes to usecases or handlers.

---

## 🧪 Testing

### Run unit tests
```bash
go test ./...
```

Mock the provider by implementing `StorageProvider` and injecting it in tests.

---

## 🛠 Contribution Guide

1. Follow Go clean architecture principles
2. Keep main.go simple
3. Never import concrete providers inside usecases or handlers
4. Log only meaningful information
5. Format code before pushing:
```bash
go fmt ./...
```

---

## 🚀 Deployment Notes

### Build binary
```bash
go build -o upload-service ./cmd/server
```

Deploy to:
- AWS ECS / Fargate
- Kubernetes
- Render / Railway / Fly.io
- Lambda (via API Gateway)

The service is stateless — safe for horizontal scaling.

---

## 📄 License
MIT

---

## 👤 Author
Developed by **Joel Ajide** — scalable infrastructure enthusiast and Senior backend engineer.

