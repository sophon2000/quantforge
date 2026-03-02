package main

import (
	"fmt"
	"time"

	"github.com/sophon2000/quantforge/broker/ibkr"

	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "运行市场扫描（IBKR）",
	Long:  `连接 TWS/IB Gateway 并执行市场扫描，如涨幅榜等。`,
	RunE:  runScan,
}

var (
	scanHost string
	scanPort int
	scanID   int
)

func init() {
	scanCmd.Flags().StringVar(&scanHost, "host", "10.30.8.154", "TWS/IB Gateway 地址")
	scanCmd.Flags().IntVar(&scanPort, "port", 4001, "端口")
	scanCmd.Flags().IntVar(&scanID, "client-id", 0, "客户端 ID")
}

func runScan(_ *cobra.Command, _ []string) error {
	client := ibkr.NewClient()
	client.IB().SetClientLogLevel(0)

	config := &ibkr.Config{
		Host:     scanHost,
		Port:     scanPort,
		ClientID: scanID,
		Timeout:  30 * time.Second,
	}
	if err := client.Connect(config); err != nil {
		return fmt.Errorf("连接失败: %w", err)
	}
	defer client.Disconnect()

	scanParams, err := client.ReqScannerParameters()
	if err != nil {
		return fmt.Errorf("获取扫描器参数: %w", err)
	}
	_ = scanParams

	scanSub := ibkr.NewScannerSubscription()
	scanSub.Instrument = "STK"
	scanSub.LocationCode = "STK.US.MAJOR"
	scanSub.ScanCode = "TOP_PERC_GAIN"
	opts := ibkr.ScannerSubscriptionOptions{
		FilterOptions: []ibkr.TagValue{
			{Tag: "changePercAbove", Value: "10"},
			{Tag: "priceAbove", Value: "5"},
		},
	}

	scanData, err := client.ReqScannerSubscription(scanSub, opts)
	if err != nil {
		return fmt.Errorf("市场扫描: %w", err)
	}

	fmt.Printf("扫描结果数量: %d\n", len(scanData))
	for _, item := range scanData {
		cd := item.ContractDetails
		c := cd.Contract
		se := ibkr.NewStock(c.Symbol, c.Exchange, c.Currency)
		contractDetails, err := client.ReqContractDetails(se)
		if err != nil {
			fmt.Printf("查询合约失败 %s: %v\n", c.Symbol, err)
			continue
		}
		if len(contractDetails) == 0 {
			continue
		}
		cd2 := contractDetails[0]
		c2 := cd2.Contract
		fmt.Printf("Rank=%d  %s  %s  %s  %s  MinTick=%.4f  全称=%s  投影=%s\n",
			item.Rank, c2.Symbol, c2.Exchange, c2.Currency, cd2.MarketName, cd2.MinTick, string(cd2.LongName), item.Projection)
	}

	time.Sleep(2 * time.Second) // 等待结果收齐
	return nil
}
