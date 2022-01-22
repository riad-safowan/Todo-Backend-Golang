package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/riad-safowan/JWT-GO-MongoDB/controllers"
	"github.com/riad-safowan/JWT-GO-MongoDB/middleware"
)

func Task(r *gin.Engine) {
r.Use(middleware.Authenticate())
r.GET("/user/task", controllers.GetTask())
r.GET("/user/tasks", controllers.GetTasks())
r.POST("/user/task", controllers.AddTask())
r.POST("/user/tasks", controllers.AddTasks())
}