package utils

import (
	"sync"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type client map[string]*mongo.Client

var (
	once sync.Once
	instance client
)

func GetDBClient() *mongo.Client {
	config := GetConfig()
	once.Do(func(){
		clientOptions := options.Client().ApplyURI("mongodb://"+config["mongoAddress"]+"/?connect=direct")
		conn, err := mongo.Connect(context.TODO(), clientOptions)

		if err != nil {
			log.Fatal(err)
		}

		instance = make(client)
		instance["db"] = conn
	})

	return instance["db"]
}
