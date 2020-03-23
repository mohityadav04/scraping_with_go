package scraper

import(
	"net/http"
	"encoding/json"
	"github.com/gocolly/colly/v2"
	"strings"
	"strconv"
	"github.com/scraping_with_go/utils"
	// "time"
	"bytes"
	_ "go.mongodb.org/mongo-driver/bson/primitive"
	_ "log"
	"errors"
)

func ScrapeHandler(w http.ResponseWriter, r *http.Request){

	var rec Record
	_ = json.NewDecoder(r.Body).Decode(&rec)

	urlToProcess := rec.Url
	c := colly.NewCollector()

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

	w.Header().Set("Content-Type", "application/json")
	err := POSTRecordAPI(rec)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "Error": "Not able to process request" }`))
		return
	}
	json.NewEncoder(w).Encode(rec)

}

func POSTRecordAPI(record Record) error {
	client := &http.Client{}
	config := utils.GetConfig()
	req, err := http.NewRequest("POST", "http://localhost:"+ config["appPort"] +"/products/", bytes.NewBuffer(record.ToJSON()))
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return err
	}
	
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return err
}