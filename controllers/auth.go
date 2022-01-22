package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/riad-safowan/JWT-GO-MongoDB/helpers"
	"github.com/riad-safowan/JWT-GO-MongoDB/models"
	"github.com/riad-safowan/JWT-GO-MongoDB/models/response"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}
func VerifyPassword(userPass string, providedPass string) (passIsValid bool, msg string) {
	err := bcrypt.CompareHashAndPassword([]byte(userPass), []byte(providedPass))
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(userPass, providedPass)
		return false, fmt.Sprint("email or password is incorrect")
	} else {
		return true, ""
	}

}
func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			defer cancel()
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			defer cancel()
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while chacking email existance"})
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		count, err = userCollection.CountDocuments(ctx, bson.M{"phone_number": user.Phone_number})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while chacking phone number existance"})
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "This email or phone number already exist"})
			return
		}
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		accessToken, refreshToken, _ := helpers.GenerateAllToken(*user.Email, *user.First_name, *user.Last_name, *user.User_type, user.User_id)
		user.Access_token = &accessToken
		user.Refresh_token = &refreshToken

		_, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprint("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, getLoginResponse(user))
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			defer cancel()
			return
		}
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email or password is incorrect"})
			defer cancel()
			return
		}

		passIsValid, msg := VerifyPassword(*foundUser.Password, *user.Password)
		defer cancel()

		if !passIsValid {
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}

		accessToken, refreshToken, _ := helpers.GenerateAllToken(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, foundUser.User_id)
		helpers.UpdateAllTokens(accessToken, refreshToken, foundUser.User_id)

		err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, getLoginResponse(foundUser))
	}
}

type TokenModel struct {
	Access_token  string `json:"access_token"`
	Refresh_token string `json:"refresh_token"`
}

func RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)

		clientToken := c.Request.Header.Get("Authorization")
		if clientToken == "" {
			clientToken = c.Request.Header.Get("token")
		} else if strings.HasPrefix(clientToken, "Bearer ") {
			reqToken := c.Request.Header.Get("Authorization")
			splitToken := strings.Split(reqToken, "Bearer ")
			clientToken = splitToken[1]
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid authorization token"})
			c.Abort()
			return
		}

		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no Authorization header provided"})
			c.Abort()
			return
		}
		// handle access token
		claims, err := helpers.ValidateToken(clientToken)

		if err != "" || claims.Token_type!="refresh_token"{
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid authorization token"})
			c.Abort()
			return
		}

		var user models.User
		id, _ := primitive.ObjectIDFromHex(claims.Uid)
		e := userCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
		if e != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
			defer cancel()
			return
		}

		accessToken, refreshToken, _ := helpers.GenerateAllToken(*user.Email, *user.First_name, *user.Last_name, *user.User_type, user.User_id)
		helpers.UpdateAllTokens(accessToken, refreshToken, user.User_id)

		var tokenModel = TokenModel{Access_token: accessToken, Refresh_token: refreshToken}
		c.JSON(http.StatusOK, tokenModel)
		defer cancel()
	}
}

func getLoginResponse(user models.User) response.LoginResponse {
	var loginResponse = response.LoginResponse{}
	b, _ := json.Marshal(&user)
	json.Unmarshal(b, &loginResponse)
	return loginResponse
}
