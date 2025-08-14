package di

import (
	logemitter "10.1.20.130/dropping/file-service/config/log_emitter"
	logger "10.1.20.130/dropping/file-service/config/logger"
	messagequeue "10.1.20.130/dropping/file-service/config/message-queue"
	"10.1.20.130/dropping/file-service/config/router"
	minioCon "10.1.20.130/dropping/file-service/config/storage"
	"10.1.20.130/dropping/file-service/internal/domain/repository"
	"10.1.20.130/dropping/file-service/internal/domain/service"
	_logger "10.1.20.130/dropping/file-service/internal/infrastructure/logger"
	minioStorage "10.1.20.130/dropping/file-service/internal/infrastructure/storage"
	"go.uber.org/dig"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	if err := container.Provide(logger.New); err != nil {
		panic("Failed to provide logger: " + err.Error())
	}
	if err := container.Provide(messagequeue.New); err != nil {
		panic("Failed to provide nats connection: " + err.Error())
	}
	if err := container.Provide(messagequeue.NewJetstream); err != nil {
		panic("Failed to provide jetstream: " + err.Error())
	}
	if err := container.Provide(minioCon.NewMinioConnection); err != nil {
		panic("Failed to provide minio Connection: " + err.Error())
	}
	if err := container.Provide(logemitter.NewLogEmitter); err != nil {
		panic("Failed to provide log emitter: " + err.Error())
	}
	if err := container.Provide(_logger.NewLoggerInfra); err != nil {
		panic("Failed to provide logger infra: " + err.Error())
	}
	if err := container.Provide(minioStorage.NewMinioStorage); err != nil {
		panic("Failed to provide minio storage interface: " + err.Error())
	}
	if err := container.Provide(repository.NewUserRepository); err != nil {
		panic("Failed to provide user repository: " + err.Error())
	}
	if err := container.Provide(service.NewUserService); err != nil {
		panic("Failed to provide user repository: " + err.Error())
	}
	if err := container.Provide(router.NewGRPC); err != nil {
		panic("Failed to provide gRPC Server: " + err.Error())
	}
	return container
}
