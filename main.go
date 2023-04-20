package main

import (
	"fmt"
	"github.com/ArtiomO/oauth/auth"
	"github.com/ArtiomO/oauth/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
	"net/url"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type TokenIn struct {
	Code        string `json:"code"`
	GrantType   string `json:"grant_type"`
	RedirectUri string `json:"redirect_uri"`
}

func main() {

	requestsContext := db.Requests{}
	clients := []db.Client{{
		ClientId:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURI:  "http://localhost:3000/api/oauthcallback",
	}}

	codes := db.Codes{}

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

	r.GET("/api/v1/authorize", func(c *gin.Context) {
		responseType := c.Query("response_type")
		clientId := c.Query("client_id")
		redirectUri := c.Query("redirect_uri")
		client, err := db.GetClient(clients, clientId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if client.RedirectURI != redirectUri {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid redirect URI, incident reported."})
			return
		}

		if responseType == "code" {
			state := c.Query("state")
			code := randStringRunes(8)
			codes[code] = db.CodeClient{
				ClientId: clientId,
			}
			code = url.QueryEscape(code)
			redirect := fmt.Sprintf("http://localhost:3000/api/oauthcallback?code=%s&state=%s", code, state)
			c.Redirect(http.StatusFound, redirect)
			return
		}

	})
	r.GET("/api/v1/login", func(c *gin.Context) {

		reqId := randStringRunes(10)
		clientId, _ := url.QueryUnescape(c.Query("client_id"))
		redirectUri, _ := url.QueryUnescape(c.Query("redirect_uri"))
		state, _ := url.QueryUnescape(c.Query("state"))

		requestsContext[reqId] = db.Client{
			ClientId:    clientId,
			RedirectURI: redirectUri,
			State:       state,
		}

		fmt.Println(reqId)
		fmt.Println(requestsContext)

		c.HTML(http.StatusOK, "login.tmpl", gin.H{"requestId": reqId})
		return
	})

	r.POST("/api/v1/login", func(c *gin.Context) {

		login := c.PostForm("login")
		password := c.PostForm("password")
		requestId := c.PostForm("reqid")

		clientent := requestsContext[requestId]

		if login == "vasya" || password == "123" {
			redirectUri := fmt.Sprintf("http://localhost:8090/api/v1/authorize?client_id=%s&response_type=code&redirect_uri=%s&state=%s",
				url.QueryEscape(clientent.ClientId),
				url.QueryEscape(clientent.RedirectURI),
				url.QueryEscape(clientent.State),
			)
			fmt.Println(redirectUri)
			c.Redirect(http.StatusFound, redirectUri)
			return
		} else {

			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid creds."})
			return
		}

	})

	r.POST("/api/v1/token", func(c *gin.Context) {

		var token TokenIn

		if c.ShouldBind(&token) != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request."})
			return
		}

		if _, fk := codes[token.Code]; !fk {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid grant."})
			return
		} else {
			delete(codes, token.Code)
		}

		authHeader := c.GetHeader("Authorization")
		clientId, secret := auth.GetCreds(authHeader)
		client, err := db.GetClient(clients, clientId)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid grant."})
			return
		}

		if client.ClientSecret != secret {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid grant."})
			return
		}

		if client.RedirectURI != token.RedirectUri {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid grant."})
			return
		}

		if token.GrantType == "authorization_code" {
			token := auth.GenerateJWT(
				auth.Header{Alg: "SHA256", Typ: "JWT"},
				auth.Payload{Username: "test_user", Exp: 123123123},
				"test",
			)
			c.JSON(http.StatusOK, gin.H{"token": token})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid grand type."})
			return
		}

	})

	if err := r.Run(":8090"); err != nil {
		log.Printf("Error: %v", err)
	}
}
