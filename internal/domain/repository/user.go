package repository

import (
	"context"
	"fmt"
	"io"

	"github.com/micros-template/file-service/internal/domain/dto"
	"github.com/micros-template/file-service/internal/infrastructure/logger"
	"github.com/micros-template/file-service/internal/infrastructure/storage"
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
		minio      storage.MinioStorage
		logger     zerolog.Logger
		logEmitter logger.LoggerInfra
	}
)

func NewUserRepository(minio storage.MinioStorage, logEmitter logger.LoggerInfra, logger zerolog.Logger) UserRepository {
	return &userRepository{
		minio:      minio,
		logEmitter: logEmitter,
		logger:     logger,
	}
}

func (u *userRepository) SaveProfileImage(ctx context.Context, bucketName string, objectPath string, reader io.Reader, objectSize int64) error {
	err := u.minio.Set(ctx, bucketName, objectPath, reader, objectSize)
	if err != nil {
		if err := u.logEmitter.EmitLog("ERR", fmt.Sprintf("failed to save profile image. image path: %s. error:%s", objectPath, err.Error())); err != nil {
			u.logger.Error().Err(err).Msg("failed to emit log")
		}
		return status.Error(codes.Internal, dto.Err_INTERNAL_SAVE_PROFILE_IMAGE.Error())
	}
	return nil
}
func (u *userRepository) RemoveProfileImage(ctx context.Context, bucketName string, objectPath string) error {
	err := u.minio.Delete(ctx, bucketName, objectPath)
	if err != nil {
		if err := u.logEmitter.EmitLog("ERR", fmt.Sprintf("failed to save remove image. image path: %s. error:%s", objectPath, err.Error())); err != nil {
			u.logger.Error().Err(err).Msg("failed to emit log")
		}
		return status.Error(codes.Internal, dto.Err_INTERNAL_REMOVE_PROFILE_IMAGE.Error())
	}
	return nil
}
