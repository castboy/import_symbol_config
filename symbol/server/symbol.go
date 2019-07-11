/*******************************************************************************
 * // Copyright AnchyTec Corp. All Rights Reserved.
 * // SPDX-License-Identifier: Apache-2.0
 * // Author: shaozhiming
 ******************************************************************************/

package server

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/juju/errors"
	"github.com/shopspring/decimal"
	"import_symbol_config/config"
	"time"
)

// Symbol represents a instance of symbol
type Symbol struct {
	// common settings
	ID            int             `json:"id" xorm:"id"`
	Symbol        string          `json:"symbol" xorm:"symbol"`
	Source        string          `json:"source" xorm:"source"`
	SymbolType    SymbolType      `json:"symbol_type" xorm:"symbol_type"`
	SecurityID    int             `json:"security_id" xorm:"security_id"`
	Digits        int             `json:"digits" xorm:"digits"`
	Point         decimal.Decimal `json:"point" xorm:"point"`
	Multiply      decimal.Decimal `json:"multiply" xorm:"multiply"`
	ContractSize  decimal.Decimal `json:"contract_size" xorm:"contract_size"`
	StopsLevel    int             `json:"stops_level" xorm:"stops_level"`
	MarginInitial decimal.Decimal `json:"margin_initial" xorm:"margin_initial"`
	MarginDivider decimal.Decimal `json:"margin_divider" xorm:"margin_divider"`
	Percentage    decimal.Decimal `json:"percentage" xorm:"percentage"`
	// profit settings
	ProfitMode     ProfitMode `json:"profit_mode" xorm:"profit_mode"`
	ProfitCurrency string     `json:"profit_currency" xorm:"profit_currency"`
	// margin settings
	MarginMode     MarginMode      `json:"margin_mode" xorm:"margin_mode"`
	MarginCurrency string          `json:"margin_currency" xorm:"margin_currency"`
	// swap settings
	SwapType     SwapType        `json:"swap_type" xorm:"swap_type"`
	SwapLong     decimal.Decimal `json:"swap_long" xorm:"swap_long"`
	SwapShort    decimal.Decimal `json:"swap_short" xorm:"swap_short"`
	Swap3Day     SwapWeekday     `json:"swap_3_day" xorm:"swap_3_day"`
	SwapCurrency string          `josn:"swap_currency" xorm:"swap_currency"`
	// session settings
	QuoteSession map[time.Weekday]string `json:"quote_session" xorm:"-"`
	TradeSession map[time.Weekday]string `json:"trade_session" xorm:"-"`
	Index        int                     `json:"index" xorm:"index"`
}

