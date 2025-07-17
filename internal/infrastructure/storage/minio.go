package storage

import (
	"context"
	"fmt"
	"io"

	"10.1.20.130/dropping/file-service/pkg/constant"
	"10.1.20.130/dropping/file-service/pkg/utils"
	"github.com/minio/minio-go/v7"
)

type (
	MinioStorage interface {
		InitBucket(context context.Context, bucketName string) error
		Set(ctx context.Context, bucketName, objectPath string, reader io.Reader, objectSize int64) error
		Get(ctx context.Context, bucketName, objectPath string) (io.ReadCloser, error)
		Delete(ctx context.Context, bucketName, fileName string) error
		CreateBucketIfNotExist(ctx context.Context, bucketName string) error
		SetPolicy(ctx context.Context, bucketName, policy string) error
		GetPolicy(ctx context.Context, bucketName string) (string, error)
	}
	minioStorage struct {
		client *minio.Client
	}
)

func NewMinioStorage(client *minio.Client) MinioStorage {
	return &minioStorage{
		client: client,
	}
}
func (m *minioStorage) InitBucket(context context.Context, bucketName string) error {
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

func (m *minioStorage) Set(ctx context.Context, bucketName, objectPath string, reader io.Reader, objectSize int64) error {
	_, err := m.client.PutObject(ctx, bucketName, objectPath, reader, objectSize, minio.PutObjectOptions{})
	return err
}

func (m *minioStorage) Get(ctx context.Context, bucketName, objectPath string) (io.ReadCloser, error) {
	object, err := m.client.GetObject(ctx, bucketName, objectPath, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return object, nil
}

func (m *minioStorage) Delete(ctx context.Context, bucketName, fileName string) error {
	return m.client.RemoveObject(ctx, bucketName, fileName, minio.RemoveObjectOptions{})
}

func (m *minioStorage) CreateBucketIfNotExist(ctx context.Context, bucketName string) error {
	exists, err := m.client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}
	if !exists {
		return m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	}
	return nil
}
func (m *minioStorage) SetPolicy(ctx context.Context, bucketName, policy string) error {
	return m.client.SetBucketPolicy(ctx, bucketName, policy)
}

func (m *minioStorage) GetPolicy(ctx context.Context, bucketName string) (string, error) {
	return m.client.GetBucketPolicy(ctx, bucketName)
}
