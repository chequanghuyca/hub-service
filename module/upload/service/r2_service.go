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

// UploadToR2 uploads a file to Cloudflare R2 and returns the public Worker URL with unique UID filename
func UploadToR2(fileHeader *multipart.FileHeader) (string, error) {
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
		return "", fmt.Errorf("failed to init minio client: %w", err)
	}

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

	_, err = minioClient.PutObject(context.Background(), bucketName, objectName, file, fileHeader.Size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", fmt.Errorf("failed to upload to R2: %w", err)
	}

	url := workerURL + objectName
	return url, nil
}
