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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderItemPack struct {
	Table_id    *string
	Order_items []models.OrderItem
}

func GetOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		result, err := collections.OrderItemCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {

			errorCode := http.StatusInternalServerError
			c.JSON(errorCode, gin.H{
				"statusCode": errorCode,
				"status":     "error",
				"message":    "error while listing order items...!",
			})
		}

		var orderItems []bson.M
		if err = result.All(ctx, &orderItems); err != nil {
			log.Fatal(err)
		}
		defer cancel()
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"orderItems": orderItems,
			},
		})
	}
}
func GetOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		orderItemId := c.Param("order_item_id")
		var orderItem models.OrderItem
		err := collections.OrderItemCollection.FindOne(ctx, bson.M{"order_item_id": orderItemId}).Decode(&orderItem)
		defer cancel()
		if err != nil {
			errorCode := http.StatusInternalServerError
			c.JSON(errorCode, gin.H{
				"statusCode": errorCode,
				"status":     "error",
				"message":    "error occured while fetching order item",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"order_item": orderItem,
			},
		})
	}
}
func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var orderItemPack OrderItemPack
		var order models.Order
		if err := c.BindJSON(&orderItemPack); err != nil {
			errorCode := http.StatusBadRequest
			c.JSON(errorCode, gin.H{
				"statusCode": errorCode,
				"status":     "error",
				"message":    err.Error(),
			})
			defer cancel()
			return
		}
		order.Order_date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		orderItemsToInsert := []interface{}{}
		order.Table_id = orderItemPack.Table_id
		order_id := OrderItemOrderCreator(order)
		for _, orderItem := range orderItemPack.Order_items {
			orderItem.Order_id = order_id
			if validationErr := validate.Struct(orderItem); validationErr != nil {
				errorCode := http.StatusBadRequest
				c.JSON(errorCode, gin.H{
					"statusCode": errorCode,
					"status":     "error",
					"message":    validationErr.Error(),
				})
				defer cancel()
				return
			}
			orderItem.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.ID = primitive.NewObjectID()
			orderItem.Order_item_id = orderItem.ID.Hex()
			var num = toFixed(*orderItem.Unit_price, 2)
			orderItem.Unit_price = &num
			orderItemsToInsert = append(orderItemsToInsert, orderItem)
		}
		result, insertionErr := collections.InvoiceCollection.InsertMany(ctx, orderItemsToInsert)
		if insertionErr != nil {

			errorCode := http.StatusInternalServerError
			c.JSON(errorCode, gin.H{
				"statusCode": errorCode,
				"status":     "error",
				"message":    "error creating invoice record",
			})
			defer cancel()
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"order_items": result,
			},
		})
	}
}
func OrderItemOrderCreator(models.Order) string {
	return ""
}
func UpdateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var orderItem models.OrderItem
		orderItemId := c.Param("order_item_id")
		var updateObj primitive.D

		if orderItem.Unit_price != nil {
			updateObj = append(updateObj, bson.E{Key: "unit_price", Value: orderItem.Unit_price})

		}
		if orderItem.Quantity != nil {
			updateObj = append(updateObj, bson.E{Key: "quantity", Value: orderItem.Quantity})

		}
		if orderItem.Food_id != nil {
			updateObj = append(updateObj, bson.E{Key: "food_id", Value: orderItem.Food_id})

		}

		orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: orderItem.Updated_at})

		upsert := true
		filter := bson.M{"order_item_id": orderItemId}
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		result, updateErr := collections.OrderItemCollection.UpdateOne(
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
				"message":    "order item updation failed",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"order_item": result,
			},
		})
		defer cancel()

	}
}
func GetOrderItemsByOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderId := c.Param("order_id")
		allOrderItems, err := ItemsByOrder(orderId)

		if err != nil {
			errorCode := http.StatusInternalServerError
			c.JSON(errorCode, gin.H{
				"statusCode": errorCode,
				"status":     "error",
				"message":    "error occured while listing order items by order id",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"order_items": allOrderItems,
			},
		})
	}
}

func ItemsByOrder(id string) (OrderItems []primitive.M, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "order_id", Value: id}}}}
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "food"},
		{Key: "localField", Value: "food_id"},
		{Key: "foreignField", Value: "food_id"},
		{Key: "as", Value: "food"},
	}}}
	unwindStage := bson.D{{
		Key: "$unwind", Value: bson.D{{Key: "path", Value: "$food"}, {Key: "preserveNullAndEmptyArrays", Value: true}},
	}}
	lookupOrderStage := bson.D{{
		Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "order"},
			{Key: "localField", Value: "order_id"},
			{Key: "foreignField", Value: "order_id"},
			{Key: "as", Value: "order"},
		},
	}}
	unwindOrderStage := bson.D{{
		Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$order"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		},
	}}
	lookupTableStage := bson.D{{
		Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "table"},
			{Key: "localField", Value: "order.table_id"},
			{Key: "foreignField", Value: "table_id"},
			{Key: "as", Value: "table"},
		},
	}}
	unwindTableStage := bson.D{{
		Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$table"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		},
	}}

	projectStage := bson.D{
		{
			Key: "$projeValue: ct", Value: bson.D{
				{Key: "id", Value: 0},
				{Key: "amouValue: nt", Value: "$food.price"},
				{Key: "total_count", Value: 1},
				{Key: "food_name", Value: "$food.name"},
				{Key: "food_image", Value: "$food.food_image"},
				{Key: "food_image", Value: "$food.food_image"},
				{Key: "table_number", Value: "$table.table_number"},
				{Key: "table_id", Value: "$table.table_id"},
				{Key: "order_id", Value: "$table.order_id"},
				{Key: "price", Value: "$food.price"},
				{Key: "quantity", Value: 1},
			},
		},
	}
	groupStage := bson.D{
		{
			Key: "_id", Value: bson.D{
				{Key: "order_id", Value: "$order_id"},
				{Key: "table_id", Value: "$table_id"},
				{Key: "table_number", Value: "$table_number"},
			},
		},
		{
			Key: "payment_due", Value: bson.D{{Key: "$sum", Value: "$amount"}},
		},
		{
			Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}},
		},
		{
			Key: "order_items", Value: bson.D{{Key: "$push", Value: "$$ROOT"}},
		},
	}
	projectStage2 := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "id", Value: 0},
			{Key: "payment_due", Value: 1},
			{Key: "total_count", Value: 1},
			{Key: "table_number", Value: "$_id.table_number"},
			{Key: "order_items", Value: 1},
		}},
	}
	result, err := collections.OrderItemCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage,
		lookupStage,
		unwindStage,
		lookupOrderStage,
		unwindOrderStage,
		lookupTableStage,
		unwindTableStage,
		projectStage,
		groupStage,
		projectStage2,
	})
	defer cancel()
	if err != nil {
		panic(err)
	}
	if err = result.All(ctx, &OrderItems); err != nil {
		panic(err)
	}
	defer cancel()
	return OrderItems, err
}
