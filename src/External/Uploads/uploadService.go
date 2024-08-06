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
	IsDirectory  bool
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

type updateImageAttempt struct {
	OldFileName      string
	OldFileExtension string
	Repository       string
	Form             models.UploadImageForm
	NewImageData     models.ImageData
	Status           error
	Done             chan struct{}
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

func (US *UploadService) CheckImageExists(repositoryPath, fileName, fileExtension string) bool {
	if repositoryPath == "" || fileName == "" || fileExtension == "" {
		return false
	}

	completeFileName := filepath.Join(repositoryPath, fileName+"."+fileExtension)

	_, err := os.Stat(completeFileName)
	if !os.IsExist(err) {
		return false
	}

	return true
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

func (US *UploadService) updateImageWorker(numAttempt int, reqChan chan updateImageAttempt) {

	var i int

	for {
		if i == numAttempt {
			break
		}

		select {
		case req := <-reqChan:
			resourcePath := filepath.Join(req.Repository, req.OldFileName+"."+req.OldFileExtension)
			err := US.delete(resourcePath)
			if err != nil {
				req.Status = err
				req.Done <- struct{}{}
			}

			imageData, err := US.upload(req.Repository, req.Form.FileName, req.Form.FileExtension, req.Form.ImageData, req.Form.UserID)
			if err != nil {
				req.Status = err
				req.Done <- struct{}{}
			}

			req.NewImageData = imageData
			i++

			req.Done <- struct{}{}
		}
	}
}

func (US *UploadService) uploadWorker(requestChan chan *uploadImageSingleAttempt) {
	for {
		select {
		case request := <-requestChan:
			imageData, err := US.upload(request.DirectoryPath, request.Image.FileName, request.Image.FileExtension, request.Image.ImageData, request.Image.UserID)
			if err != nil {
				request.Status = err
				request.Done <- struct{}{}
			}
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
			if !req.IsDirectory {

				req.Status = US.delete(req.ResourcePath)

				req.Done <- struct{}{}
				requestChan <- req
				i++
			}

			req.Status = US.deleteDirectory(req.ResourcePath)
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
