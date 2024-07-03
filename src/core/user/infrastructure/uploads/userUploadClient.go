package userUploads

import (
	"bytes"
	"errors"
	"github.com/google/uuid"
	"shopperia/src/Uploads"
	"shopperia/src/common/models"
)

type uploadClient interface {
	MakeNewDirectory(fatherDirectory, NewDirName string) (string, error)
	MakeNewMediaRepositoryForUser(userId uuid.UUID) (string, error)

	CheckUserHasAMediaRepository(userId uuid.UUID) bool
	GetUserRepositoryPath(userId uuid.UUID) (string, error)

	UploadMedia(image models.UploadImageForm) (models.ImageData, error)
	UploadMultipleMediaResources(images []models.UploadImageForm) ([]models.ImageData, error)

	GetMedia(filePath string) (bytes.Buffer, error)
}

var userPath string = "./src/core/user/infrastructure/uploads/main/repository"

type uploadsClient struct {
	uploadClient
}

func NewUploadsClient() *uploadsClient {

	client := Uploads.NewUploadService(userPath)

	return &uploadsClient{
		uploadClient: client,
	}
}

func (C *uploadsClient) UploadMedia(image models.UploadImageForm) (models.ImageData, error) {
	if image.FileName == "" {
		return models.ImageData{}, errors.New("no file name provided")
	}

	imageData, err := C.uploadClient.UploadMedia(image)
	if err != nil {
		return models.ImageData{}, err
	}

	return imageData, nil
}
