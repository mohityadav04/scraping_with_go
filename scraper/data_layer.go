package scraper

import (
	"encoding/json"
)


type Record struct{
	Url				string  `json:"url" bson:"url,omitempty"`
	ProductName 	string	`json:"name" bson:"name,omitempty"`
	ImageUrl		string 	`json:"imageurl" bson:"imageurl,omitempty"`
	Description 	string 	`json:"description" bson:"description,omitempty"`
	ProductPrice 	string	`json:"price" bson:"price,omitempty"`
	TotalReviews 	int  	`json:"totalreviews" bson:"totalreviews,omitempty"`
}

func (r *Record) ToJSON() []byte {
	data, err := json.Marshal(r)
	if err != nil {
		panic("Not able to convert the record to json")
	}
	return data
}