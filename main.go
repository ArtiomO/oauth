package main

import (
	"fmt"
	"github.com/ArtiomO/oauth/auth"
	"github.com/ArtiomO/oauth/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var redisExpiration = time.Duration(time.Duration.Seconds(5))

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

type LoginIn struct {
	Login    string `form:"login"`
	Password string `form:"password"`
	ReqId    string `form:"reqid"`
}

type LoginFormIn struct {
	ClientId    string `form:"client_id"`
	RedirectUri string `form:"redirect_uri"`
	State       string `form:"state"`
}

type AuthorizeIn struct {
	ResponseType string `form:"response_type"`
	CliendId     string `form:"client_id"`
	RedirectUri  string `form:"redirect_uri"`
	State        string `form:"state"`
}

func main() {

	requestsContext := db.Requests{}
	clients := []db.Client{{
		ClientId:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURI:  "http://localhost:3000/api/oauthcallback",
	}}

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	r.GET("/api/v1/authorize", func(c *gin.Context) {

		var authIn AuthorizeIn

		if c.ShouldBind(&authIn) != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request."})
			return
		}
		client, err := db.GetClient(clients, authIn.CliendId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if client.RedirectURI != authIn.RedirectUri {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid redirect URI, incident reported."})
			return
		}
		if authIn.ResponseType == "code" {
			code := randStringRunes(8)
			err := rdb.Set(c.Request.Context(), fmt.Sprintf("oauth_code_%s", code), "", redisExpiration).Err()
			if err != nil {
				panic(err)
			}
			redirect := fmt.Sprintf("http://localhost:3000/api/oauthcallback?code=%s&state=%s", code, authIn.State)
			c.Redirect(http.StatusFound, redirect)
			return
		}

	})
	r.GET("/api/v1/login", func(c *gin.Context) {

		var loginFormIn LoginFormIn

		reqId := randStringRunes(10)

		if c.ShouldBind(&loginFormIn) != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request."})
			return
		}

		requestsContext[reqId] = db.Client{
			ClientId:    loginFormIn.ClientId,
			RedirectURI: loginFormIn.RedirectUri,
			State:       loginFormIn.State,
		}

		c.HTML(http.StatusOK, "login.tmpl", gin.H{"requestId": reqId})
		return
	})

	r.POST("/api/v1/login", func(c *gin.Context) {

		var loginIn LoginIn

		if c.ShouldBind(&loginIn) != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request."})
			return
		}

		clientent, found := requestsContext[loginIn.ReqId]

		if !found {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No auth request."})
			return
		}

		if loginIn.Login == "vasya" || loginIn.Password == "123" {
			redirectUri := fmt.Sprintf("http://localhost:8090/api/v1/authorize?client_id=%s&response_type=code&redirect_uri=%s&state=%s",
				url.QueryEscape(clientent.ClientId),
				url.QueryEscape(clientent.RedirectURI),
				url.QueryEscape(clientent.State),
			)
			c.Redirect(http.StatusFound, redirectUri)
			delete(requestsContext, loginIn.ReqId)
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid creds."})
			delete(requestsContext, loginIn.ReqId)
			return
		}

	})

	r.POST("/api/v1/token", func(c *gin.Context) {

		var token TokenIn

		if c.ShouldBind(&token) != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request."})
			return
		}

		_, err := rdb.Get(c.Request.Context(), fmt.Sprintf("oauth_code_%s", token.Code)).Result()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid grant."})
			return
		}
		rdb.Del(c.Request.Context(), fmt.Sprintf("oauth_code_%s", token.Code))

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
