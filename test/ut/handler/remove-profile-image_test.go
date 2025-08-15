package handler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/micros-template/file-service/internal/domain/handler"
	"github.com/micros-template/file-service/test/mocks"
	"github.com/micros-template/proto-file/pkg/fpb"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

type RemoveProfileImageHandlerSuite struct {
	suite.Suite
	userHandler handler.UserGrpcHandler
	userService *mocks.UserServiceMock
}

func (r *RemoveProfileImageHandlerSuite) SetupSuite() {
	mockedUserService := new(mocks.UserServiceMock)
	r.userService = mockedUserService
	grpcServer := grpc.NewServer()
	handler.RegisterUserService(grpcServer, mockedUserService)
	r.userHandler = *handler.NewUserGrpcHandler(mockedUserService)
}

func (r *RemoveProfileImageHandlerSuite) SetupTest() {
	r.userService.ExpectedCalls = nil
	r.userService.Calls = nil
}

func TestRemoveProfileImageHandlerSuite(t *testing.T) {
	suite.Run(t, &RemoveProfileImageHandlerSuite{})
}

func (r *RemoveProfileImageHandlerSuite) TestUserHandler_RemoveProfileImageHandler_Success() {
	ctx := context.Background()
	imageName := "profile.jpg"
	fakeImageName := &fpb.ImageName{Name: imageName}

	r.userService.On("RemoveProfileImage", ctx, imageName).Return(nil)

	status, err := r.userHandler.RemoveProfileImage(ctx, fakeImageName)

	r.NoError(err)
	r.NotNil(status)
	r.True(status.Status)
	r.userService.AssertExpectations(r.T())
}

func (r *RemoveProfileImageHandlerSuite) TestUserHandler_RemoveProfileImageHandler_Error() {
	ctx := context.Background()
	imageName := "profile.jpg"
	fakeImageName := &fpb.ImageName{Name: imageName}
	expectedErr := errors.New("remove failed")

	r.userService.On("RemoveProfileImage", ctx, imageName).Return(expectedErr)

	status, err := r.userHandler.RemoveProfileImage(ctx, fakeImageName)

	r.Error(err)
	r.NotNil(status)
	r.False(status.Status)
	r.userService.AssertExpectations(r.T())
}
