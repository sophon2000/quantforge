package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
)

//go:embed web/*
var webFS embed.FS

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "启动回测可视化 HTTP 服务",
	Long:  `启动 HTTP 服务，提供 /api/backtest 与前端 K 线图表页面。`,
	RunE:  runServe,
}

var servePort int

func init() {
	serveCmd.Flags().IntVarP(&servePort, "port", "p", 8080, "监听端口")
	serveCmd.Flags().StringP("symbol", "s", "AAPL", "默认回测标的")
	serveCmd.Flags().StringP("strategy", "S", "bollinger", "默认策略")
	serveCmd.Flags().Float64P("cash", "c", 100000, "默认初始资金")
	serveCmd.Flags().IntP("quantity", "q", 100, "默认每笔数量")
}

func runServe(cmd *cobra.Command, _ []string) error {
	http.HandleFunc("/api/backtest", handleBacktest)
	sub, _ := fs.Sub(webFS, "web")
	http.Handle("/", http.FileServer(http.FS(sub)))
	addr := ":" + strconv.Itoa(servePort)
	log.Printf("回测可视化服务: http://localhost%s", addr)
	return http.ListenAndServe(addr, nil)
}

func handleBacktest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	q := r.URL.Query()
	symbol := q.Get("symbol")
	if symbol == "" {
		symbol = "AAPL"
	}
	strategy := q.Get("strategy")
	if strategy == "" {
		strategy = "bollinger"
	}~
	cash := 100000.0
	if s := q.Get("cash"); s != "" {
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			cash = f
		}
	}
	quantity := 100
	if s := q.Get("quantity"); s != "" {
		if n, err := strconv.Atoi(s); err == nil {
			quantity = n
		}
	}

	res, err := RunBacktest(symbol, strategy, cash, quantity)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	_ = json.NewEncoder(w).Encode(res)
}
