package it

import (
	"context"
	"log"
	"testing"
	"time"

	"10.1.20.130/dropping/file-service/test/helper"
	_helper "10.1.20.130/dropping/sharedlib/test/helper"
	"github.com/dropboks/proto-file/pkg/fpb"
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
	defer conn.Close()

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
