package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"restaurant-management/database/collections"
	"restaurant-management/models"
	"restaurant-management/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		recordsPerPage, err := strconv.Atoi(c.Query("recordsPerPage"))
		if err != nil || recordsPerPage < 1 {
			recordsPerPage = 10
		}
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}
		// startIndex := (page - 1) * recordsPerPage
		// startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
		skipStage := bson.D{{Key: "$skip", Value: page - 1}}
		limitStage := bson.D{{Key: "$limit", Value: recordsPerPage}}

		result, aggregationErr := collections.UserCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage,
			skipStage,
			limitStage,
		})

		if aggregationErr != nil {
			errorCode := http.StatusInternalServerError
			c.JSON(errorCode, gin.H{
				"statusCode": errorCode,
				"status":     "error",
				"message":    "error while listing users...!",
			})
			defer cancel()
			return
		}
		var allUsers []bson.M
		if err = result.All(ctx, &allUsers); err != nil {
			log.Fatal(err)
		}
		fmt.Println(allUsers)

		defer cancel()
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"users": allUsers,
			},
		})
	}
}
func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		userId := c.Param("user_id")
		var user models.User
		err := collections.UserCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()
		if err != nil {
			errorCode := http.StatusInternalServerError
			c.JSON(errorCode, gin.H{
				"statusCode": errorCode,
				"status":     "error",
				"message":    "error occured while fetching user",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"user": user,
			},
		})
	}
}
func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		// convert data to golang native format
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			errorCode := http.StatusBadRequest
			c.JSON(errorCode, gin.H{
				"statusCode": errorCode,
				"status":     "error",
				"message":    err.Error(),
			})
			defer cancel()
			return
		}
		// validate incoming data based on user struct model
		if validationErr := validate.Struct(user); validationErr != nil {
			errorCode := http.StatusBadRequest
			c.JSON(errorCode, gin.H{
				"statusCode": errorCode,
				"status":     "error",
				"message":    validationErr.Error(),
			})
			defer cancel()
			return
		}

		emailCount, err := collections.UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			errorCode := http.StatusInternalServerError
			c.JSON(errorCode, gin.H{
				"statusCode": errorCode,
				"status":     "error",
				"message":    "error while checking for email",
			})
			return
		}
		// hash password
		password := HashPassword(*user.Password)
		user.Password = &password
		phoneCount, err := collections.UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			errorCode := http.StatusInternalServerError
			c.JSON(errorCode, gin.H{
				"statusCode": errorCode,
				"status":     "error",
				"message":    "error while checking for phone.no",
			})
			return
		}
		if emailCount > 0 || phoneCount > 0 {
			errorCode := http.StatusInternalServerError
			c.JSON(errorCode, gin.H{
				"statusCode": errorCode,
				"status":     "error",
				"message":    "This email/phone already exists",
			})
			return
		}
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		// generate jwt token
		token, refreshToken, err := utils.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, user.User_id, c)
		user.Token = &token
		user.Refresh_Token = &refreshToken

		//create new user (insert user in DB)
		insertResult, insertionErr := collections.UserCollection.InsertOne(ctx, &user)
		if insertionErr != nil {

			errorCode := http.StatusInternalServerError
			c.JSON(errorCode, gin.H{
				"statusCode": errorCode,
				"status":     "error",
				"message":    "error signing up the user",
			})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"user": insertResult,
			},
		})
	}
}
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		// convert data to golang native format
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			errorCode := http.StatusBadRequest
			c.JSON(errorCode, gin.H{
				"statusCode": errorCode,
				"status":     "error",
				"message":    err.Error(),
			})
			defer cancel()
			return
		}
		err := collections.UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			errorCode := http.StatusInternalServerError
			c.JSON(errorCode, gin.H{
				"statusCode": errorCode,
				"status":     "error",
				"message":    "user not found, please sign up with this email",
			})
			return
		}
		// verify password
		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if !passwordIsValid {
			errorCode := http.StatusInternalServerError

			c.JSON(errorCode, gin.H{
				"statusCode": errorCode,
				"status":     "error",
				"message":    msg,
			})
			return
		}
		// generate jwt token
		token, refreshToken, err := utils.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, foundUser.User_id, c)
		user.Token = &token
		user.Refresh_Token = &refreshToken

		// update token and refresh token
		utils.UpdateAllTokens(token, refreshToken, foundUser.User_id)

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"jwt":  token,
				"user": foundUser,
			},
		})
	}
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""
	if err != nil {
		msg = "invalid email / password combination"
		check = false

	}
	return check, msg
}
