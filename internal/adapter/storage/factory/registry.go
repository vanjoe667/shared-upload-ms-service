package factory

import (
	"context"

	s3adapter "go.peniremit.dev/peniremit-shared-upload-service/internal/adapter/storage/s3"
	"go.peniremit.dev/peniremit-shared-upload-service/internal/port"
)

func LoadProviders(ctx context.Context, cfg RegistryConfig) (ProviderMap, error) {
    providers := ProviderMap{}

    if cfg.S3.Enabled {
        client, err := s3adapter.S3Client(ctx, s3adapter.Config{
            Region: cfg.S3.Region,
        })
        if err != nil {
            return nil, err
        }
        providers[port.ProviderS3] = s3adapter.New(client, cfg.S3.Bucket)
    }

    return providers, nil
}