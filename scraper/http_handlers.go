package scraper

import(
	"net/http"
	"encoding/json"
	"log"
	"bytes"
	"strings"
	"strconv"
	
	"github.com/gocolly/colly/v2"
	"go.mongodb.org/mongo-driver/mongo"
	// "time"
	// _ "go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/scraping_with_go/utils"

)

var dbClient *mongo.Client = utils.GetDBClient()
var dBLayer DBLayer = MongoDBLayer{dbClient}

func GetDBRecordsHandler(w http.ResponseWriter, r *http.Request){
	log.Println("recieved request to read records")
	records, err := dBLayer.ReadAllRecords("amazon", "products")
	
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Panicln("Error getting records from database" + err.Error())
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(records)
}

func CreateDBRecordHandler(w http.ResponseWriter, r *http.Request){
	log.Println("recieved request to create record")
	var rec Record
	_ = json.NewDecoder(r.Body).Decode(&rec)

	err := dBLayer.CreateRecord(rec, "amazon", "products")
	
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Panicln("Error creating record in database" + err.Error())
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rec)


}

func ScraperHandler(w http.ResponseWriter, r *http.Request){

	var rec Record
	_ = json.NewDecoder(r.Body).Decode(&rec)

	urlToProcess := rec.Url
	c := colly.NewCollector()

	isUrlValid, err := utils.ValidateURL(urlToProcess)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Panicln("Problem parsing URL")
		w.Write([]byte(`{ "Error": "can not parse URL" }`))		
		return
	}
	if isUrlValid == false {
		w.WriteHeader(http.StatusInternalServerError)
		log.Panicln("Invalid URL Found")
		w.Write([]byte(`{ "Error": "invalid URL" }`))	
		return
	}

	c.OnHTML("div#ppd", func(e *colly.HTMLElement){
		// product name
		rec.ProductName = e.ChildText(`span#productTitle`)

		// image URL
		productImagesMap := make(map[string] interface{})
		allImagesOfProduct := e.ChildAttr(`img#landingImage`, "data-a-dynamic-image")
		_ = json.Unmarshal([]byte(allImagesOfProduct), &productImagesMap)

		for k, _ := range productImagesMap{
			rec.ImageUrl = k
			break
		}

		// price
		// TODO: missing product price handle
		rec.ProductPrice = e.ChildText(`span#priceblock_ourprice`)
		if rec.ProductPrice == "" {
			rec.ProductPrice = "-1"
		}

		// reviews
		customerReviewes := strings.Split(e.ChildText(`span#acrCustomerReviewText`), " ")
		if len(customerReviewes) != 0 {
			rev, err := strconv.Atoi(strings.Trim(customerReviewes[0], " "))
			if err != nil {
				rec.TotalReviews = 0
			} else {
				rec.TotalReviews = rev
			}
		} else {
			rec.TotalReviews = 0
		}

		// decription
		var pdesc []string
		e.ForEach("div#featurebullets_feature_div", func(_ int, descDiv *colly.HTMLElement){
			descDiv.ForEach("span.a-list-item", func(_ int, elem *colly.HTMLElement) {
				cleanedText := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(elem.Text, "\t", ""), "\n", ""),"  ", "")
				pdesc = append(pdesc, cleanedText)
			})
		})
		rec.Description = strings.Join(pdesc[1:], "")

	})
	c.Visit(urlToProcess)

	log.Println("Scraping complete. Result is:")
	log.Println(rec)

	w.Header().Set("Content-Type", "application/json")
	err = POSTRecordAPI(rec)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Panicln("Error sending request to create record in database" + err.Error())
		w.Write([]byte(`{ "Error": "Not able to process request" }`))
		return
	}
	json.NewEncoder(w).Encode(rec)

}

func POSTRecordAPI(record Record) error {
	client := &http.Client{}
	config := utils.GetConfig()

	log.Print("Preparaing to send request to create database record after scraping")
	req, err := http.NewRequest("POST", "http://localhost:"+ config["appPort"] +"/products/", bytes.NewBuffer(record.ToJSON()))
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		log.Panicln("Error in creating request")
		return err
	}
	
	resp, err := client.Do(req)
	if err != nil {
		log.Panicln("couldn't send sending request")
		return err
	}
	defer resp.Body.Close()
	return err
}