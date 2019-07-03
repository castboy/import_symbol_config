/*******************************************************************************
 * // Copyright AnchyTec Corp. All Rights Reserved.
 * // SPDX-License-Identifier: Apache-2.0
 * // Author: shaozhiming
 ******************************************************************************/

package server

import (
	"errors"
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/shopspring/decimal"
	"sync"
	"time"
	"import_symbol_config/config"
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
	Leverage       decimal.Decimal `json:"leverage" xorm:"leverage"`
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
	GetSymbolInfoByName(symbolName string) (symbol *Symbol, err error)
	GetSymbolInfoByID(ID int) (symbol *Symbol, err error)
	GetSymbols() (symbols []Symbol, err error)

	//create new record, update existing record, delete existing record.
	InsertSymbol(symbol *Symbol) error
	InsertSession(sess []*Session) error
	GetIDByName(symbolName string) (ID int, err error)

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
	ValidSymbolNameID(symbolName string, symbolID int) (valid bool, err error)
	ValidSymbolSecurity(symbolID int, securityID int) (valid bool, err error)

	GetSymbolsNameBySecurityID(securityID int) (symbols []string, err error)
	GetSymbolsName() (symbols []string, err error)

	SecurityHoldSymbols(securityID int) (hold bool, err error)
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
		ss.importSymbols(symbols)

		securities, err := parseSecurity()
		if err != nil {
			panic(err)
		}
		ss.insertSecurityInfo(securities)

		ss.setSymbolSecurity(securities)

		holidays, err := parseHoliday()
		if err != nil {
			panic(err)
		}
		importHolidays(holidays)
	}

	date := time.Now().UTC().Format("2006-01-02")
	LoadHolidayCacheByDate(date)
	LoadSymbolCache()
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
	_, err := ho.holidayRepo.Insert(holiday)
	return err
}

func (ss *symbolOperator) insertSecurityInfo(securities []Security) {
	so := GetSecurityOperator()
	for i, _ := range securities {
		if err := so.securityRepo.InsertSecurity(&securities[i]); err != nil {
			panic(err)
		}
	}
}

func (ss *symbolOperator) setSymbolSecurity(securities []Security) {
	so := GetSecurityOperator()
	for i, _ := range securities {
		securityID, err := so.GetIDByName(securities[i].SecurityName)
		if err != nil {
			panic(err)
		}

		for j, _ := range securities[i].Symbols {
			symbolID, err := symbolOp.GetIDByName(securities[i].Symbols[j])
			if err != nil {
				fmt.Println(err) // TODO
			}

			if err := ss.SetSymbolSecurity(symbolID, securityID); err != nil {
				//fmt.Println(err) // TODO
			}
		}
	}
}

func (ss *symbolOperator) importSymbols(symbols []Symbol) {
	Len := len(symbols)
	for i := 0; i < Len; i++ {
		err := ss.importSymbol(&symbols[i])
		if err != nil {
			panic(err)
		}
	}
}

func (ss *symbolOperator) importSymbol(symbol *Symbol) error {
	err := ss.symbolRepo.InsertSymbol(symbol)
	if err != nil {
		return err
	}

	id, err := ss.symbolRepo.GetIDByName(symbol.Symbol)
	if err != nil {
		return err
	}

	symbol.ID = id

	quote := DecodeSession(symbol.ID, symbol.Symbol, "quote", symbol.QuoteSession)
	trade := DecodeSession(symbol.ID, symbol.Symbol, "trade", symbol.TradeSession)
	quote = append(quote, trade...)

	return GetSessionOperator().InsertSessions(quote...)
}

func (ss *symbolOperator) GetSymbolInfoByName(symbolName string) (symbol *Symbol) {
	// get symbol from cache firstly.
	symbol = getSymbolByNameFromCache(symbolName)
	if symbol != nil {
		return
	}
	// get symbol from mysql if get symbol failed from cache.
	symbol, err := ss.symbolRepo.GetSymbolInfoByName(symbolName)
	if err != nil || symbol == nil {
		return nil // TODO optimize
	}

	sess, err := GetSessionOperator().GetSessionsByName(symbolName)
	if err != nil || sess == nil {
		return symbol // TODO ??
	}

	sessions, err := EncodeSession(sess)
	if err != nil {
		return nil
	}

	symbol.QuoteSession, symbol.QuoteSession = make(map[time.Weekday]string), make(map[time.Weekday]string)
	symbol.QuoteSession, symbol.TradeSession = sessions["quote"], sessions["trade"]

	return
}

