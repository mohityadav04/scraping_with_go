package scraper

import(
	"net/http"
	"encoding/json"
	"log"
	"bytes"
	"strings"
	"strconv"

	"github.com/gocolly/colly/v2"

	"github.com/scraping_with_go/utils"

)

var dbClient = utils.GetDBClient()
var dBLayer = MongoDBLayer{dbClient}

func GetDBRecordsHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("content-type", "application/json")

	log.Println("recieved request to read records")

	records, err := dBLayer.ReadAllRecords("products")
	
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}


	w.WriteHeader(http.StatusOK)
	if len(records) == 0 {
		w.Write([]byte("[]"))
		return
	}

	json.NewEncoder(w).Encode(records)
}

func CreateDBRecordHandler(w http.ResponseWriter, r *http.Request){
	log.Println("recieved request to create record")
	var rec Record
	_ = json.NewDecoder(r.Body).Decode(&rec)

	err := dBLayer.CreateOrUpdateRecord(rec, "products")
	
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rec)

}

func ScraperHandler(w http.ResponseWriter, r *http.Request){

	var rec Record
	_ = json.NewDecoder(r.Body).Decode(&rec)

	urlToProcess := rec.Url
	c := colly.NewCollector()

	urlMatchForProductPage := utils.ValidateURL(urlToProcess)
	rec.ProductId = urlMatchForProductPage[4:]

	if urlMatchForProductPage == "" {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Invalid URL Found")
		w.Write([]byte(`{ "Error": "invalid URL" }`))
	
		return
	}
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set(
			"User-Agent", 
			"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.132 Safari/537.36")
	})

	c.OnHTML("div#ppd", func(e *colly.HTMLElement){
		// product name
		rec.ProductName = e.ChildText(`span#productTitle`)

		// image URL
		productImagesMap := make(map[string] interface{})
		allImagesOfProduct := e.ChildAttr("img#landingImage", "data-a-dynamic-image")
		_ = json.Unmarshal([]byte(allImagesOfProduct), &productImagesMap)

		for pImage, _ := range productImagesMap{
			rec.ImageUrl = pImage
			break
		}

		// price
		if e.ChildText("#priceblock_ourprice") == "" {

			if e.ChildText("#priceblock_dealprice") == "" {

				e.ForEach("div#olp_feature_div", func(_ int, usedPriceDiv *colly.HTMLElement){

					if usedPriceDiv.ChildText(".a-color-price") == "" {

						e.ForEach("div#buybox", func(_ int, buyBoxPriceDiv *colly.HTMLElement){
							if buyBoxPriceDiv.ChildText(".a-color-price") == "" {
								rec.ProductPrice = "-1"
							} else {
								rec.ProductPrice = buyBoxPriceDiv.ChildText(".a-color-price")
							}
						})
					} else {
						rec.ProductPrice = usedPriceDiv.ChildText(".a-color-price")
					}
				})
			} else {
				rec.ProductPrice = e.ChildText("span#priceblock_dealprice")
			}
		} else {
			rec.ProductPrice = e.ChildText("span#priceblock_ourprice")
		}
		//check if price from HTML got only string and replace it with missing label
		if utils.MatchPrice(rec.ProductPrice) == "" {
			rec.ProductPrice = "-1"
		}

		// reviews
		customerReviewes := strings.Split(e.ChildText("span#acrCustomerReviewText"), " ")
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
		// remove unnecessary details
		if len(pdesc) > 1 {
			rec.Description = strings.Join(pdesc[1:], "")
		}
	})
	c.Visit(urlToProcess)

	w.Header().Set("Content-Type", "application/json")
	err := POSTRecordAPI(rec)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "Error": "Not able to process request to save in DB" }`))
		return
	}

	log.Println("Scraping complete. Result is:")
	log.Println(rec)

	dbRecord, err := dBLayer.ReadARecordByProductId(rec.ProductId, "products")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "Error": "there seems to be some problem conncting to database"`))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dbRecord)

}

func POSTRecordAPI(record Record) error {
	client := &http.Client{}
	config := utils.GetConfig()

	log.Print("Preparaing to send request to create database record after scraping")
	req, err := http.NewRequest(
					"POST", 
					"http://" + config["appHost"] + ":"+ config["appPort"] +"/products/", 
					bytes.NewBuffer(record.ToJSON()))
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		log.Panicln("Error in creating request" + err.Error())
		return err
	}
	
	resp, err := client.Do(req)
	if err != nil {
		log.Panicln("couldn't send sending request" + err.Error())
		return err
	}
	defer resp.Body.Close()
	return err
}