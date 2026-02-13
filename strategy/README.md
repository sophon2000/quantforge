# 策略库

基于常用技术指标的交易策略实现。

## 策略列表

### 1. 布林带策略

#### BollingerBandStrategy - 布林带突破策略
- **买入信号**: 价格从下轨下方上穿下轨（超卖反弹）
- **卖出信号**: 价格从上轨上方下穿上轨（超买回落）
- **适用场景**: 震荡市场
- **参数**: period=20, stdDev=2.0

#### BollingerBandMeanReversionStrategy - 布林带均值回归策略
- **买入信号**: 价格触及下轨
- **卖出信号**: 价格回到中轨或触及上轨
- **适用场景**: 震荡市场，价格趋向均值回归
- **参数**: period=20, stdDev=2.0

### 2. MACD 策略

#### MACDCrossoverStrategy - MACD 金叉死叉策略
- **买入信号**: MACD 线上穿信号线（金叉）
- **卖出信号**: MACD 线下穿信号线（死叉）
- **适用场景**: 趋势市场
- **参数**: fast=12, slow=26, signal=9

#### MACDHistogramStrategy - MACD 柱状图策略
- **买入信号**: 柱状图由负转正（上穿零轴）
- **卖出信号**: 柱状图由正转负（下穿零轴）
- **适用场景**: 趋势转折点捕捉
- **参数**: fast=12, slow=26, signal=9

### 3. RSI 策略

#### RSIStrategy - RSI 超买超卖策略
- **买入信号**: RSI 上穿超卖线（默认 30）
- **卖出信号**: RSI 下穿超买线（默认 70）
- **适用场景**: 震荡市场，捕捉超买超卖反转
- **参数**: period=14, oversold=30, overbought=70

#### RSIDivergenceStrategy - RSI 背离策略（简化版）
- **买入信号**: RSI < 30 且开始上升
- **卖出信号**: RSI > 70 且开始下降
- **适用场景**: 寻找价格与指标的背离信号
- **参数**: period=14

### 4. KDJ 策略

#### KDJCrossoverStrategy - KDJ 金叉死叉策略
- **买入信号**: K 线上穿 D 线，且 J < 20
- **卖出信号**: K 线下穿 D 线，且 J > 80
- **适用场景**: 短线交易，捕捉快速反转
- **参数**: period=9, smoothK=3, smoothD=3

#### KDJOversoldOverboughtStrategy - KDJ 超买超卖策略
- **买入信号**: J 值上穿 0
- **卖出信号**: J 值下穿 100
- **适用场景**: 极端超买超卖情况
- **参数**: period=9, smoothK=3, smoothD=3

### 5. 组合策略

#### MultiIndicatorStrategy - 多指标组合策略
- **买入信号**: MACD 金叉 + RSI < 50 + 价格接近布林带下轨
- **卖出信号**: MACD 死叉或 RSI > 70
- **适用场景**: 多指标信号共振，提高准确率
- **参数**: 综合使用 MACD(12,26,9) + RSI(14) + BOLL(20,2)

#### TrendFollowingStrategy - 趋势跟踪策略
- **买入信号**: MACD > 0 + RSI > 50 + 价格在布林带中轨上方
- **卖出信号**: MACD < 0 或价格跌破布林带中轨
- **适用场景**: 强趋势市场，跟随主趋势
- **参数**: 综合使用 MACD(12,26,9) + RSI(14) + BOLL(20,2)

## 使用示例

### 基础使用

```go
import (
    "github.com/sdcoffey/techan"
    "itick/strategy"
)

// 创建时间序列
series := techan.NewTimeSeries()
// ... 添加蜡烛图数据

// 使用单一策略
bollStrategy := strategy.BuildBollingerStrategy(series, 20, 2.0)
macdStrategy := strategy.BuildMACDStrategy(series, 12, 26, 9)
rsiStrategy := strategy.BuildRSIStrategy(series, 14, 30, 70)
kdjStrategy := strategy.BuildKDJCrossoverStrategy(series, 9, 3, 3)

// 使用组合策略
multiStrategy := strategy.BuildMultiIndicatorStrategy(series)
trendStrategy := strategy.BuildTrendFollowingStrategy(series)
```

### 批量获取所有默认策略

```go
// 获取所有默认策略配置
strategies := strategy.DefaultStrategies(series)

// 遍历使用
for name, strat := range strategies {
    // 使用策略进行回测或实盘交易
    record := techan.NewTradingRecord()
    for i := range series.Candles {
        if strat.ShouldEnter(i, record) {
            // 执行买入
        }
        if strat.ShouldExit(i, record) {
            // 执行卖出
        }
    }
}
```

### 自定义策略

如果需要自定义策略参数，可以直接调用策略函数：

```go
import (
    "github.com/sdcoffey/techan"
    "itick/indicator"
    "itick/strategy"
)

series := techan.NewTimeSeries()
closePrices := techan.NewClosePriceIndicator(series)

// 自定义布林带参数
upper, middle, lower := indicator.BollingerBands(closePrices, 30, 3.0)
customBollStrategy := strategy.BollingerBandStrategy(closePrices, upper, middle, lower)

// 自定义 RSI 参数
rsi := indicator.RSI(closePrices, 21)
customRsiStrategy := strategy.RSIStrategy(rsi, 20, 80)
```

## 注意事项

1. **不稳定期**: 所有策略都有 `UnstablePeriod`，在此期间指标未完全初始化，不应进行交易
2. **参数优化**: 默认参数适用于大多数情况，但针对特定市场和标的可能需要优化
3. **风险管理**: 策略不包含止损止盈逻辑，实际使用时应添加风险控制
4. **回测验证**: 在实盘前务必进行充分的历史数据回测
5. **市场适应**: 不同策略适用于不同市场环境（趋势/震荡），应根据市场状态选择

## 策略选择建议

- **震荡市**: 布林带策略、RSI 策略
- **趋势市**: MACD 策略、趋势跟踪策略
- **短线交易**: KDJ 策略
- **稳健交易**: 多指标组合策略（降低假信号）
- **激进交易**: 单一指标策略（信号更频繁）