func (ss *symbolOperator) GetSymbolInfoByID(ID int) (symbol *Symbol) {
	// get symbol from cache firstly.
	symbol = getSymbolByIDFromCache(ID)
	if symbol != nil {
		return
	}
	// get symbol from mysql if get symbol failed from cache.
	symbol, err := ss.symbolRepo.GetSymbolInfoByID(ID)
	if err != nil {
		return nil // TODO optimize
	}

	sess, err := GetSessionOperator().GetSessionsByID(ID)
	if err != nil || sess == nil {
		return symbol // TODO ??
	}

	sessions, err := EncodeSession(sess)
	if err != nil {
		return nil
	}

	symbol.QuoteSession, symbol.QuoteSession = make(map[time.Weekday]string), make(map[time.Weekday]string)
	symbol.QuoteSession, symbol.TradeSession = sessions["quote"], sessions["trade"]

	return
}

func (ss *symbolOperator) GetSymbols() (symbols []Symbol) {
	// if symbols were cached.
	if symbCache != nil {
		return getSymbolsFromCache()
	}

	// if symbols were not cached.
	symbols, err := ss.symbolRepo.GetSymbols()
	if err != nil || len(symbols) == 0 {
		return nil // TODO optimize
	}

	for k, _ := range symbols {
		ss, err := GetSessionOperator().GetSessionsByID(symbols[k].ID)
		if err != nil || ss == nil {
			continue // TODO ??
		}
		sessions, err := EncodeSession(ss)
		if err != nil {
			return nil // ??
		}

		symbols[k].QuoteSession = sessions["quote"]
		symbols[k].TradeSession = sessions["trade"]
	}
	return
}

func symbolFormatCheck(symbol *Symbol) error {
	if symbol.SymbolType != SymbolFx && symbol.SymbolType != SymbolMetal && symbol.SymbolType != SymbolEnergy &&
		symbol.SymbolType != SymbolIndex && symbol.SymbolType != SymbolCrypto {
		return errors.New(fmt.Sprintf("invalid symbol type: %d", symbol.SymbolType))
	}

	valid, err := GetSecurityOperator().ValidSecurityID(symbol.SecurityID)
	if err != nil {
		return err
	}
	if !valid {
		return errors.New(fmt.Sprintf("invalid security id: %d", symbol.SecurityID))
	}

	if symbol.ProfitMode != ProfitForex && symbol.ProfitMode != ProfitCfd && symbol.ProfitMode != ProfitFutures {
		return errors.New(fmt.Sprintf("invalid profit mode: %d", symbol.ProfitMode))
	}

	if symbol.MarginMode != MarginForex && symbol.MarginMode != MarginCfd && symbol.MarginMode != MarginFutures &&
		symbol.MarginMode != MarginCfdIndex && symbol.MarginMode != MarginCfdLeverage {
		return errors.New(fmt.Sprintf("invalid margin mode: %d", symbol.MarginMode))
	}

	if symbol.SwapType != ByPoints && symbol.SwapType != ByMoney && symbol.SwapType != ByInterest &&
		symbol.SwapType != ByMoneyInMarginCurrency && symbol.SwapType != ByInterestOfCfds &&
		symbol.SwapType != ByInterestOfFutures {
		return errors.New(fmt.Sprintf("invalid swaptype: %d", symbol.SwapType))
	}

	if symbol.Swap3Day != SwapMonday && symbol.Swap3Day != SwapTuesday && symbol.Swap3Day != SwapWednesday &&
		symbol.Swap3Day != SwapThursday && symbol.Swap3Day != SwapFriday && symbol.Swap3Day != SwapSaturday &&
		symbol.Swap3Day != SwapSunday {
		return errors.New(fmt.Sprintf("invalid swap3day: %d", symbol.Swap3Day))
	}

	return nil
}

func (ss *symbolOperator) InsertSymbol(symbol *Symbol) error {
	// insert mysql firstly.
	err := symbolFormatCheck(symbol)
	if err != nil {
		return err
	}

	err = ss.symbolRepo.InsertSymbol(symbol)
	if err != nil {
		return err
	}

	if symbol.QuoteSession == nil && symbol.TradeSession == nil {
		return nil
	}

	id, err := ss.symbolRepo.GetIDByName(symbol.Symbol)
	if err != nil {
		return err
	}

	symbol.ID = id

	if symbol.QuoteSession != nil {
		quote := DecodeSession(symbol.ID, symbol.Symbol, "quote", symbol.QuoteSession)
		err := ss.symbolRepo.InsertSession(quote)
		if err != nil {
			return err
		}
	}

	if symbol.TradeSession == nil {
		trade := DecodeSession(symbol.ID, symbol.Symbol, "quote", symbol.TradeSession)
		err := ss.symbolRepo.InsertSession(trade)
		if err != nil {
			return err
		}
	}

	// then insert cache.
	insertSymbolToCache(symbol)

	return nil
}

