package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"github.com/riad-safowan/JWT-GO-MongoDB/routes"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading .env file")
	}
	port := os.Getenv("PORT")

	if port == "" {
		port = "9999"
	}

	// testMongo()

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("api-1", func(c *gin.Context) { c.JSON(200, gin.H{"success": "access granted for api-1"}) })
	router.GET("api-2", func(c *gin.Context) { c.JSON(200, gin.H{"success": "access granted for api-2 "}) })

	routes.Auth(router)
	routes.User(router)

	router.Run("127.0.0.1:" + port)
}

func testMongo() {
	clientOptions := options.Client().
		ApplyURI("mongodb+srv://riadsafowan:12345**asdf@clusterrs.fajlk.mongodb.net/jwt?retryWrites=true&w=majority")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(ctx)

	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		log.Fatal(err)
	}
	database, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(database)
}
