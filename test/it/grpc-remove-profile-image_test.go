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
	minioContainer       *_helper.MinioContainer
	natsContainer        *_helper.NatsContainer
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

	mContainer, err := _helper.StartMinioContainer(s.ctx, s.network.Name, viper.GetString("container.minio_version"))
	if err != nil {
		log.Fatalf("failed starting minio container: %s", err)
	}
	s.minioContainer = mContainer
	// spawn nats
	nContainer, err := _helper.StartNatsContainer(s.ctx, s.network.Name, viper.GetString("container.nats_version"))
	if err != nil {
		log.Fatalf("failed starting minio container: %s", err)
	}
	s.natsContainer = nContainer

	fContainer, err := _helper.StartFileServiceContainer(s.ctx, s.network.Name, viper.GetString("container.file_service_version"))
	if err != nil {
		log.Println("make sure the image is exist")
		log.Fatalf("failed starting file service container: %s", err)
	}
	s.fileServiceContainer = fContainer

	gatewayContainer, err := _helper.StartGatewayContainer(s.ctx, s.network.Name, viper.GetString("container.gateway_version"))
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
