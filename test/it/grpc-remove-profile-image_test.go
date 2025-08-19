package it

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/micros-template/file-service/test/helper"
	"github.com/micros-template/proto-file/pkg/fpb"
	_helper "github.com/micros-template/sharedlib/test/helper"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
)

type GRPCSaveProfileImageHandlerSuite struct {
	suite.Suite
	ctx context.Context

	network              *testcontainers.DockerNetwork
	gatewayContainer     *_helper.GatewayContainer
	minioContainer       *_helper.StorageContainer
	natsContainer        *_helper.MessageQueueContainer
	fileServiceContainer *_helper.FileServiceContainer
}

func (s *GRPCSaveProfileImageHandlerSuite) SetupSuite() {
	log.Println("Setting up integration test suite for GRPCSaveProfileImageHandlerSuite")
	s.ctx = context.Background()

	viper.SetConfigName("config.test")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../../")
	if err := viper.ReadInConfig(); err != nil {
		panic("failed to read config")
	}

	s.network = _helper.StartNetwork(s.ctx)

	mContainer, err := _helper.StartStorageContainer(_helper.StorageParameterOption{
		Context:       s.ctx,
		SharedNetwork: s.network.Name,
		ImageName:     viper.GetString("container.minio_image"),
		ContainerName: "test-minio",
		WaitingSignal: "API:",
		Cmd:           []string{"server", "/data"},
		Env: map[string]string{
			"MINIO_ROOT_USER":     viper.GetString("minio.credential.user"),
			"MINIO_ROOT_PASSWORD": viper.GetString("minio.credential.password"),
		},
	})
	if err != nil {
		log.Fatalf("failed starting minio container: %s", err)
	}
	s.minioContainer = mContainer

	// spawn nats
	nContainer, err := _helper.StartMessageQueueContainer(_helper.MessageQueueParameterOption{
		Context:            s.ctx,
		SharedNetwork:      s.network.Name,
		ImageName:          viper.GetString("container.nats_image"),
		ContainerName:      "test_nats",
		MQConfigPath:       viper.GetString("script.nats_server"),
		MQInsideConfigPath: "/etc/nats/nats.conf",
		WaitingSignal:      "Server is ready",
		MappedPort:         []string{"4221:4221/tcp"},
		Cmd: []string{
			"-c", "/etc/nats/nats.conf",
			"--name", "nats",
			"-p", "4221",
		},
		Env: map[string]string{
			"NATS_USER":     viper.GetString("nats.credential.user"),
			"NATS_PASSWORD": viper.GetString("nats.credential.password"),
		},
	})
	if err != nil {
		log.Fatalf("failed starting minio container: %s", err)
	}
	s.natsContainer = nContainer

	fContainer, err := _helper.StartFileServiceContainer(_helper.FileServiceParameterOption{
		Context:       s.ctx,
		SharedNetwork: s.network.Name,
		ImageName:     viper.GetString("container.file_service_image"),
		ContainerName: "test_file_service",
		WaitingSignal: "gRPC server running in port",
		Cmd:           []string{"/file_service"},
		Env:           map[string]string{"ENV": "test"},
	})

	if err != nil {
		log.Fatalf("failed starting file service container: %s", err)
	}
	s.fileServiceContainer = fContainer

	gatewayContainer, err := _helper.StartGatewayContainer(_helper.GatewayParameterOption{
		Context:                   s.ctx,
		SharedNetwork:             s.network.Name,
		ImageName:                 viper.GetString("container.gateway_image"),
		ContainerName:             "test_gateway",
		NginxConfigPath:           viper.GetString("script.nginx"),
		NginxInsideConfigPath:     "/etc/nginx/conf.d/default.conf",
		GrpcErrorConfigPath:       viper.GetString("script.grpc_error"),
		GrpcErrorInsideConfigPath: "/etc/nginx/conf.d/errors.grpc_conf",
		WaitingSignal:             "Configuration complete; ready for start up",
		MappedPort:                []string{"9090:80/tcp", "50051:50051/tcp"},
	})
	if err != nil {
		log.Fatalf("failed starting gateway container: %s", err)
	}
	s.gatewayContainer = gatewayContainer

	time.Sleep(time.Second)
}
func (s *GRPCSaveProfileImageHandlerSuite) TearDownSuite() {

	if err := s.minioContainer.Terminate(s.ctx); err != nil {
		log.Fatalf("error terminating minio container: %s", err)
	}

	if err := s.natsContainer.Terminate(s.ctx); err != nil {
		log.Fatalf("error terminating nats container: %s", err)
	}
	if err := s.fileServiceContainer.Terminate(s.ctx); err != nil {
		log.Fatalf("error terminating file service container: %s", err)
	}

	if err := s.gatewayContainer.Terminate(s.ctx); err != nil {
		log.Fatalf("error terminating gateway container: %s", err)
	}
	log.Println("Tear Down integration test suite for GRPCSaveProfileImageHandlerSuite")

}
func TestGRPCSaveProfileImageHandlerSuite(t *testing.T) {
	suite.Run(t, &GRPCSaveProfileImageHandlerSuite{})
}

func (s *GRPCSaveProfileImageHandlerSuite) TestUserHandler_SaveProfileImageHandler_Success() {

	conn, err := helper.ConnectGRPC("localhost:50051")
	s.Require().NoError(err, "Failed to connect to gRPC server")
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("failed to close gRPC connection: %v", err)
		}
	}()

	fileServiceClient := fpb.NewFileServiceClient(conn)
	imageName, err := fileServiceClient.SaveProfileImage(s.ctx, &fpb.Image{Image: []byte{}, Ext: "jpg"})

	s.NotEmpty(imageName)
	s.NoError(err)
}
