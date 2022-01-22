package helpers

import (
	"context"

	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/riad-safowan/JWT-GO-MongoDB/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email      string
	First_name string
	Last_name  string
	Uid        string
	User_type  string
	Token_type string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func GenerateAllToken(email string, firstName string, lastName string, userType string, uid string) (signedAccessToken string, signedRefreshToken string, err error) {
	accessClaims := &SignedDetails{
		Email:      email,
		First_name: firstName,
		Last_name:  lastName,
		User_type:  userType,
		Uid:        uid,
		Token_type :"access_token",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(100)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		Uid: uid,
		Token_type :"refresh_token",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24*7)).Unix(),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		log.Panic(err)
		return
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		log.Panic(err)
		return
	}

	return accessToken, refreshToken, err

}

func UpdateAllTokens(signedAccessToken string, signedRefreshToken string, userId string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	// var updateObj primitive.D

	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	id, _ := primitive.ObjectIDFromHex(userId)
	_, err := userCollection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.D{
			{"$set", bson.D{{"access_token", signedAccessToken}, {"refresh_token", signedRefreshToken}, {"updated_at", updated_at}}},
		},
	)

	defer cancel()

	if err != nil {
		log.Panic(err)
	}

}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = err.Error()
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = err.Error()
		return
	}

	return claims, msg
}

