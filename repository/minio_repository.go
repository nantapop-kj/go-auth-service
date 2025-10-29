package repository

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/chai2010/webp"
	"github.com/minio/minio-go/v7"
	"github.com/nfnt/resize"
)

type MinioRepository interface {
	GetFileURL(ctx context.Context, filename string) (string, error)
	UploadImage(ctx context.Context, filename string, content []byte, maxWidth, maxHeight uint) (string, error)
	DeleteFile(ctx context.Context, filename string) error
}

type MinioRepo struct {
	Client     *minio.Client
	BucketName string
}

func NewMinioRepo(client *minio.Client, bucket string) MinioRepository {
	return &MinioRepo{
		Client:     client,
		BucketName: bucket,
	}
}

func decodeImage(content []byte, ext string) (image.Image, error) {
	reader := bytes.NewReader(content)
	ext = strings.ToLower(ext)

	switch ext {
	case ".jpg", ".jpeg":
		return jpeg.Decode(reader)
	case ".png":
		return png.Decode(reader)
	default:
		img, _, err := image.Decode(bytes.NewReader(content))
		if err != nil {
			return nil, fmt.Errorf("unsupported image format: %s", ext)
		}
		return img, nil
	}
}

func changeExtension(filename, newExt string) string {
	name := strings.TrimSuffix(filename, filepath.Ext(filename))
	return name + newExt
}

func (m *MinioRepo) GetFileURL(ctx context.Context, filename string) (string, error) {
	url, err := m.Client.PresignedGetObject(ctx, m.BucketName, filename, 10*time.Minute, nil)
	if err != nil {
		log.Printf("Failed to generate presigned URL for %s: %v", filename, err)
		return "", err
	}

	return url.String(), nil
}

func (m *MinioRepo) UploadImage(ctx context.Context, filename string, content []byte, maxWidth, maxHeight uint) (string, error) {
	img, err := decodeImage(content, filepath.Ext(filename))
	if err != nil {
		log.Printf("[MinIO] Failed to decode image %s: %v", filename, err)
		return "", fmt.Errorf("invalid image format: %w", err)
	}

	if maxWidth > 0 || maxHeight > 0 {
		img = resize.Thumbnail(maxWidth, maxHeight, img, resize.Lanczos3)
	}

	webpData, err := webp.EncodeRGBA(img, 85)
	if err != nil {
		log.Printf("[MinIO] Failed to encode image to WebP: %v", err)
		return "", fmt.Errorf("failed to encode webp: %w", err)
	}

	webpFilename := changeExtension(filename, ".webp")

	res, err := m.Client.PutObject(
		ctx,
		m.BucketName,
		webpFilename,
		bytes.NewReader(webpData),
		int64(len(webpData)),
		minio.PutObjectOptions{
			ContentType: "image/webp",
		},
	)
	if err != nil {
		log.Printf("[MinIO] Failed to upload file %s: %v", webpFilename, err)
		return "", err
	}
	return res.Key, nil
}

func (m *MinioRepo) DeleteFile(ctx context.Context, filename string) error {
	err := m.Client.RemoveObject(ctx, m.BucketName, filename, minio.RemoveObjectOptions{})
	if err != nil {
		log.Printf("[MinIO] Failed to delete file %s: %v", filename, err)
	}
	return err
}
