package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBSet() *mongo.Client{
	MongoDB:="mongodb://localhost:27017"
	fmt.Println("MongoDB URL: ", MongoDB)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MongoDB))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB Database")
	return client
}

var Client *mongo.Client = DBSet()


func UserData(client *mongo.Client, collectionName string ) *mongo.Collection {
	var userCollection *mongo.Collection = client.Database("ginapp").Collection(collectionName)
	return userCollection
}

func ProductData(client *mongo.Client, collectionName string) *mongo.Collection {
	var productCollection *mongo.Collection = client.Database("ginapp").Collection(collectionName)
	return productCollection
}

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("ginapp").Collection(collectionName)
	return collection
}