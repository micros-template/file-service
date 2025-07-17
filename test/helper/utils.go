package helper

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ConnectGRPC(grpcURL string) (*grpc.ClientConn, error) {
	return grpc.NewClient(grpcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
