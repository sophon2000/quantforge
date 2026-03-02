package main

import (
	"fmt"
	"time"

	"github.com/sophon2000/quantforge/broker/ibkr"

	"github.com/spf13/cobra"
)

var liveCmd = &cobra.Command{
	Use:   "live",
	Short: "连接实盘/模拟（IBKR）",
	Long:  `连接 TWS/IB Gateway，展示账户、持仓与盈亏等。`,
	RunE:  runLive,
}

var (
	liveHost string
	livePort int
	liveID   int
)

func init() {
	liveCmd.Flags().StringVar(&liveHost, "host", "10.30.8.154", "TWS/IB Gateway 地址")
	liveCmd.Flags().IntVar(&livePort, "port", 4001, "端口")
	liveCmd.Flags().IntVar(&liveID, "client-id", 0, "客户端 ID")
}

func runLive(_ *cobra.Command, _ []string) error {
	client := ibkr.NewClient()
	client.IB().SetClientLogLevel(0)

	config := &ibkr.Config{
		Host:     liveHost,
		Port:     livePort,
		ClientID: liveID,
		Timeout:  30 * time.Second,
	}
	if err := client.Connect(config); err != nil {
		return fmt.Errorf("连接失败: %w", err)
	}
	defer client.Disconnect()

	fmt.Println("连接成功")
	accounts := client.ManagedAccounts()
	fmt.Printf("管理的账户: %v\n", accounts)

	accountValues := client.AccountValues()
	fmt.Printf("账户值数量: %d\n", len(accountValues))
	portfolio := client.Portfolio()
	fmt.Printf("投资组合项数: %d\n", len(portfolio))

	if len(accounts) > 0 {
		account := accounts[0]
		client.ReqPnL(account, "")
		go func() {
			pnlChan := client.PnlChan(account, "")
			for pnl := range pnlChan {
				fmt.Printf("盈亏更新: 每日盈亏=%v, 未实现盈亏=%v, 已实现盈亏=%v\n",
					pnl.DailyPNL, pnl.UnrealizedPnl, pnl.RealizedPNL)
			}
		}()
	}

	fmt.Println("按 Ctrl+C 退出...")
	select {}
}
