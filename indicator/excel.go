package indicator

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/sdcoffey/big"
	"github.com/sdcoffey/techan"
)

const csvPath = "/home/baixiao/workspace/quantforge/datasource/csv/S&P 500 Stock Prices 2014-2017.csv"

// row 表示 CSV 的一行，用于按 code-date 排序
type row struct {
	symbol string
	date   string
	open   string
	high   string
	low    string
	close  string
	volume string
}

func GetData() (map[string][]row, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, fmt.Errorf("打开 CSV: %w", err)
	}
	defer file.Close()

	r := csv.NewReader(file)
	all, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("读取 CSV: %w", err)
	}
	if len(all) < 2 {
		return nil, fmt.Errorf("CSV 无数据行")
	}

	rows := make([]row, 0, len(all)-1)
	for _, rec := range all[1:] {
		if len(rec) < 7 {
			continue
		}
		rows = append(rows, row{
			symbol: rec[0],
			date:   rec[1],
			open:   rec[2],
			high:   rec[3],
			low:    rec[4],
			close:  rec[5],
			volume: rec[6],
		})
	}

	// 按股票 code、再按时间聚合（排序）
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].symbol != rows[j].symbol {
			return rows[i].symbol < rows[j].symbol
		}
		return rows[i].date < rows[j].date
	})

	groupedRows := make(map[string][]row, 0)
	for _, r := range rows {
		groupedRows[r.symbol] = append(groupedRows[r.symbol], r)
	}
	return groupedRows, nil
}

func GenerateSeries(rows []row) (*techan.TimeSeries, error) {
	series := techan.NewTimeSeries()
	for _, r := range rows {
		date, err := time.ParseInLocation(time.DateOnly, r.date, time.UTC)
		if err != nil {
			return nil, fmt.Errorf("解析日期: %w", err)
		}
		period := techan.NewTimePeriod(date, time.Hour*24)
		candle := techan.NewCandle(period)
		candle.OpenPrice = big.NewFromString(r.open)
		candle.ClosePrice = big.NewFromString(r.close)
		candle.MaxPrice = big.NewFromString(r.high)
		candle.MinPrice = big.NewFromString(r.low)
		candle.Volume = big.NewFromString(r.volume)
		series.AddCandle(candle)
	}
	return series, nil
}
