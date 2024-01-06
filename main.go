package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"encoding/json"
	"time"
	"github.com/fatih/color"
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

	portfolioList := []Portfolio{}
	err = json.Unmarshal(body, &portfolioList)
	if err != nil {
		log.Fatal(err)
	}

	// Sort portfolio by Ppl

	// read from private/data.json file and unmarshal into portfolioList


	// today := time.Now().Format("YYMMDD")
	today := time.Now().Format("060102")
	
	fileName := "data-" + today + ".json"
	f, err := os.Open("private/" + fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// portfolioList := []Portfolio{}
	err = json.NewDecoder(f).Decode(&portfolioList)
	if err != nil {
		log.Fatal(err)
	}

	sortPortfolioByPpl(portfolioList)
	// Print portfolio with Ticker, Current Price, and Ppl (Ppl displayed in yellow)
	for _, p := range portfolioList {
		fmt.Printf("%s  \t  %.2f \t", p.Ticker, p.CurrentPrice)
		if p.Ppl > 0 {
			color.Set(color.FgGreen)
		} else if p.Ppl < 0 {
			color.Set(color.FgRed)
		} else {
			color.Set(color.FgWhite)
		}
		fmt.Printf("Ppl: £%.2f\n", p.Ppl)
		color.Unset() // Resets to default colors

	}

	fmt.Printf("Total Ppl: $%.2f\n", sumPpl(portfolioList))
	// save body to file json
	f2, err2 := os.Create("private/"+fileName)
	if err2 != nil {
		panic(err2)
	}

	err2 = json.NewEncoder(f2).Encode(portfolioList)
	if err2 != nil {
		log.Fatal(err2)
	}

	defer f2.Close()
}

// Calculate the sum of Ppl values in the portfolioList
func sumPpl(portfolioList []Portfolio) float64 {
	sum := 0.0
	for _, p := range portfolioList {
		sum += p.Ppl
	}
	return sum
}

// Sort portfolio by Ppl
func sortPortfolioByPpl(portfolio []Portfolio) {
	if len(portfolio) <= 1 {
		return
	}

	mid := len(portfolio) / 2
	left := make([]Portfolio, mid)
	right := make([]Portfolio, len(portfolio)-mid)

	copy(left, portfolio[:mid])
	copy(right, portfolio[mid:])

	sortPortfolioByPpl(left)
	sortPortfolioByPpl(right)

	merge(portfolio, left, right)
}

func merge(portfolio, left, right []Portfolio) {
	i, j, k := 0, 0, 0

	for i < len(left) && j < len(right) {
		if left[i].Ppl > right[j].Ppl {
			portfolio[k] = left[i]
			i++
		} else {
			portfolio[k] = right[j]
			j++
		}
		k++
	}

	for i < len(left) {
		portfolio[k] = left[i]
		i++
		k++
	}

	for j < len(right) {
		portfolio[k] = right[j]
		j++
		k++
	}
}
