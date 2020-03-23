package utils

import (
	"sync"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type client map[string]*mongo.Client

var (
	once sync.Once
	instance client
)

func GetDBClient() *mongo.Client {
	once.Do(func(){
		clientOptions := options.Client().ApplyURI("mongodb://db:27017/?connect=direct")
		conn, err := mongo.Connect(context.TODO(), clientOptions)

		if err != nil {
			log.Fatal(err)
		}

		instance = make(client)
		instance["db"] = conn
	})

	return instance["db"]
}
