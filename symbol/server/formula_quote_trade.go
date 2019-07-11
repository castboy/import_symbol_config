package server

import (
	"fmt"
	"github.com/shopspring/decimal"
	"strings"
	"time"
)

const (
	SystemLeverage   = 100
	SymbolLeverage = 100
	DaysInYear = 360
)

func (symb *Symbol) MarginFormula() func(lots, marketPrice decimal.Decimal) decimal.Decimal {
	switch symb.MarginMode {
	case MarginForex:
		return func(lots, marketPrice decimal.Decimal) decimal.Decimal {
			// TODO : symbol.Leverage
			return lots.Mul(symb.ContractSize).Div(decimal.NewFromFloat(SymbolLeverage)).Mul(symb.Percentage).Div(decimal.NewFromFloat(SystemLeverage))
		}

	case MarginCfd:
		return func(lots, marketPrice decimal.Decimal) decimal.Decimal {
			return lots.Mul(symb.ContractSize).Mul(marketPrice).Mul(symb.Percentage).Div(decimal.NewFromFloat(SystemLeverage))
		}

	case MarginFutures:
		return func(lots, marketPrice decimal.Decimal) decimal.Decimal {
			return lots.Mul(symb.MarginInitial).Mul(symb.Percentage).Div(decimal.NewFromFloat(SystemLeverage))
		}

	case MarginCfdIndex:
		// TODO

	case MarginCfdLeverage:
		return func(lots, marketPrice decimal.Decimal) decimal.Decimal {
			// TODO : symbol.Leverage
			return lots.Mul(symb.ContractSize).Mul(marketPrice).Div(decimal.NewFromFloat(SymbolLeverage)).Mul(symb.Percentage).Div(decimal.NewFromFloat(SystemLeverage))
		}

	default:
		panic(fmt.Sprintf("invalid margin mode: %v", symb.MarginMode))
	}

	return nil
}

func (symb *Symbol) ProfitFormula() func(lots, openPrice, closePrice decimal.Decimal) decimal.Decimal {
	switch symb.ProfitMode {
	case ProfitForex, ProfitCfd:
		return func(lots, openPrice, closePrice decimal.Decimal) decimal.Decimal {
			return closePrice.Sub(openPrice).Mul(symb.ContractSize).Mul(lots)
		}

	case ProfitFutures:
		// TODO

	default:
		panic(fmt.Sprintf("invalid profit mode: %v", symb.ProfitMode))
	}

	return nil
}

func (symb *Symbol) SwapFormula() func(lots, longOrShort decimal.Decimal, price ...decimal.Decimal) decimal.Decimal {
	switch symb.SwapType {
	case ByPoints:
		return func(lots, longOrShort decimal.Decimal, price ...decimal.Decimal) decimal.Decimal {
			divider := decimal.New(1, int32(-symb.Digits))
			return lots.Mul(longOrShort).Mul(symb.ContractSize).Mul(divider)
		}

	case ByMoney:
		// TODO

	case ByInterest:
		return func(lots, longOrShort decimal.Decimal, price ...decimal.Decimal) decimal.Decimal {
			return lots.Mul(longOrShort).Mul(symb.ContractSize).Div(decimal.NewFromFloat(SystemLeverage)).Div(decimal.NewFromFloat(DaysInYear))
		}

	case ByMoneyInMarginCurrency:
		// TODO

	case ByInterestOfCfds:
		return func(lots, longOrShort decimal.Decimal, price ...decimal.Decimal) decimal.Decimal {
			return lots.Mul(longOrShort).Mul(symb.ContractSize).Mul(price[0]).Div(decimal.NewFromFloat(SystemLeverage)).Div(decimal.NewFromFloat(DaysInYear))
		}

	case ByInterestOfFutures:
		// TODO

	default:
		panic(fmt.Sprintf("invalid swap type: %v", symb.SwapType))
	}

	return nil
}

//

func (symb *Symbol) IsQuotable() bool {
	return GetHolidayOperator().IsTradable(symb.Symbol) && HitSession(symb.QuoteSession)
}

func (symb *Symbol) IsTradable() bool {
	return GetHolidayOperator().IsTradable(symb.Symbol) && HitSession(symb.TradeSession)
}

func HitSession(allSessions map[time.Weekday]string) bool {
	// GMT is equal to UTC
	gmt := time.Now().UTC()
	nowStr := gmt.Format("15:04:05")
	weekday := gmt.Weekday()
	weekdaySessions := strings.Split(allSessions[weekday], ",")

	for _, session := range weekdaySessions {
		beginEnd := strings.Split(session, "-")
		if beginEnd[0] <= nowStr && nowStr < beginEnd[1] {
			return true
		}
	}

	return false
}
