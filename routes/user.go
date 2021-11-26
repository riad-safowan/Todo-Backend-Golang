package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/riad-safowan/JWT-GO-MongoDB/controllers"
	"github.com/riad-safowan/JWT-GO-MongoDB/middleware"
)

func User(r gin.Engine) {
	r.Use(middleware.Authenticate())
	r.GET("/users", controllers.GetUser())
	r.GET("/users/:user_id", controllers.GetUsers())
}
