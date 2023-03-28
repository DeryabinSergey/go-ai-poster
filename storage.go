package aiposter

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
)

type Storage struct {
	bucket *storage.BucketHandle
}

func NewStorage(ctx context.Context, bucket string) (storageClient *Storage, err error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return storageClient, fmt.Errorf("storage.NewClient: %w", err)
	}

	storageClient.bucket = client.Bucket(bucket)

	return
}

func (s *Storage) GetReader(ctx context.Context, object string) (reader *storage.Reader, err error) {
	reader, err = s.bucket.Object(object).NewReader(ctx)
	if err != nil {
		return reader, fmt.Errorf("read error: %w", err)
	}

	return
}

func (s *Storage) GetWriter(ctx context.Context, object string) *storage.Writer {
	return s.bucket.Object(object).NewWriter(ctx)
}
