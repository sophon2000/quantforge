package main

import (
	"fmt"
	"time"

	"quantforge/datasource"

	"github.com/sdcoffey/techan"
	"github.com/spf13/cobra"
)

var backtestCmd = &cobra.Command{
	Use:   "backtest",
	Short: "运行策略回测",
	Long:  `使用历史数据与指标运行回测。`,
	RunE:  runBacktest,
}

func init() {
	backtestCmd.Flags().StringP("symbol", "s", "AAPL", "回测标的代码")
}

func runBacktest(cmd *cobra.Command, _ []string) error {
	symbol, _ := cmd.Flags().GetString("symbol")

	rows, err := datasource.GetData("")
	if err != nil {
		return fmt.Errorf("获取数据: %w", err)
	}
	searchRows, ok := rows[symbol]
	if !ok {
		return fmt.Errorf("未找到标的: %s", symbol)
	}

	series, err := datasource.GenerateSeries(searchRows)
	if err != nil {
		return fmt.Errorf("生成序列: %w", err)
	}

	closePrice := techan.NewClosePriceIndicator(series)
	upper, middle, lower := datasource.BollingerBands(closePrice, 20, 2)
	k, d, j := datasource.KDJ(series, 9, 3, 3)

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
	return nil
}
