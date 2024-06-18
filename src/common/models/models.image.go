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
	UserID             uuid.UUID
	FileName           string
	ImageData          bytes.Buffer
	UserRepositoryPath string
}

type ImageData struct {
	UserId              uuid.UUID
	Image_id            uuid.UUID
	UserMediaRepository string
	FileName            string
	ImagePath           string
}

type GetImage struct {
	FileName    string
	ImageBuffer bytes.Buffer
}

type Collection struct {
	ID             uuid.UUID `json:"id"`
	UserId         uuid.UUID `json:"user_id"`
	CollectionName string    `json:"collection_name"`
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

type GetImageData struct {
	FilePath string
	Data     bytes.Buffer
}
