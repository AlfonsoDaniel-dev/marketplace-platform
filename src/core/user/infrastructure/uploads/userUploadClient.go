package userUploads

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"path/filepath"
	"shopperia/src/External/Uploads"
	"shopperia/src/common/models"
	UserDTO "shopperia/src/core/user/domain/DTO"
)

type uploadClient interface {
	MakeNewDirectory(fatherDirectory, NewDirName string) (string, error)
	MakeNewMediaRepositoryForUser(userId uuid.UUID) (string, error)

	CheckUserHasAMediaRepository(userId uuid.UUID) bool
	CheckImageExists(repository, fileName, fileExtension string) bool
	GetUserRepositoryPath(userId uuid.UUID) (string, error)

	UploadMedia(repositoryPath string, image models.UploadImageForm) (models.ImageData, error)
	UploadProfileImage(repositoryPath string, image models.UploadImageForm) (models.ImageData, error)
	UpdateProfileImage(repositoryPath, oldFileName, oldFileExtension string, image models.UploadImageForm) (models.ImageData, error)
	UploadMultipleMediaResourcesOnRepository(repositoryPath string, images []models.UploadImageForm) ([]models.ImageData, error)

	GetMedia(repositoryPath, fileName, fileExtension string) (bytes.Buffer, error)
	GetProfileImage(repositoryPath, fileName, fileExtension string, userId uuid.UUID) (bytes.Buffer, error)
}

type collectionsInterface interface {
	CreateCollection(form UserDTO.CreateCollectionForm) (models.CollectionData, error)
	DeleteCollection(collectionPath string) error
	InsertImageOnCollection(collectionPath string, image models.UploadImageForm) (models.ImageData, error)
	InsertMultipleImagesOnCollection(collectionPath string, forms []models.UploadImageForm) ([]models.ImageData, error)
	GetAllMediaFromCollection(repositoryPath, collectionPath string, forms []models.GetImageForm) ([]models.GetImage, error)
	UpdateImageOnCollection(collectionPath, fileName, fileExtension string, form models.UploadImageForm) (models.ImageData, error)
	DeleteImageOnCollection(request models.DeleteOnCollectionRequest) error
	DeleteMultipleImagesOnCollection(requests []models.DeleteOnCollectionRequest) error
	CheckIfImageIsOnCollection(collectionPath, fileName, fileExtension string) bool
	UpdateCollectionName(userRepository, collectionName, newCollectionName string, userID uuid.UUID) (string, string, error)
}

type postsInterface interface {
	MakeNewPostsDir(userRepository string) (models.CollectionData, error)
	CheckIfUserHasPostsDir(userRepository string) (bool, error)
	CheckIfPostExists(PostDir string) bool
	NewPost(postsDir, postName string) (string, error)
	UpdatePostName(postsDir, OldPostName, NewPostName string, userId uuid.UUID) (string, string, error)
	DeleteImageOnPost(filePath string) error
	UploadImagesOnPost(postDirectory string, images []models.UploadImageForm) ([]models.ImageData, error)
	UpdatePostImage(postDir, fileName, fileExtension string, NewImage models.UploadImageForm) (models.ImageData, error)
}

var userPath string = "./src/core/user/infrastructure/uploads/main/repository"

type UploadsClient struct {
	uploadClient
	collectionsInterface
	postsInterface
}

func NewUploadsClient() *UploadsClient {

	client := Uploads.NewUploadService(userPath)

	return &UploadsClient{
		uploadClient: client,
	}
}

/*
func (C *UploadsClient) CreateUserRepository(userId uuid.UUID) (string, error) {

} */

