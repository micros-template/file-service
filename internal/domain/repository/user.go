package repository

import (
	"context"
	"io"

	"10.1.20.130/dropping/file-service/internal/domain/dto"
	"10.1.20.130/dropping/file-service/internal/infrastructure/storage"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	UserRepository interface {
		SaveProfileImage(ctx context.Context, bucketName, objectPath string, reader io.Reader, objectSize int64) error
		RemoveProfileImage(ctx context.Context, bucketName, objectPath string) error
	}
	userRepository struct {
		miniio *storage.MinioStorage
		logger zerolog.Logger
	}
)

func NewUserRepository(miniio *storage.MinioStorage, logger zerolog.Logger) UserRepository {
	return &userRepository{
		miniio: miniio,
		logger: logger,
	}
}

func (u *userRepository) SaveProfileImage(ctx context.Context, bucketName string, objectPath string, reader io.Reader, objectSize int64) error {
	err := u.miniio.Set(ctx, bucketName, objectPath, reader, objectSize)
	if err != nil {
		u.logger.Error().Err(err).Str("imagePath", objectPath).Msg("failed to save profile image")
		return status.Error(codes.Internal, dto.Err_INTERNAL_SAVE_PROFILE_IMAGE.Error())
	}
	return nil
}
func (u *userRepository) RemoveProfileImage(ctx context.Context, bucketName string, objectPath string) error {
	err := u.miniio.Delete(ctx, bucketName, objectPath)
	if err != nil {
		u.logger.Error().Err(err).Str("imagePath", objectPath).Msg("failed to remove profile image")
		return status.Error(codes.Internal, dto.Err_INTERNAL_REMOVE_PROFILE_IMAGE.Error())
	}
	return nil
}
