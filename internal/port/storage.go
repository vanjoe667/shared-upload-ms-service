package port

import (
	"context"
	"time"
)

// UploadInput describes a single file the client intends to upload.
type UploadInput struct {
	Key         string            // object key in storage (e.g., uploads/users/123/file.pdf)
	ContentType string            // expected content type, e.g., "application/pdf", "image/png"
	Metadata    map[string]string // optional custom metadata stored with object
	MaxSize     int64             // optional size limit in bytes (0 = no limit)
}

// PresignResponse is returned for a presigned operation.
type PresignResponse struct {
	URL     string
	Fields  map[string]string // for POST (form) style uploads, empty for PUT presigns
	Expires time.Duration
	Key     string 
}

// UploadPart represents a multipart upload part when completing the upload.
type UploadPart struct {
	PartNumber int32
	ETag       string
}

// MultipartSession represents the state for a multipart upload
type MultipartSession struct {
	UploadID string
	Key      string
	Parts    []UploadPart
	Expires  time.Time
}

type StorageProvider interface {
	// Single object presign (PUT style). Good for small files.
	PresignUpload(ctx context.Context, in UploadInput, expires time.Duration) (PresignResponse, error)

	// Generate presigned URL for download/getting an object.
	PresignDownload(ctx context.Context, key string, expires time.Duration) (PresignResponse, error)

	// Delete an object by key.
	DeleteObject(ctx context.Context, key string) error

	// Multi-file presign: accept many UploadInput items and return corresponding presigns.
	// Useful to reduce round trips: generate N presigns in one call.
	PresignMultiUpload(ctx context.Context, inputs []UploadInput, expires time.Duration) ([]PresignResponse, error)

	// --- Multipart (large file) flow ---
	// InitiateMultipart starts a multipart upload and returns an uploadID.
	InitiateMultipart(ctx context.Context, in UploadInput) (MultipartSession, error)

	// PresignMultipartPart generates a presigned URL for uploading a single part
	// for a multipart session (uploadID + partNumber).
	PresignMultipartPart(ctx context.Context, key, uploadID string, partNumber int32, expires time.Duration) (PresignResponse, error)

	// CompleteMultipart completes the multipart upload with the collected parts (ETags).
	CompleteMultipart(ctx context.Context, key, uploadID string, parts []UploadPart) error

	// AbortMultipart aborts an in-progress multipart upload.
	AbortMultipart(ctx context.Context, key, uploadID string) error
}