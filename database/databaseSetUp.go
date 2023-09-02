package database

import (
	"context"
	"fmt"
	"log"
	"time"
	"os"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBSet() *mongo.Client{
	err := godotenv.Load(".env")
	if err != nil {
    log.Println("Error loading .env file")
  }
	MongoDB:= os.Getenv("MONGO_URL")
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

// OPEN COLLECTION 


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