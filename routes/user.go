package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/riad-safowan/JWT-GO-MongoDB/controllers"
	"github.com/riad-safowan/JWT-GO-MongoDB/middleware"
)

func User(r *gin.Engine) {
	r.GET("/user/refresh_token", controllers.RefreshToken())
	r.Use(middleware.Authenticate())
	r.GET("/users", controllers.GetUsers())
	r.GET("/user/:user_id", controllers.GetUser())
}
