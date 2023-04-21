package db

import (
	"encoding/json"
	"errors"
)

type Client struct {
	ClientId     string
	ClientSecret string
	RedirectURI  string
}

type ClientRequest struct {
	RedirectURI         string
	State               string
	ClientId            string
	CodeChallenge       string
	CodeChallengeMethod string
}

func (c ClientRequest) String() string {
	out, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return string(out)
}

func (c Client) String() string {
	out, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return string(out)
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
