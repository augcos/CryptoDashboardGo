package binancePipeline

import (
	"fmt"
	"sort"
)

// twoSlices is the auxiliar struct type used to order the results
type twoSlices struct {
    stringSlice  []string
    floatSlice  []float32
}

// an interface is created to modify the behaviour of the sort package and order
// a slice of strings according to a slice of float32
type sortByOther twoSlices
func (sbo sortByOther) Len() int {
    return len(sbo.stringSlice)
}
func (sbo sortByOther) Swap(i, j int) {
    sbo.stringSlice[i], sbo.stringSlice[j] = sbo.stringSlice[j], sbo.stringSlice[i]
    sbo.floatSlice[i], sbo.floatSlice[j] = sbo.floatSlice[j], sbo.floatSlice[i]
}
func (sbo sortByOther) Less(i, j int) bool {
    return sbo.floatSlice[i] > sbo.floatSlice[j] 
}

// BoardResult is the auxiliar struct type used to process the results
type BoardResult struct {
	NumAltcoin float32
	InvestedBTC float32
	CurrentBTC float32
	AvgPrice float32
	CurrentPrice float32
	GainedBTC float32
	GainedRatio float32
	LastPurchasingPrice float32
	LastSellingPrice float32
	ConsolidatedGains float32
}

// OrderPairs(): returns a slice with the altcoin pairs oredered according to the selected parameter
func OrderPairs(orderType string, dashboard map[string]*BoardResult, tradedPairs []string) []string{
	values := make([]float32, 0, len(tradedPairs))
	switch orderType {
		case "Invested BTC":
			for  _,pair := range tradedPairs {
				values = append(values, dashboard[pair].InvestedBTC)
			}
		case "Current BTC":
			for  _,pair := range tradedPairs {
				values = append(values, dashboard[pair].CurrentBTC)
			}
		case "Avg Price":
			for  _,pair := range tradedPairs {
				values = append(values, dashboard[pair].AvgPrice)
			}
		case "Current Price":
			for  _,pair := range tradedPairs {
				values = append(values, dashboard[pair].CurrentPrice)
			}
		case "Gained BTC":
			for  _,pair := range tradedPairs {
				values = append(values, dashboard[pair].GainedBTC)
			}
		case "Gained %":
			for  _,pair := range tradedPairs {
				values = append(values, dashboard[pair].GainedRatio)
			}
		default:
			for  _,pair := range tradedPairs {
				values = append(values, dashboard[pair].InvestedBTC)
			}
	}
	pairedSlices := twoSlices{stringSlice: tradedPairs, floatSlice: values}
	sort.Sort(sortByOther(pairedSlices))
	return pairedSlices.stringSlice
}

// GetInterface(): returns an interface used to create the printed table
func GetInterface(result *BoardResult, pair string) ([]interface{}) {
    row := []interface{}{
		pair,
		fmt.Sprintf("%.8f",result.InvestedBTC),
		fmt.Sprintf("%.8f",result.CurrentBTC),
		fmt.Sprintf("%.8f",result.AvgPrice),
		fmt.Sprintf("%.8f",result.CurrentPrice),
		fmt.Sprintf("%.8f",result.GainedBTC),
		fmt.Sprint(fmt.Sprintf("%.2f",result.GainedRatio) + " %"),
    }
	return row
}