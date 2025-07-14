package storage

import (
	"context"
	"fmt"
	"io"

	"10.1.20.130/dropping/file-service/pkg/constant"
	"10.1.20.130/dropping/file-service/pkg/utils"
	"github.com/minio/minio-go/v7"
)

type MinioStorage struct {
	Client *minio.Client
}

func NewMinioStorage(client *minio.Client) *MinioStorage {
	return &MinioStorage{
		Client: client,
	}
}
func (m *MinioStorage) InitBucket(context context.Context, bucketName string) error {
	err := m.CreateBucketIfNotExist(context, bucketName)
	if err != nil {
		return err
	}
	policyStr, err := m.GetPolicy(context, bucketName)
	if err != nil {
		return err
	}
	if policyStr != "" {
		public, err := utils.IsBucketPublic(policyStr)
		if err != nil {
			return err
		}
		if !public {
			policy := fmt.Sprintf(constant.PUBLIC_PERMISSION, bucketName)
			err = m.SetPolicy(context, bucketName, policy)
			if err != nil {
				return err
			}
		}
	} else {
		policy := fmt.Sprintf(constant.PUBLIC_PERMISSION, bucketName)
		err = m.SetPolicy(context, bucketName, policy)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MinioStorage) Set(ctx context.Context, bucketName, objectPath string, reader io.Reader, objectSize int64) error {
	_, err := m.Client.PutObject(ctx, bucketName, objectPath, reader, objectSize, minio.PutObjectOptions{})
	return err
}

func (m *MinioStorage) Get(ctx context.Context, bucketName, objectPath string) (io.ReadCloser, error) {
	object, err := m.Client.GetObject(ctx, bucketName, objectPath, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return object, nil
}

func (m *MinioStorage) Delete(ctx context.Context, bucketName, fileName string) error {
	return m.Client.RemoveObject(ctx, bucketName, fileName, minio.RemoveObjectOptions{})
}

func (m *MinioStorage) CreateBucketIfNotExist(ctx context.Context, bucketName string) error {
	exists, err := m.Client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}
	if !exists {
		return m.Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	}
	return nil
}
func (m *MinioStorage) SetPolicy(ctx context.Context, bucketName, policy string) error {
	return m.Client.SetBucketPolicy(ctx, bucketName, policy)
}

func (m *MinioStorage) GetPolicy(ctx context.Context, bucketName string) (string, error) {
	return m.Client.GetBucketPolicy(ctx, bucketName)
}
