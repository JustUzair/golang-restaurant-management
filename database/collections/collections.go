package collections

import (
	"go.mongodb.org/mongo-driver/mongo"
	"restaurant-management/database"
)

var FoodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")
var MenuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")
var OrderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")
var TableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")
var InvoiceCollection *mongo.Collection = database.OpenCollection(database.Client, "invoice")
