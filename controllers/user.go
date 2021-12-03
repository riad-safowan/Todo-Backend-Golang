package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/riad-safowan/JWT-GO-MongoDB/database"
	"github.com/riad-safowan/JWT-GO-MongoDB/helpers"
	"github.com/riad-safowan/JWT-GO-MongoDB/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		if err := helpers.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User

		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, user)
	}
}
func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
		if err := helpers.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		cursor, err := userCollection.Find(ctx, bson.M{})
		if err != nil {
			log.Fatal(err)
		}
		var users []bson.M
		if err = cursor.All(ctx, &users); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, users)


		//LOGIC WITH PAGINATION___
		// fmt.Fprintf(c.Writer, c.GetString("user_type"))
		// if err := helpers.CheckUserType(c, "ADMIN"); err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 	return
		// }
		// fmt.Fprintf(c.Writer, "1")
		// var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		// recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))

		// if err != nil || recordPerPage < 1 {
		// 	recordPerPage = 10
		// }
		// page, err1 := strconv.Atoi(c.Query("page"))
		// if err1 != nil || page < 1 {
		// 	page = 1
		// }
		// fmt.Fprintf(c.Writer, "2")
		// startIndex := (page - 1) * recordPerPage
		// startIndex, err = strconv.Atoi(c.Query("startIndex"))

		// matchStage := bson.D{{"$match", bson.D{{}}}}
		// groupStage := bson.D{{"$group", bson.D{{"_id", bson.D{{"_id", "null"}}}, {"total_count", bson.D{{"$sum", 1}}}, {"data", bson.D{{"$push", "$$ROOT"}}}}}}
		// projectStage := bson.D{
		// 	{"$project", bson.D{{"_id", 0}}},
		// 	{"total_count", 1},
		// 	{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, startIndex}}}},
		// }
		// fmt.Fprintf(c.Writer, "3")
		// result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
		// 	matchStage, groupStage, projectStage,
		// })
		// defer cancel()
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "error occuerd while listing user items"})
		// }

		// var allUsers []bson.M
		// if err = result.All(ctx, &allUsers); err != nil {
		// 	log.Fatal(err)
		// }
		// fmt.Fprintf(c.Writer, "4")
		// c.JSON(http.StatusOK, allUsers[0])
	}
}
