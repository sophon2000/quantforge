package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "quantforge",
	Short: "QuantForge 量化交易框架",
	Long:  `QuantForge 支持回测、实盘与市场扫描等能力。`,
}

// Execute 执行根命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(backtestCmd)
	rootCmd.AddCommand(liveCmd)
	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(serveCmd)
}
