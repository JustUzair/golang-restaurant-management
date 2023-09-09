package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"restaurant-management/database/collections"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email string

	First_name string
	Last_name  string
	UID        string
	jwt.StandardClaims
}

func GenerateAllTokens(email string, firstname string, lastname string, uid string, c *gin.Context) (signedToken string, signedRefreshToken string, err error) {
	err = godotenv.Load("config.env")
	if err != nil {
		log.Fatalln(err)
	}
	jwt_secret := os.Getenv("JWT_SECRET")
	jwt_expires_in := os.Getenv("JWT_EXPIRES_IN")
	jwt_validity, _ := strconv.Atoi(jwt_expires_in)
	if jwt_expires_in == "" {
		jwt_validity = 1
	}

	if jwt_secret == "" {
		panic("JWT Signature is required")
	}

	claims := &SignedDetails{
		Email:      email,
		First_name: firstname,
		Last_name:  lastname,
		UID:        uid,
		StandardClaims: jwt.StandardClaims{

			ExpiresAt: time.Now().Local().Add(time.Hour * 24 * time.Duration(jwt_validity)).Unix(),
		},
	}
	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{

			ExpiresAt: time.Now().Local().Add(time.Hour * 24 * time.Duration(jwt_validity)).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(jwt_secret))
	if err != nil {
		log.Fatalln(err)
		return
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(jwt_secret))
	if err != nil {
		log.Fatalln(err)
		return
	}
	c.SetCookie("jwt", token, jwt_validity*24*60*60*1000, "/", "localhost", false, true)
	c.SetCookie("jwt_refresh", refreshToken, jwt_validity*24*60*60*1000, "/", "localhost", false, true)

	return token, refreshToken, err
}
func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var updateObj primitive.D
	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})
	updateObj = append(updateObj, bson.E{Key: "refreshToken", Value: signedRefreshToken})
	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: Updated_at})

	upsert := true
	filter := bson.M{"user_id": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	_, updateErr := collections.UserCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{Key: "$set", Value: updateObj},
		},
		&opt,
	)
	defer cancel()
	if updateErr != nil {
		log.Panic(updateErr)
		return
	}

}
func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatalln(err)
	}
	jwt_secret := os.Getenv("JWT_SECRET")
	if jwt_secret == "" {
		panic("JWT Signature is required")
	}

	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwt_secret), nil
	})
	if err != nil {
		msg = err.Error()
		return
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = fmt.Sprintln("the token is invalid")
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprintln("token has expired")
		return
	}

	return claims, msg
}
