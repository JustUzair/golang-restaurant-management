package controllers

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func GetUsers() gin.HandlerFunc {
	log.Println("user route : /")
	fmt.Println("user route : /")
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "success",
		})
	}
}
func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func HashPassword(password string) string {
	return ""
}
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	return false, ""
}