func (ss *symbolOperator) UpdateSymbolByID(ID int, symbol *Symbol) error {
	valid, err := ss.symbolRepo.ValidSymbolID(ID)
	if err != nil {
		return err
	}
	if !valid {
		return errors.New(fmt.Sprintf("invalid symbol id: %d", ID))
	}

	err = symbolFormatCheck(symbol)
	if err != nil {
		return err
	}

	err = ss.symbolRepo.UpdateByID(ID, symbol)
	if err != nil {
		return err
	}

	// update cache.
	deleteSymbolCacheByID(ID)
	symbol = ss.GetSymbolInfoByID(ID)
	if symbol == nil {
		return errors.New("update symbol cache failed by id")
	}
	insertSymbolToCache(symbol)

	return nil
}

func (ss *symbolOperator) UpdateSymbolByName(symbolName string, symbol *Symbol) error {
	valid, err := ss.symbolRepo.ValidSymbolName(symbolName)
	if err != nil {
		return err
	}
	if !valid {
		return errors.New(fmt.Sprintf("invalid symbol name: %s", symbolName))
	}

	err = symbolFormatCheck(symbol)
	if err != nil {
		return err
	}

	err = ss.symbolRepo.UpdateByName(symbolName, symbol)
	if err != nil {
		return err
	}

	// update cache.
	deleteSymbolCacheByName(symbolName)
	symbol = ss.GetSymbolInfoByName(symbolName)
	if symbol == nil {
		return errors.New("update symbol cache failed by name")
	}
	insertSymbolToCache(symbol)

	return nil
}

func (ss *symbolOperator) DeleteSymbolByID(ID int) error {
	trac := ss.symbolRepo.NewTransaction()
	defer trac.Close()

	if err := trac.Begin(); err != nil {
		return err
	}

	hit, err := ss.symbolRepo.TransactionDeleteSymbolByID(trac, ID)
	if hit == 0 && err == nil {
		return errors.New(fmt.Sprintf("invalid symbol id: %d", ID))
	}
	if err != nil {
		return err
	}

	if err := ss.symbolRepo.TransactionDeleteSessionByID(trac, ID); err != nil {
		trac.Rollback()
		return err
	}

	err = trac.Commit()
	if err != nil {
		return err
	}

	// delete cache.
	deleteSymbolCacheByID(ID)

	return nil
}

func (ss *symbolOperator) DeleteSymbolByName(symbolName string) error {
	trac := ss.symbolRepo.NewTransaction()
	defer trac.Close()

	if err := trac.Begin(); err != nil {
		return err
	}

	hit, err := ss.symbolRepo.TransactionDeleteByName(trac, Symbol{}, symbolName)
	if hit == 0 && err == nil {
		return errors.New(fmt.Sprintf("invalid symbol name: %s", symbolName))
	}
	if err != nil {
		return err
	}

	if _, err := ss.symbolRepo.TransactionDeleteByName(trac, Session{}, symbolName); err != nil {
		trac.Rollback()
		return err
	}

	err = trac.Commit()
	if err != nil {
		return err
	}

	// delete cache.
	deleteSymbolCacheByName(symbolName)

	return nil
}

func (ss *symbolOperator) GetAllSecuritySymbols() (map[string]string, error) {
	ses, err := ss.symbolRepo.GetAllSecuritySymbols()
	if err != nil {
		return nil, err
	}

	if ses == nil {
		// TODO
	}

	return ses[0], nil
}

func (ss *symbolOperator) GetSecuritySymbol(securityID int) ([]string, error) {
	return ss.symbolRepo.GetSecuritySymbols(securityID)
}

func (ss *symbolOperator) SetSymbolSecurity(symbolID int, securityID int) error {
	valid, err := ss.symbolRepo.ValidSymbolID(symbolID)
	if err != nil {
		return err
	}
	if !valid {
		return errors.New(fmt.Sprintf("invalid symbol id: %d", symbolID))
	}

	valid, err = GetSecurityOperator().ValidSecurityID(securityID)
	if err != nil {
		return err
	}
	if !valid {
		return errors.New(fmt.Sprintf("invalid security id: %d", securityID))
	}

	_, err = ss.symbolRepo.UpdateSymbolSecurity(symbolID, securityID)
	if err != nil {
		return err
	}

	// set cache symbol-security.
	setCacheSymbolSecurity(symbolID, securityID)

	return nil
}

