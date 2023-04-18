package db

import "errors"

type Client struct {
	ClientId     string
	ClientSecret string
	RedirectURI  string
}

type CodeClient struct {
	ClientId string
}

type Requests map[string]Client
type Codes map[string]CodeClient

func GetClient(arr []Client, id string) (*Client, error) {
	for _, cl := range arr {
		if cl.ClientId == id {
			return &cl, nil
		}
	}
	return nil, errors.New("Client not registered")
}
