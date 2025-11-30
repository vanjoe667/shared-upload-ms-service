package factory

import (
	"os"
)

type RegistryConfig struct {
    S3 struct {
        Enabled bool
        Region  string
        Bucket  string
    }
}

func LoadRegistryConfigFromEnv() RegistryConfig {
    var cfg RegistryConfig

    if os.Getenv("S3_ENABLED") == "true" {
        cfg.S3.Enabled = true
        cfg.S3.Region = os.Getenv("AWS_REGION")
        cfg.S3.Bucket = os.Getenv("S3_BUCKET")
    }

    return cfg
}