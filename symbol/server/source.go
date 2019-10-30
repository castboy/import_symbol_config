package server

import (
	"github.com/shopspring/decimal"
	"time"
)

type Source struct {
	ID             int                     `json:"id" xorm:"id"`
	Source         string                  `json:"source" xorm:"source"`
	SourceType     SourceType              `json:"source_type" xorm:"source_type"`
	Digits         int                     `json:"digits" xorm:"digits"`
	Multiply       decimal.Decimal         `json:"multiply" xorm:"multiply"`
	ContractSize   decimal.Decimal         `json:"contract_size" xorm:"contract_size"`
	StopsLevel     int                     `json:"stops_level" xorm:"stops_level"`
	ProfitMode     ProfitMode              `json:"profit_mode" xorm:"profit_mode"`
	ProfitCurrency string                  `json:"profit_currency" xorm:"profit_currency"`
	MarginMode     MarginMode              `json:"margin_mode" xorm:"margin_mode"`
	MarginCurrency string                  `json:"margin_currency" xorm:"margin_currency"`
	SwapType       SwapType                `json:"swap_type" xorm:"swap_type"`
	SwapLong       decimal.Decimal         `json:"swap_long" xorm:"swap_long"`
	SwapShort      decimal.Decimal         `json:"swap_short" xorm:"swap_short"`
	SwapCurrency   string                  `json:"swap_currency" xorm:"swap_currency"`
	Swap3Day       time.Weekday            `json:"swap_3_day" xorm:"swap_3_day"`
	QuoteSession   map[time.Weekday]string `json:"quote_session" xorm:"-"`
	TradeSession   map[time.Weekday]string `json:"trade_session" xorm:"-"`
	Symbols        []string                `json:"symbols" xorm:"-"`
}

type (
	// ProfitMode
	ProfitMode int
	// SwapType
	SwapType int
	// MarginMode
	MarginMode int
	// SourceType
	SourceType int
)

// ProfitForex: 0 =>(closePrice - openPrice ) * contractSize * lots
//
// ProfitCfd: 1 => (closePrice - openPrice ) * contractSize * lots
//
// ProfitFutures: 2
const (
	ProfitForex ProfitMode = iota
	ProfitCfd
	ProfitFutures
)

// ByPoints: 0 => lots * longOrShort points * pointsSize
//
// ByMoney: 1
//
// ByInterest: 2 => lots * contractSize * longOrShort points /100 /360
//
// ByMoneyInMarginCurrency: 3
//
// ByInterestOfCfds: 4 => lots * contractSize * price * longOrShort points /100 /360
//
// ByInterestOfFutures: 5
const (
	ByPoints SwapType = iota
	ByMoney
	ByInterest
	ByMoneyInMarginCurrency
	ByInterestOfCfds
	ByInterestOfFutures
)

// MarginForex: 0 => lots * contractSize / leverage * percentage / 100
//
// MarginCfd: 1 => lots * contractSize * marketPrice * percentage / 100
//
// MarginFutures: 2 => lots * marginInitial * percentage / 100
//
// MarginCfdIndex: 3
//
// MarginCfdLeverage: 4 => lots * contractSize * marketPrice / leverage * percentage / 100
const (
	MarginForex MarginMode = iota
	MarginCfd
	MarginFutures
	MarginCfdIndex
	MarginCfdLeverage
)

// SourceFx: 0 => Currency Pair
//
// SourceMetal: 1 => Precious Metals, Gold, Silver, etc.
//
// SourceEnergy: 2 => Oil or NAT GAS
//
// SourceIndex: 3 => Index
//
// SourceCrypto: 4 => Visual Coin
const (
	SourceFx SourceType = iota
	SourceMetal
	SourceEnergy
	SourceIndex
	SourceCrypto
)

type SourceOperator interface {
	GetIDByName(source string) (ID int, exist bool, err error)
	InsertSource(source *Source) error
}

type sourceOperator struct {
	sourceRepo SourceOperator
}

var sourceOp *sourceOperator

func GetSourceOperator() *sourceOperator {
	return sourceOp
}

func InitSourceOperator(sourceRepo SourceOperator) *sourceOperator {
	if sourceOp == nil {
		sourceOp = &sourceOperator{
			sourceRepo,
		}
	}
	return sourceOp
}

func (ss *sourceOperator) GetIDByName(source string) (ID int, exist bool, err error) {
	return ss.sourceRepo.GetIDByName(source)
}

func (ss *sourceOperator) InsertSource(source *Source) error {
	return ss.sourceRepo.InsertSource(source)
}
