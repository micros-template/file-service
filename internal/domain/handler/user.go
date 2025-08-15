package handler

import (
	"context"

	"github.com/micros-template/file-service/internal/domain/service"
	"github.com/micros-template/proto-file/pkg/fpb"
	"google.golang.org/grpc"
)

type UserGrpcHandler struct {
	userService service.UserService
	fpb.UnimplementedFileServiceServer
}

func NewUserGrpcHandler(userService service.UserService) *UserGrpcHandler {
	return &UserGrpcHandler{
		userService: userService,
	}
}

func RegisterUserService(grpc *grpc.Server, userService service.UserService) {
	grpcHandler := NewUserGrpcHandler(userService)
	fpb.RegisterFileServiceServer(grpc, grpcHandler)
}

func (u *UserGrpcHandler) SaveProfileImage(ctx context.Context, imageByte *fpb.Image) (*fpb.ImageName, error) {
	imageName, err := u.userService.SaveProfileImage(ctx, imageByte.GetImage(), imageByte.GetExt())
	if err != nil {
		return nil, err
	}
	return &fpb.ImageName{
		Name: imageName,
	}, nil
}

func (u *UserGrpcHandler) RemoveProfileImage(ctx context.Context, imageName *fpb.ImageName) (*fpb.Status, error) {
	err := u.userService.RemoveProfileImage(ctx, imageName.GetName())
	if err != nil {
		return &fpb.Status{
			Status: false,
		}, err
	}
	return &fpb.Status{
		Status: true,
	}, nil
}
