package repository

import (
	"errors"
	"github.com/ArtiomO/oauth/internal/models"
	"sync"
)

type MemoryClientRepository struct {
	clients map[string]models.Client
	lock    *sync.RWMutex
}

func (r MemoryClientRepository) InitClientRepo() *MemoryClientRepository {

	clients := map[string]models.Client{
		"test-client-id": {
			ClientId:     "test-client-id",
			ClientSecret: "test-client-secret",
			RedirectURI:  "http://localhost:3000/api/oauthcallback",
		},
	}
	r.clients = clients
	r.lock = &sync.RWMutex{}
	return &r
}

func (r MemoryClientRepository) GetClient(clientId string) (*models.Client, error) {

	r.lock.Lock()
	defer r.lock.Unlock()

	client, ok := r.clients[clientId]

	if !ok {
		return nil, errors.New("no such client")
	}

	return &client, nil
}
