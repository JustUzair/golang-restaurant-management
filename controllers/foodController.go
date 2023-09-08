package controllers

import (
	"context"
	"fmt"
	"net/http"
	"restaurant-management/database/collections"
	"restaurant-management/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var validate = validator.New()

func GetFoods() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
func GetFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		foodId := c.Param("food_id")
		var food models.Food
		err := collections.FoodCollection.FindOne(ctx, bson.M{"food_id": foodId}).Decode(&food)

		if err != nil {
			errorCode := http.StatusInternalServerError
			c.JSON(errorCode, gin.H{
				"statusCode": errorCode,
				"status":     "error",
				"message":    "Internal Server Error",
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"food": food,
			},
		})
	}
}
func CreateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var food models.Food
		var menu models.Menu
		if err := c.BindJSON(&food); err != nil {
			errorCode := http.StatusBadRequest
			c.JSON(errorCode, gin.H{
				"statusCode": errorCode,
				"status":     "error",
				"message":    err.Error(),
			})
			return
		}

		if validationErr := validate.Struct(food); validationErr != nil {
			errorCode := http.StatusBadRequest
			c.JSON(errorCode, gin.H{
				"statusCode": errorCode,
				"status":     "error",
				"message":    validationErr.Error(),
			})
		}

		err := collections.MenuCollection.FindOne(ctx, bson.M{
			"menu_id": food.Menu_id,
		}).Decode(&menu)
		defer cancel()
		if err != nil {
			errorCode := http.StatusInternalServerError
			c.JSON(errorCode, gin.H{
				"statusCode": errorCode,
				"status":     "error",
				"message":    "No menu with that menu ID exists!",
			})
			return
		}
		food.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.ID = primitive.NewObjectID()
		food.Food_id = food.ID.Hex()
		var num = toFixed(*food.Price, 2)
		food.Price = &num
		result, insertionErr := collections.FoodCollection.InsertOne(ctx, &food)
		if insertionErr != nil {
			msg := fmt.Sprint("error creating food item")
			errorCode := http.StatusInternalServerError
			c.JSON(errorCode, gin.H{
				"statusCode": errorCode,
				"status":     "error",
				"message":    msg,
			})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"food": result,
			},
		})
	}
}
func UpdateFood() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func round(num float64) int {

}

func toFixed(num float64, precision int) float64 {

}
