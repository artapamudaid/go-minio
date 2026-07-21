package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"

	"go-minio/internal/config"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func (s *storageService) getScheme() string {
	if s.cfg.UseSSL {
		return "https"
	}
	return "http"
}

type StorageService interface {
	UploadAsync(file io.Reader, filename string, fileSize int64, contentType string, client, folder string) string
	ListFiles(ctx context.Context, client, folder string) ([]string, error)
	GetFileStat(ctx context.Context, fileURL string) (map[string]interface{}, error)
	DeleteFile(ctx context.Context, fileURL string) error
}

type storageService struct {
	cfg         *config.Config
	minioClient *minio.Client
}

func NewStorageService(cfg *config.Config) (StorageService, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	return &storageService{
		cfg:         cfg,
		minioClient: client,
	}, nil
}

func (s *storageService) UploadAsync(file io.Reader, filename string, fileSize int64, contentType string, client, folder string) string {
	ext := ""
	if parts := strings.Split(filename, "."); len(parts) > 1 {
		ext = "." + parts[len(parts)-1]
	}
	randomName := uuid.New().String() + ext

	folderPath := ""
	if folder != "" {
		folderPath = strings.Trim(folder, "/") + "/"
	}
	objectKey := fmt.Sprintf("%s/uploads/%s%s", client, folderPath, randomName)

	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Goroutine untuk simulasi queue / async upload
	go func(key string, size int64, cType string) {
		ctx := context.Background()
		_, err := s.minioClient.PutObject(ctx, s.cfg.BucketName, key, file, size, minio.PutObjectOptions{
			ContentType: cType,
		})
		if err != nil {
			log.Printf("[ERROR] Upload background gagal (%s): %v\n", key, err)
		} else {
			log.Printf("[INFO] Berhasil upload background: %s\n", key)
		}
	}(objectKey, fileSize, contentType)

	return fmt.Sprintf("%s://%s/%s/%s", s.getScheme(), s.cfg.Endpoint, s.cfg.BucketName, objectKey)
}

func (s *storageService) ListFiles(ctx context.Context, client, folder string) ([]string, error) {
	folderPath := ""
	if folder != "" {
		folderPath = strings.Trim(folder, "/") + "/"
	}
	prefix := fmt.Sprintf("%s/uploads/%s", client, folderPath)

	objectCh := s.minioClient.ListObjects(ctx, s.cfg.BucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	var files []string
	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}
		files = append(files, fmt.Sprintf("%s://%s/%s/%s", s.getScheme(), s.cfg.Endpoint, s.cfg.BucketName, object.Key))
	}

	if files == nil {
		files = []string{}
	}
	return files, nil
}

func (s *storageService) GetFileStat(ctx context.Context, fileURL string) (map[string]interface{}, error) {
	key := extractKeyFromURL(fileURL, s.cfg.BucketName)
	stat, err := s.minioClient.StatObject(ctx, s.cfg.BucketName, key, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"url":           fileURL,
		"key":           key,
		"size":          stat.Size,
		"type":          stat.ContentType,
		"last_modified": stat.LastModified.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *storageService) DeleteFile(ctx context.Context, fileURL string) error {
	key := extractKeyFromURL(fileURL, s.cfg.BucketName)
	return s.minioClient.RemoveObject(ctx, s.cfg.BucketName, key, minio.RemoveObjectOptions{})
}

func extractKeyFromURL(rawURL, bucket string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	path := strings.TrimPrefix(u.Path, "/")
	prefixToRemove := bucket + "/"
	if strings.HasPrefix(path, prefixToRemove) {
		return strings.TrimPrefix(path, prefixToRemove)
	}
	return path
}
