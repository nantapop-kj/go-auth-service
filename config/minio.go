package config

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func InitMinio() *minio.Client {
	endpoint := Getenv("MINIO_ENDPOINT", "")
	accessKeyID := Getenv("MINIO_ACCESS_KEY", "")
	secretAccessKey := Getenv("MINIO_SECRET_KEY", "")
	useSSL := true

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	})

	if err != nil {
		log.Fatalf("Failed to Connected Minio client: %v", err)
	}

	log.Println("\033[32m✅ [Minio] Connected successfully.\033[0m")

	return client
}
