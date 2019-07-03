/*******************************************************************************
 * // Copyright AnchyTec Corp. All Rights Reserved.
 * // SPDX-License-Identifier: Apache-2.0
 * // Author: shaozhiming
 ******************************************************************************/

package server

import (
	"github.com/shopspring/decimal"
	"testing"
	"time"
)

func TestSymbolService_SetSymbol(t *testing.T) {

}

func TestSymbolService_GetSymbolInfoByID(t *testing.T) {

}

func TestSymbolService_GetSymbolInfoByName(t *testing.T) {

}

func TestSwapFormula(t *testing.T) {
	lots := decimal.NewFromFloat(1)
	longOrShort := decimal.NewFromFloat(-3.86)
	expectedForexRes := decimal.NewFromFloat(-3.86)

	symbol := &Symbol{
		Digits:       5,
		ContractSize: decimal.NewFromFloat(100000.0),
	}

	swapFunc := symbol.SwapFormula()
	res := swapFunc(lots, longOrShort)
	if res.Equal(expectedForexRes) {
		t.Log("SwapByPointsCalc Right")
	} else {
		t.Error("SwapByPointsCalc Wrong!")
	}
}

func TestMarginFormula(t *testing.T) {
	lots := decimal.NewFromFloat(0.01)
	marketPrice := decimal.NewFromFloat(100)
	expectedForexRes := decimal.NewFromFloat(10)

	symbol := &Symbol{
		ContractSize: decimal.NewFromFloat(100000.0),
		Leverage:     decimal.NewFromFloat(100.0),
		Percentage:   decimal.NewFromFloat(100.0),
	}

	marginFunc := symbol.MarginFormula()
	res := marginFunc(lots, marketPrice)
	if res.Equal(expectedForexRes) {
		t.Log("MarginForexCalc Right")
	} else {
		t.Error("MarginForexCalc Wrong!")
	}
}

func TestProfitFormula(t *testing.T) {
	closePrice := decimal.NewFromFloat(100)
	openPrice := decimal.NewFromFloat(90)
	lots := decimal.NewFromFloat(0.01)
	exceptedRes := decimal.NewFromFloat(10000)

	symbol := &Symbol{
		ContractSize: decimal.NewFromFloat(100000.0),
	}

	profitFunc := symbol.ProfitFormula()
	res := profitFunc(lots, openPrice, closePrice)
	if res.Equal(exceptedRes) {
		t.Log("ProfitForexCalc Right")
	} else {
		t.Error("ProfitForexCalc Wrong!")
	}

}

var sessionSymbols = []*Symbol{
	&Symbol{
		Symbol: "AUDCAD",
		QuoteSession: map[time.Weekday]string{
			time.Sunday:    "06:05-07:05,00:00-00:00,00:00-00:00",
			time.Monday:    "00:00-02:55,06:05-07:05,00:00-00:00",
			time.Tuesday:   "00:00-02:55,06:05-07:05,00:00-00:00",
			time.Wednesday: "00:00-02:55,06:05-07:05,00:00-00:00",
			time.Thursday:  "00:00-02:55,06:05-07:05,00:00-00:00",
			time.Friday:    "00:00-02:55,06:05-07:05,00:00-00:00",
			time.Saturday:  "00:00-00:00,00:00-00:00,00:00-00:00",
		},
		TradeSession: map[time.Weekday]string{
			time.Sunday:    "06:05-07:05,00:00-00:00,00:00-00:00",
			time.Monday:    "00:00-02:55,06:05-07:05,00:00-00:00",
			time.Tuesday:   "00:00-02:55,06:05-07:05,00:00-00:00",
			time.Wednesday: "00:00-02:55,06:05-07:05,00:00-00:00",
			time.Thursday:  "00:00-02:55,06:05-07:05,00:00-00:00",
			time.Friday:    "00:00-02:55,06:05-07:05,00:00-00:00",
			time.Saturday:  "00:00-00:00,00:00-00:00,00:00-00:00",
		},
	},
	&Symbol{
		Symbol: "AUDUSD",
		QuoteSession: map[time.Weekday]string{
			time.Sunday:    "06:05-07:05,00:00-00:00,00:00-00:00",
			time.Monday:    "02:55-06:05,07:05-24:00,00:00-00:00",
			time.Tuesday:   "02:55-06:05,07:05-24:00,00:00-00:00",
			time.Wednesday: "02:55-06:05,07:05-24:00,00:00-00:00",
			time.Thursday:  "02:55-06:05,07:05-24:00,00:00-00:00",
			time.Friday:    "02:55-06:05,07:05-24:00,00:00-00:00",
			time.Saturday:  "00:00-00:00,00:00-00:00,00:00-00:00",
		},
		TradeSession: map[time.Weekday]string{
			time.Sunday:    "07:05-24:00,00:00-00:00,00:00-00:00",
			time.Monday:    "02:55-06:05,07:05-24:00,00:00-00:00",
			time.Tuesday:   "02:55-06:05,07:05-24:00,00:00-00:00",
			time.Wednesday: "02:55-06:05,07:05-24:00,00:00-00:00",
			time.Thursday:  "02:55-06:05,07:05-24:00,00:00-00:00",
			time.Friday:    "02:55-06:05,07:05-24:00,00:00-00:00",
			time.Saturday:  "00:00-00:00,00:00-00:00,00:00-00:00",
		},
	},
}

func TestIsQuotable(t *testing.T) {
	for _, symbol := range sessionSymbols {
		if symbol.IsQuotable() {
			t.Logf("%s is quotable", symbol.Symbol)
		} else {
			t.Logf("%s is not quotable", symbol.Symbol)
		}
	}
}

func TestIsTradable(t *testing.T) {
	for _, symbol := range sessionSymbols {
		if symbol.IsTradable() {
			t.Logf("%s is tradable", symbol.Symbol)
		} else {
			t.Logf("%s is not tradable", symbol.Symbol)
		}
	}
}
