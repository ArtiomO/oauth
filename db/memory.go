package db

import "errors"

type Client struct {
	ClientId     string
	ClientSecret string
	RedirectURI  string
}

type Requests map[string]Client

func GetClient(arr []Client, id string) (*Client, error) {
	for _, cl := range arr {
		if cl.ClientId == id {
			return &cl, nil
		}
	}
	return nil, errors.New("Client not registered")
}
