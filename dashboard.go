package main

// Package imports
import (
	"os"
	"fmt"
	"flag"
	"time"
	"sort"
	"encoding/csv"

	pl "github.com/augcos/CryptoDashboard/binancePipeline"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func main() {
	// defines all the flags used in the code
	csvFilename := flag.String("csv", "myOrders.csv", "csv file with the Binance order history")
	refresh := flag.Int("refresh", 60, "refresh time")
	orderType := flag.String("order", "Invested BTC", "altcoin display order of the results")
	flag.Parse()

	// opens the csv file
	file, err := os.Open(*csvFilename)
	if err != nil {
		pl.Exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilename))
	}

	// reads the data from the csv
	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		pl.Exit("Failed to parse the provided CSV file.\n")
	}

	// parses the data and returns the data in map format
	myParsedOrders := pl.ParseOrders(lines)
	tradedPairs := pl.UniquePairs(myParsedOrders)

	// initializes the dashboard values to zero
	dashboard := make(map[string]*pl.BoardResult)	
	for _,pair := range tradedPairs {
		dashboard[pair] = &pl.BoardResult{0,0,0,0,0,0,0,0,0,0}
	}
	
	// analyzes the orders in chronological order
	for _, nextOrder := range myParsedOrders {
		if nextOrder.Side=="BUY" {
			dashboard[nextOrder.Pair].NumAltcoin += nextOrder.Altcoin
			dashboard[nextOrder.Pair].InvestedBTC += nextOrder.Bitcoin
			dashboard[nextOrder.Pair].AvgPrice = dashboard[nextOrder.Pair].InvestedBTC/dashboard[nextOrder.Pair].NumAltcoin
			dashboard[nextOrder.Pair].LastPurchasingPrice = nextOrder.Bitcoin/nextOrder.Altcoin
		} else if nextOrder.Side=="SELL" {
			ratioAlt := nextOrder.Altcoin/dashboard[nextOrder.Pair].NumAltcoin
			dashboard[nextOrder.Pair].ConsolidatedGains += nextOrder.Bitcoin-ratioAlt*dashboard[nextOrder.Pair].InvestedBTC

			dashboard[nextOrder.Pair].NumAltcoin -= nextOrder.Altcoin
			dashboard[nextOrder.Pair].InvestedBTC -= ratioAlt*dashboard[nextOrder.Pair].InvestedBTC
			dashboard[nextOrder.Pair].GainedBTC = dashboard[nextOrder.Pair].CurrentBTC-dashboard[nextOrder.Pair].InvestedBTC
			dashboard[nextOrder.Pair].AvgPrice = dashboard[nextOrder.Pair].InvestedBTC/dashboard[nextOrder.Pair].NumAltcoin
			dashboard[nextOrder.Pair].LastSellingPrice = nextOrder.Bitcoin/nextOrder.Altcoin
		}
	}

	// zeros the values of the altcoin pairs that are out of the scope of the code (e.g. Binance Launchpad tokens)
	for _,pair := range tradedPairs {
		if dashboard[pair].InvestedBTC != dashboard[pair].InvestedBTC || dashboard[pair].InvestedBTC<0 {
			dashboard[pair] = &pl.BoardResult{0,0,0,0,0,0,0,0,0,0}
		}
	}

	// stats the loop to update the price
	var investedBTC,currentBTC,currentGains,ratioGain,consolidatedGains float32
	for {
		// downloads the price data 
		priceData := pl.GetData()
		fmt.Println("\033[2J")

		// updates the fields dependant on price and the total values
		investedBTC,currentBTC,currentGains,ratioGain,consolidatedGains = 0,0,0,0,0
		for _,pair := range tradedPairs {
			dashboard[pair].CurrentPrice = priceData[pair]
			dashboard[pair].CurrentBTC = dashboard[pair].NumAltcoin*dashboard[pair].CurrentPrice
			dashboard[pair].GainedBTC = dashboard[pair].CurrentBTC-dashboard[pair].InvestedBTC
			dashboard[pair].GainedRatio = float32(100)*dashboard[pair].GainedBTC/dashboard[pair].InvestedBTC
			
			investedBTC += dashboard[pair].InvestedBTC
			currentBTC += dashboard[pair].CurrentBTC
			currentGains += dashboard[pair].GainedBTC
			consolidatedGains += dashboard[pair].ConsolidatedGains
			ratioGain = float32(100)*currentGains/investedBTC
		}
		
		// prints a summary of the altcoin balances
		fmt.Printf("\n**************************************** CryptoBoard ****************************************\n")
		fmt.Printf("BTC Price: %.0f EUR / %.0f USD\n", priceData["BTCEUR"], priceData["BTCUSDT"])
		fmt.Printf("Invested BTC: %.4f BTC - Current BTC: %.4f BTC - Current gains: %.4f BTC - Current gain: %.2f ", 
			investedBTC,currentBTC,currentGains,ratioGain)
		fmt.Println("%")
		fmt.Printf("Consolidated Gains: %.4f BTC\n", consolidatedGains)
		fmt.Print("\n")

		// creates an empty table to print the results
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Pair", "Invested BTC", "Current BTC", "Avg Price", "Current Price", "Gained BTC", "Gained %"})

		// orders the results according to the selected parameter in the flags
		var orderedPairs []string = tradedPairs
		if *orderType=="Pair" {
			sort.Strings(orderedPairs)
		} else {
			orderedPairs = pl.OrderPairs(*orderType, dashboard, tradedPairs)
		}

		// adds the altcoin pairs (over 10eur) to the table
		for _,pair := range orderedPairs {
			if dashboard[pair].InvestedBTC>10/priceData["BTCEUR"] {
				row := pl.GetInterface(dashboard[pair], pair)
				t.AppendRow(row)
			}
		}
		// modifies the settings of the table style and prints
		t.AppendSeparator()
		t.SetColumnConfigs([]table.ColumnConfig{
			{Number: 1, Align: text.AlignLeft, AlignHeader: text.AlignCenter},
			{Number: 2, Align: text.AlignRight, AlignHeader: text.AlignCenter},
			{Number: 3, Align: text.AlignRight, AlignHeader: text.AlignCenter},
			{Number: 4, Align: text.AlignRight, AlignHeader: text.AlignCenter},
			{Number: 5, Align: text.AlignRight, AlignHeader: text.AlignCenter},
			{Number: 6, Align: text.AlignRight, AlignHeader: text.AlignCenter},
			{Number: 7, Align: text.AlignRight, AlignHeader: text.AlignCenter},
		})
		t.Render()
		// the table refreshes according to the refresh parameter in the flags
		time.Sleep(time.Duration(*refresh) * time.Second)
	}
}