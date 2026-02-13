package main

import (
	"fmt"
	"quantforge/indicator"
	"time"

	"github.com/sdcoffey/techan"
)

func main() {

	rows, err := indicator.GetData()
	if err != nil {
		fmt.Println(err)
	}
	searchSymbol := "AAPL"
	searchRows := rows[searchSymbol]

	series, err := indicator.GenerateSeries(searchRows)
	if err != nil {
		fmt.Println(err)
	}

	closePrice := techan.NewClosePriceIndicator(series)
	upper, middle, lower := indicator.BollingerBands(closePrice, 20, 2)

	k, d, j := indicator.KDJ(series, 9, 3, 3)
	for i := 0; i < 1000; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("k", k.Calculate(i).FormattedString(2))
		fmt.Println("d", d.Calculate(i).FormattedString(2))
		fmt.Println("j", j.Calculate(i).FormattedString(2))
		fmt.Println("closePrice", closePrice.Calculate(i).FormattedString(2))
		fmt.Println("upper", upper.Calculate(i).FormattedString(2))
		fmt.Println("middle", middle.Calculate(i).FormattedString(2))
		fmt.Println("lower", lower.Calculate(i).FormattedString(2))
		fmt.Println("--------------------------------")
	}

}
