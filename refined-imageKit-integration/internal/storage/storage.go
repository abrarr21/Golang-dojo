package storage

import (
	"imakit-practice/internal/models"
	"mime/multipart"
)

type Storage interface {
	UploadImage(file multipart.File, filename string) (*models.Image, error)
	DeleteImage(fileID string) error
}
