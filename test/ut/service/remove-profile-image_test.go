package service_test

import (
	"context"
	"errors"
	"testing"

	"10.1.20.130/dropping/file-service/internal/domain/service"
	"10.1.20.130/dropping/file-service/pkg/constant"
	"10.1.20.130/dropping/file-service/test/mocks"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RemoveProfileImageServiceSuite struct {
	suite.Suite
	userService    service.UserService
	userRepository *mocks.UserRepositoryMock
}

func (r *RemoveProfileImageServiceSuite) SetupSuite() {

	mockUserRepo := new(mocks.UserRepositoryMock)
	mockLogEmitter := new(mocks.LoggerInfraMock)

	logger := zerolog.Nop()
	r.userRepository = mockUserRepo
	r.userService = service.NewUserService(r.userRepository, mockLogEmitter, logger)
}

func (r *RemoveProfileImageServiceSuite) SetupTest() {
	r.userRepository.ExpectedCalls = nil
	r.userRepository.Calls = nil
}

func TestRemoveProfileImageSuite(t *testing.T) {
	suite.Run(t, &RemoveProfileImageServiceSuite{})
}
func (r *RemoveProfileImageServiceSuite) TestUserService_RemoveProfileImage_Success_Default() {
	ctx := context.Background()

	imageName := "image name"
	r.userRepository.On("RemoveProfileImage", ctx, constant.APP_BUCKET, mock.AnythingOfType("string")).Return(nil)

	err := r.userService.RemoveProfileImage(ctx, imageName)
	r.NoError(err)
	r.userRepository.AssertExpectations(r.T())
}

func (r *RemoveProfileImageServiceSuite) TestUserService_RemoveProfileImage_Error() {
	ctx := context.Background()
	imageName := "image name"

	r.userRepository.On("RemoveProfileImage", ctx, constant.APP_BUCKET, mock.AnythingOfType("string")).Return(errors.New("save error"))

	err := r.userService.RemoveProfileImage(ctx, imageName)
	r.Error(err)
	r.userRepository.AssertExpectations(r.T())
}
