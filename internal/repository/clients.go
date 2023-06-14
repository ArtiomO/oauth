package repository

import (
	"github.com/ArtiomO/oauth/internal/models"
)

type ClientsRepository interface {
	GetClient(clientId string) (*models.Client, error)
}
