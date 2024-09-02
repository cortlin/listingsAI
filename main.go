package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func PrettyString(str string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}

func main() {
	url := "https://www.talktotucker.com/api/v1/search/search"

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
	body, _ := io.ReadAll(response.Body)

	jsonBody, err := PrettyString(string(body))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("response Body:", jsonBody)
}
