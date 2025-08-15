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

type GRPCRemoveProfileImageHandlerSuite struct {
	suite.Suite
	ctx context.Context

	network              *testcontainers.DockerNetwork
	gatewayContainer     *_helper.GatewayContainer
	minioContainer       *_helper.MinioContainer
	natsContainer        *_helper.NatsContainer
	fileServiceContainer *_helper.FileServiceContainer
}

func (r *GRPCRemoveProfileImageHandlerSuite) SetupSuite() {
	log.Println("Setting up integration test suite for GRPCRemoveProfileImageHandlerSuite")
	r.ctx = context.Background()

	viper.SetConfigName("config.test")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../../")
	if err := viper.ReadInConfig(); err != nil {
		panic("failed to read config")
	}

	r.network = _helper.StartNetwork(r.ctx)

	mContainer, err := _helper.StartMinioContainer(r.ctx, r.network.Name, viper.GetString("container.minio_version"))
	if err != nil {
		log.Fatalf("failed starting minio container: %r", err)
	}
	r.minioContainer = mContainer

	// spawn nats
	nContainer, err := _helper.StartNatsContainer(r.ctx, r.network.Name, viper.GetString("container.nats_version"))
	if err != nil {
		log.Fatalf("failed starting minio container: %r", err)
	}
	r.natsContainer = nContainer

	fContainer, err := _helper.StartFileServiceContainer(r.ctx, r.network.Name, viper.GetString("container.file_service_version"))
	if err != nil {
		log.Println("make sure the image is exist")
		log.Fatalf("failed starting file service container: %r", err)
	}
	r.fileServiceContainer = fContainer

	gatewayContainer, err := _helper.StartGatewayContainer(r.ctx, r.network.Name, viper.GetString("container.gateway_version"))
	if err != nil {
		log.Fatalf("failed starting gateway container: %r", err)
	}
	r.gatewayContainer = gatewayContainer
	time.Sleep(time.Second)
}
func (r *GRPCRemoveProfileImageHandlerSuite) TearDownSuite() {

	if err := r.minioContainer.Terminate(r.ctx); err != nil {
		log.Fatalf("error terminating minio container: %r", err)
	}

	if err := r.natsContainer.Terminate(r.ctx); err != nil {
		log.Fatalf("error terminating nats container: %s", err)
	}

	if err := r.fileServiceContainer.Terminate(r.ctx); err != nil {
		log.Fatalf("error terminating file service container: %r", err)
	}

	if err := r.gatewayContainer.Terminate(r.ctx); err != nil {
		log.Fatalf("error terminating gateway container: %r", err)
	}
	log.Println("Tear Down integration test suite for GRPCRemoveProfileImageHandlerSuite")

}
func TestGRPCRemoveProfileImageHandlerSuite(t *testing.T) {
	suite.Run(t, &GRPCRemoveProfileImageHandlerSuite{})
}

func (r *GRPCRemoveProfileImageHandlerSuite) TestUserHandler_RemoveProfileImageHandler_Success() {
	conn, err := helper.ConnectGRPC("localhost:50051")
	r.Require().NoError(err, "Failed to connect to gRPC server")
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("failed to close gRPC connection: %v", err)
		}
	}()

	fileServiceClient := fpb.NewFileServiceClient(conn)
	imageName, err := fileServiceClient.SaveProfileImage(r.ctx, &fpb.Image{Image: []byte{}, Ext: "jpg"})
	r.NotEmpty(imageName)
	r.NoError(err)

	n := imageName.GetName()
	status, err := fileServiceClient.RemoveProfileImage(r.ctx, &fpb.ImageName{
		Name: n,
	})
	r.NotEmpty(status)
	r.NoError(err)
	r.Equal(status.GetStatus(), true)
}
