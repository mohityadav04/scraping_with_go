package main

import(
	"fmt"
	"log"
	"github.com/gorilla/mux"
	"github.com/scraping_with_go/scraper"
	"net/http"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/scraping_with_go/utils"

)

func main(){
	// fmt.Println("Hello")
	fmt.Printf("%T", utils.GetDBClient())


	r := mux.NewRouter()

	// r.HandleFunc("/fetch", scraper.SHandlerMock).Methods("GET")

	/*
		Req:
			POST 
			URL: localhost:8000/
			Header: 
			Body: 
				'url': 'https://www.amazon.com/LEGO-Building-Awesome-Playset-Creative
					/dp/B07WJJKNRG/ref=gbps_tit_m-9_475e_a8951db7?smid=ATVPDKIKX0DER
					&pf_rd_p=5d86def2-ec10-4364-9008-8fbccf30475e
					&pf_rd_s=merchandised-search-9&pf_rd_t=101&pf_rd_i=15529609011
					&pf_rd_m=ATVPDKIKX0DER&pf_rd_r=E5FN5TEKCG3CDYMAR5CC'

		Resp:
			Header
			REST Response Code
			REST response body
	*/
	// r.HandleFunc("/resources/", scraper.SHandler).Methods("POST")
	r.HandleFunc("/products/", scraper.CreateDBRecordHandler).Methods("POST")
	r.HandleFunc("/products/", scraper.GetDBRecordsHandler).Methods("GET")
	r.HandleFunc("/resources/", scraper.ScrapeHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(":5000", r))

	// Connect to MongoDB
}