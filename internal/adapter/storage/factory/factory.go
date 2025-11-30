package factory

import (
	"fmt"

	"go.peniremit.dev/peniremit-shared-upload-service/internal/port"
)

type ProviderFactory struct {
    providers ProviderMap
}

func FileUploadProviderFactory(providers ProviderMap) *ProviderFactory {
    return &ProviderFactory{
        providers: providers,
    }
}

func (f *ProviderFactory) GetProvider(slug port.ProviderSlug) (port.StorageProvider, error) {
    provider, ok := f.providers[slug]
    if !ok {
        return nil, fmt.Errorf("storage provider '%s' not registered", slug)
    }
    return provider, nil
}