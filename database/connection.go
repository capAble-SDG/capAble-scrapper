package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Connection() *mongo.Client {
	options := options.Client().ApplyURI("mongodb://user:password@localhost:27017")
	client, err := mongo.Connect(context.TODO(), options)
	if err != nil {
		fmt.Println(err)
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	return client
}
