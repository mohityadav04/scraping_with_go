package scraper

import (
	"log"
	"context"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/scraping_with_go/utils"
)

var config = utils.GetConfig()

type DBLayer interface{
	CreateOrUpdateRecord(record Record, collection string) error
	ReadAllRecords(collection string) ([]Record, error)
	ReadARecordByProductId(productId string, collection string) (Record, error)
}

type MongoDBLayer struct{
	connection *mongo.Client
}

func (mdb MongoDBLayer) CreateOrUpdateRecord(record Record, collection string) error {
	dbCollection := mdb.connection.Database("amazon").Collection(collection)

	// Check if record exists in db and update it
	pId := record.ProductId
	filter := bson.M{
		"productid": bson.M{
			"$eq": pId,
		},
	}
	update := bson.M{
		"$set": bson.M{
			"url": record.Url,
			"name": record.ProductName,
			"imageurl": record.ImageUrl,
			"description": record.Description,
			"price": record.ProductPrice,
			"totalreviews": record.TotalReviews,
			"lastupdatedat": time.Now().String(),
		},
	}

	updateResult, err := dbCollection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		log.Panicln("error when updating in database")
	}

	if updateResult.MatchedCount == 0 {
		record.CreatedAt = time.Now().String()
		_, err := dbCollection.InsertOne(context.TODO(), record)
		if err != nil {
			log.Panicln("Database write error")
		}
		return err
	}
	return err
}


func (mdb MongoDBLayer) ReadAllRecords(collection string) ([]Record, error) {
	var products []Record

	dbCollection := mdb.connection.Database("amazon").Collection(collection)
	cursor, err := dbCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Panicln("Database read error")
		return products, err
	}

	for cursor.Next(context.TODO()) {
		var rec Record
		cursor.Decode(&rec)
		products = append(products, rec)
	}
	defer cursor.Close(context.TODO())
	return products, err
}

func (mdb MongoDBLayer) ReadARecordByProductId(productId string, collection string) (Record,error) {
	var record Record
	filter := bson.M{
		"productid": bson.M{
			"$eq": productId,
		},
	}
	dbCollection := mdb.connection.Database("amazon").Collection(collection)
	err := dbCollection.FindOne(context.Background(), filter).Decode(&record)
	if err != nil{
		log.Panicln("Error reading from database with given filter")
	}
	return record, err
}
