package storage

import (
    "mingda_cloud_service/internal/pkg/config"
    "github.com/minio/minio-go/v7"
    "github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinioClient(cfg *config.MinioConfig) (*minio.Client, error) {
    return minio.New(cfg.Endpoint, &minio.Options{
        Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
        Secure: cfg.UseSSL,
    })
}
