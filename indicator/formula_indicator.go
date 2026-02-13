package indicator

import "github.com/sdcoffey/big"

type FormulaIndicator struct {
	formula func(int) big.Decimal
}

func (f *FormulaIndicator) Calculate(index int) big.Decimal {
	return f.formula(index)
}