func (C *UploadsClient) UploadMedia(userRepository string, image models.UploadImageForm) (models.ImageData, error) {
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

func (C *UploadsClient) GetMedia(repositoryPath string, form models.GetImageForm) (bytes.Buffer, error) {
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

func (C *UploadsClient) UploadProfilePicture(repositoryPath string, form models.UploadImageForm) (models.ImageData, error) {
	if form.FileName == "" || form.FileExtension == "" {
		return models.ImageData{}, errors.New("no file name provided")
	} else if repositoryPath == "" {
		var err error

		okChan := make(chan bool)

		go func(condition chan<- bool) {

			ok := C.uploadClient.CheckUserHasAMediaRepository(form.UserID)
			condition <- ok

		}(okChan)

		ok := <-okChan

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

func (C *UploadsClient) GetProfileImage(repositoryPath, fileName, fileExtension string, userId uuid.UUID) (models.GetImage, error) {

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

func (C *UploadsClient) UpdateProfilePicture(repositoryPath, oldFileName, oldFileExtension string, NewImage models.UploadImageForm) (models.ImageData, error) {
	if repositoryPath == "" || oldFileName == "" || oldFileExtension == "" {
		return models.ImageData{}, errors.New("no parameters provide")
	}

	okChan := make(chan bool)

	go func(condition chan<- bool) {

		ok := C.uploadClient.CheckImageExists(repositoryPath, oldFileName, oldFileExtension)

		condition <- ok
	}(okChan)

	ok := <-okChan
	if !ok {
		return models.ImageData{}, errors.New("image does not exist")
	}

	NewData, err := C.uploadClient.UpdateProfileImage(repositoryPath, oldFileName, oldFileExtension, NewImage)
	if err != nil {
		return models.ImageData{}, err
	}

	return NewData, nil
}

func (C *UploadsClient) CreateCollection(userId uuid.UUID, collectionName, repositoryPath string) (models.CollectionData, error) {
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

				collectionData, err := C.collectionsInterface.CreateCollection(form)
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

		collectionData, err := C.collectionsInterface.CreateCollection(form)
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

func (C *UploadsClient) UpdateCollectionName(userRespository, collectionName, NewCollectionName string, userID uuid.UUID) (models.CollectionData, error) {
	if userRespository == "" || collectionName == "" || NewCollectionName == "" {
		return models.CollectionData{}, errors.New("no parameters provide")
	}

	collectionPath, collectionName, err := C.collectionsInterface.UpdateCollectionName(userRespository, collectionName, NewCollectionName, userID)
	if err != nil {
		return models.CollectionData{}, err
	}

	data := models.CollectionData{
		CollectionName: collectionName,
		UserRepository: userRespository,
		CollectionPath: collectionPath,
	}

	return data, nil
}

func (C *UploadsClient) UploadImageIntoCollection(collectionPath string, form models.UploadImageForm) (models.ImageData, error) {
	if form.FileName == "" || form.FileExtension == "" || collectionPath == "" {
		return models.ImageData{}, errors.New("no parameters provide")
	}

	Data, err := C.collectionsInterface.InsertImageOnCollection(collectionPath, form)
	if err != nil {
		return models.ImageData{}, err
	}

	return Data, nil
}

func (C *UploadsClient) GetImagesOnCollection(userRepository, collectionName string, forms []models.GetImageForm) ([]models.GetImage, error) {
	if userRepository == "" || len(forms) == 0 {
		return nil, errors.New("no parameters provide")
	}

	Images, err := C.collectionsInterface.GetAllMediaFromCollection(userRepository, collectionName, forms)
	if err != nil {
		return nil, err
	}

	return Images, nil
}

func (C *UploadsClient) UpdateImageOnCollection(collectionPath, fileName, fileExtension string, form models.UploadImageForm) (models.ImageData, error) {
	if collectionPath == "" || fileName == "" || fileExtension == "" {
		return models.ImageData{}, errors.New("no parameters provide")
	}

	okChan := make(chan bool)

	go func(ok chan<- bool) {

		exits := C.collectionsInterface.CheckIfImageIsOnCollection(collectionPath, fileName, fileExtension)

		ok <- exits

	}(okChan)

	ok := <-okChan
	if !ok {
		close(okChan)
		return models.ImageData{}, errors.New("image is not on collection")
	}

	close(okChan)

	ImageData, err := C.collectionsInterface.UpdateImageOnCollection(collectionPath, fileName, fileExtension, form)
	if err != nil {
		return models.ImageData{}, err
	}

	return ImageData, nil

}

func (C *UploadsClient) NewPost(userID uuid.UUID, postsDir, userRepository, postName string) (models.PostData, error) {
	if postName == "" {
		return models.PostData{}, errors.New("new post name is required")
	} else if postsDir == "" || userRepository == "" {

		okChan := make(chan bool)

		go func(condition chan<- bool) {

			ok := C.uploadClient.CheckUserHasAMediaRepository(userID)

			condition <- ok

		}(okChan)

		ok := <-okChan
		close(okChan)
		if !ok {
			var err error

			userRepository, err = C.uploadClient.MakeNewMediaRepositoryForUser(userID)
			if err != nil {
				return models.PostData{}, err
			}

			postsDirData, err := C.postsInterface.MakeNewPostsDir(userRepository)
			if err != nil {
				return models.PostData{}, err
			}

			postPath, err := C.postsInterface.NewPost(postsDirData.CollectionPath, postName)
			if err != nil {
				return models.PostData{}, err
			}

			Data := models.PostData{
				ID:             uuid.New(),
				UserRepository: userRepository,
				UserPostsDir:   postsDirData.CollectionPath,
				Path:           postPath,
			}

			return Data, nil

		}
	}

	var err error

	NewPostPath, err := C.postsInterface.NewPost(postsDir, postName)
	if err != nil {
		return models.PostData{}, err
	}

	userRepository, err = C.uploadClient.GetUserRepositoryPath(userID)
	if err != nil {
		errStr := fmt.Sprintf("no user repository found, error while creating it: %v", err)
		return models.PostData{}, errors.New(errStr)
	}

	Data := models.PostData{
		ID:             uuid.New(),
		UserRepository: userRepository,
		UserPostsDir:   postsDir,
		Path:           NewPostPath,
	}

	return Data, nil
}

func (C *UploadsClient) UpdateImageOnPost(postDir, fileName, fileDescription string, NewImage models.UploadImageForm) (models.ImageData, error) {
	if postDir == "" || fileName == "" || fileDescription == "" {
		return models.ImageData{}, errors.New("no parameters provide")
	}

	ok := C.postsInterface.CheckIfPostExists(postDir)
	if !ok {
		return models.ImageData{}, errors.New("post does not exits")
	}

	Data, err := C.postsInterface.UpdatePostImage(postDir, fileName, fileDescription, NewImage)
	if err != nil {
		return models.ImageData{}, err
	}

	return Data, nil
}

func (C *UploadsClient) DeleteImageOnPost(postDir, fileName, fileExtension string) error {
	if postDir == "" || fileName == "" || fileExtension == "" {
		return errors.New("no parameters provide")
	}

	path := filepath.Join(postDir, fileName+"."+fileExtension)

	err := C.postsInterface.DeleteImageOnPost(path)
	if err != nil {
		return err
	}

	return nil
}

func (C *UploadsClient) UpdatePostName(postsDir, postName, NewPostName string, userId uuid.UUID) (models.PostData, error) {
	if postsDir == "" || postName == "" || NewPostName == "" || userId == uuid.Nil {
		return models.PostData{}, errors.New("no parameters provide")
	}

	NewPostPath, NewPostName, err := C.postsInterface.UpdatePostName(postsDir, postName, NewPostName, userId)
	if err != nil {
		return models.PostData{}, err
	}

	userRepository, err := C.GetUserRepositoryPath(userId)
	if err != nil {
		return models.PostData{}, err
	}

	Data := models.PostData{
		ID:             uuid.UUID{},
		UserRepository: userRepository,
		UserPostsDir:   postsDir,
		Path:           NewPostPath,
	}

	return Data, nil
}
