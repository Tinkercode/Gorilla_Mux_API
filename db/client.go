package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//DB is collection struct
type DB struct {
	Collection *mongo.Collection
}

//GetClient returns db struct for user and employee,mongo client and context
func GetClient() (*DB, *DB, *mongo.Client, context.Context) {
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017/")
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer cancel()

	EmployeeCollection := client.Database("assignment").Collection("EmployeeCollection")
	dbEmployee := &DB{Collection: EmployeeCollection}
	UserCollection := client.Database("assignment").Collection("UserCollection")
	dbUser := &DB{Collection: UserCollection}
	return dbUser, dbEmployee, client, ctx
}
