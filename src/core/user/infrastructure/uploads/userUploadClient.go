package userUploads

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"shopperia/src/External/Uploads"
	"shopperia/src/common/models"
	UserDTO "shopperia/src/core/user/domain/DTO"
	"sync"
)

type uploadClient interface {
	MakeNewDirectory(fatherDirectory, NewDirName string) (string, error)
	MakeNewMediaRepositoryForUser(userId uuid.UUID) (string, error)

	CheckUserHasAMediaRepository(userId uuid.UUID) bool
	GetUserRepositoryPath(userId uuid.UUID) (string, error)

	CreateCollection(form UserDTO.CreateCollectionForm) (models.CollectionData, error)
	InsertImageOnCollection(repositoryPath, CollectionPath string, image models.UploadImageForm) (models.ImageData, error)
	InsertMultipleImagesOnCollection(collectionPath string, forms []models.UploadImageForm) ([]models.ImageData, error)
	GetAllMediaFromCollection(repositoryPath, collectionPath string, forms []models.GetImageForm) ([]models.GetImage, error)
	UpdateImageOnCollection(collectionPath, fileName, fileExtension string, form models.UploadImageForm) (models.ImageData, error)
	DeleteImageOnCollection(request models.DeleteOnCollectionRequest) error
	DeleteMultipleImagesOnCollection(requests []models.DeleteOnCollectionRequest) error
	CheckIfImageIsOnCollection(userRepository, collectionName, fileName, fileExtension string) bool
	UpdateCollectionName(userRepository, collectionName, newCollectionName string, userID uuid.UUID) (string, string, error)

	UploadMedia(repositoryPath string, image models.UploadImageForm) (models.ImageData, error)
	UploadProfileImage(repositoryPath string, image models.UploadImageForm) (models.ImageData, error)
	UploadMultipleMediaResourcesOnRepository(repositoryPath string, images []models.UploadImageForm) ([]models.ImageData, error)

	GetMedia(repositoryPath, fileName, fileExtension string) (bytes.Buffer, error)
	GetProfileImage(repositoryPath, fileName, fileExtension string, userId uuid.UUID) (bytes.Buffer, error)
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

		return
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

		return
	}()

	err := <-errchan
	data := <-datachan
	if err != nil {
		errStr := fmt.Sprintf("Error getting image: %s", err)
		return bytes.Buffer{}, errors.New(errStr)
	}

	return data, nil
}

func (C *uploadsClient) UploadProfilePicture(repositoryPath string, form models.UploadImageForm) (models.ImageData, error) {
	if form.FileName == "" || form.FileExtension == "" {
		return models.ImageData{}, errors.New("no file name provided")
	} else if repositoryPath == "" {
		var ok bool
		var err error

		wg := sync.WaitGroup{}

		wg.Add(1)

		go func(w *sync.WaitGroup) {

			ok = C.uploadClient.CheckUserHasAMediaRepository(form.UserID)

			w.Done()

		}(&wg)

		wg.Wait()

		if !ok {
			repositoryPath, err = C.uploadClient.MakeNewMediaRepositoryForUser(form.UserID)
			if err != nil {
				return models.ImageData{}, err
			}
		}

		data, err := C.uploadClient.UploadProfileImage(repositoryPath, form)
		if err != nil {
			return models.ImageData{}, err
		}

		return data, nil
	}

	data, err := C.uploadClient.UploadProfileImage(repositoryPath, form)
	if err != nil {
		return models.ImageData{}, err
	}

	return data, nil
}

