package models

import (
	"bytes"
	"github.com/google/uuid"
)

type Image struct {
	ID       uuid.UUID `json:"id"`
	UserID   uuid.UUID `json:"user_id"`
	FileName string    `json:"file_name"`
	FilePath string    `json:"filename"`
	Data     []byte    `json:"data"`
}

type UploadImageForm struct {
	ImageID       uuid.UUID
	UserID        uuid.UUID
	UserEmail     string
	FileName      string
	FileExtension string
	ImageData     bytes.Buffer
}

type ImageData struct {
	UserId              uuid.UUID
	Image_id            uuid.UUID
	UserMediaRepository string
	FileName            string
	FileExtension       string
	ImagePath           string
}

type GetImageForm struct {
	FileName      string
	FileExtension string
}

type GetImage struct {
	FileName      string
	FileExtension string
	ImageBuffer   bytes.Buffer
}

type Collection struct {
	ID             uuid.UUID `json:"id"`
	UserId         uuid.UUID `json:"user_id"`
	UserName       string    `json:"user_name"`
	UserEmail      string    `json:"user_email"`
	Name           string    `json:"collection_name"`
	Cover          Image     `json:"cover_photo"`
	Description    string    `json:"description"`
	UserRepository string    `json:"user_repository"`
	Path           string    `json:"path"`
	Content        []Image   `json:"content"`
}

type CollectionPath struct {
	UserRepositoryPath string
	CollectionName     string
}

type CollectionData struct {
	CollectionId   uuid.UUID
	CollectionName string
	UserRepository string
	CollectionPath string
}
