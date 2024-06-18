package userUploads

import (
	"bytes"
	"errors"
	"shopperia/src/Uploads"
	"shopperia/src/common/models"
)

type uploadClient interface {
	UploadMedia(image models.UploadImageForm) (models.ImageData, error)
	UploadProfileImage(image models.UploadImageForm) (models.ImageData, error)
	GetMedia(filepath string) (bytes.Buffer, error)
	SaveMultipleMediaResources(images []models.UploadImageForm) ([]models.ImageData, error)
	CreateCollection(form models.CreateCollectionForm) (models.CollectionData, error)
	GetMultipleMediaResources(repositoryPath string, fileNames []string) ([]models.GetImage, error)
	InsertMultipleMediaIntoCollection(collectionPath models.CollectionData, images []models.UploadImageForm) ([]models.ImageData, error)
	GetMultipleDataFromPath(Fatherpath string, filenames []string) ([]models.GetImageData, error)
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
