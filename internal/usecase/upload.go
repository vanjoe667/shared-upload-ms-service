package usecase

import (
	"context"
	"time"

	"go.peniremit.dev/peniremit-shared-upload-service/internal/adapter/storage/factory"
	"go.peniremit.dev/peniremit-shared-upload-service/internal/port"
)

type UploadUsecase struct {
	storage port.StorageProvider
}

func (u *UploadUsecase) PresignDownload(ctx context.Context, key string) (port.PresignResponse, error) {
	return u.storage.PresignDownload(ctx, key, 15*time.Minute)
}

func (u *UploadUsecase) DeleteObject(ctx context.Context, key string) error {
	return u.storage.DeleteObject(ctx, key)
}

func FileUploadUsecase(factory *factory.ProviderFactory, slug port.ProviderSlug) (*UploadUsecase, error) {
	provider, err := factory.GetProvider(slug)
	if err != nil {
		return nil, err
	}
	return &UploadUsecase{storage: provider}, nil
}

func (u *UploadUsecase) PresignSingle(ctx context.Context, key, contentType string) (port.PresignResponse, error) {
	in := port.UploadInput{
		Key:         key,
		ContentType: contentType,
	}
	return u.storage.PresignUpload(ctx, in, 15*time.Minute)
}

func (u *UploadUsecase) PresignMany(ctx context.Context, inputs []port.UploadInput) ([]port.PresignResponse, error) {
	return u.storage.PresignMultiUpload(ctx, inputs, 15*time.Minute)
}

func (u *UploadUsecase) StartMultipart(ctx context.Context, key, contentType string, maxSize int64) (port.MultipartSession, error) {
	in := port.UploadInput{
		Key:         key,
		ContentType: contentType,
		MaxSize:     maxSize,
	}
	return u.storage.InitiateMultipart(ctx, in)
}

func (u *UploadUsecase) GetMultipartPartURL(ctx context.Context, key, uploadID string, partNumber int32) (port.PresignResponse, error) {
	return u.storage.PresignMultipartPart(ctx, key, uploadID, partNumber, 30*time.Minute)
}

func (u *UploadUsecase) CompleteMultipart(ctx context.Context, key, uploadID string, parts []port.UploadPart) error {
	return u.storage.CompleteMultipart(ctx, key, uploadID, parts)
}
