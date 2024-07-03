package Uploads

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"shopperia/src/common/models"
)

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

/*
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
*/
