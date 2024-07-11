package userUploads

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"shopperia/src/Uploads"
	"shopperia/src/common/models"
)

type uploadClient interface {
	MakeNewDirectory(fatherDirectory, NewDirName string) (string, error)
	MakeNewMediaRepositoryForUser(userId uuid.UUID) (string, error)

	CheckUserHasAMediaRepository(userId uuid.UUID) bool
	GetUserRepositoryPath(userId uuid.UUID) (string, error)

	UploadMedia(repositoryPath string, image models.UploadImageForm) (models.ImageData, error)
	UploadProfileImage(repositoryPath string, image models.UploadImageForm) (models.ImageData, error)
	UploadMultipleMediaResourcesOnRepository(repositoryPath string, images []models.UploadImageForm) ([]models.ImageData, error)

	GetMedia(repositoryPath, fileName, fileExtension string) (bytes.Buffer, error)
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

/*
func (C *uploadsClient) CreateUserRepository(userId uuid.UUID) (string, error) {

} */

func (C *uploadsClient) UploadMedia(userRepository string, image models.UploadImageForm) (models.ImageData, error) {
	if image.FileName == "" {
		return models.ImageData{}, errors.New("no file name provided")
	}

	imageDataChan := make(chan models.ImageData)
	errchan := make(chan error)

	go func() {
		if userRepository == "" {

			var err error

			userRepositoryExists := C.uploadClient.CheckUserHasAMediaRepository(image.UserID)

			if !userRepositoryExists {

				userRepository, err = C.uploadClient.MakeNewMediaRepositoryForUser(image.UserID)
				if err != nil {
					errchan <- err
					imageDataChan <- models.ImageData{}
					return
				}
			}

			userRepository, err = C.uploadClient.GetUserRepositoryPath(image.UserID)
			if err != nil {
				errchan <- err
				imageDataChan <- models.ImageData{}
				return
			}

		}

		imageData, err := C.uploadClient.UploadMedia(userRepository, image)
		if err != nil {
			errchan <- err
			return
		}

		imageDataChan <- imageData
		errchan <- nil
	}()

	err := <-errchan
	if err != nil {
		errStr := fmt.Sprintf("failed to upload media: %s", err)
		return models.ImageData{}, errors.New(errStr)
	}

	imageData := <-imageDataChan

	return imageData, nil
}

func (C *uploadsClient) GetMedia(repositoryPath string, form models.GetImageForm) (bytes.Buffer, error) {
	if repositoryPath == "" || form.FileName == "" || form.FileExtension == "" {
		return bytes.Buffer{}, errors.New("no parameters provide")
	}

	datachan := make(chan bytes.Buffer)
	errchan := make(chan error)
	go func() {
		imageData, err := C.uploadClient.GetMedia(repositoryPath, form.FileName, form.FileExtension)
		if err != nil {
			datachan <- bytes.Buffer{}
			errchan <- err
			return
		}

		datachan <- imageData
		errchan <- nil
	}()

	err := <-errchan
	data := <-datachan
	if err != nil {
		errStr := fmt.Sprintf("Error getting image: %s", err)
		return bytes.Buffer{}, errors.New(errStr)
	}

	return data, nil
}
