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
		defer cancel()
		var task models.Task
		var tasks []models.Task
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
		task.ID = primitive.NewObjectID()
		task.Task_id = task.ID.Hex()

		projection := bson.D{
			{"_id", 0},
			{"tasks", 1},
		}

		var response models.Tasks
		err := taskCollection.FindOne(ctx, bson.M{"user_id": user_id}, options.FindOne().SetProjection(projection)).Decode(&response)
		tasks = response.Tasks
		if err != nil {
			if strings.Contains(err.Error(), "no documents") {

			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

		}

		_, insertErr := taskCollection.UpdateOne(ctx, bson.M{"user_id": user_id}, bson.D{{"$push", bson.D{{"tasks", task}}}}, options.Update().SetUpsert(true))
		if insertErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": insertErr.Error()})
			return
		}
		tasks = append(tasks, task)
		c.JSON(http.StatusOK, tasks)
	}
}

func AddTasks() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var newTasks []models.Task
		var user_id = c.GetString("user_id")

		if err := c.BindJSON(&newTasks); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			defer cancel()
			return
		}

		var validatedTask []models.Task
		for _, task := range newTasks {
			validationErr := validate.Struct(task)
			if validationErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
				defer cancel()
				return
			}
			task.ID = primitive.NewObjectID()
			task.Task_id = task.ID.Hex()
			validatedTask = append(validatedTask, task)
		}
		projection := bson.D{
			{"_id", 0},
			{"tasks", 1},
		}

		var response models.Tasks
		err := taskCollection.FindOne(ctx, bson.M{"user_id": user_id}, options.FindOne().SetProjection(projection)).Decode(&response)
		oldTasks := response.Tasks
		if err != nil {
			if strings.Contains(err.Error(), "no documents") {

			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}

		oldTasks = append(oldTasks, validatedTask...)

		_, insertErr := taskCollection.UpdateOne(ctx, bson.M{"user_id": user_id}, bson.D{{"$push", bson.D{{"tasks", bson.D{{"$each", validatedTask}}}}}}, options.Update().SetUpsert(true))
		if insertErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": insertErr.Error()})
			return
		}
		c.JSON(http.StatusOK, oldTasks)
	}
}

func GetTask() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		taskId := c.Param("task_id")
		id, _ := primitive.ObjectIDFromHex(taskId)
		println(taskId)

		projection := bson.D{
			{"_id", 0},
			{"user_id", 0},
			{"tasks", bson.D{
				{"$elemMatch", bson.D{{"_id", id}}},
			}},
		}

		var response models.Tasks
		err := taskCollection.FindOne(ctx, bson.M{"tasks._id": id},
			options.FindOne().SetProjection(projection),
		).Decode(&response)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		task := response.Tasks[0]
		c.JSON(http.StatusOK, task)
	}
}

func GetTasks() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var user_id = c.GetString("user_id")

		projection := bson.D{
			{"_id", 0},
			{"tasks", 1},
		}

		var response models.Tasks
		err := taskCollection.FindOne(ctx, bson.M{"user_id": user_id}, options.FindOne().SetProjection(projection)).Decode(&response)
		tasks := response.Tasks
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, tasks)
	}

}
