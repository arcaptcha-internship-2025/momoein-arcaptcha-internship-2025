package storage

import (
	"context"
	"testing"

	appctx "github.com/arcaptcha-internship-2025/momoein-apartment/pkg/context"
	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/logger"
	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/minio"
	"github.com/stretchr/testify/assert"
)

var (
	log = logger.NewConsoleZapLogger(logger.ModeDevelopment)
	ctx = appctx.New(context.Background(), appctx.WithLogger(log))
)

func TestBillObjectStorage_FPut(t *testing.T) {
	client := minio.MustNewClient("localhost:9000", "minioadmin", "minioadmin")
	obj := MustNewBillObjectStorage(client)

	err := obj.FPut(ctx, "test-2", "../../../test/data/my-testfile")
	assert.Nil(t, err)
}

func TestBillObjectStorage_FGet(t *testing.T) {
	client := minio.MustNewClient("localhost:9000", "minioadmin", "minioadmin")

	obj := MustNewBillObjectStorage(client)

	err := obj.FGet(ctx, "test-1", "../../../test/data/my-testfile1")
	assert.Nil(t, err)
}
