package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"10.1.20.130/dropping/file-service/internal/domain/handler"
	"10.1.20.130/dropping/file-service/internal/domain/service"
	_logger "10.1.20.130/dropping/file-service/internal/infrastructure/logger"
	"10.1.20.130/dropping/file-service/internal/infrastructure/storage"
	"10.1.20.130/dropping/file-service/pkg/constant"
	"github.com/rs/zerolog"
	"go.uber.org/dig"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	Container   *dig.Container
	ServerReady chan bool
	Address     string
}

func (s *GRPCServer) Run(ctx context.Context) {
	err := s.Container.Invoke(func(
		grpcServer *grpc.Server,
		logger zerolog.Logger,
		svc service.UserService,
		minio storage.MinioStorage,
		logEmitter _logger.LoggerInfra,
	) {
		err := minio.InitBucket(ctx, constant.APP_BUCKET)
		if err != nil {
			go func() {
				if err := logEmitter.EmitLog("ERR", fmt.Sprintf("Failed to init bucket: %v", err)); err != nil {
					logger.Error().Err(err).Msg("failed to emit log")
				}
			}()
			log.Fatalf("Failed to init bucket: %v", err)
		}

		listen, err := net.Listen("tcp", s.Address)
		if err != nil {
			go func() {
				if err := logEmitter.EmitLog("ERR", fmt.Sprintf("failed to listen:%v", err)); err != nil {
					logger.Error().Err(err).Msg("failed to emit log")
				}
			}()
			logger.Fatal().Msgf("failed to listen:%v", err)
		}
		handler.RegisterUserService(grpcServer, svc)

		go func() {
			if serveErr := grpcServer.Serve(listen); serveErr != nil {
				go func() {
					if err := logEmitter.EmitLog("ERR", fmt.Sprintf("gRPC server error: %v", serveErr)); err != nil {
						logger.Error().Err(err).Msg("failed to emit log")
					}
				}()
				logger.Fatal().Msgf("gRPC server error: %v", serveErr)
			}
		}()
		if s.ServerReady != nil {
			for range 50 {
				conn, err := net.DialTimeout("tcp", s.Address, 100*time.Millisecond)
				if err == nil {
					if err := conn.Close(); err != nil {
						go func() {
							if err := logEmitter.EmitLog("ERR", "establish check connection failed to close"); err != nil {
								logger.Error().Err(err).Msg("failed to emit log")
							}
						}()
						logger.Fatal().Err(err).Msg("establish check connection failed to close")
					}
					s.ServerReady <- true
					break
				}
				time.Sleep(100 * time.Millisecond)
			}
		}
		go func() {
			if err := logEmitter.EmitLog("INFO", fmt.Sprintf("gRPC server running in port %s", s.Address)); err != nil {
				logger.Error().Err(err).Msg("failed to emit log")
			}
		}()
		logger.Info().Msg("gRPC server running in port " + s.Address)

		<-ctx.Done()
		go func() {
			if err := logEmitter.EmitLog("INFO", "Shutting down gRPC server..."); err != nil {
				logger.Error().Err(err).Msg("failed to emit log")
			}
		}()
		logger.Info().Msg("Shutting down gRPC server...")
		grpcServer.GracefulStop()
		go func() {
			if err := logEmitter.EmitLog("INFO", "gRPC server stopped gracefully."); err != nil {
				logger.Error().Err(err).Msg("failed to emit log")
			}
		}()
		logger.Info().Msg("gRPC server stopped gracefully.")
	})
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}
}
