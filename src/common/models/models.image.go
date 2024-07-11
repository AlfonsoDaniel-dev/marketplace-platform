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
	FileName    string
	ImageBuffer bytes.Buffer
}

type Collection struct {
	ID             uuid.UUID `json:"id"`
	UserId         uuid.UUID `json:"user_id"`
	CollectionName string    `json:"collection_name"`
	Cover          Image     `json:"cover_photo"`
	Description    string    `json:"description"`
	Path           string    `json:"path"`
	Content        []Image   `json:"content"`
}

type CreateCollectionForm struct {
	UserID             uuid.UUID `json:"user_id"`
	CollectionName     string    `json:"collection_name"`
	UserRepositoryPath string    `json:"user_repository_path"`
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
