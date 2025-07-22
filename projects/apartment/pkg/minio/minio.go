package minio

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

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
