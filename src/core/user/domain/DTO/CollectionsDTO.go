package UserDTO

import (
	"github.com/google/uuid"
	"shopperia/src/common/models"
)

type CollectionDTO struct {
	Name        string `json:"collection_name"`
	Description string `json:"description"`
}

type GetCollection struct {
	Name        string            `json:"collection_name"`
	Description string            `json:"description"`
	UserEmail   string            `json:"user_email"`
	CoverPhoto  models.GetImage   `json:"cover_photo"`
	Content     []models.GetImage `json:"content"`
}

type CreateCollection struct {
	Email string
	CreateCollectionForm
}

type DbCreateCollection struct {
	Id                 uuid.UUID
	UserId             uuid.UUID
	UserName           string
	CollectionName     string
	Description        string
	UserRepositoryPath string
}

type CreateCollectionForm struct {
	CollectionName     string `json:"collection_name"`
	Description        string `json:"description"`
	UserRepositoryPath string `json:"user_repository_path"`
}
