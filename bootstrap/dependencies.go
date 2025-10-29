package bootstrap

import (
	"github.com/nantapop-kj/go-auth-service/config"
	"github.com/nantapop-kj/go-auth-service/db"
	"github.com/nantapop-kj/go-auth-service/repository"
	"gorm.io/gorm"
)

type Dependencies struct {
	DB    *gorm.DB
	Minio repository.MinioRepository
	Redis repository.RedisRepository
}

func SetupDependencies() (*Dependencies, error) {
	database := db.ConnectDB()
	db.AutoMigrate(database)

	bucketName := config.Getenv("MINIO_BUCKET", "")
	minio := config.InitMinio()
	minioRepo := repository.NewMinioRepo(minio, bucketName)

	redis := config.InitRedis()
	redisRepo := repository.NewRedisRepo(redis)

	return &Dependencies{
		DB:    database,
		Redis: redisRepo,
		Minio: minioRepo,
	}, nil
}
