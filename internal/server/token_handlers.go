package server

import (
	"fmt"
	"net/http"

	"github.com/ArtiomO/oauth/internal/encode"
	"github.com/ArtiomO/oauth/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func (s *Server) PostTokenHandler(c *gin.Context) {

	var tokenIn models.TokenIn

	if c.ShouldBind(&tokenIn) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request."})
		return
	}

	rdsKey := fmt.Sprintf("oauth_code_%s", tokenIn.Code)

	clientStr, err := s.Cache.GetCacheKey(c.Request.Context(), rdsKey)

	if err == redis.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid grant."})
		return
	} else if err != nil {
		panic(err)
	}

	s.Cache.DelCacheKey(c.Request.Context(), rdsKey)

	storedCodeReq := models.LoginInFromString(clientStr)

	authHeader := c.GetHeader("Authorization")
	clientId, secret := encode.GetCreds(authHeader)
	client, err := s.Clients.GetClient(clientId)

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

		calculatedChallenge := encode.Sha256SumHex(tokenIn.CodeVerifier)

		if calculatedChallenge != storedCodeReq.CodeChallenge {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request."})
			return
		}
		tokenJWT := encode.GenerateJWT("test_user")
		c.JSON(http.StatusOK, gin.H{"token": tokenJWT})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid grand type."})
		return
	}

}
