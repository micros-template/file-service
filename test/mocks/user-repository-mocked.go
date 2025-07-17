package mocks

import (
	"context"
	"io"

	m "github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	m.Mock
}

func (m *UserRepositoryMock) SaveProfileImage(ctx context.Context, bucketName, objectPath string, reader io.Reader, objectSize int64) error {
	args := m.Called(ctx, bucketName, objectPath, reader, objectSize)
	return args.Error(0)
}

func (m *UserRepositoryMock) RemoveProfileImage(ctx context.Context, bucketName, objectPath string) error {
	args := m.Called(ctx, bucketName, objectPath)
	return args.Error(0)
}