func (C *uploadsClient) GetProfilePicture(repositoryPath, fileName, fileExtension string, userId uuid.UUID) (models.GetImage, error) {

	if userId == uuid.Nil || repositoryPath == "" || fileName == "" || fileExtension == "" {
		return models.GetImage{}, errors.New("no parameters provide")
	}

	dataChan := make(chan models.GetImage)
	errchan := make(chan error)

	go func() {

		img, err := C.uploadClient.GetProfileImage(repositoryPath, fileName, fileExtension, userId)
		if err != nil {
			errchan <- err
			dataChan <- models.GetImage{}
			return
		}

		imageData := models.GetImage{
			FileName:    fileName + "." + fileExtension,
			ImageBuffer: img,
		}

		errchan <- nil
		dataChan <- imageData

		return
	}()

	err := <-errchan
	if err != nil {
		return models.GetImage{}, nil
	}

	imageData := <-dataChan

	return imageData, nil
}

func (C *uploadsClient) CreateCollection(userId uuid.UUID, collectionName, repositoryPath string) (models.CollectionData, error) {
	if collectionName == "" || userId == uuid.Nil {
		return models.CollectionData{}, errors.New("no collection name provide")
	}

	collectionDataChan := make(chan models.CollectionData)
	errChan := make(chan error)

	go func(channel chan<- models.CollectionData, errorChan chan<- error) {

		var err error
		if repositoryPath == "" {

			exists := C.uploadClient.CheckUserHasAMediaRepository(userId)

			if !exists {

				repositoryData, err := C.uploadClient.MakeNewMediaRepositoryForUser(userId)
				if err != nil {
					errChan <- err
					return
				}

				form := UserDTO.CreateCollectionForm{
					CollectionName:     collectionName,
					UserRepositoryPath: repositoryData,
				}

				collectionData, err := C.uploadClient.CreateCollection(form)
				if err != nil {
					errChan <- err
					return
				}

				channel <- collectionData
				return
			}

			repositoryPath, err = C.GetUserRepositoryPath(userId)
			if err != nil {
				errChan <- err
				return
			}

		}

		form := UserDTO.CreateCollectionForm{
			CollectionName:     collectionName,
			UserRepositoryPath: repositoryPath,
		}

		collectionData, err := C.uploadClient.CreateCollection(form)
		if err != nil {
			errChan <- err
			return
		}

		channel <- collectionData

		return
	}(collectionDataChan, errChan)

	err := <-errChan
	if err != nil {
		return models.CollectionData{}, err
	}

	collectionData := <-collectionDataChan

	return collectionData, nil
}

func (C *uploadsClient) UploadImageIntoCollection(UserRepository, collectionName string, form models.UploadImageForm) (models.ImageData, error) {
	if form.FileName == "" || form.FileExtension == "" || collectionName == "" || UserRepository == "" {
		return models.ImageData{}, errors.New("no parameters provide")
	}

	Data, err := C.uploadClient.InsertImageOnCollection(UserRepository, collectionName, form)
	if err != nil {
		return models.ImageData{}, err
	}

	return Data, nil
}

func (C *uploadsClient) GetImagesOnCollection(userRepository, collectionName string, forms []models.GetImageForm) ([]models.GetImage, error) {
	if userRepository == "" || len(forms) == 0 {
		return nil, errors.New("no parameters provide")
	}

	Images, err := C.uploadClient.GetAllMediaFromCollection(userRepository, collectionName, forms)
	if err != nil {
		return nil, err
	}

	return Images, nil
}

func (C *uploadsClient) UpdateImageOnCollection(collectionPath, fileName, fileExtension string, form models.UploadImageForm) (models.ImageData, error) {
	if collectionPath == "" || fileName == "" || fileExtension == "" {
		return models.ImageData{}, errors.New("no parameters provide")
	}

	okChan := make(chan bool)

	go func(ok chan<- bool) {

		exits := C.uploadClient.CheckIfImageIsOnCollection(request.UserRepositoryPath, request.CollectionName, request.FileName, request.FileExtension)

		ok <- exits

	}(okChan)

	ok := <-okChan
	if !ok {
		close(okChan)
		return models.ImageData{}, errors.New("image is not on collection")
	}

	close(okChan)

	ImageData, err := C.uploadClient.UpdateImageOnCollection(request, form)
	if err != nil {
		return models.ImageData{}, err
	}

	return ImageData, nil

}
