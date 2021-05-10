package main

import (
	"context"
	"net/url"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	expires = 10 * time.Minute
)

func GetMinioClient() (minioClient *minio.Client, err error) {
	endpoint := os.Getenv("MINIO_HOST")
	accessKeyID := os.Getenv("MINIO_ROOT_USER")
	secretAccessKey := os.Getenv("MINIO_ROOT_PASSWORD")
	useSSL := false

	// Initialize minio client object.
	minioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	return
}

func GetUploadLink(minioClient *minio.Client, bucketName, objectName string) (*url.URL, error) {
	return minioClient.PresignedPutObject(context.Background(), bucketName, objectName, expires)
}
