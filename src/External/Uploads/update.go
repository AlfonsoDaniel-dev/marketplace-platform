package Uploads

import (
	"errors"
	"shopperia/src/common/models"
)

func (US *UploadService) UpdateProfileImage(repositoryPath, oldFileName, oldFileExtension string, form models.UploadImageForm) (models.ImageData, error) {
	if repositoryPath == "" || form.FileName == "" || form.FileExtension == "" {
		return models.ImageData{}, errors.New("no parameters provide")
	}

	attempt := updateImageAttempt{
		OldFileName:      oldFileName,
		OldFileExtension: oldFileExtension,
		Repository:       repositoryPath,
		Form:             form,
		Done:             make(chan struct{}),
	}

	attemptChan := make(chan updateImageAttempt)
	go US.updateImageWorker(1, attemptChan)

	<-attempt.Done
	if attempt.Status != nil {
		return models.ImageData{}, attempt.Status
	}

	return attempt.NewImageData, nil
}
