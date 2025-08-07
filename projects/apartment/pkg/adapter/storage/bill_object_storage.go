package storage

import (
	"context"

	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/bill/port"
	appctx "github.com/arcaptcha-internship-2025/momoein-apartment/pkg/context"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

type billObjectStorage struct {
	client *minio.Client
}

const BillImageBucket = "bill-image-bucket"

func MustNewBillObjectStorage(c *minio.Client) port.ObjectStorage {
	storage, err := NewBillObjectStorage(c)
	if err != nil {
		panic(err)
	}
	return storage
}

func NewBillObjectStorage(c *minio.Client) (port.ObjectStorage, error) {
	storage := &billObjectStorage{client: c}
	err := storage.CreateBucket(context.Background(), BillImageBucket)
	if err != nil {
		return nil, err
	}
	return storage, nil
}

func (s *billObjectStorage) CreateBucket(ctx context.Context, bucketName string) error {
	found, err := s.client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}
	if found {
		return nil
	}
	opts := minio.MakeBucketOptions{}
	return s.client.MakeBucket(ctx, bucketName, opts)
}

func (s *billObjectStorage) Set(key string, val any) error {
	panic("unimplemented")
}

func (s *billObjectStorage) Get(key string) any {
	panic("unimplemented")
}

func (s *billObjectStorage) FPut(ctx context.Context, key, filename string) error {
	log := appctx.Logger(ctx)
	opts := minio.PutObjectOptions{ContentType: "application/octet-stream"}
	info, err := s.client.FPutObject(ctx, BillImageBucket, key, filename, opts)
	if err != nil {
		log.Error("ObjectStorage.FPut", zap.Error(err))
		return err
	}
	log.Info("ObjectStorage.FPut", zap.Any("upload info", info))
	return nil
}

func (s *billObjectStorage) FGet(ctx context.Context, key, path string) error {
	opts := minio.GetObjectOptions{}
	return s.client.FGetObject(ctx, BillImageBucket, key, path, opts)
}

func (s *billObjectStorage) Del(ctx context.Context, key string) error {
	opts := minio.RemoveObjectOptions{
		GovernanceBypass: true,
	}
	return s.client.RemoveObject(ctx, BillImageBucket, key, opts)
}
