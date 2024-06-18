package Uploads

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"shopperia/src/common/models"
	"time"
)

type UploadService struct {
	OriginPath string
}

type userMediaRepository struct {
	path string
}

func (US *UploadService) checkUserHasAMediaRepository(userId uuid.UUID) bool {
	repoPath := US.OriginPath + "User" + userId.String() + "_" + "repository"
	_, err := os.Stat(repoPath)
	if os.IsExist(err) {
		return true
	}

	return false
}

func (US *UploadService) checkMediaRepository(repositoryPath, collectionNameOpcional string) bool {
	if collectionNameOpcional != "" {
		path := filepath.Join(US.OriginPath, repositoryPath, collectionNameOpcional)

		_, err := os.Stat(path)
		if !os.IsExist(err) {
			return false
		}

		return true
	}

	path := filepath.Join(US.OriginPath, repositoryPath)

	_, err := os.Stat(path)
	if !os.IsExist(err) {
		return false
	}

	return true
}

func (US *UploadService) MakeNewDirectory(internalPath, NewdirName string) (string, error) {
	if internalPath != "" {
		path := filepath.Join(US.OriginPath, internalPath, NewdirName)
		err := os.Mkdir(path, 0755)
		if err != nil {
			return "", err
		}

		return path, nil
	}

	path := filepath.Join(US.OriginPath, NewdirName)
	err := os.Mkdir(path, 0755)
	if err != nil {
		return "", err
	}

	return path, nil
}

func (US *UploadService) makeNewCollectionRepository(userRepositoryPath, collectionName string) (string, error) {
	collectionPath, collErr := US.MakeNewDirectory(userRepositoryPath, collectionName)
	if collErr != nil {
		return "", collErr
	}

	return collectionPath, nil
}

func (US *UploadService) makeNewMediaRepositoryForUser(userId uuid.UUID) error {
	mediaPath := "/user" + "_" + "user_id"

	mediaPath, err := US.MakeNewDirectory("", mediaPath)
	if err != nil {
		return err
	}

	return nil
}

func (US *UploadService) upload(repositoryPath, fileName string, image bytes.Buffer, userId uuid.UUID) (models.ImageData, error) {

	if repositoryPath == "" {
		ok := US.checkUserHasAMediaRepository(userId)
		if !ok {
			err := US.makeNewMediaRepositoryForUser(userId)
			if err != nil {
				errStr := fmt.Sprintf("Error al crear el repositorio del usuario")
				return models.ImageData{}, errors.New(errStr)
			}

		}

	}

	fileName = fileName + ".jpg"

	dest := filepath.Join(US.OriginPath, repositoryPath, fileName)

	imageBytes := image.Bytes()
	err := os.WriteFile(dest, imageBytes, 0644)
	if err != nil {
		return models.ImageData{}, err
	}

	imageId := uuid.New()

	imageData := models.ImageData{
		UserId:              userId,
		Image_id:            imageId,
		UserMediaRepository: repositoryPath,
		FileName:            fileName,
		ImagePath:           dest,
	}

	return imageData, nil
}

func (US *UploadService) UploadMedia(image models.UploadImageForm) (models.ImageData, error) {
	image.FileName = image.FileName + image.UserID.String() + time.Now().String()

	imageData, err := US.upload(image.UserRepositoryPath, image.FileName, image.ImageData, image.UserID)
	if err != nil {
		return models.ImageData{}, err
	}

	return imageData, nil
}

func (US *UploadService) UploadProfileImage(image models.UploadImageForm) (models.ImageData, error) {
	image.FileName = image.UserID.String() + "_" + image.FileName + image.UserID.String() + "profilePic" + time.Now().String()

	imageData, err := US.upload(image.UserRepositoryPath, image.FileName, image.ImageData, image.UserID)
	if err != nil {
		return models.ImageData{}, err
	}

	return imageData, nil
}

func (US *UploadService) GetMedia(filePath string) (bytes.Buffer, error) {
	if filePath == "" {
		return bytes.Buffer{}, errors.New("No file path provided")
	}
	image, err := os.ReadFile(filePath)
	if err != nil {
		return bytes.Buffer{}, err
	}

	var buf bytes.Buffer

	_, err = buf.Write(image)
	if err != nil {
		return bytes.Buffer{}, err
	}

	return buf, nil
}

func (US *UploadService) SaveMultipleMediaResources(images []models.UploadImageForm) ([]models.ImageData, error) {
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

		ImageData, err := US.upload(image.UserRepositoryPath, image.FileName, image.ImageData, image.UserID)
		if err != nil {
			errStr := fmt.Sprintf("Hubo un error al guardar la imagen %v, ERR: %v", image.FileName, err)
			return nil, errors.New(errStr)
		}

		imageData = append(imageData, ImageData)
	}

	return imageData, nil
}

func (US *UploadService) CreateCollection(form models.CreateCollectionForm) (models.CollectionData, error) {
	if form.UserRepositoryPath == "" || form.CollectionName == "" {
		errStr := fmt.Sprint("Please provide all fields")
		return models.CollectionData{}, errors.New(errStr)
	}

	collection, err := US.makeNewCollectionRepository(form.UserRepositoryPath, form.CollectionName)
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

func (US *UploadService) GetMultipleMediaResources(repositoryPath string, filenames []string) ([]models.GetImage, error) {
	if repositoryPath == "" || len(filenames) == 0 {
		return nil, errors.New("no path to search")
	}

	var images []models.GetImage = make([]models.GetImage, len(filenames))
	for _, filename := range filenames {

		totalPath := filepath.Join(repositoryPath, filename)

		img, err := US.GetMedia(totalPath)
		if err != nil {
			return nil, err
		}

		image := models.GetImage{
			FileName:    filename,
			ImageBuffer: img,
		}

		images = append(images, image)
	}

	return images, nil
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

	imagesData, err := US.SaveMultipleMediaResources(images)
	if err != nil {
		return []models.ImageData{}, err
	}

	return imagesData, nil
}

func (US *UploadService) GetMultipleDataFromPath(Fatherpath string, filenames []string) ([]models.GetImageData, error) {
	var buffers []bytes.Buffer

	ImagesData := []models.GetImageData{}

	for _, name := range filenames {
		path := filepath.Join(Fatherpath, name)
		imageDataBuffer, err := US.GetMedia(path)
		if err != nil {
			return []models.GetImageData{}, err
		}

		buffers = append(buffers, imageDataBuffer)

		Image := models.GetImageData{
			FilePath: path,
			Data:     imageDataBuffer,
		}

		ImagesData = append(ImagesData, Image)
	}

	return ImagesData, nil
}

func NewUploadService(FileRepositoryPath string) *UploadService {
	return &UploadService{
		OriginPath: FileRepositoryPath,
	}
}
