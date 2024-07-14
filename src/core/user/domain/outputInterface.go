package user_domain

import (
	"bytes"
	"github.com/google/uuid"
	models "shopperia/src/common/models"
	UserDTO "shopperia/src/core/user/domain/DTO"
)

type OutputInterface interface {
	// User getters
	PsqlCreateUserWithOutAddress(user models.User) error
	PsqlGetUserIdByEmail(email string) (uuid.UUID, error)
	PsqlGetUserNameByEmail(email string) (string, error)
	PsqlInsertAddressData(address models.Address) error

	// Login methods
	PsqlGetHashPassword(email string) ([]byte, error)
	PsqlVerifyEmailExists(email string) (string, error)
	PsqlCheckTwoStepsVerificationIsTrue(email string) (bool, error)
	PsqlInsertTsvCode(email, code string) (string, error)
	PsqlCleanAccessToken(email string) error
	PsqlGetAccessToken(email string) (string, error)

	//images
	PsqlInsertRepositoryPathOnUser(userId uuid.UUID, repositoryPath string) error
	PsqlGetUserRepositoryPath(userId uuid.UUID) (string, error)
	PsqlInsertImageData(imageID, userId uuid.UUID, userRepositoryPath, filename, fileextension, filepath string) error
	PsqlGetUserProfilePictureData(email string) (models.ImageData, error)
	PsqlCreateCollection(CollectionPath string, form UserDTO.DbCreateCollection) error

	// Update
	PsqlChangeUserName(newUserNAme, email string) error
	PsqlChangeUserLastName(newLastName, email string) error
	PsqlChangeUserFirstName(newFirstName, email string) error
	PsqlChangeUserEmail(newEmail, userId string) error
	PsqlChangeUserPassword(newPassword, email, oldPassword string) error
	PsqlChangeUserTsvStatus(email string, value bool) error
}

type EmailInterface interface {
	SendWelcomeEmail(emailContent models.WelcomeEmail, email models.EmailDto) error
	SendLoginConfirmationEmail(emailContent models.SendTSVLoginEmail, email, name string) error
	SendPasswordChangeConfirmationEmail(content models.PasswordChangeEmail, destEmail, destName string) error
	SendTsvChangeConfirmation(content models.TsvChangeEmail, DestEmail, DestName string) error
}

type UploadsInterface interface {
	UploadMedia(userRepository string, image models.UploadImageForm) (models.ImageData, error)
	CreateCollection(userId uuid.UUID, collectionName, repositoryPath string) (models.CollectionData, error)
	GetMedia(repositoryPath string, form models.GetImageForm) (bytes.Buffer, error)
	GetProfilePicture(repositoryPath, fileName, fileExtension string, userId uuid.UUID) (models.GetImage, error)
}
