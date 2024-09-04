package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type jsonResponse struct {
	Hits struct {
		Hits []jsonResponseInner `json:"hits"`
	} `json:"hits"`
}

type jsonResponseInner struct {
	Source Source `json:"_source"`
}

type Source struct {
	Zip           string  `json:"Zip"`
	StreetName    string  `json:"StreetName"`
	FullBathrooms string  `json:"FullBathrooms"`
	Latitude      float64 `json:"Latitude"`
	Longitude     float64 `json:"Longitude"`
	StreetSuffix  string
	StreetAddress string
	City          string
	State         string
	County        string
	ListingPrice  int
	HalfBathrooms int
	Bedrooms      int
	SquareFeet    int
}

type Listing struct {
	gorm.Model
	zip            string
	street_name    string
	street_suffix  string
	city           string
	state          string
	county         string
	lat            float64
	lng            float64
	listing_price  int
	square_feet    int
	bedrooms       uint8
	bathrooms      string
	half_bathrooms uint8
}

var db, err = gorm.Open(sqlite.Open("mls-ai.db"), &gorm.Config{})

func StringToListings(str string) jsonResponse {
	data := jsonResponse{}
	err := json.Unmarshal([]byte(str), &data)

	if err != nil {
		fmt.Println(err.Error())
	}

	return data
}

func mlsToDB(mls jsonResponse) []Listing {
	arr := mls.Hits.Hits
	r := []Listing{}
	for _, d := range arr {
		r = append(r, Listing{
			zip:            d.Source.Zip,
			street_name:    d.Source.StreetName,
			street_suffix:  d.Source.StreetSuffix,
			city:           d.Source.City,
			state:          d.Source.State,
			county:         d.Source.County,
			lat:            d.Source.Latitude,
			lng:            d.Source.Longitude,
			listing_price:  d.Source.ListingPrice,
			square_feet:    d.Source.SquareFeet,
			bedrooms:       uint8(d.Source.Bedrooms),
			bathrooms:      d.Source.FullBathrooms,
			half_bathrooms: uint8(d.Source.HalfBathrooms),
		})
	}

	return r
}

func main() {
	url := "https://www.talktotucker.com/api/v1/search/search"

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Listing{})

	fmt.Println("HTTP JSON POST URL:", url)

	requestBody := []byte(`{
		"array": {
			"PropertyType":["Residential","Condominiums"],
			"Status":["Active","Pending"],
			"TransType":["Sale"],
			"stateOrProvince":["IN"]
		},
		"gte":{},
		"range":{},
		"shouldarray":{},
		"mustNotarray":{},
		"data":true,
		"getall":true,
		"browser":true
		}
	`)

	client := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	response, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, err := io.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err.Error())
	}

	jsonBody := StringToListings(string(body))

	localData := mlsToDB(jsonBody)
	fmt.Println("INSERTING INTO TABLE")
	if res := db.CreateInBatches(localData, 100); res.Error != nil {
		panic(res.Error)
	}

	fmt.Println("TOTAL LISTINGS: ", len(localData))
}
