package datasource

import "github.com/sdcoffey/big"

// FormulaIndicator 公式指标：由闭包计算每个索引的值（用于 KDJ J 等）
type FormulaIndicator struct {
	Formula func(int) big.Decimal
}

// Calculate 实现 techan.Indicator
func (f *FormulaIndicator) Calculate(index int) big.Decimal {
	return f.Formula(index)
}
