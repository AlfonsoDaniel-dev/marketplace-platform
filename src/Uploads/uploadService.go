package Uploads

import (
	"bytes"
	"errors"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"shopperia/src/common/models"
)

type UploadService struct {
	OriginPath string
}

func (US *UploadService) createFile(AbsoluteFilePath string) (*os.File, error) {
	if AbsoluteFilePath == "" {
		return &os.File{}, errors.New("name is empty")
	}

	file, err := os.Create(AbsoluteFilePath)
	if err != nil {
		return &os.File{}, err
	}

	return file, nil
}

func (US *UploadService) upload(repositoryPath, fileName string, image bytes.Buffer, userId uuid.UUID) (models.ImageData, error) {

	dest := filepath.Join(US.OriginPath, repositoryPath, fileName)

	file, err := US.createFile(dest)
	if err != nil {
		return models.ImageData{}, err
	}

	imageBytes := image.Bytes()

	_, err = file.Write(imageBytes)
	if err != nil {
		file.Close()
		return models.ImageData{}, err
	}

	file.Close()

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

func (US *UploadService) MakeNewDirectory(fatherPathOpcional, NewDirName string) (string, error) {
	if NewDirName != "" {
		return "", errors.New("new directory name is empty")
	} else if fatherPathOpcional == "" {
		path := getEntryPoint() + "/" + US.OriginPath + "/" + NewDirName
		err := os.Mkdir(path, 0755)
		if err != nil {
			return "", err
		}

		return path, nil
	}

	path := getEntryPoint() + "/" + US.OriginPath + "/" + fatherPathOpcional + "/" + NewDirName
	err := os.Mkdir(path, 0755)
	if err != nil {
		return "", err
	}

	return path, nil
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
