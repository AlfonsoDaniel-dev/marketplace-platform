package user_domain

import (
	"errors"
	"shopperia/src/common/models"
)

func (u *UserDomain) GetProfilePicture(email string) (models.GetImage, error) {
	if email == "" {
		return models.GetImage{}, errors.New("no email provide")
	}

	ProfilePicData, err := u.OutputInterface.PsqlGetUserProfilePictureData(email)
	if err != nil {
		return models.GetImage{}, err
	}

	form := models.GetImageForm{
		FileName:      ProfilePicData.FileName,
		FileExtension: ProfilePicData.FileExtension,
	}

	image, err := u.UploadsInterface.GetMedia(ProfilePicData.UserMediaRepository, form)
	if err != nil {
		return models.GetImage{}, err
	}

	Data := models.GetImage{
		FileName:    ProfilePicData.FileName + "." + ProfilePicData.FileExtension,
		ImageBuffer: image,
	}

	return Data, nil
}

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