func (ss *symbolOperator) UpdateSymbolSecurity(symbolID int, oldSecurityID int, newSecurityID int) error {
	valid, err := ss.symbolRepo.ValidSymbolSecurity(symbolID, oldSecurityID)
	if err != nil {
		return err
	}
	if !valid {
		return errors.New(fmt.Sprintf("invalid symbolId: %d or oldSecurityID: %d", symbolID, oldSecurityID))
	}

	valid, err = GetSecurityOperator().ValidSecurityID(newSecurityID)
	if err != nil {
		return err
	}
	if !valid {
		return errors.New(fmt.Sprintf("invalid newSecurityID: %d", newSecurityID))
	}

	_, err = ss.symbolRepo.UpdateSymbolSecurity(symbolID, newSecurityID)
	if err != nil {
		return err
	}

	// set cache symbol-security.
	setCacheSymbolSecurity(symbolID, newSecurityID)

	return nil
}

func (ss *symbolOperator) GetIDByName(symbolName string) (int, error) {
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

func (ss *symbolOperator) ValidSymbolNameID(symbolName string, symbolID int) (valid bool, err error) {
	return ss.symbolRepo.ValidSymbolNameID(symbolName, symbolID)
}

func (ss *symbolOperator) SecurityHoldSymbols(securityID int) (hold bool, err error) {
	return ss.symbolRepo.SecurityHoldSymbols(securityID)
}

//cache
type symbolCache struct {
	id2Name map[int]string
	info    map[string]*Symbol
	sync.RWMutex
}

var symbCache = &symbolCache{
	id2Name: make(map[int]string),
	info:    make(map[string]*Symbol),
}

func LoadSymbolCache() {
	symbols := GetSymbolOperator().GetSymbols()
	if symbols == nil {
		panic("cache symbol failed")
	}

	symbCache.Lock()
	defer symbCache.Unlock()

	for i := range symbols {
		ss, err := GetSessionOperator().GetSessionsByID(symbols[i].ID)
		if err != nil || ss == nil {
			panic("get session err when load symbol cache")
		}
		sessions, err := EncodeSession(ss)
		if err != nil {
			panic("encode session err when load symbol cache")
		}

		symbols[i].QuoteSession = sessions["quote"]
		symbols[i].TradeSession = sessions["trade"]

		symbCache.id2Name[symbols[i].ID] = symbols[i].Symbol
		symbCache.info[symbols[i].Symbol] = &symbols[i]
	}
}

func getSymbolByIDFromCache(ID int) *Symbol {
	symbCache.RLock()
	defer symbCache.RUnlock()

	symbolName := symbCache.id2Name[ID]
	symbol, ok := symbCache.info[symbolName]
	if !ok {
		return nil
	}
	symb := *symbol

	return &symb
}

func getSymbolByNameFromCache(symbolName string) *Symbol {
	symbCache.RLock()
	defer symbCache.RUnlock()

	symbol, ok := symbCache.info[symbolName]
	if !ok {
		return nil
	}
	symb := *symbol

	return &symb
}

func getSymbolsFromCache() []Symbol {
	symbCache.RLock()
	defer symbCache.RUnlock()

	symbs := make([]Symbol, 0)
	for i := range symbCache.info {
		symbs = append(symbs, *symbCache.info[i])
	}

	return symbs
}

func insertSymbolToCache(symbol *Symbol) {
	symbCache.Lock()
	defer symbCache.Unlock()

	symbCache.id2Name[symbol.ID] = symbol.Symbol
	symbCache.info[symbol.Symbol] = symbol
}

func deleteSymbolCacheByID(ID int) {
	symbCache.Lock()
	defer symbCache.Unlock()

	symb := symbCache.id2Name[ID]
	delete(symbCache.id2Name, ID)
	delete(symbCache.info, symb)
}

func deleteSymbolCacheByName(symbolName string) {
	symbCache.Lock()
	defer symbCache.Unlock()

	for id := range symbCache.id2Name {
		if symbCache.id2Name[id] == symbolName {
			delete(symbCache.id2Name, id)
			break
		}
	}

	delete(symbCache.info, symbolName)
}

func setCacheSymbolSecurity(symbolID int, securityID int) {
	symbCache.Lock()
	defer symbCache.Unlock()

	symb := symbCache.id2Name[symbolID]
	symbCache.info[symb].SecurityID = securityID
}