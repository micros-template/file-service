package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"10.1.20.130/dropping/file-service/internal/domain/repository"
	"10.1.20.130/dropping/file-service/test/mocks"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RemoveProfileImageRepositorySuite struct {
	suite.Suite
	userRepository repository.UserRepository
	minioStorage   *mocks.MinioStorageMock
	logEmitter     *mocks.LoggerInfraMock
}

func (r *RemoveProfileImageRepositorySuite) SetupSuite() {

	mockMinioStorage := new(mocks.MinioStorageMock)
	mockLogEmitter := new(mocks.LoggerInfraMock)

	logger := zerolog.Nop()
	r.minioStorage = mockMinioStorage
	r.logEmitter = mockLogEmitter
	r.userRepository = repository.NewUserRepository(r.minioStorage, mockLogEmitter, logger)
}

func (r *RemoveProfileImageRepositorySuite) SetupTest() {
	r.minioStorage.ExpectedCalls = nil
	r.logEmitter.ExpectedCalls = nil
	r.minioStorage.Calls = nil
	r.logEmitter.Calls = nil
}

func TestRemoveProfileImageSuite(t *testing.T) {
	suite.Run(t, &RemoveProfileImageRepositorySuite{})
}
func (r *RemoveProfileImageRepositorySuite) TestUserRepository_RemoveProfileImage_Success() {
	ctx := context.Background()
	bucketName := "test-bucket"
	objectPath := "profile-images/user123.jpg"

	r.minioStorage.On("Delete", ctx, bucketName, objectPath).Return(nil)

	err := r.userRepository.RemoveProfileImage(ctx, bucketName, objectPath)
	r.NoError(err)
	r.minioStorage.AssertCalled(r.T(), "Delete", ctx, bucketName, objectPath)
}

func (r *RemoveProfileImageRepositorySuite) TestUserRepository_RemoveProfileImage_Error() {
	ctx := context.Background()
	bucketName := "test-bucket"
	objectPath := "profile-images/user123.jpg"

	r.minioStorage.On("Delete", ctx, bucketName, objectPath).Return(errors.New("error to remove profile image"))
	r.logEmitter.On("EmitLog", mock.Anything, mock.Anything).Return(nil).Once()

	err := r.userRepository.RemoveProfileImage(ctx, bucketName, objectPath)
	r.Error(err)
	r.minioStorage.AssertCalled(r.T(), "Delete", ctx, bucketName, objectPath)

	time.Sleep(time.Second)
	r.logEmitter.AssertExpectations(r.T())
}
