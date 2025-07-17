package mocks

import (
	"context"

	m "github.com/stretchr/testify/mock"
)

type UserServiceMock struct {
	m.Mock
}

func (m *UserServiceMock) SaveProfileImage(ctx context.Context, imageByte []byte, imageExt string) (string, error) {
	args := m.Called(ctx, imageByte, imageExt)
	return args.String(0), args.Error(1)
}

func (m *UserServiceMock) RemoveProfileImage(ctx context.Context, imageName string) error {
	args := m.Called(ctx, imageName)
	return args.Error(0)
}
