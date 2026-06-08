package services

import (
	"context"
	"log"
	"mime/multipart"

	"github.com/abrarr21/auth-practice-3/internal/models"
	imagekit "github.com/imagekit-developer/imagekit-go/v2"
	"github.com/imagekit-developer/imagekit-go/v2/option"
)

var Client imagekit.Client

func InitImageKit() {
	log.Println("initializing imagekit")
	Client = imagekit.NewClient(
		option.WithPrivateKey("private_84DrMMeyPYz+5fKsyhdd46UNQ38="),
	)
}

func UploadImage(
	file multipart.File,
	filename string,
) (*models.Image, error) {

	resp, err := Client.Files.Upload(
		context.Background(),
		imagekit.FileUploadParams{
			File:     file,
			FileName: filename,
		},
	)

	if err != nil {
		return nil, err
	}

	return &models.Image{
		URL:    resp.URL,
		Name:   resp.Name,
		FileID: resp.FileID,
	}, nil
}
