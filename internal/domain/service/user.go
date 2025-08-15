package service

import (
	"bytes"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/micros-template/file-service/internal/domain/repository"
	"github.com/micros-template/file-service/internal/infrastructure/logger"
	"github.com/micros-template/file-service/pkg/constant"
	"github.com/micros-template/file-service/pkg/utils"
	"github.com/rs/zerolog"
)

type (
	UserService interface {
		SaveProfileImage(context context.Context, imageByte []byte, imageExt string) (string, error)
		RemoveProfileImage(context context.Context, imageName string) error
	}
	userService struct {
		userRepository repository.UserRepository
		logEmitter     logger.LoggerInfra
		logger         zerolog.Logger
	}
)

func NewUserService(userRepository repository.UserRepository, logEmitter logger.LoggerInfra, logger zerolog.Logger) UserService {
	return &userService{
		userRepository: userRepository,
		logEmitter:     logEmitter,
		logger:         logger,
	}
}

func (u *userService) SaveProfileImage(context context.Context, imageBytes []byte, imageExt string) (string, error) {
	imageName := fmt.Sprintf("%s.%s", uuid.New().String(), imageExt)
	imagePath := fmt.Sprintf("%s/%s", constant.PROFILE_IMAGE_FOLDER, imageName)
	compressedBytes, err := utils.CompressImage(imageBytes, imageExt)
	if err != nil {
		if err := u.logEmitter.EmitLog("ERR", "failed to compress the file, use default instead"); err != nil {
			u.logger.Error().Err(err).Msg("failed to emit log")
		}
		err = u.userRepository.SaveProfileImage(context, constant.APP_BUCKET, imagePath, bytes.NewReader(imageBytes), int64(len(imageBytes)))
	} else {
		err = u.userRepository.SaveProfileImage(context, constant.APP_BUCKET, imagePath, bytes.NewReader(compressedBytes), int64(len(compressedBytes)))
	}
	if err != nil {
		return "", err
	}
	return imagePath, nil
}

func (u *userService) RemoveProfileImage(context context.Context, imageName string) error {
	if err := u.userRepository.RemoveProfileImage(context, constant.APP_BUCKET, imageName); err != nil {
		return err
	}
	return nil
}
