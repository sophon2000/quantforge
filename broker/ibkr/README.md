# IBKR 模块

这是对 Interactive Brokers（盈透证券）API 的 Go 语言封装，基于 [ibsync](https://github.com/scmhub/ibsync) 库。

## 功能特性

### 1. 连接管理
- 简化的连接配置
- 自动重连支持
- 账户管理

### 2. 账户操作 (`account.go`)
- ✅ 账户值查询
- ✅ 账户摘要
- ✅ 投资组合
- ✅ 持仓查询
- ✅ 交易记录

### 3. 合约操作 (`contract.go`)
- ✅ 股票合约
- ✅ 外汇合约
- ✅ 期权合约
- ✅ 期货合约
- ✅ 合约详情查询

### 4. 历史数据 (`bar_data.go`)
- ✅ 历史 K 线数据
- ✅ 实时更新的历史数据
- ✅ 实时 Bar 数据（5秒）
- ✅ 辅助函数获取数据切片

### 5. 实时行情 (`ticket_data.go`)
- ✅ 市场快照
- ✅ 逐笔数据（Tick by Tick）
- ✅ 历史逐笔数据
- ✅ 实时市场数据订阅

### 6. 订单管理 (`order.go`)
- ✅ 限价单
- ✅ 市价单
- ✅ 止损单
- ✅ 止损限价单
- ✅ 下单
- ✅ 撤单
- ✅ 全局撤单

### 7. 市场扫描 (`scanner.go`)
- ✅ 市场扫描器参数
- ✅ 市场扫描订阅
- ✅ 过滤选项

### 8. 盈亏监控 (`pnl.go`)
- ✅ 账户盈亏订阅
- ✅ 单个持仓盈亏
- ✅ 实时盈亏更新

## 安装

```bash
go get github.com/scmhub/ibapi
go get github.com/scmhub/ibsync
```

## 快速开始

### 1. 基础连接

```go
import "itick/ibkr"

func main() {
    // 创建客户端
    client := ibkr.NewClient()
    
    // 配置连接
    config := &ibkr.Config{
        Host:     "127.0.0.1",
        Port:     7497,        // TWS: 7497, IB Gateway: 4001/4002
        ClientID: 0,
        Timeout:  30 * time.Second,
    }
    
    // 连接
    err := client.Connect(config)
    if err != nil {
        panic(err)
    }
    defer client.Disconnect()
    
    // 获取账户
    accounts := client.ManagedAccounts()
    fmt.Println("账户:", accounts)
}
```

### 2. 查询合约

```go
// 股票
amd := ibkr.NewStock("AMD", "SMART", "USD")
details, err := client.ReqContractDetails(amd)

// 外汇
eurusd := ibkr.NewForex("EUR", "IDEALPRO", "USD")

// 期权
spy_call := ibkr.NewOption("SPY", "SMART", "20240315", "C", 400.0, "USD")
```

### 3. 获取历史数据

```go
// 获取最近 1 天的 1 分钟 K 线
bars, err := client.GetHistoricalBars(
    amd,
    time.Now(),
    "1 D",      // 持续时间: "1 D", "1 W", "1 M", "1 Y"
    "1 min",    // K线大小: "1 min", "5 mins", "1 hour", "1 day"
    "TRADES",   // 数据类型: "TRADES", "MIDPOINT", "BID", "ASK"
    true,       // 只使用常规交易时间
)

for _, bar := range bars {
    fmt.Printf("时间=%v, 开=%v, 高=%v, 低=%v, 收=%v, 量=%v\n",
        bar.Date, bar.Open, bar.High, bar.Low, bar.Close, bar.Volume)
}
```

### 4. 实时数据

```go
// 获取快照
snapshot, err := client.Snapshot(amd)
fmt.Printf("当前价格: %v\n", snapshot.MarketPrice())

// 订阅实时数据
ticker := client.ReqMktData(amd)
// ... 使用 ticker
client.CancelMktData(amd)
```

### 5. 下单

```go
// 创建限价单
order := ibkr.LimitOrder("BUY", ibkr.StringToDecimal("100"), 150.0)

// 下单
trade := client.PlaceOrder(amd, order)

// 等待完成
<-trade.Done()
fmt.Println("交易完成!")

// 撤单
client.CancelOrder(order, ibkr.NewOrderCancel())

// 撤销所有订单
client.ReqGlobalCancel()
```

### 6. 账户信息

```go
// 账户值
accountValues := client.AccountValues()

// 投资组合
portfolio := client.Portfolio()

// 持仓订阅
client.ReqPositions()
posChan := client.PositionChan()
for pos := range posChan {
    fmt.Printf("持仓: %s, 数量: %v\n", pos.Contract.Symbol, pos.Position)
}
```

### 7. 盈亏监控

```go
account := "DU123456"

// 订阅盈亏
client.ReqPnL(account, "")

// 接收更新
pnlChan := client.PnlChan(account, "")
for pnl := range pnlChan {
    fmt.Printf("每日盈亏=%v, 未实现=%v, 已实现=%v\n",
        pnl.DailyPnL, pnl.UnrealizedPnL, pnl.RealizedPnL)
}
```

### 8. 市场扫描

```go
// 创建扫描器
scanSub := ibkr.NewScannerSubscription()
scanSub.Instrument = "STK"
scanSub.LocationCode = "STK.US.MAJOR"
scanSub.ScanCode = "TOP_PERC_GAIN"

// 添加过滤器
opts := ibkr.ScannerSubscriptionOptions{
    FilterOptions: []ibkr.TagValue{
        {Tag: "changePercAbove", Value: "10"},
        {Tag: "priceAbove", Value: "5"},
    },
}

// 执行扫描
scanData, err := client.ReqScannerSubscription(scanSub, opts)
```

## API 文档

### Client 结构体

```go
type Client struct {
    // 私有字段
}
```

### 配置结构体

```go
type Config struct {
    Host     string        // TWS/Gateway 地址
    Port     int          // 端口号
    ClientID int          // 客户端 ID
    Timeout  time.Duration // 超时时间
}
```

### 主要方法

#### 连接管理
- `NewClient() *Client` - 创建新客户端
- `Connect(config *Config) error` - 连接
- `Disconnect()` - 断开连接
- `ManagedAccounts() []string` - 获取账户列表

#### 账户操作
- `AccountValues() []ibsync.AccountValue`
- `AccountSummary() []ibsync.AccountValue`
- `Portfolio() []ibsync.PortfolioItem`
- `ReqPositions()`
- `PositionChan() <-chan ibsync.Position`
- `Trades() []ibsync.Trade`
- `OpenTrades() []ibsync.Trade`

#### 合约操作
- `ReqContractDetails(*ibapi.Contract) ([]ibsync.ContractDetails, error)`
- `NewStock(symbol, exchange, currency string) *ibapi.Contract`
- `NewForex(base, exchange, quote string) *ibapi.Contract`
- `NewOption(...) *ibapi.Contract`
- `NewFuture(...) *ibapi.Contract`

#### 数据获取
- `GetHistoricalBars(...) ([]ibsync.Bar, error)`
- `ReqHistoricalData(...) (<-chan ibsync.Bar, error)`
- `ReqRealTimeBars(...) (<-chan ibsync.RealTimeBar, context.CancelFunc)`
- `Snapshot(*ibapi.Contract) (*ibsync.Ticker, error)`
- `ReqMktData(*ibapi.Contract) *ibsync.Ticker`

#### 订单管理
- `PlaceOrder(*ibapi.Contract, *ibapi.Order) *ibsync.Trade`
- `CancelOrder(*ibapi.Order, *ibapi.OrderCancel)`
- `ReqGlobalCancel()`
- `LimitOrder(action string, quantity Decimal, limitPrice float64) *ibapi.Order`
- `MarketOrder(action string, quantity Decimal) *ibapi.Order`

#### 盈亏监控
- `ReqPnL(account, modelCode string)`
- `PnlChan(account, modelCode string) <-chan ibsync.PnL`
- `Pnl(account, modelCode string) *ibsync.PnL`

## 注意事项

1. **连接前提**
   - 确保 TWS 或 IB Gateway 正在运行
   - 在 TWS 中启用 API 连接（文件 -> 全局配置 -> API -> 设置）
   - 配置允许的 IP 地址和端口

2. **端口说明**
   - TWS 默认端口: 7497（实盘）、7496（模拟）
   - IB Gateway 默认端口: 4001（实盘）、4002（模拟）

3. **ClientID**
   - 每个连接需要唯一的 ClientID
   - ClientID=0 可以接收手动下单的订单

4. **数据权限**
   - 某些数据需要订阅市场数据
   - 实时数据可能有延迟（非订阅用户）

5. **风险提示**
   - 这是真实交易 API，请谨慎操作
   - 建议先在模拟账户测试
   - 生产环境务必添加错误处理和风控逻辑

## 完整示例

参考项目根目录的 `ibkr_example_main.go` 文件，包含所有功能的完整示例。

## 依赖

- [github.com/scmhub/ibapi](https://github.com/scmhub/ibapi) - IBKR API Go 绑定
- [github.com/scmhub/ibsync](https://github.com/scmhub/ibsync) - IBKR API 同步封装

## 许可证

与主项目相同
