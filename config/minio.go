package config

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/nantapop-kj/go-auth-service/enums"
)

func CreateMinioBucket(client *minio.Client, bucketName string) {
	ctx := context.Background()

	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		log.Fatalf("❌ [MinIO] Failed to check bucket %s: %v", bucketName, err)
	}

	if exists {
		return
	}

	err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		log.Fatalf("❌ [MinIO] Failed to create bucket %s: %v", bucketName, err)
	}

	log.Println("\033[32m✅ [Minio] Created new bucket successfully.\033[0m")
}

func InitMinio() *minio.Client {
	endpoint := GetEnv("MINIO_ENDPOINT", "")
	accessKeyID := GetEnv("MINIO_ACCESS_KEY", "")
	secretAccessKey := GetEnv("MINIO_SECRET_KEY", "")
	bucketName := GetEnv("MINIO_BUCKET", "mybucket")
	isProd := GetEnv("APP_ENV", "") == string(enums.EnvProd)
	useSSL := isProd

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

	CreateMinioBucket(client, bucketName)

	return client
}
