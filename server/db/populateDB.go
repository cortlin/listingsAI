package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
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
	Zip           string      `json:"Zip"`
	StreetName    string      `json:"StreetName"`
	FullBathrooms interface{} `json:"FullBathrooms"`
	Latitude      float64     `json:"Latitude"`
	Longitude     float64     `json:"Longitude"`
	StreetSuffix  string
	StreetAddress string
	City          string
	State         string
	County        string
	ListingPrice  int
	HalfBathrooms interface{}
	Bedrooms      int
	SquareFeet    interface{}
}

type Listing struct {
	gorm.Model
	Zip            string
	Street_name    string
	Street_suffix  string
	City           string
	State          string
	County         string
	Lat            float64
	Lng            float64
	Listing_price  int
	Square_feet    string
	Bedrooms       uint8
	Bathrooms      string
	Half_bathrooms string
}

var db, err = gorm.Open(sqlite.Open("mls-ai.db"), &gorm.Config{})

// PropToString function
func PropToString(value interface{}) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	case int:
		return strconv.Itoa(v), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	default:
		fmt.Println(reflect.TypeOf(v))
		return "", fmt.Errorf("unsupported type: %T", value)
	}
}

func StringValidated(v interface{}) string {
	c, err := PropToString(v)
	if err != nil {
		panic(err.Error)
	}
	return c
}

func StringToListings(str string) jsonResponse {
	data := jsonResponse{}
	err := json.Unmarshal([]byte(str), &data)

	if err != nil {
		fmt.Println(err.Error())
	}

	return data
}

func mlsToDB(mls jsonResponse) (r []Listing) {
	arr := mls.Hits.Hits
	for _, d := range arr {
		r = append(r, Listing{
			Zip:            d.Source.Zip,
			Street_name:    d.Source.StreetName,
			Street_suffix:  d.Source.StreetSuffix,
			City:           d.Source.City,
			State:          d.Source.State,
			County:         d.Source.County,
			Lat:            d.Source.Latitude,
			Lng:            d.Source.Longitude,
			Listing_price:  d.Source.ListingPrice,
			Square_feet:    StringValidated(d.Source.SquareFeet),
			Bedrooms:       uint8(d.Source.Bedrooms),
			Bathrooms:      StringValidated(d.Source.FullBathrooms),
			Half_bathrooms: StringValidated(d.Source.HalfBathrooms),
		})
	}
	return
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
