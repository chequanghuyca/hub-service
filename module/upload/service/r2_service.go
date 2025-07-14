package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"hub-service/common"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// R2Service struct để quản lý R2 operations
type R2Service struct {
	client     *minio.Client
	bucketName string
	workerURL  string
}

// NewR2Service tạo instance mới của R2Service
func NewR2Service() (*R2Service, error) {
	endpoint := os.Getenv("R2_ENDPOINT")
	accessKeyID := os.Getenv("R2_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	bucketName := os.Getenv("R2_BUCKET")
	workerURL := os.Getenv("R2_WORKER_URL")

	endpoint = strings.TrimPrefix(endpoint, "https://")

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init minio client: %w", err)
	}

	return &R2Service{
		client:     minioClient,
		bucketName: bucketName,
		workerURL:  workerURL,
	}, nil
}

// UploadToR2 uploads a file to Cloudflare R2 and returns the public Worker URL with unique UID filename
func (s *R2Service) UploadToR2(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	uid := common.NewUID(uint32(time.Now().UnixNano()), 1, 0).String()

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext == "" {
		ext = ".bin"
	}

	objectName := uid + ext
	contentType := fileHeader.Header.Get("Content-Type")

	_, err = s.client.PutObject(context.Background(), s.bucketName, objectName, file, fileHeader.Size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", fmt.Errorf("failed to upload to R2: %w", err)
	}

	url := s.workerURL + objectName
	return url, nil
}

func (s *R2Service) DeleteFile(fileName string) error {
	ctx := context.Background()

	err := s.client.RemoveObject(ctx, s.bucketName, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

func UploadToR2(fileHeader *multipart.FileHeader) (string, error) {
	service, err := NewR2Service()
	if err != nil {
		return "", err
	}
	return service.UploadToR2(fileHeader)
}
