package handler_test

import (
	"context"
	"errors"
	"testing"

	"10.1.20.130/dropping/file-service/internal/domain/handler"
	"10.1.20.130/dropping/file-service/test/mocks"
	"github.com/dropboks/proto-file/pkg/fpb"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

type SaveProfileImageHandlerSuite struct {
	suite.Suite
	userHandler handler.UserGrpcHandler
	userService *mocks.UserServiceMock
}

func (s *SaveProfileImageHandlerSuite) SetupSuite() {
	mockedUserService := new(mocks.UserServiceMock)
	s.userService = mockedUserService
	grpcServer := grpc.NewServer()
	handler.RegisterUserService(grpcServer, mockedUserService)
	s.userHandler = *handler.NewUserGrpcHandler(mockedUserService)
}

func (s *SaveProfileImageHandlerSuite) SetupTest() {
	s.userService.ExpectedCalls = nil
	s.userService.Calls = nil
}

func TestSaveProfileImageHandlerSuite(t *testing.T) {
	suite.Run(t, &SaveProfileImageHandlerSuite{})
}

func (s *SaveProfileImageHandlerSuite) TestUserHandler_SaveProfileImageHandler_Success() {
	ctx := context.Background()
	imageBytes := []byte{0x01, 0x02, 0x03}
	ext := ".jpg"
	expectedImageName := "profile_123.jpg"

	s.userService.On("SaveProfileImage", ctx, imageBytes, ext).Return(expectedImageName, nil)

	image := &fpb.Image{
		Image: imageBytes,
		Ext:   ext,
	}

	result, err := s.userHandler.SaveProfileImage(ctx, image)
	s.NoError(err)
	s.NotNil(result)
	s.Equal(expectedImageName, result.Name)
	s.userService.AssertExpectations(s.T())
}

func (s *SaveProfileImageHandlerSuite) TestUserHandler_SaveProfileImageHandler_Error() {
	ctx := context.Background()
	imageBytes := []byte{0x01, 0x02, 0x03}
	ext := ".png"
	expectedErr := errors.New("failed to save image")

	s.userService.On("SaveProfileImage", ctx, imageBytes, ext).Return("", expectedErr)

	image := &fpb.Image{
		Image: imageBytes,
		Ext:   ext,
	}

	result, err := s.userHandler.SaveProfileImage(ctx, image)
	s.Error(err)
	s.Nil(result)
	s.Equal(expectedErr, err)
	s.userService.AssertExpectations(s.T())
}
