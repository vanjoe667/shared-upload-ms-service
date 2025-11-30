package bootstrap

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.peniremit.dev/peniremit-shared-upload-service/internal/adapter/storage/factory"
	"go.peniremit.dev/peniremit-shared-upload-service/internal/port"
	"go.peniremit.dev/peniremit-shared-upload-service/internal/usecase"
)

func InitUploadService(ctx context.Context) (*usecase.UploadUsecase, error) {
	  err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found, relying on environment variables")
    }

    active := port.ProviderSlug(os.Getenv("UPLOAD_PROVIDER"))
    if active == "" {
        return nil, fmt.Errorf("UPLOAD_PROVIDER env var is required")
    }

    cfg := factory.LoadRegistryConfigFromEnv()

    providerMap, err := factory.LoadProviders(ctx, cfg)
    if err != nil {
        return nil, err
    }

    providerFactory := factory.FileUploadProviderFactory(providerMap)

    return usecase.FileUploadUsecase(providerFactory, active)
}