package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var backtestCmd = &cobra.Command{
	Use:   "backtest",
	Short: "运行策略回测",
	Long:  `使用历史 K 线 + TechanStrategy 回测，统计权益与成交。`,
	RunE:  runBacktest,
}

func init() {
	backtestCmd.Flags().StringP("symbol", "s", "AAPL", "回测标的")
	backtestCmd.Flags().StringP("strategy", "S", "bollinger", "策略: bollinger, bollinger-mean-reversion, macd, macd-histogram, rsi, rsi-divergence, kdj-crossover, kdj-extreme, multi-indicator, trend-following")
	backtestCmd.Flags().Float64P("cash", "c", 100000, "初始资金")
	backtestCmd.Flags().IntP("quantity", "q", 100, "每笔信号下单数量")
}

func runBacktest(cmd *cobra.Command, _ []string) error {
	symbol, _ := cmd.Flags().GetString("symbol")
	strategyName, _ := cmd.Flags().GetString("strategy")
	initialCash, _ := cmd.Flags().GetFloat64("cash")
	quantity, _ := cmd.Flags().GetInt("quantity")

	res, err := RunBacktest(symbol, strategyName, initialCash, quantity)
	if err != nil {
		return err
	}

	s := res.Summary
	fmt.Printf("回测: %s | 策略=%s | 初始资金=%.2f | 每笔=%d\n", res.Symbol, res.Strategy, s.InitialCash, quantity)
	fmt.Println("----------------------------------------")
	for _, sig := range res.Signals {
		fmt.Printf("  [信号] %s %s @ %.2f x %d\n", res.Symbol, sig.Type, sig.Price, sig.Quantity)
	}
	fmt.Println("----------------------------------------")
	fmt.Printf("K 线数: %d | 信号/成交笔数: %d\n", len(res.CategoryData), s.TradeCount)
	fmt.Printf("期末现金: %.2f | 收益率: %.2f%% | 持仓 %s: %d\n", s.FinalCash, s.ReturnPct, res.Symbol, s.Position)
	if s.TradeCount > 0 {
		fmt.Printf("期末权益(简化): %.2f (未计持仓市值)\n", s.FinalCash)
	}
	return nil
}
