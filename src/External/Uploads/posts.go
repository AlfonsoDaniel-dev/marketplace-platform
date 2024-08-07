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

func (US *UploadService) CheckIfPostExists(PostDir string) bool {
	if PostDir == "" {
		return false
	}

	okChan := make(chan bool)

	go func(condition chan<- bool) {
		_, err := os.Stat(PostDir)
		if !os.IsExist(err) {
			condition <- false
		}

		condition <- true

	}(okChan)

	ok := <-okChan

	return ok
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

func (US *UploadService) UpdatePostName(postsDir, oldPostName, NewPostName string, userId uuid.UUID) (string, string, error) {

	NewPostPath, NewPostName, err := US.UpdateCollectionName(postsDir, oldPostName, NewPostName, userId)

	return NewPostPath, NewPostName, err
}

func (US *UploadService) DeletePost() {}

func (US *UploadService) DeleteImageOnPost(filePath string) error {
	if filePath == "" {
		return errors.New("filepath is required")
	}

	reqChan := make(chan *deleteRequest)

	go US.deleteWorker(1, reqChan)

	req := &deleteRequest{
		IsDirectory:  false,
		ResourcePath: filePath,
		Done:         make(chan struct{}),
	}

	reqChan <- req

	<-req.Done
	close(reqChan)

	return req.Status
}