type (
	// ProfitMode
	ProfitMode int
	// SwapType
	SwapType int
	// SwapWeekday
	SwapWeekday int
	// SessionWeekday
	SessionWeekday int
	// MarginMode
	MarginMode int
	// SymbolType
	SymbolType int
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

// SwapMonday: 0
//
// SwapTuesday: 1
//
// SwapWednesday: 2
//
// SwapThursday: 3
//
// SwapFriday: 4
//
// SwapSaturday: 5
//
// SwapSunday: 6
const (
	SwapMonday SwapWeekday = iota
	SwapTuesday
	SwapWednesday
	SwapThursday
	SwapFriday
	SwapSaturday
	SwapSunday
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

// SymbolFx: 0 => Currency Pair
//
// SymbolMetal: 1 => Precious Metals, Gold, Silver, etc.
//
// SymbolEnergy: 2 => Oil or NAT GAS
//
// SymbolIndex: 3 => Index
//
// SymbolCrypto: 4 => Visual Coin
const (
	SymbolFx SymbolType = iota
	SymbolMetal
	SymbolEnergy
	SymbolIndex
	SymbolCrypto
)

type SymbolRepository interface {
	GetSymbolInfoByName(symbolName string) (symbol *Symbol, exist bool, err error)
	GetSymbolInfoByID(ID int) (symbol *Symbol, exist bool, err error)
	GetSymbols() (symbols []Symbol, err error)

	//create new record, update existing record, delete existing record.
	InsertSymbol(symbol *Symbol) error
	InsertSession(sess []*Session) error
	GetIDByName(symbolName string) (ID int, exist bool, err error)

	UpdateByID(ID int, symbol *Symbol) error
	UpdateByName(symbolName string, symbol *Symbol) error

	NewTransaction() *xorm.Session
	TransactionDeleteByName(ss *xorm.Session, tableName interface{}, symbolName string) (num int64, err error)
	TransactionDeleteSymbolByID(ss *xorm.Session, ID int) (num int64, err error)
	TransactionDeleteSessionByID(ss *xorm.Session, ID int) error

	GetAllSecuritySymbols() ([]map[string]string, error)
	GetSecuritySymbols(securityID int) (symbols []string, err error)

	UpdateSymbolSecurity(symbolID int, securityID int) (num int64, err error)
	ValidSymbolID(ID int) (valid bool, err error)
	ValidSymbolName(symbolName string) (valid bool, err error)
	ValidSymbolSecurity(symbolID int, securityID int) (valid bool, err error)

	GetSymbolsNameBySecurityID(securityID int) (symbols []string, err error)
	GetSymbolsName() (symbols []string, err error)

	SecurityHoldSymbols(securityID int) (hold bool, err error)

	GetSymbolLeverage(symbolSource string) (symbols []string, err error)
}

type symbolOperator struct {
	symbolRepo SymbolRepository
}

var symbolOp *symbolOperator

func GetSymbolOperator() *symbolOperator {
	return symbolOp
}

func InitSymbolOperator(symbolRepo SymbolRepository) *symbolOperator {
	if symbolOp == nil {
		symbolOp = &symbolOperator{
			symbolRepo,
		}
	}
	return symbolOp
}


func (ss *symbolOperator) Start() {
	if config.GetConfigService("symbol").GetBool("import_from_config") {
		symbols, err := parseSymbols()
		if err != nil {
			//panic(err)
		}
		importSymbols(symbols)

		securities, err := parseSecurity()
		if err != nil {
			panic(err)
		}
		insertSecurityInfo(securities)

		setSymbolSecurity(securities)

		holidays, err := parseHoliday()
		if err != nil {
			panic(err)
		}
		importHolidays(holidays)
	}
}

func importHolidays(holidays []Holiday) {
	Len := len(holidays)
	for i := 0; i < Len; i++ {
		err := importHoliday(&holidays[i])
		if err != nil {
			panic(err)
		}
	}
}

func importHoliday(holiday *Holiday) error {
	ho := GetHolidayOperator()
	err := ho.InsertHoliday(holiday)
	return err
}

func insertSecurityInfo(securities []Security) {
	so := GetSecurityOperator()
	for i, _ := range securities {
		if err := so.InsertSecurityInfo(&securities[i]); err != nil {
			panic(err)
		}
	}
}

func setSymbolSecurity(securities []Security) {
	so := GetSecurityOperator()
	for i, _ := range securities {
		securityID, exist, err := so.GetIDByName(securities[i].SecurityName)
		if err != nil {
			panic(err)
		}

		if !exist {
			err = fmt.Errorf("invalid security name: %s", securities[i].SecurityName)
			panic(err)
		}

		ss := GetSymbolOperator()
		for j, _ := range securities[i].Symbols {
			symbolID, exist, err := symbolOp.GetIDByName(securities[i].Symbols[j])
			if err != nil {
				fmt.Println(err) // TODO
			}

			if !exist {
				err = fmt.Errorf("sdfasf")
				panic(err)
			}

			if err := ss.SetSymbolSecurity(symbolID, securityID); err != nil {
				//fmt.Println(err) // TODO
			}
		}
	}
}

func importSymbols(symbols []Symbol) {
	Len := len(symbols)
	for i := 0; i < Len; i++ {
		err := importSymbol(&symbols[i])
		if err != nil {
			panic(err)
		}
	}
}

func importSymbol(symbol *Symbol) error {
	so := GetSymbolOperator()
	err := so.InsertSymbol(symbol)
	if err != nil {
		return err
	}

	id, exist, err := so.GetIDByName(symbol.Symbol)
	if err != nil {
		return err
	}

	if !exist {
		err = fmt.Errorf("invalid symbol name: %s", symbol.Symbol)
		return err
	}

	symbol.ID = id

	quote := DecodeSession(symbol.ID, Quote, symbol.QuoteSession)
	trade := DecodeSession(symbol.ID, Trade, symbol.TradeSession)
	quote = append(quote, trade...)

	return GetSessionOperator().InsertSessions(quote...)
}


func (ss *symbolOperator) GetSymbolInfoByName(symbolName string) (symbol *Symbol) {
	// get symbol from mysql if get symbol failed from cache.
	symbol, exist, err := ss.symbolRepo.GetSymbolInfoByName(symbolName)
	if err != nil {
		err = errors.Annotate(err, "sql exec")
		return nil
	}

	if !exist {
		err = errors.NotFoundf("symbol name %s", symbolName)
		return nil
	}

	sess, err := GetSessionOperator().GetSessionsByName(symbolName)
	if err != nil {
		err = errors.Annotate(err, "sql exec")
		return nil
	}

	if len(sess) == 0 {
		err = errors.NotFoundf("symbol name %s", symbolName)
		return nil
	}

	sessions, err := EncodeSession(sess)
	if err != nil {
		err = errors.Annotate(err, "encode session")
		return nil
	}

	symbol.QuoteSession, symbol.QuoteSession = make(map[time.Weekday]string), make(map[time.Weekday]string)
	symbol.QuoteSession, symbol.TradeSession = sessions[Quote], sessions[Trade]

	// cache

	return
}

func (ss *symbolOperator) GetSymbolInfoByID(ID int) (symbol *Symbol) {
	// get symbol from mysql if get symbol failed from cache.
	symbol, exist, err := ss.symbolRepo.GetSymbolInfoByID(ID)
	if err != nil {
		err = errors.Annotate(err, "sql exec")
		return nil
	}

	if !exist {
		err = errors.NotFoundf("symbol id %d", ID)
		return nil
	}

	sess, err := GetSessionOperator().GetSessionsByID(ID)
	if err != nil {
		err = errors.Annotate(err, "sql exec")
		return nil
	}

	if len(sess) == 0 {
		err = errors.NotFoundf("symbol id %d", ID)
		return nil
	}

	sessions, err := EncodeSession(sess)
	if err != nil {
		err = errors.Annotate(err, "encode session")
		return nil
	}

	symbol.QuoteSession, symbol.QuoteSession = make(map[time.Weekday]string), make(map[time.Weekday]string)
	symbol.QuoteSession, symbol.TradeSession = sessions[Quote], sessions[Trade]

	// cache

	return
}

func (ss *symbolOperator) GetSymbols() (symbols []Symbol) {

	// if symbols were not cached.
	symbols, err := ss.symbolRepo.GetSymbols()
	if err != nil {
		err = errors.Annotate(err, "sql exec")
		return nil
	}

	if len(symbols) == 0 {
		err = errors.NotFoundf("symbols")
		return nil
	}

	for k, _ := range symbols {
		ss, err := GetSessionOperator().GetSessionsByID(symbols[k].ID)
		if err != nil {
			err = errors.Annotate(err, "sql exec")
			continue
		}
		if len(ss) == 0 {
			err = errors.NotFoundf("symbol id %d", symbols[k].ID)
			continue
		}

		sessions, err := EncodeSession(ss)
		if err != nil {
			err = errors.Annotate(err, "encode session")
			continue // TODO
		}

		symbols[k].QuoteSession, symbols[k].QuoteSession = make(map[time.Weekday]string), make(map[time.Weekday]string)
		symbols[k].QuoteSession, symbols[k].TradeSession = sessions[Quote], sessions[Trade]
	}
	return
}

func symbolFormatCheck(symbol *Symbol) error {
	if symbol.SymbolType != SymbolFx && symbol.SymbolType != SymbolMetal && symbol.SymbolType != SymbolEnergy &&
		symbol.SymbolType != SymbolIndex && symbol.SymbolType != SymbolCrypto {
		return errors.NotValidf("symbol type %d", symbol.SymbolType)
	}

	valid, err := GetSecurityOperator().ValidSecurityID(symbol.SecurityID)
	if err != nil {
		err = errors.Annotatef(err, "security id %d", symbol.SecurityID)
		return err
	}
	if !valid {
		return errors.NotValidf("security id %d", symbol.SecurityID)
	}

	if symbol.ProfitMode != ProfitForex && symbol.ProfitMode != ProfitCfd && symbol.ProfitMode != ProfitFutures {
		return errors.NotValidf("profit mode %d", symbol.ProfitMode)
	}

	if symbol.MarginMode != MarginForex && symbol.MarginMode != MarginCfd && symbol.MarginMode != MarginFutures &&
		symbol.MarginMode != MarginCfdIndex && symbol.MarginMode != MarginCfdLeverage {
		return errors.NotValidf("margin mode %d", symbol.MarginMode)
	}

	if symbol.SwapType != ByPoints && symbol.SwapType != ByMoney && symbol.SwapType != ByInterest &&
		symbol.SwapType != ByMoneyInMarginCurrency && symbol.SwapType != ByInterestOfCfds &&
		symbol.SwapType != ByInterestOfFutures {
		return errors.NotValidf("swaptype %d", symbol.SwapType)
	}

	if symbol.Swap3Day != SwapMonday && symbol.Swap3Day != SwapTuesday && symbol.Swap3Day != SwapWednesday &&
		symbol.Swap3Day != SwapThursday && symbol.Swap3Day != SwapFriday && symbol.Swap3Day != SwapSaturday &&
		symbol.Swap3Day != SwapSunday {
		return errors.NotValidf("swap3day: %d", symbol.Swap3Day)
	}

	return nil
}

func (ss *symbolOperator) InsertSymbol(symbol *Symbol) error {
	// insert mysql firstly.
	//err := symbolFormatCheck(symbol)
	//if err != nil {
	//	err = errors.NewNotValid(err, "validation failed")
	//	return err
	//}

	err := ss.symbolRepo.InsertSymbol(symbol)
	if err != nil {
		err = errors.Annotate(err, "sql exec")
		return err
	}

	if symbol.QuoteSession == nil && symbol.TradeSession == nil {
		return nil
	}

	id, exist, err := ss.symbolRepo.GetIDByName(symbol.Symbol)
	if err != nil {
		err = errors.Annotate(err, "sql exec")
		return err
	}

	if !exist {
		err = errors.NotFoundf("symbol name %s", symbol.Symbol)
		return err
	}

	symbol.ID = id

	if symbol.QuoteSession != nil {
		quote := DecodeSession(symbol.ID, Quote, symbol.QuoteSession)
		err := ss.symbolRepo.InsertSession(quote)
		if err != nil {
			err = errors.Annotate(err, "sql exec: insert quote session")
			return err
		}
	}

	if symbol.TradeSession == nil {
		trade := DecodeSession(symbol.ID, Trade, symbol.TradeSession)
		err := ss.symbolRepo.InsertSession(trade)
		if err != nil {
			err = errors.Annotate(err, "sql exec: insert trade session")
			return err
		}
	}

	// then insert cache.

	return nil
}

func (ss *symbolOperator) UpdateSymbolByID(ID int, symbol *Symbol) error {
	valid, err := ss.symbolRepo.ValidSymbolID(ID)
	if err != nil {
		err = errors.Annotate(err, "sql exec")
		return err
	}
	if !valid {
		err = errors.NotValidf("symbol id %d", ID)
		return err
	}

	err = symbolFormatCheck(symbol)
	if err != nil {
		err = errors.NewNotValid(err, "validation failed")
		return err
	}

	err = ss.symbolRepo.UpdateByID(ID, symbol)
	if err != nil {
		err = errors.Annotatef(err, "sql exec, id %d, symbol %+v", ID, symbol)
		return err
	}

	// delete cache to re-cache.

	return nil
}

func (ss *symbolOperator) UpdateSymbolByName(symbolName string, symbol *Symbol) error {
	valid, err := ss.symbolRepo.ValidSymbolName(symbolName)
	if err != nil {
		err = errors.Annotate(err, "sql exec")
		return err
	}
	if !valid {
		err = errors.NotValidf("symbol name %s", symbolName)
		return err
	}

	err = symbolFormatCheck(symbol)
	if err != nil {
		err = errors.NewNotValid(err, "validation failed")
		return err
	}

	err = ss.symbolRepo.UpdateByName(symbolName, symbol)
	if err != nil {
		err = errors.Annotatef(err, "sql exec, name %d, symbol %+v", symbolName, symbol)
		return err
	}

	// delete cache to re-cache.

	return nil
}

func (ss *symbolOperator) DeleteSymbolByID(ID int) error {
	trac := ss.symbolRepo.NewTransaction()
	defer trac.Close()

	if err := trac.Begin(); err != nil {
		err = errors.Annotatef(err, "sql exec")
		return err
	}

	hit, err := ss.symbolRepo.TransactionDeleteSymbolByID(trac, ID)
	if err != nil {
		err = errors.Annotatef(err, "sql exec")
		return err
	}
	if hit == 0 {
		err = errors.NotFoundf("symbol id %d", ID)
		return err
	}

	if err := ss.symbolRepo.TransactionDeleteSessionByID(trac, ID); err != nil {
		trac.Rollback()
		err = errors.Annotatef(err, "sql exec")
		return err
	}

	err = trac.Commit()
	if err != nil {
		err = errors.Annotatef(err, "sql exec")
		return err
	}

	// delete cache.

	return nil
}

func (ss *symbolOperator) DeleteSymbolByName(symbolName string) error {
	trac := ss.symbolRepo.NewTransaction()
	defer trac.Close()

	if err := trac.Begin(); err != nil {
		err = errors.Annotatef(err, "sql exec")
		return err
	}

	hit, err := ss.symbolRepo.TransactionDeleteByName(trac, Symbol{}, symbolName)
	if err != nil {
		err = errors.Annotatef(err, "sql exec")
		return err
	}
	if hit == 0 {
		err = errors.NotFoundf("symbol name %s", symbolName)
		return err
	}

	if _, err := ss.symbolRepo.TransactionDeleteByName(trac, Session{}, symbolName); err != nil {
		trac.Rollback()
		err = errors.Annotatef(err, "sql exec")
		return err
	}

	err = trac.Commit()
	if err != nil {
		err = errors.Annotatef(err, "sql exec")
		return err
	}

	// delete cache.

	return nil
}

func (ss *symbolOperator) GetAllSecuritySymbols() ([]map[string]string, error) {
	return ss.symbolRepo.GetAllSecuritySymbols()
}

func (ss *symbolOperator) GetSecuritySymbol(securityID int) ([]string, error) {
	return ss.symbolRepo.GetSecuritySymbols(securityID)
}

func (ss *symbolOperator) SetSymbolSecurity(symbolID int, securityID int) error {
	valid, err := ss.symbolRepo.ValidSymbolID(symbolID)
	if err != nil {
		err = errors.Annotatef(err, "sql exec")
		return err
	}
	if !valid {
		err = errors.NotValidf("symbol id %d", symbolID)
		return err
	}

	valid, err = GetSecurityOperator().ValidSecurityID(securityID)
	if err != nil {
		err = errors.Annotatef(err, "sql exec")
		return err
	}
	if !valid {
		err = errors.NotValidf("security id: %d", securityID)
		return err
	}

	_, err = ss.symbolRepo.UpdateSymbolSecurity(symbolID, securityID)
	if err != nil {
		err = errors.Annotatef(err, "sql exec")
		return err
	}

	// set cache symbol-security.

	return nil
}

func (ss *symbolOperator) UpdateSymbolSecurity(symbolID int, oldSecurityID int, newSecurityID int) error {
	valid, err := ss.symbolRepo.ValidSymbolSecurity(symbolID, oldSecurityID)
	if err != nil {
		err = errors.Annotatef(err, "sql exec")
		return err
	}
	if !valid {
		err = errors.NotFoundf("symbol id %d, security id %d", symbolID, oldSecurityID)
		return err
	}

	valid, err = GetSecurityOperator().ValidSecurityID(newSecurityID)
	if err != nil {
		err = errors.Annotatef(err, "sql exec")
		return err
	}
	if !valid {
		err = errors.NotValidf("security id %d", newSecurityID)
		return err
	}

	_, err = ss.symbolRepo.UpdateSymbolSecurity(symbolID, newSecurityID)
	if err != nil {
		err = errors.Annotatef(err, "sql exec")
		return err
	}

	// set cache symbol-security.

	return nil
}

func (ss *symbolOperator) GetIDByName(symbolName string) (int, bool, error) {
	return ss.symbolRepo.GetIDByName(symbolName)
}

func (ss *symbolOperator) GetSymbolsNameBySecurityID(securityID int) (symbols []string, err error) {
	return ss.symbolRepo.GetSymbolsNameBySecurityID(securityID)
}

func (ss *symbolOperator) GetSymbolsName() (symbols []string, err error) {
	return ss.symbolRepo.GetSymbolsName()
}

func (ss *symbolOperator) ValidSymbolName(symbolName string) (valid bool, err error) {
	return ss.symbolRepo.ValidSymbolName(symbolName)
}

func (ss *symbolOperator) ValidSymbolID(symbolID int) (valid bool, err error) {
	return ss.symbolRepo.ValidSymbolID(symbolID)
}

func (ss *symbolOperator) SecurityHoldSymbols(securityID int) (hold bool, err error) {
	return ss.symbolRepo.SecurityHoldSymbols(securityID)
}

func (ss *symbolOperator) GetSymbolLeverage(symbolSource string) (symbols []string, err error) {
	symbols, err = ss.symbolRepo.GetSymbolLeverage(symbolSource)
	if err != nil {
		err = errors.Annotatef(err, "sql exec")
		return nil, err
	}

	if len(symbols) == 0 {
		err = errors.NotFoundf("source %s", symbolSource)
		return nil, err
	}

	return symbols, nil
}
