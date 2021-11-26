package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/riad-safowan/JWT-GO-MongoDB/controllers"
)

func Auth(r *gin.Engine) {
	r.POST("users/signup", controller.Signup())
	r.POST("users/login", controller.Login())
}
