package db

import "errors"

type Character struct {
	Id  int `json:"-"`
	Int int `json:"int" binding:"required"`
	Cha int `json:"cha" binding:"required"`
	Str int `json:"str" binding:"required"`
	Wis int `json:"wis" binding:"required"`
	Dex int `json:"dex" binding:"required"`
	Con int `json:"con" binding:"required"`
}

type Characters struct {
	List      []Character
	IdCounter int
}

func (c *Characters) Add(char Character) {
	c.IdCounter += 1
	char.Id = c.IdCounter
	c.List = append(c.List, char)
}

type Client struct {
	ClientId     string
	ClientSecret string
	RedirectURI  string
}

type Requests map[string]Client

func FilterId(arr []Character, id int) (*Character, error) {
	for _, char := range arr {
		if char.Id == id {
			return &char, nil
		}
	}
	return nil, errors.New("Not found")
}

func GetClient(arr []Client, id string) (*Client, error) {
	for _, cl := range arr {
		if cl.ClientId == id {
			return &cl, nil
		}
	}
	return nil, errors.New("Client not registered")
}
