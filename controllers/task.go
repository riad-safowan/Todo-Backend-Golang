package controllers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/riad-safowan/JWT-GO-MongoDB/database"
	"github.com/riad-safowan/JWT-GO-MongoDB/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var taskCollection *mongo.Collection = database.OpenCollection(database.Client, "tasks")

func AddTask() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		var task models.Task
		task.ID = primitive.NewObjectID()
		task.Task_id = task.ID.Hex()
		var tasks models.Tasks
		var user_id = c.GetString("user_id")

		if err := c.BindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			defer cancel()
			return
		}

		validationErr := validate.Struct(task)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			defer cancel()
			return
		}

		err := taskCollection.FindOne(ctx, bson.M{"user_id": user_id}).Decode(&tasks)
		if err != nil {
			if strings.Contains(err.Error(), "no documents") {
				tasks.UserId = user_id
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				defer cancel()
				return
			}
		}

		tasks.Tasks = append(tasks.Tasks, task)
		_, insertErr := taskCollection.UpdateOne(ctx, bson.M{"user_id": user_id}, bson.D{{"$set", bson.D{{"tasks", tasks.Tasks}}}}, options.Update().SetUpsert(true))
		if insertErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": insertErr.Error()})
			defer cancel()
			return
		}
		c.JSON(http.StatusOK, tasks.Tasks)
		defer cancel()
	}
}

func AddTasks() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetTask() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetTasks() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
