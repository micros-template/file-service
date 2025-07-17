package repository_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"10.1.20.130/dropping/file-service/internal/domain/repository"
	"10.1.20.130/dropping/file-service/test/mocks"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type SaveProfileImageRepositorySuite struct {
	suite.Suite
	userRepository repository.UserRepository
	minioStorage   *mocks.MinioStorageMock
}

func (s *SaveProfileImageRepositorySuite) SetupSuite() {

	mockMinioStorage := new(mocks.MinioStorageMock)
	logger := zerolog.Nop()
	s.minioStorage = mockMinioStorage
	s.userRepository = repository.NewUserRepository(s.minioStorage, logger)
}

func (s *SaveProfileImageRepositorySuite) SetupTest() {
	s.minioStorage.ExpectedCalls = nil
	s.minioStorage.Calls = nil
}

func TestSaveProfileImageSuite(t *testing.T) {
	suite.Run(t, &SaveProfileImageRepositorySuite{})
}
func (s *SaveProfileImageRepositorySuite) TestUserRepository_SaveProfileImage_Success() {
	ctx := context.Background()
	bucketName := "test-bucket"
	objectPath := "profile-images/user123.jpg"
	imageContent := []byte("fake image data")
	objectSize := int64(len(imageContent))
	reader := bytes.NewReader(imageContent)

	s.minioStorage.On("Set", ctx, bucketName, objectPath, mock.Anything, objectSize).Return(nil)

	err := s.userRepository.SaveProfileImage(ctx, bucketName, objectPath, reader, objectSize)
	s.NoError(err)
	s.minioStorage.AssertCalled(s.T(), "Set", ctx, bucketName, objectPath, mock.Anything, objectSize)
}

func (s *SaveProfileImageRepositorySuite) TestUserRepository_SaveProfileImage_Error() {
	ctx := context.Background()
	bucketName := "test-bucket"
	objectPath := "profile-images/user123.jpg"
	imageContent := []byte("fake image data")
	objectSize := int64(len(imageContent))
	reader := bytes.NewReader(imageContent)

	s.minioStorage.On("Set", ctx, bucketName, objectPath, mock.Anything, objectSize).Return(errors.New("failed to insert"))

	err := s.userRepository.SaveProfileImage(ctx, bucketName, objectPath, reader, objectSize)
	s.Error(err)
	s.minioStorage.AssertCalled(s.T(), "Set", ctx, bucketName, objectPath, mock.Anything, objectSize)
}
