package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/riad-safowan/JWT-GO-MongoDB/controllers"
)

func Auth(r *gin.Engine) {
	r.POST("users/signup", controllers.Signup())
	r.POST("users/login", controllers.Login())
}
