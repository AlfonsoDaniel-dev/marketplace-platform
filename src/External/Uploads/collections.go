package Uploads

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"os"
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

func (US *UploadService) InsertImageOnCollection(collectionPath string, image models.UploadImageForm) (models.ImageData, error) {
	if collectionPath == "" {
		return models.ImageData{}, errors.New("collection path is empty")
	}

	imageData, err := US.upload(collectionPath, image.FileName, image.FileExtension, image.ImageData, image.UserID)
	if err != nil {
		return models.ImageData{}, err
	}

	return imageData, nil
}

type getImageAttempt struct {
	completeFilePath string
	fileName         string
	fileExtension    string
	status           error
	data             models.GetImage
	done             chan struct{}
}

func searchWorker(numAttemps int, attempt chan *getImageAttempt) {
	for numAttemps != 0 {
		request := <-attempt

		imgBytes, err := os.ReadFile(request.completeFilePath)
		if err != nil {
			request.status = err
			request.done <- struct{}{}
		}

		var buf bytes.Buffer

		_, err = buf.Write(imgBytes)
		if err != nil {
			request.status = err
			request.done <- struct{}{}
		}

		data := models.GetImage{
			FileName:      request.fileName,
			FileExtension: request.fileExtension,
			ImageBuffer:   buf,
		}

		request.data = data

		numAttemps--

		request.done <- struct{}{}
	}
}

func (US *UploadService) GetAllMediaFromCollection(repositoryPath, collectionPath string, forms []models.GetImageForm) ([]models.GetImage, error) {
	if repositoryPath == "" || collectionPath == "" {
		return nil, errors.New("repository path or collection path is empty")
	}

	imgsData := []models.GetImage{}

	requestChan := make(chan *getImageAttempt)
	go searchWorker(len(forms), requestChan)

	for _, file := range forms {

		if file.FileName == "" || file.FileExtension == "" {
			return nil, errors.New("file name or file extension is empty")
		}

		searchFileName := US.OriginPath + "/" + repositoryPath + "/" + collectionPath + "/" + file.FileName + "." + file.FileExtension

		attempt := getImageAttempt{
			completeFilePath: searchFileName,
			fileName:         file.FileName,
			fileExtension:    file.FileExtension,
		}

		requestChan <- &attempt

		<-attempt.done
	}

	for attempt := range requestChan {
		if attempt.status != nil {
			return imgsData, attempt.status
		}

		imgsData = append(imgsData, attempt.data)
	}

	close(requestChan)

	return imgsData, nil

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
