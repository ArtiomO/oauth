package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"
	"os"

	"github.com/ArtiomO/oauth/auth"
	"github.com/ArtiomO/oauth/db"
	"github.com/ArtiomO/oauth/encdec"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var redisExpiration = time.Duration(time.Duration.Seconds(300))

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type TokenIn struct {
	Code         string `form:"code"`
	GrantType    string `form:"grant_type"`
	RedirectUri  string `form:"redirect_uri"`
	CodeVerifier string `form:"code_verifier"`
}

type LoginIn struct {
	Email    string `form:"email"`
	Password string `form:"password"`
	ReqId    string `form:"reqid"`
}

type LoginFormIn struct {
	ClientId            string `form:"client_id"`
	RedirectUri         string `form:"redirect_uri"`
	State               string `form:"state"`
	CodeChallenge       string `form:"code_challenge"`
	CodeChallengeMethod string `form:"code_challenge_method"`
}

func (c LoginFormIn) String() string {
	out, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return string(out)
}

func LoginInFromString(s string) LoginFormIn {
	var loginFormReg LoginFormIn
	err := json.Unmarshal([]byte(s), &loginFormReg)
	if err != nil {
		panic(err)
	}
	return loginFormReg
}

func main() {

	clients := []db.Client{{
		ClientId:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURI:  "http://rssscraperweb:3000/api/oauthcallback",
	}}

	r := gin.Default()
	r.LoadHTMLGlob("./templates/*")
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))
	r.Static(os.Getenv("STATIC_URI"), "./static/")

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	r.GET("/api/v1/login", func(c *gin.Context) {

		var loginFormIn LoginFormIn

		reqId := randStringRunes(10)

		if c.ShouldBind(&loginFormIn) != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request."})
			return
		}

		loginFormIn.RedirectUri, _ = url.QueryUnescape(loginFormIn.RedirectUri)
		loginFormIn.ClientId, _ = url.QueryUnescape(loginFormIn.ClientId)
		loginFormIn.State, _ = url.QueryUnescape(loginFormIn.State)
		loginFormIn.CodeChallenge, _ = url.QueryUnescape(loginFormIn.CodeChallenge)
		client, err := db.GetClient(clients, loginFormIn.ClientId)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown client."})
			return
		}

		if client.RedirectURI != loginFormIn.RedirectUri {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid redirect URI."})
			return
		}

		rdsKey := fmt.Sprintf("oauth_request_%s", reqId)

		err = rdb.Set(c.Request.Context(), rdsKey, loginFormIn.String(), redisExpiration).Err()
		if err != nil {
			panic(err)
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

		rdsKeyReq := fmt.Sprintf("oauth_request_%s", loginIn.ReqId)
		requestStr, err := rdb.Get(c.Request.Context(), rdsKeyReq).Result()

		if err == redis.Nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request."})
			return
		} else if err != nil {
			panic(err)
		}

		loginInReq := LoginInFromString(requestStr)

		if loginIn.Email == "vasya@vasya.com" || loginIn.Password == "123" {

			_, err := db.GetClient(clients, loginInReq.ClientId)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			code := randStringRunes(8)
			rdsKeyCode := fmt.Sprintf("oauth_code_%s", code)
			err = rdb.Set(c.Request.Context(), rdsKeyCode, requestStr, redisExpiration).Err()
			if err != nil {
				panic(err)
			}
			redirect := fmt.Sprintf("http://localhost:3000/api/oauthcallback?code=%s&state=%s", code, loginInReq.State)
			c.Redirect(http.StatusFound, redirect)
			rdb.Del(c.Request.Context(), rdsKeyReq)
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid creds."})
			rdb.Del(c.Request.Context(), rdsKeyReq)
			return
		}

	})

	r.POST("/api/v1/token", func(c *gin.Context) {

		var tokenIn TokenIn

		if c.ShouldBind(&tokenIn) != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request."})
			return
		}

		rdsKey := fmt.Sprintf("oauth_code_%s", tokenIn.Code)

		clientStr, err := rdb.Get(c.Request.Context(), rdsKey).Result()

		if err == redis.Nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid grant."})
			return
		} else if err != nil {
			panic(err)
		}

		rdb.Del(c.Request.Context(), rdsKey)

		storedCodeReq := LoginInFromString(clientStr)

		authHeader := c.GetHeader("Authorization")
		clientId, secret := encdec.GetCreds(authHeader)
		client, err := db.GetClient(clients, clientId)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown client."})
			return
		}

		if client.ClientSecret != secret {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid grant."})
			return
		}

		if tokenIn.GrantType == "authorization_code" {

			if tokenIn.RedirectUri != storedCodeReq.RedirectUri {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid grant."})
				return
			}

			if clientId != storedCodeReq.ClientId {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid grant."})
				return
			}

			calculatedChallenge := encdec.Sha256SumB64(tokenIn.CodeVerifier)

			if calculatedChallenge != storedCodeReq.CodeChallenge {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request."})
				return
			}

			tokenJWT := auth.GenerateJWT(
				auth.Header{Alg: "SHA256", Typ: "JWT"},
				auth.Payload{Username: "test_user", Exp: 123123123},
				"test2",
			)
			c.JSON(http.StatusOK, gin.H{"token": tokenJWT})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid grand type."})
			return
		}

	})

	if err := r.Run("0.0.0.0:8090"); err != nil {
		log.Printf("Error: %v", err)
	}
}
