package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/riad-safowan/JWT-GO-MongoDB/helpers"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {

		clientToken := c.Request.Header.Get("Authorization")
		if clientToken == "" {
			clientToken = c.Request.Header.Get("token")
		} else if strings.HasPrefix(clientToken, "Bearer ") {
			reqToken := c.Request.Header.Get("Authorization")
			splitToken := strings.Split(reqToken, "Bearer ")
			clientToken = splitToken[1]
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid authorization token"})
			c.Abort()
			return
		}

		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no Authorization header provided"})
			c.Abort()
			return
		}
		// handle access token
		claims, err := helpers.ValidateToken(clientToken)

		if err != "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err})
			c.Abort()
			return
		}

		if claims.Token_type == "access_token" {
			c.Set("email", claims.Email)
			c.Set("first_name", claims.First_name)
			c.Set("last_name", claims.Last_name)
			c.Set("user_id", claims.Uid)
			c.Set("user_type", claims.User_type)
			c.Next()
		} else if claims.Token_type == "refresh_token" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid access_token"})
			c.Abort()
			return
		}

	}
}
