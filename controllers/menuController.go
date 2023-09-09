package controllers

import (
	"context"
	"log"
	"net/http"
	"restaurant-management/database/collections"
	"restaurant-management/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		result, err := collections.MenuCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {

			errorCode := http.StatusInternalServerError
			c.JSON(errorCode, gin.H{
				"statusCode": errorCode,
				"status":     "error",
				"message":    "error while listing menu items...!",
			})

		}
		var menus []bson.M
		if err = result.All(ctx, &menus); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"menus": menus,
			},
		})
	}
}
func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		menuId := c.Param("menu_id")
		var menu models.Menu
		err := collections.MenuCollection.FindOne(ctx, bson.M{"menu_id": menuId}).Decode(&menu)

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
				"menu": menu,
			},
		})
	}
}
func CreateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var menu models.Menu
		if err := c.BindJSON(&menu); err != nil {
			errorCode := http.StatusBadRequest
			c.JSON(errorCode, gin.H{
				"statusCode": errorCode,
				"status":     "error",
				"message":    err.Error(),
			})
			defer cancel()
			return
		}

		if validationErr := validate.Struct(menu); validationErr != nil {
			errorCode := http.StatusBadRequest
			c.JSON(errorCode, gin.H{
				"statusCode": errorCode,
				"status":     "error",
				"message":    validationErr.Error(),
			})
			defer cancel()
			return
		}

		menu.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.ID = primitive.NewObjectID()
		menu.Menu_id = menu.ID.Hex()

		result, err := collections.MenuCollection.InsertOne(ctx, menu)
		defer cancel()
		if err != nil {
			errorCode := http.StatusBadRequest
			c.JSON(errorCode, gin.H{
				"statusCode": errorCode,
				"status":     "error",
				"message":    "menu creation failed",
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"menu": result,
			},
		})
		defer cancel()
	}
}
func UpdateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var menu models.Menu
		if err := c.BindJSON(&menu); err != nil {
			errorCode := http.StatusBadRequest
			c.JSON(errorCode, gin.H{
				"statusCode": errorCode,
				"status":     "error",
				"message":    err.Error(),
			})
			return
		}
		menuId := c.Param("menu_id")
		filter := bson.M{"menu_id": menuId}
		var updateObj primitive.D
		if menu.Start_Date != nil && menu.End_Date != nil {
			if !inTimeSpan(*menu.Start_Date, *menu.End_Date, time.Now()) {
				msg := "Invalid Time"
				errorCode := http.StatusInternalServerError
				c.JSON(errorCode, gin.H{
					"statusCode": errorCode,
					"status":     "error",
					"message":    msg,
				})
				defer cancel()
				return
			}
			updateObj = append(updateObj, bson.E{Key: "start_date", Value: menu.Start_Date})
			updateObj = append(updateObj, bson.E{Key: "end_date", Value: menu.End_Date})

			if menu.Name != "" {
				updateObj = append(updateObj, bson.E{Key: "name", Value: menu.Name})
			}
			if menu.Category != "" {
				updateObj = append(updateObj, bson.E{Key: "name", Value: menu.Category})
			}
			menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			updateObj = append(updateObj, bson.E{Key: "updated_at", Value: menu.Updated_at})
			upsert := true

			opt := options.UpdateOptions{
				Upsert: &upsert,
			}
			result, updateErr := collections.MenuCollection.UpdateOne(ctx, filter, bson.D{
				{Key: "$set", Value: updateObj},
			}, &opt)
			defer cancel()
			if updateErr != nil {
				errorCode := http.StatusBadRequest
				c.JSON(errorCode, gin.H{
					"statusCode": errorCode,
					"status":     "error",
					"message":    "menu updation failed",
				})
			}
			c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"data": gin.H{
					"menu": result,
				},
			})
			defer cancel()

		}
	}
}

func inTimeSpan(start, end, check time.Time) bool {
	return (start.After(time.Now())) && end.After(start)
}
