package strategy

import (
	"time"

	"github.com/sdcoffey/big"
	"github.com/sdcoffey/techan"
)

func StrategyOne(indicator techan.Indicator) techan.RuleStrategy {
	entryConstant := techan.NewConstantIndicator(30)
	exitConstant := techan.NewConstantIndicator(10)

	entryRule := techan.And(
		techan.NewCrossUpIndicatorRule(entryConstant, indicator),
		techan.PositionNewRule{})

	exitRule := techan.And(
		techan.NewCrossDownIndicatorRule(indicator, exitConstant),
		techan.PositionOpenRule{})

	return techan.RuleStrategy{
		UnstablePeriod: 10,
		EntryRule:      entryRule,
		ExitRule:       exitRule,
	}
}

func Test1(ruleStratgy techan.RuleStrategy, trade *techan.TradingRecord, index int) {
	if ruleStratgy.ShouldEnter(index, trade) {
		entranceOrder := techan.Order{
			Side:          techan.BUY,
			Security:      "APPL",
			Amount:        big.ONE,
			Price:         big.NewFromString("2"),
			ExecutionTime: time.Now(),
		}
		trade.Operate(entranceOrder)
	}

	if ruleStratgy.ShouldExit(index, trade) {
		entranceOrder := techan.Order{
			Side:          techan.SELL,
			Security:      "APPL",
			Amount:        big.ONE,
			Price:         big.NewFromString("2"),
			ExecutionTime: time.Now(),
		}
		trade.Operate(entranceOrder)
	}
}
