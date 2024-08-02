package Uploads

import (
	"bytes"
	"errors"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"shopperia/src/common/models"
	"strings"
)

type UploadService struct {
	OriginPath string
}

type file struct {
	CompleteName string
	Name         string
	Extension    string
}

type attemptchangeCollectionName struct {
	UserID            uuid.UUID
	OldPath           string
	NewPath           string
	NewCollectionName string
	Status            error
	Files             []file
	FilesData         []bytes.Buffer
	NewCollectionPath string
	ImagesUpdated     []models.ImageData
	Done              chan struct{}
}

type deleteRequest struct {
	ResourcePath string
	Status       error
	Done         chan struct{}
}

type uploadImageSingleAttempt struct {
	Data          models.ImageData
	Image         models.UploadImageForm
	DirectoryPath string
	Status        error
	Done          chan struct{}
}

func (US *UploadService) createFile(AbsoluteFilePath string) (*os.File, error) {
	if AbsoluteFilePath == "" {
		return &os.File{}, errors.New("name is empty")
	}

	NewFile, err := os.Create(AbsoluteFilePath)
	if err != nil {
		return &os.File{}, err
	}

	return NewFile, nil
}

func (US *UploadService) countFilesOnDirectory(oldPath string, count chan<- int, errChan chan<- error) {

	var numberOfFiles int

	err := filepath.Walk(oldPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			numberOfFiles++
		}
		return nil
	})

	count <- numberOfFiles
	errChan <- err
}

func (US *UploadService) upload(repositoryPath, fileName, fileExtension string, image bytes.Buffer, userId uuid.UUID) (models.ImageData, error) {
	imageId := uuid.New()

	fileName = strings.ReplaceAll(fileName, " ", "_")

	name := imageId.String() + "_" + fileName + "." + fileExtension

	dest := filepath.Join(US.OriginPath, repositoryPath, name)

	newFile, err := US.createFile(dest)
	if err != nil {
		return models.ImageData{}, err
	}

	imageBytes := image.Bytes()

	_, err = newFile.Write(imageBytes)
	if err != nil {
		newFile.Close()
		return models.ImageData{}, err
	}

	newFile.Close()

	imageData := models.ImageData{
		UserId:              userId,
		Image_id:            imageId,
		UserMediaRepository: repositoryPath,
		FileName:            fileName,
		FileExtension:       fileExtension,
		ImagePath:           dest,
	}

	return imageData, nil
}

func (US *UploadService) delete(RelativeResourcePath string) error {
	if RelativeResourcePath == "" {
		return errors.New("resource path is empty")
	}

	completePath := filepath.Join(US.OriginPath, RelativeResourcePath)

	err := os.Remove(completePath)
	if err != nil {
		return err
	}

	return nil
}

func (US *UploadService) deleteDirectory(directoryPath string) error {
	if directoryPath == "" {
		return errors.New("no directory path provide")
	}

	resourcePath := filepath.Join(directoryPath)
	err := os.RemoveAll(resourcePath)
	if err != nil {
		return err
	}

	return nil
}

func (US *UploadService) MakeNewDirectory(fatherPath, NewDirName string) (string, error) {
	if NewDirName != "" {
		return "", errors.New("new directory name is empty")
	} else if fatherPath == "" {
		path := US.OriginPath + "/" + NewDirName
		err := os.Mkdir(path, 0755)
		if err != nil {
			return "", err
		}

		return path, nil
	}

	path := fatherPath + "/" + NewDirName
	err := os.Mkdir(path, 0755)
	if err != nil {
		return "", err
	}

	return path, nil
}

func (US *UploadService) uploadWorker(requestChan chan *uploadImageSingleAttempt) {
	for {
		select {
		case request := <-requestChan:
			imageData, err := US.upload(request.DirectoryPath, request.Image.FileName, request.Image.FileExtension, request.Image.ImageData, request.Image.UserID)
			request.Data = imageData
			request.Status = err

			requestChan <- request
			request.Done <- struct{}{}
		}
	}
}

func (US *UploadService) deleteWorker(numAttemps int, requestChan chan *deleteRequest) {

	var i int
	for {

		if i == numAttemps {
			return
		}

		select {
		case req := <-requestChan:
			err := US.delete(req.ResourcePath)
			req.Status = err
			req.Done <- struct{}{}
			requestChan <- req
			i++
		}
	}
}

func getEntryPoint() string {
	entryPoint := os.Getenv("APP_ENTRY_POINT")
	if entryPoint == "" {
		return ""
	}

	return entryPoint
}

func NewUploadService(fatherRepositoryName string) *UploadService {

	originPath := getEntryPoint() + "/" + fatherRepositoryName
	return &UploadService{
		OriginPath: originPath,
	}
}
