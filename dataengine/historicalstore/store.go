package historicalstore

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/sdcoffey/big"
	"github.com/sdcoffey/techan"
)

// CSVRow 一行历史 K 线（CSV）
type CSVRow struct {
	Symbol string
	Date   string
	Open   string
	High   string
	Low    string
	Close  string
	Volume string
}

// HistoricalStore 历史数据存储接口
type HistoricalStore interface {
	// LoadCSV 从 CSV 加载并按标的聚合，csvPath 为空时用默认路径
	LoadCSV(csvPath string) (map[string][]CSVRow, error)
	// TimeSeries 将 rows 转为 techan.TimeSeries
	TimeSeries(rows []CSVRow) (*techan.TimeSeries, error)
}

// CSVStore 基于 CSV 的历史存储实现
type CSVStore struct {
	DefaultPath string
}

// NewCSVStore 构造，defaultPath 为空时使用 datasource/csv/ 下默认文件
func NewCSVStore(defaultPath string) *CSVStore {
	if defaultPath == "" {
		defaultPath = "dataengine/historicalstore/S&P 500 Stock Prices 2014-2017.csv"
	}
	return &CSVStore{DefaultPath: defaultPath}
}

// LoadCSV 实现 HistoricalStore
func (s *CSVStore) LoadCSV(csvPath string) (map[string][]CSVRow, error) {
	if csvPath == "" {
		csvPath = s.DefaultPath
	}
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

	rows := make([]CSVRow, 0, len(all)-1)
	for _, rec := range all[1:] {
		if len(rec) < 7 {
			continue
		}
		rows = append(rows, CSVRow{
			Symbol: rec[0],
			Date:   rec[1],
			Open:   rec[2],
			High:   rec[3],
			Low:    rec[4],
			Close:  rec[5],
			Volume: rec[6],
		})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Symbol != rows[j].Symbol {
			return rows[i].Symbol < rows[j].Symbol
		}
		return rows[i].Date < rows[j].Date
	})

	grouped := make(map[string][]CSVRow)
	for _, r := range rows {
		grouped[r.Symbol] = append(grouped[r.Symbol], r)
	}
	return grouped, nil
}

// TimeSeries 实现 HistoricalStore
func (s *CSVStore) TimeSeries(rows []CSVRow) (*techan.TimeSeries, error) {
	series := techan.NewTimeSeries()
	for _, r := range rows {
		date, err := time.ParseInLocation(time.DateOnly, r.Date, time.UTC)
		if err != nil {
			return nil, fmt.Errorf("解析日期: %w", err)
		}
		period := techan.NewTimePeriod(date, time.Hour*24)
		candle := techan.NewCandle(period)
		candle.OpenPrice = big.NewFromString(r.Open)
		candle.ClosePrice = big.NewFromString(r.Close)
		candle.MaxPrice = big.NewFromString(r.High)
		candle.MinPrice = big.NewFromString(r.Low)
		candle.Volume = big.NewFromString(r.Volume)
		series.AddCandle(candle)
	}
	return series, nil
}
