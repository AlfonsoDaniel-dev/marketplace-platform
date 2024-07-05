package Uploads

import (
	"bytes"
	"errors"
	"os"
	"shopperia/src/common/models"
)

func (US *UploadService) GetMedia(repositoryPath, fileName, fileExtension string) (bytes.Buffer, error) {
	if fileName == "" || repositoryPath == "" {
		return bytes.Buffer{}, errors.New("No path or name provided")
	}

	completePath := US.OriginPath + "/" + repositoryPath + "/" + fileName + "." + fileExtension
	image, err := os.ReadFile(completePath)
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

func (US *UploadService) GetMultipleMediaResources(repositoryPath string, filenames []models.GetImageForm) ([]models.GetImage, error) {
	if repositoryPath == "" || len(filenames) == 0 {
		return nil, errors.New("no path to search")
	}

	var images []models.GetImage = make([]models.GetImage, len(filenames))
	for _, filename := range filenames {

		img, err := US.GetMedia(repositoryPath, filename.FileName, filename.FileExtension)
		if err != nil {
			return nil, err
		}

		name := filename.FileName + "." + filename.FileExtension

		image := models.GetImage{
			FileName:    name,
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
