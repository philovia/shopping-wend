package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBSet() *mongo.Client {

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Println("failed to connect the mongodb")
		return nil
	}

	fmt.Println("Successfully connected to mongodb")
	return client

}

var Client *mongo.Client = DBSet()

func init() {
	Client = DBSet()
}

func UserData(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("Shopping").Collection(collectionName)
	return collection

}

func ProductData(client *mongo.Client, collectionName string) *mongo.Collection {
	var productCollection *mongo.Collection = client.Database("Shopping").Collection(collectionName)
	return productCollection

}
