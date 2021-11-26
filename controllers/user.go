package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/riad-safowan/JWT-GO-MongoDB/database"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func GetUser()
func GetUsers()
