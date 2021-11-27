package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/riad-safowan/JWT-GO-MongoDB/helpers"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientAccessToken := c.Request.Header.Get("access_token")
		if clientAccessToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprint("No Authorization header provided")})
			c.Abort()
			return
		}

		claims , err := helpers.ValidateToken(clientAccessToken)

		if err!= "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error":err})
			c.Abort()
			return 
		}

		c.Set("email", claims.Email)
		c.Set("first_name", claims.First_name)
		c.Set("last_name", claims.Last_name)
		c.Set("user_id", claims.Uid)
		c.Set("user_type", claims.User_type)
		c.Next()
	}
}
