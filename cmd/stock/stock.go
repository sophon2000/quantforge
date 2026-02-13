package main

import (
	"fmt"
	"time"

	"quantforge/broker/ibkr"
)

func main() {
	// 创建 IBKR 客户端
	client := ibkr.NewClient()
	client.IB().SetClientLogLevel(0)

	// 配置连接参数
	config := &ibkr.Config{
		Host:     "10.30.8.154", // TWS/IB Gateway 地址
		Port:     4001,          // TWS 默认端口（IB Gateway 为 4001/4002）
		ClientID: 0,             // 客户端 ID，0 表示接收手动订单
		Timeout:  30 * time.Second,
	}
	// 连接到 IBKR
	err := client.Connect(config)
	if err != nil {
		fmt.Printf("连接失败: %v\n", err)
		return
	}
	defer client.Disconnect()

	fmt.Println("连接成功")
	// 获取管理的账户列表
	accounts := client.ManagedAccounts()
	fmt.Printf("管理的账户: %v\n", accounts)

	// ===== 账户信息示例 =====
	accountValues := client.AccountValues()
	fmt.Printf("账户值数量: %d\n", len(accountValues))
	// for _, item := range accountValues {
	// 	fmt.Printf("账户值: %s, 值: %v, 标签: %s, 货币: %s\n", item.Account, item.Value, item.Tag, item.Currency)
	// }

	portfolio := client.Portfolio()
	fmt.Printf("投资组合项数: %d\n", len(portfolio))
	// for _, item := range portfolio {
	// 	fmt.Printf("投资组合: %s, 数量: %d, 市值: %f\n", item.Contract.Symbol, item.Position, item.MarketValue)
	// }

	// ===== 合约查询示例 =====
	// 查询 SE 股票
	// se := ibkr.NewStock("SE", "SMART", "USD")
	// contractDetails, err := client.ReqContractDetails(se)
	// if err != nil {
	// 	fmt.Printf("查询合约失败: %v\n", err)
	// } else {
	// 	fmt.Printf("找到 SE 合约数量: %d\n", len(contractDetails))
	// }

	// ===== 历史数据示例 =====
	// 获取最近 1 天的 1 分钟 K 线
	// bars := client.GetHistoricalBars(
	// 	se,
	// 	"",       // 空字符串表示当前时间
	// 	"1 D",    // 1 天
	// 	"1 min",  // 1 分钟 K 线
	// 	"TRADES", // 交易数据
	// 	true,     // 只使用常规交易时间
	// )
	// fmt.Printf("获取到 %d 根 K 线\n", len(bars))
	// if len(bars) > 0 {
	// 	lastBar := bars[len(bars)-1]
	// 	fmt.Printf("最新 K 线: 时间=%v, 开=%v, 高=%v, 低=%v, 收=%v, 量=%v\n",
	// 		lastBar.Date, lastBar.Open, lastBar.High, lastBar.Low, lastBar.Close, lastBar.Volume)
	// }

	// // ===== 市场扫描示例 =====
	// // 获取扫描器参数
	// scanParams, err := client.ReqScannerParameters()
	// if err != nil {
	// 	fmt.Printf("获取扫描器参数失败: %v\n", err)
	// }
	// fmt.Printf("扫描器参数: %s\n", scanParams)
	// // 创建扫描器订阅
	// scanSub := ibkr.NewScannerSubscription()
	// scanSub.Instrument = "STK"
	// scanSub.LocationCode = "STK.US.MAJOR"
	// scanSub.ScanCode = "TOP_PERC_GAIN"

	// // 带过滤选项
	// opts := ibkr.ScannerSubscriptionOptions{
	// 	FilterOptions: []ibkr.TagValue{
	// 		{Tag: "changePercAbove", Value: "10"},
	// 		{Tag: "priceAbove", Value: "5"},
	// 	},
	// }

	// scanData, err := client.ReqScannerSubscription(scanSub, opts)
	// if err != nil {
	// 	fmt.Printf("市场扫描失败: %v\n", err)
	// } else {
	// 	fmt.Printf("扫描结果数量: %d\n", len(scanData))
	// }
	// for _, item := range scanData {
	// 	cd := item.ContractDetails
	// 	c := cd.Contract

	// 	se := ibkr.NewStock(c.Symbol, c.Exchange, c.Currency)
	// 	contractDetails, err := client.ReqContractDetails(se)
	// 	if err != nil {
	// 		fmt.Printf("查询合约失败: %v\n", err)
	// 	} else {
	// 		cd2 := contractDetails[0]
	// 		c2 := cd2.Contract

	// 		fmt.Printf("扫描结果: Rank=%d, 股票=%s, 交易所=%s, 货币=%s, 市场=%s, 最小变动=%f, 价格倍数=%d, 全称=%s, 投影=%s\n",
	// 			item.Rank,
	// 			c2.Symbol, c2.Exchange, c2.Currency, cd2.MarketName, cd2.MinTick, cd2.PriceMagnifier, string(cd2.LongName), item.Projection)
	// 	}
	// }

	// ===== 盈亏监控示例 =====
	if len(accounts) > 0 {
		account := accounts[0]

		// 订阅盈亏更新
		client.ReqPnL(account, "")

		// 启动协程接收盈亏更新
		go func() {
			pnlChan := client.PnlChan(account, "")
			for pnl := range pnlChan {
				fmt.Printf("盈亏更新: 每日盈亏=%v, 未实现盈亏=%v, 已实现盈亏=%v\n",
					pnl.DailyPNL, pnl.UnrealizedPnl, pnl.RealizedPNL)
			}
		}()
	}

	select {}
}
