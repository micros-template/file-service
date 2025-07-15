package storage

import (
	"fmt"
	"log"
	"net"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
)

func NewMinioConnection() *minio.Client {

	hostPort := viper.GetString("minio.host") + ":" + viper.GetString("minio.port")
	fmt.Print("Checking MinIO connectivity: ", hostPort, "\n")

	// Connectivity check
	conn, err := net.Dial("tcp", hostPort)
	if err != nil {
		log.Fatalf("Cannot connect to MinIO at %s: %v", hostPort, err)
	}
	conn.Close()

	minioClient, err := minio.New(hostPort, &minio.Options{
		Creds: credentials.NewStaticV4(
			viper.GetString("minio.credential.user"),
			viper.GetString("minio.credential.password"),
			"",
		),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}
	return minioClient
}
