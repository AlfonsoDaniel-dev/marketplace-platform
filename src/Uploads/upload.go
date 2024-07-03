package Uploads

import (
	"errors"
	"fmt"
	"shopperia/src/common/models"
	"time"
)

func (US *UploadService) UploadMedia(image models.UploadImageForm) (models.ImageData, error) {
	image.FileName = image.FileName + image.UserID.String() + time.Now().String()

	imageData, err := US.upload(image.UserRepositoryPath, image.FileName, image.FileExtension, image.ImageData, image.UserID)
	if err != nil {
		return models.ImageData{}, err
	}

	return imageData, nil
}

func (US *UploadService) UploadProfileImage(image models.UploadImageForm) (models.ImageData, error) {
	image.FileName = image.UserID.String() + "_" + image.FileName + image.UserID.String() + "profilePic" + time.Now().String()

	imageData, err := US.upload(image.UserRepositoryPath, image.FileName, image.FileExtension, image.ImageData, image.UserID)
	if err != nil {
		return models.ImageData{}, err
	}

	return imageData, nil
}

func (US *UploadService) UploadMultipleMediaResources(images []models.UploadImageForm) ([]models.ImageData, error) {
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

		ImageData, err := US.upload(image.UserRepositoryPath, image.FileName, image.FileExtension, image.ImageData, image.UserID)
		if err != nil {
			errStr := fmt.Sprintf("Hubo un error al guardar la imagen %v, ERR: %v", image.FileName, err)
			return nil, errors.New(errStr)
		}

		imageData = append(imageData, ImageData)
	}

	return imageData, nil
}
