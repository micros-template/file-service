package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/micros-template/file-service/internal/domain/service"
	"github.com/micros-template/file-service/pkg/constant"
	"github.com/micros-template/file-service/test/mocks"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type SaveProfileImageServiceSuite struct {
	suite.Suite
	userService    service.UserService
	userRepository *mocks.UserRepositoryMock
	logEmitter     *mocks.LoggerInfraMock
}

func (r *SaveProfileImageServiceSuite) SetupSuite() {

	mockUserRepo := new(mocks.UserRepositoryMock)
	mockLogEmitter := new(mocks.LoggerInfraMock)
	logger := zerolog.Nop()
	r.userRepository = mockUserRepo
	r.logEmitter = mockLogEmitter
	r.userService = service.NewUserService(r.userRepository, mockLogEmitter, logger)
}

func (r *SaveProfileImageServiceSuite) SetupTest() {
	r.userRepository.ExpectedCalls = nil
	r.logEmitter.ExpectedCalls = nil
	r.userRepository.Calls = nil
	r.logEmitter.Calls = nil
}

func TestSaveProfileImageSuite(t *testing.T) {
	suite.Run(t, &SaveProfileImageServiceSuite{})
}
func (r *SaveProfileImageServiceSuite) TestUserService_SaveProfileImage_Success_Default() {
	ctx := context.Background()
	imageBytes := []byte{1, 2, 3}
	imageExt := "png"

	r.userRepository.On("SaveProfileImage", ctx, constant.APP_BUCKET, mock.AnythingOfType("string"), mock.Anything, int64(len(imageBytes))).Return(nil)
	r.logEmitter.On("EmitLog", mock.Anything, mock.Anything).Return(nil).Once()

	path, err := r.userService.SaveProfileImage(ctx, imageBytes, imageExt)
	r.NoError(err)
	r.NotEmpty(path)
	r.userRepository.AssertExpectations(r.T())

	time.Sleep(time.Second)
	r.logEmitter.AssertExpectations(r.T())
}

func (r *SaveProfileImageServiceSuite) TestUserService_SaveProfileImage_Error() {
	ctx := context.Background()
	imageBytes := []byte{1, 2, 3}
	imageExt := "jpg"

	r.userRepository.On("SaveProfileImage", ctx, constant.APP_BUCKET, mock.AnythingOfType("string"), mock.Anything, int64(3)).Return(errors.New("save error"))
	r.logEmitter.On("EmitLog", mock.Anything, mock.Anything).Return(nil).Once()

	path, err := r.userService.SaveProfileImage(ctx, imageBytes, imageExt)
	r.Error(err)
	r.Empty(path)
	r.userRepository.AssertExpectations(r.T())

	time.Sleep(time.Second)
	r.logEmitter.AssertExpectations(r.T())
}
