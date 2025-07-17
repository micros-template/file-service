package mocks

import (
	"context"
	"io"

	"github.com/stretchr/testify/mock"
)

type MinioStorageMock struct {
	mock.Mock
}

func (m *MinioStorageMock) InitBucket(ctx context.Context, bucketName string) error {
	args := m.Called(ctx, bucketName)
	return args.Error(0)
}

func (m *MinioStorageMock) Set(ctx context.Context, bucketName, objectPath string, reader io.Reader, objectSize int64) error {
	args := m.Called(ctx, bucketName, objectPath, reader, objectSize)
	return args.Error(0)
}

func (m *MinioStorageMock) Get(ctx context.Context, bucketName, objectPath string) (io.ReadCloser, error) {
	args := m.Called(ctx, bucketName, objectPath)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *MinioStorageMock) Delete(ctx context.Context, bucketName, fileName string) error {
	args := m.Called(ctx, bucketName, fileName)
	return args.Error(0)
}

func (m *MinioStorageMock) CreateBucketIfNotExist(ctx context.Context, bucketName string) error {
	args := m.Called(ctx, bucketName)
	return args.Error(0)
}

func (m *MinioStorageMock) SetPolicy(ctx context.Context, bucketName, policy string) error {
	args := m.Called(ctx, bucketName, policy)
	return args.Error(0)
}

func (m *MinioStorageMock) GetPolicy(ctx context.Context, bucketName string) (string, error) {
	args := m.Called(ctx, bucketName)
	return args.String(0), args.Error(1)
}
