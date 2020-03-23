package scraper

import (
	"log"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	// "net/http"
	// "github.com/scraping_with_go/utils"
	// "encoding/json"
	// "fmt"
)

type DBLayer interface{
	CreateRecord(record Record, database string, collection string) error
	ReadAllRecords(database string, collection string) ([]Record, error)
}

type MongoDBLayer struct{
	connection *mongo.Client
}

func (mdb MongoDBLayer) CreateRecord(record Record, database string, collection string) error {
	dbCollection := mdb.connection.Database(database).Collection(collection)
	
	_, err := dbCollection.InsertOne(context.TODO(), record)
	if err != nil {
		log.Panicln("Database write error")
	}
	return err
}

func (mdb MongoDBLayer) ReadAllRecords(database string, collection string) ([]Record, error) {
	var products []Record

	dbCollection := mdb.connection.Database(database).Collection(collection)
	cursor, err := dbCollection.Find(context.TODO(), bson.M{})
	defer cursor.Close(context.TODO())
	if err != nil {
		log.Panicln("Database read error")
		return products, err
	}

	for cursor.Next(context.TODO()) {
		var rec Record
		cursor.Decode(&rec)
		products = append(products, rec)
	}
	return products, err
}
