package main

import(
	// "fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/scraping_with_go/scraper"
	"github.com/scraping_with_go/utils"

)

func main(){

	r := mux.NewRouter()
	config := utils.GetConfig()

	r.HandleFunc("/products/", scraper.CreateDBRecordHandler).Methods("POST")
	r.HandleFunc("/products/", scraper.GetDBRecordsHandler).Methods("GET")
	r.HandleFunc("/resources/", scraper.ScraperHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(":" + config["appPort"], r))

	// Connect to MongoDB
}