package binancePipeline

import (
	"os"
	"fmt"
	"strconv"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

// pairPrice is the auxiliar struct type used to process the JSON data
type pairPrice struct {
	Pair string `json:"symbol"`
	CurrentPrice string `json:"lastPrice"`
}

// GetData(): downloads and processes the price data from Binance's API
func GetData() map[string]float32{
	// does a GET request to the Binance API
	url := "https://api.binance.com/api/v3/ticker/24hr"
	resp, err := http.Get(url)
	if err != nil {
	   Exit("Unable to connect to Binance. Please check your internet connection.\n")
	}

	// reads the body of the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Exit("An error has occurred with Binance's data.\n")
	}
	
	// transforms the JSON data to an array of maps
	var priceSlice []pairPrice
	err = json.Unmarshal(body, &priceSlice)
	if err != nil {
		Exit("An error has occurred with Binance's data.\n")
	}

	// transforms the array of maps to a map between each pair and its price
	priceMap := make(map[string]float32)
	for _, pair := range priceSlice {
		valuePrice, _ := strconv.ParseFloat(pair.CurrentPrice, 32)
		priceMap[pair.Pair] = float32(valuePrice)
	}
	return priceMap
}



// Exit(): prints a messsage and exits the execution
func Exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}