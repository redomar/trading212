package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"encoding/json"

	"github.com/joho/godotenv"
)

type Asset struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type Portfolio struct {
	AveragePrice    float64 `json:"averagePrice"`
	CurrentPrice    float64 `json:"currentPrice"`
	Frontend        string  `json:"frontend"`
	FxPpl           float64 `json:"fxPpl"`
	InitialFillDate string  `json:"initialFillDate"`
	MaxBuy          float64 `json:"maxBuy"`
	MaxSell         float64 `json:"maxSell"`
	PieQuantity     float64 `json:"pieQuantity"`
	Ppl             float64 `json:"ppl"`
	Quantity        float64 `json:"quantity"`
	Ticker          string  `json:"ticker"`
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	reqUrl := "https://live.trading212.com/api/v0/equity/portfolio"
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		panic(err)
	}

	apiKey := os.Getenv("TRADING212_API_KEY")
	if apiKey == "" {
		log.Fatal("API key not set in .env file")
	}

	req.Header.Add("Authorization", apiKey)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
	fmt.Println(res.Status)
	fmt.Println("")
	fmt.Println(string(body))

	portfolioList := []Portfolio{}
	err = json.Unmarshal(body, &portfolioList)
	if err != nil {
		log.Fatal(err)
	}

	// save body to file json
	f, err := os.Create("private/data.json")
	if err != nil {
		panic(err)
	}

	err = json.NewEncoder(f).Encode(portfolioList)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

}
