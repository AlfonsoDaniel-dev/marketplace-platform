package Uploads

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"shopperia/src/common/models"
	UserDTO "shopperia/src/core/user/domain/DTO"
)

func (US *UploadService) makeNewCollection(userRepositoryPath, collectionName string) (string, error) {
	if userRepositoryPath == "" || collectionName == "" {
		return "", errors.New("collection name or collection name is empty")
	}
	collectionPath, collErr := US.MakeNewDirectory(userRepositoryPath, collectionName)
	if collErr != nil {
		return "", collErr
	}

	return collectionPath, nil
}

func (US *UploadService) CreateCollection(form UserDTO.CreateCollectionForm) (models.CollectionData, error) {
	if form.UserRepositoryPath == "" || form.CollectionName == "" {
		errStr := fmt.Sprint("Please provide all fields")
		return models.CollectionData{}, errors.New(errStr)
	}

	collection, err := US.makeNewCollection(form.UserRepositoryPath, form.CollectionName)
	if err != nil {
		return models.CollectionData{}, err
	}

	newCollectionId := uuid.New()

	collectionData := models.CollectionData{
		CollectionId:   newCollectionId,
		UserRepository: form.UserRepositoryPath,
		CollectionName: form.CollectionName,
		CollectionPath: collection,
	}

	return collectionData, nil
}

func (US *UploadService) InsertMultipleMediaIntoCollection(collectionPath models.CollectionData, images []models.UploadImageForm) ([]models.ImageData, error) {
	if collectionPath.CollectionPath == "" || len(images) == 0 {
		return []models.ImageData{}, errors.New("path no found or no images to upload")
	}

	for i, image := range images {
		NextImage := images[i+1]

		if image.FileName == NextImage.FileName {
			NextImage.FileName = NextImage.FileName + string(i+1)
		}
	}

	imagesData, err := US.UploadMultipleMediaResourcesOnRepository(collectionPath.CollectionPath, images)
	if err != nil {
		return []models.ImageData{}, err
	}

	return imagesData, nil
}
