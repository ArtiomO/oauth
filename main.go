package main

import (
	"fmt"
	"github.com/ArtiomO/oauth/auth"
	"github.com/ArtiomO/oauth/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
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

	var charsStorage = db.Characters{
		List: []db.Character{},
	}

	requestsContext := db.Requests{}
	charsStorage.Add(db.Character{
		Id:  1,
		Int: 11,
		Cha: 11,
		Str: 11,
		Wis: 11,
		Dex: 11,
		Con: 11,
	})

	clients := []db.Client{{
		ClientId:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURI:  "http://localhost:3000/oauthcallback",
	}}

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))
	r.GET("/character/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, _ := strconv.Atoi(idStr)

		char, err := db.FilterId(charsStorage.List, id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, char)
	})

	r.POST("/character", func(c *gin.Context) {
		var newCharacter db.Character
		if err := c.BindJSON(&newCharacter); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		charsStorage.Add(newCharacter)
		c.IndentedJSON(http.StatusCreated, newCharacter)
	})

	r.GET("/authorize", func(c *gin.Context) {
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
			redirect := fmt.Sprintf("http://localhost:3000/oauthcallback#access_token=%s&token_type=Bearer", token)
			c.Redirect(http.StatusFound, redirect)
			return
		}

	})
	r.GET("/login", func(c *gin.Context) {

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

	r.POST("/login", func(c *gin.Context) {

		login := c.PostForm("login")
		password := c.PostForm("password")
		requestId := c.PostForm("reqid")

		clientent := requestsContext[requestId]

		if login == "vasya" || password == "123" {
			redirectUri := fmt.Sprintf("http://localhost:8080/authorize?client_id=%s&response_type=token&redirect_uri=%s",
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

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
