package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/riad-safowan/JWT-GO-MongoDB/routes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading .env file")
	}
	port := os.Getenv("PORT")

	if port == "" {
		port = "9090"
	}

	// testMongo()

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("api-1", func(c *gin.Context) { c.JSON(200, gin.H{"success": "access granted for api-1"}) })
	router.GET("/ws", websocketConnection())

	routes.Auth(router)
	routes.User(router)
	routes.Task(router)

	router.Run("192.168.31.215:" + port)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func websocketConnection() gin.HandlerFunc {
	return func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			_, _ = fmt.Fprint(c.Writer, "You must use the web socket protocol to connect to this endpoint.", err)
			return
		}

		defer ws.Close()

		for {
			_, p, err := ws.ReadMessage()
			if err != nil {
				println(err)
				break
			}
			msg := string(p)
			println("massage: ", msg)
			for i := 0; i < 5; i++ {
				ws.WriteMessage(websocket.TextMessage, []byte(strconv.Itoa(i)))
				time.Sleep(5 * time.Second)
			}
			println("Riad")
		}

	}
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
