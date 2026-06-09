package storage

import (
	"context"
	"fmt"
	"imakit-practice/internal/config"
	"imakit-practice/internal/models"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	imagekit "github.com/imagekit-developer/imagekit-go/v2"
	"github.com/imagekit-developer/imagekit-go/v2/option"
)

type ImageKitStorage struct {
	client imagekit.Client
}

func NewImageKitStorage(cfg *config.ImageKitConfig) *ImageKitStorage {
	return &ImageKitStorage{
		client: imagekit.NewClient(option.WithPrivateKey(cfg.ImgPrivateKey)),
	}
}

func (s *ImageKitStorage) UploadImage(file multipart.File, filename string) (*models.Image, error) {

	// validate MIME type by sniffing first 512 byte
	buffer := make([]byte, 512)

	_, err := file.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	contentType := http.DetectContentType(buffer)

	switch contentType {
	case "image/jpeg", "image/jpg", "image/png", "image/webp":
	default:
		return nil, fmt.Errorf("unsupported image type: %s", contentType)
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("failed to reset file pointer: %w", err)
	}

	// cross-validate: extension must also be an alllowed image type
	ext := strings.ToLower(filepath.Ext(filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}
	if !allowedExts[ext] {
		return nil, fmt.Errorf("unsupported file extension: %s", ext)
	}

	uniqueFilename := fmt.Sprintf(
		"%d%s",
		time.Now().UnixNano(),
		ext,
	)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	resp, err := s.client.Files.Upload(
		ctx,
		imagekit.FileUploadParams{
			File:     file,
			FileName: uniqueFilename,
			Folder:   imagekit.String("/users"),
		},
	)

	if err != nil {
		return nil, fmt.Errorf("imagkit upload failed: %w", err)
	}

	return &models.Image{
		URL:    resp.URL,
		Name:   resp.Name,
		FileID: resp.FileID,
	}, nil
}

func (s *ImageKitStorage) DeleteImage(fileID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.client.Files.Delete(ctx, fileID)
	if err != nil {
		return fmt.Errorf("imagekit delete failed: %w", err)
	}

	return nil
}
