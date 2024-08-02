package Uploads

import (
	"errors"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"shopperia/src/common/models"
)

func (US *UploadService) MakeNewPostsDir(userRepository string) (models.CollectionData, error) {
	if userRepository == "" {
		return models.CollectionData{}, errors.New("userRepository is required")
	}

	path, err := US.MakeNewDirectory(userRepository, "posts")
	if err != nil {
		return models.CollectionData{}, err
	}

	dirID := uuid.New()

	postsDirInfo := models.CollectionData{
		CollectionId:   dirID,
		CollectionName: "posts",
		CollectionPath: path,
		UserRepository: userRepository,
	}

	return postsDirInfo, nil

}

func (US *UploadService) CheckIfUserHasPostsDir(userRepository string) (bool, error) {
	if userRepository == "" {
		return false, errors.New("userRepository is required")
	}

	pathToSearch := filepath.Join(userRepository, "posts")

	_, err := os.Stat(pathToSearch)
	if !os.IsExist(err) {
		return false, nil
	}

	return true, nil
}

func (US *UploadService) NewPost(postsDir, postName string) (string, error) {
	if postsDir == "" || postName == "" {
		return "", errors.New("userRepository, postsDir, postName is required")
	}

	postDirPath, err := US.MakeNewDirectory(postsDir, postName)
	if err != nil {
		return "", err
	}

	return postDirPath, nil

}

func (US *UploadService) UploadImagesOnPost(postDirectory string, images []models.UploadImageForm) ([]models.ImageData, error) {
	if postDirectory == "" || len(images) == 0 {
		return nil, errors.New("parameters required")
	}

	ImgsData, err := US.InsertMultipleImagesOnCollection(postDirectory, images)
	if err != nil {
		return nil, err
	}

	return ImgsData, nil
}

func (US *UploadService) UpdatePostImage(postDir, fileName, fileExtension string, NewImage models.UploadImageForm) (models.ImageData, error) {
	if postDir == "" {
		return models.ImageData{}, errors.New("postDir is required")
	}

	imgData, err := US.UpdateImageOnCollection(postDir, fileName, fileExtension, NewImage)
	if err != nil {
		return models.ImageData{}, err
	}

	return imgData, nil
}
