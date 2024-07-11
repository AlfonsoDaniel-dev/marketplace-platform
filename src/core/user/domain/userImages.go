package user_domain

import (
	"errors"
	"shopperia/src/common/models"
)

/*
func (u *UserDomain) GetProfilePic(email string) (models.GetImage, error) {
	if email == "" {
		return errors.New("no email provide")
	}

	u.UploadsInterface.

} */

func (u *UserDomain) UploadNewImage(imageform models.UploadImageForm) error {
	if imageform.FileName == "" || imageform.FileExtension == "" || imageform.UserEmail == "" {
		return errors.New("no needed fields provide")
	}

	userId, err := u.GetUserIdByEmail(imageform.UserEmail)
	if err != nil {
		return err
	}

	imageform.UserID = userId
	repositoryPath, err := u.OutputInterface.PsqlGetUserRepositoryPath(imageform.UserID)

	imageData, err := u.UploadsInterface.UploadMedia(repositoryPath, imageform)
	if err != nil {
		return err
	}

	err = u.OutputInterface.PsqlInsertImageData(imageData.Image_id, imageData.UserId, imageData.UserMediaRepository, imageData.FileName, imageform.FileExtension, imageData.ImagePath)
	if err != nil {
		return err
	}

	return nil
}
