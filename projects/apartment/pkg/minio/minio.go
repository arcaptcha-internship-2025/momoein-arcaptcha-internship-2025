package minio

import (
	"context"
	"errors"
	"time"

	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/fp"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var ErrPing = errors.New("minio failed to ping")

func MustNewClient(endpoint, accessKey, secretKey string) *minio.Client {
	c, err := NewClient(endpoint, accessKey, secretKey)
	if err != nil {
		panic(err)
	}
	return c
}

func NewClient(endpoint, accessKey, secretKey string) (*minio.Client, error) {
	return minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
}

// Ping checks whether the MinIO server is reachable and credentials are valid.
// Returns nil if successful; otherwise, returns an error.
func Ping(endpoint, accessKey, secretKey string, useSSL bool, timeout time.Duration) error {
	// Initialize MinIO client
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return fp.WrapErrors(ErrPing, err)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Try listing buckets as a lightweight reachability/auth check
	if _, err = client.ListBuckets(ctx); err != nil {
		return fp.WrapErrors(ErrPing, err)
	}
	return nil
}
