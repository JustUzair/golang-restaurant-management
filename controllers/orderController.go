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

func GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		result, err := collections.OrderCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {

			errorCode := http.StatusInternalServerError
			c.JSON(errorCode, gin.H{
				"statusCode": errorCode,
				"status":     "error",
				"message":    "error while listing order items...!",
			})
		}
		var orders []bson.M
		if err = result.All(ctx, &orders); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"orders": orders,
			},
		})
	}
}
func GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		orderId := c.Param("order_id")
		var order models.Order
		err := collections.OrderCollection.FindOne(ctx, bson.M{"order_id": orderId}).Decode(&order)
		defer cancel()
		if err != nil {
			errorCode := http.StatusInternalServerError
			c.JSON(errorCode, gin.H{
				"statusCode": errorCode,
				"status":     "error",
				"message":    "error occured while fetching order",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"order": order,
			},
		})
	}
}
func CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var table models.Table
		var order models.Order

		if err := c.BindJSON(&order); err != nil {
			errorCode := http.StatusBadRequest
			c.JSON(errorCode, gin.H{
				"statusCode": errorCode,
				"status":     "error",
				"message":    err.Error(),
			})
			defer cancel()
			return
		}

		if validationErr := validate.Struct(order); validationErr != nil {
			errorCode := http.StatusBadRequest
			c.JSON(errorCode, gin.H{
				"statusCode": errorCode,
				"status":     "error",
				"message":    validationErr.Error(),
			})
		}

		if order.Table_id != nil {
			err := collections.TableCollection.FindOne(ctx, bson.M{
				"table_id": order.Table_id,
			}).Decode(&table)
			defer cancel()
			if err != nil {
				errorCode := http.StatusInternalServerError
				c.JSON(errorCode, gin.H{
					"statusCode": errorCode,
					"status":     "error",
					"message":    "No order with that table ID exists!",
				})
				return
			}
		}
		defer cancel()
		order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.ID = primitive.NewObjectID()
		order.Order_id = order.ID.Hex()

		result, insertionErr := collections.OrderCollection.InsertOne(ctx, &order)
		if insertionErr != nil {

			errorCode := http.StatusInternalServerError
			c.JSON(errorCode, gin.H{
				"statusCode": errorCode,
				"status":     "error",
				"message":    "error creating order item",
			})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"order": result,
			},
		})
	}
}
func UpdateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var table models.Table
		var order models.Order
		orderId := c.Param("order_id")

		if err := c.BindJSON(&order); err != nil {
			errorCode := http.StatusBadRequest
			c.JSON(errorCode, gin.H{
				"statusCode": errorCode,
				"status":     "error",
				"message":    err.Error(),
			})
			defer cancel()
			return
		}
		var updateObj primitive.D
		if order.Table_id != nil {
			err := collections.OrderCollection.FindOne(ctx, bson.M{"table_id": order.Table_id}).Decode(&table)
			defer cancel()
			if err != nil {
				errorCode := http.StatusInternalServerError
				c.JSON(errorCode, gin.H{
					"statusCode": errorCode,
					"status":     "error",
					"message":    "given table_id does not exist",
				})
				return
			}
			updateObj = append(updateObj, bson.E{Key: "table", Value: order.Table_id})
		}
		order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: order.Updated_at})
		upsert := true
		filter := bson.M{"order_id": orderId}
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		result, updateErr := collections.OrderCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{Key: "$set", Value: updateObj},
			},
			&opt,
		)
		defer cancel()
		if updateErr != nil {
			errorCode := http.StatusInternalServerError
			c.JSON(errorCode, gin.H{
				"statusCode": errorCode,
				"status":     "error",
				"message":    "food item updation failed",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"food": result,
			},
		})
		defer cancel()
	}
}

func OrderIemOrderCreator(order models.Order) string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.ID = primitive.NewObjectID()
	order.Order_id = order.ID.Hex()
	defer cancel()
	collections.OrderCollection.InsertOne(ctx, &order)

	return ""
}
