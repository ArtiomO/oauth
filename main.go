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

		if responseType == "token" {
			token := auth.GenerateJWT(
				auth.Header{Alg: "SHA256", Typ: "JWT"},
				auth.Payload{Username: "test_user", Exp: 123123123},
				"test",
			)
			token = url.QueryEscape(token)
			redirect := fmt.Sprintf("http://localhost:3000/api/oauthcallback#access_token=%s&token_type=Bearer", token)
			c.Redirect(http.StatusFound, redirect)
			return
		}

		if responseType == "code" {
			state := c.Query("state")
			code := randStringRunes(8)
			codes[code] = db.CodeClient{
				ClientId: clientId,
			}
			code = url.QueryEscape(code)
			redirect := fmt.Sprintf("http://localhost:3000/api/oauthcallback#code=%s&state=%s", code, state)
			c.Redirect(http.StatusFound, redirect)
			return
		}

	})
	r.GET("/api/v1/login", func(c *gin.Context) {

		reqId := randStringRunes(10)
		clientId, _ := url.QueryUnescape(c.Query("client_id"))
		redirectUri, _ := url.QueryUnescape(c.Query("redirect_uri"))

		requestsContext[reqId] = db.Client{
			ClientId:    clientId,
			RedirectURI: redirectUri,
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
			redirectUri := fmt.Sprintf("http://localhost:8090/api/v1/authorize?client_id=%s&response_type=code&redirect_uri=%s",
				url.QueryEscape(clientent.ClientId),
				url.QueryEscape(clientent.RedirectURI),
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

		login := c.PostForm("login")
		password := c.PostForm("password")
		requestId := c.PostForm("reqid")

		clientent := requestsContext[requestId]

		if login == "vasya" || password == "123" {
			redirectUri := fmt.Sprintf("http://localhost:8090/api/v1/authorize?client_id=%s&response_type=code&redirect_uri=%s",
				url.QueryEscape(clientent.ClientId),
				url.QueryEscape(clientent.RedirectURI),
			)
			fmt.Println(redirectUri)
			c.Redirect(http.StatusFound, redirectUri)
			return
		} else {

			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid creds."})
			return
		}

	})

	if err := r.Run(":8090"); err != nil {
		log.Printf("Error: %v", err)
	}
}
