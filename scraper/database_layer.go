package scraper

import (

	"context"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"github.com/scraping_with_go/utils"
	"encoding/json"
	"fmt"
)

type DBLayer interface{
	CreateRecord() error
	ReadMultipleRecords() ([]Record, error)
}

type MongoDBLayer struct{
	connection *mongo.Client
}

func (mdb MongoDBLayer) CreateRecord(record Record) (*InsertOneResult, error) {
	collection := dbClient.Database("amazon").Collection("products")
	return collection.InsertOne(context.TODO(), rec)
}

func (mdb MongoDBLayer) ReadMultipleRecords()
