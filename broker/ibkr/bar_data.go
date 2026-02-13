package ibkr

import (
	"github.com/scmhub/ibapi"
	"github.com/scmhub/ibsync"
)

// ReqHistoricalData 请求历史数据
// endDateTime: 结束时间，格式 "20230101 12:00:00" 或空字符串表示当前时间
// duration: 持续时间，如 "1 D", "1 W", "1 M", "1 Y"
// barSize: K线大小，如 "1 min", "5 mins", "1 hour", "1 day"
// whatToShow: 数据类型，如 "TRADES", "MIDPOINT", "BID", "ASK"
// useRTH: 是否只使用常规交易时间
// formatDate: 日期格式，1 表示 yyyyMMdd{space}{space}HH:mm:ss，2 表示时间戳
func (c *Client) ReqHistoricalData(
	contract *ibapi.Contract,
	endDateTime string,
	duration string,
	barSize string,
	whatToShow string,
	useRTH bool,
	formatDate int,
) (chan ibsync.Bar, ibsync.CancelFunc) {
	return c.ib.ReqHistoricalData(contract, endDateTime, duration, barSize, whatToShow, useRTH, formatDate)
}

// ReqHistoricalDataUpToDate 请求实时更新的历史数据
func (c *Client) ReqHistoricalDataUpToDate(
	contract *ibapi.Contract,
	duration string,
	barSize string,
	whatToShow string,
	useRTH bool,
	formatDate int,
) (chan ibsync.Bar, ibsync.CancelFunc) {
	return c.ib.ReqHistoricalDataUpToDate(contract, duration, barSize, whatToShow, useRTH, formatDate)
}

// ReqRealTimeBars 请求实时 Bar 数据
// barSize: 只能是 5 秒
// whatToShow: "TRADES", "MIDPOINT", "BID", "ASK"
func (c *Client) ReqRealTimeBars(
	contract *ibapi.Contract,
	barSize int,
	whatToShow string,
	useRTH bool,
) (chan ibsync.RealTimeBar, ibsync.CancelFunc) {
	return c.ib.ReqRealTimeBars(contract, barSize, whatToShow, useRTH)
}

// GetHistoricalBars 获取历史 K 线数据（辅助函数，返回 slice）
// endDateTime: 结束时间，格式 "20230101 12:00:00" 或空字符串表示当前时间
func (c *Client) GetHistoricalBars(
	contract *ibapi.Contract,
	endDateTime string,
	duration string,
	barSize string,
	whatToShow string,
	useRTH bool,
) []ibsync.Bar {
	barChan, _ := c.ib.ReqHistoricalData(contract, endDateTime, duration, barSize, whatToShow, useRTH, 1)

	var bars []ibsync.Bar
	for bar := range barChan {
		bars = append(bars, bar)
	}

	return bars
}
