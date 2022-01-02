package binancePipeline

import (
	"strings"
	"strconv"
)

// order is the auxiliar struct type used to process the order history
type order struct {
	Pair string
	Bitcoin float32
	Altcoin float32
	Side string
	Status string
}

// ParseOrders(): processes the order history and return an order slice
func ParseOrders(lines [][]string) []order{
	// gets separately the column names and the rest of the data
	columnNames := lines[0]
	data := lines[1:]

	// returns the index for each column name in a map
	columnIdx := make(map[string]int)
	for i, name := range columnNames {
		columnIdx[name] = i
	}

	// removes the orders that are not completed or are not made from BTC
	cleanData := make([][]string,0)
	for _,row := range data {
		if strings.Contains(row[columnIdx["Pair"]],"BTC") && strings.Contains(row[columnIdx["Status"]],"FILLED") &&
			!strings.Contains(row[columnIdx["Pair"]],"USD") && !strings.Contains(row[columnIdx["Pair"]],"EUR"){
			cleanData = append(cleanData,row)
		}
	}

	// transforms the CSV rows into an array of maps in inverse order (oldest to most recent)
	orders := make([]order, len(cleanData))
	for i, row := range cleanData {
		valueBTC, errBTC := strconv.ParseFloat(strings.Replace(row[columnIdx["Trading total"]][0:12],",","",-1), 32)
		valueAlt, errAlt := strconv.ParseFloat(strings.Replace(row[columnIdx["Executed"]][0:12],",","",-1), 32)
		if errBTC!=nil || errAlt!=nil {
			Exit("Failed to convert order amount to float32. Data is not in the appropiate format.\n")
		}
		orders[len(cleanData)-1-i] = order{
			Pair: row[columnIdx["Pair"]],
			Bitcoin: float32(valueBTC),
			Altcoin: float32(valueAlt),
			Side: row[columnIdx["Side"]],
			Status: row[columnIdx["Status"]],
		}
	}
	return orders
}

// UniquePairs(): returns the unique pairs found in the order history
func UniquePairs(slice []order) []string {
    listed := make(map[string]bool)
    uniqueList := make([]string,0)
    for _,row := range slice {
        if _,ok := listed[row.Pair]; !ok {
            listed[row.Pair] = true
            uniqueList = append(uniqueList, row.Pair)
        }
    }    
    return uniqueList
}