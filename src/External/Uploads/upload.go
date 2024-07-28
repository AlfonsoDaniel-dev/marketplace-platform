package Uploads

import (
	"errors"
	"fmt"
	"shopperia/src/common/models"
	"strings"
	"time"
)

func (US *UploadService) UploadMedia(repositoryPath string, image models.UploadImageForm) (models.ImageData, error) {
	image.FileName = image.FileName + image.UserID.String() + time.Now().String()

	imageData, err := US.upload(repositoryPath, image.FileName, image.FileExtension, image.ImageData, image.UserID)
	if err != nil {
		return models.ImageData{}, err
	}

	return imageData, nil
}

func (US *UploadService) UploadProfileImage(repositoryPath string, image models.UploadImageForm) (models.ImageData, error) {
	image.FileName = strings.ReplaceAll(image.FileName, " ", "_")
	image.FileName = image.UserID.String() + "_" + image.FileName + "_" + "profilePic"

	imageData, err := US.upload(repositoryPath, image.FileName, image.FileExtension, image.ImageData, image.UserID)
	if err != nil {
		return models.ImageData{}, err
	}

	return imageData, nil
}

func (US *UploadService) UploadMultipleMediaResourcesOnRepository(repositorytPath string, images []models.UploadImageForm) ([]models.ImageData, error) {
	if len(images) == 0 {
		return []models.ImageData{}, errors.New("no images to upload")
	}

	imageData := []models.ImageData{}

	for i, image := range images {
		NextImage := images[i+1]

		addition := "0" + string(i)
		image.FileName = addition

		if image.FileName == NextImage.FileName {
			errStr := fmt.Sprintf("some files have the same name, conflic on: %s with %s", image.FileName, NextImage.FileName)
			return nil, errors.New(errStr)
		}

		ImageData, err := US.upload(repositorytPath, image.FileName, image.FileExtension, image.ImageData, image.UserID)
		if err != nil {
			errStr := fmt.Sprintf("Hubo un error al guardar la imagen %v, ERR: %v", image.FileName, err)
			return nil, errors.New(errStr)
		}

		imageData = append(imageData, ImageData)
	}

	return imageData, nil
}
