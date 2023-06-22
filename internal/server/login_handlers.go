package server

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/ArtiomO/oauth/internal/models"
	"github.com/ArtiomO/oauth/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var cacheExpiration = time.Duration(time.Duration.Seconds(300))

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (s *Server) GetLoginHandler(c *gin.Context) {

	var loginFormIn models.LoginFormIn

	reqId := randStringRunes(10)

	if c.ShouldBind(&loginFormIn) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request."})
		return
	}

	loginFormIn.RedirectUri, _ = url.QueryUnescape(loginFormIn.RedirectUri)
	loginFormIn.ClientId, _ = url.QueryUnescape(loginFormIn.ClientId)
	loginFormIn.State, _ = url.QueryUnescape(loginFormIn.State)
	loginFormIn.CodeChallenge, _ = url.QueryUnescape(loginFormIn.CodeChallenge)
	client, err := s.Clients.GetClient(loginFormIn.ClientId)

	if errors.Is(err, repository.ErrClientDoesntExists) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown client."})
		return
	}

	if client.RedirectURI != loginFormIn.RedirectUri {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid redirect URI."})
		return
	}

	cacheKey := fmt.Sprintf("oauth_request_%s", reqId)
	ok, err := s.Cache.SetCacheKey(c.Request.Context(), cacheKey, loginFormIn.String(), cacheExpiration)

	if !ok {
		panic(err)
	}

	c.HTML(http.StatusOK, "login.tmpl", gin.H{"requestId": reqId, "staticUri": os.Getenv("STATIC_URI")})
}

// @BasePath /api/v1

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /example/helloworld [get]
func (s *Server) PostLoginHandler(c *gin.Context) {

	var loginIn models.LoginIn

	if c.ShouldBind(&loginIn) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request."})
		return
	}

	cacheKeyReq := fmt.Sprintf("oauth_request_%s", loginIn.ReqId)
	requestStr, err := s.Cache.GetCacheKey(c.Request.Context(), cacheKeyReq)

	if err == redis.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request."})
		return
	} else if err != nil {
		panic(err)
	}

	loginInReq := models.LoginInFromString(requestStr)

	if loginIn.Email == "vasya@vasya.com" || loginIn.Password == "123" {

		_, err := s.Clients.GetClient(loginInReq.ClientId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		code := randStringRunes(8)
		keyCode := fmt.Sprintf("oauth_code_%s", code)
		ok, err := s.Cache.SetCacheKey(c.Request.Context(), keyCode, requestStr, cacheExpiration)
		if !ok {
			panic(err)
		}
		redirect := fmt.Sprintf("https://localhost:3000/api/oauthcallback?code=%s&state=%s", code, loginInReq.State)
		c.Redirect(http.StatusFound, redirect)
		s.Cache.DelCacheKey(c.Request.Context(), cacheKeyReq)
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid creds."})
		s.Cache.DelCacheKey(c.Request.Context(), cacheKeyReq)
		return
	}

}
