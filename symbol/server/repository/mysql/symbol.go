package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"errors"
	"import_symbol_config/symbol/server"
)

type symbolRepository struct {
	engine *xorm.Engine
}

var symbolRep *symbolRepository

func GetSymbolRepository() *symbolRepository {
	if symbolRep == nil {
		symbolRep = &symbolRepository{
			engine: xEngine,
		}
	}

	return symbolRep
}

func (sr *symbolRepository) GetSymbolInfoByName(symbolName string) (symbol *server.Symbol, err error) {
	symbol = new(server.Symbol)
	hit, err := sr.engine.Table(server.Symbol{}).Where("symbol=?", symbolName).NoAutoCondition(true).Get(symbol)
	if err == nil && !hit {
		err = errors.New(fmt.Sprintf("invalid symbol name: %s", symbolName))
	}
	return
}

func (sr *symbolRepository) GetSymbolInfoByID(ID int) (symbol *server.Symbol, err error) {
	symbol = new(server.Symbol)
	hit, err := sr.engine.Table(server.Symbol{}).Where("id=?", ID).NoAutoCondition(true).Get(symbol)
	if err == nil && !hit {
		err = errors.New(fmt.Sprintf("invalid symbol id: %d", ID))
	}
	return
}

func (sr *symbolRepository) GetSymbols() (symbols []server.Symbol, err error) {
	err = sr.engine.Table(server.Symbol{}).Find(&symbols)
	return
}

func (sr *symbolRepository) InsertSymbol(symbol *server.Symbol) error {
	_, err := sr.engine.Table(server.Symbol{}).Omit("id").Insert(symbol)
	return err
}

func (sr *symbolRepository) InsertSession(sess []*server.Session) error {
	_, err := sr.engine.Table(server.Session{}).Insert(sess)
	return err
}

func (sr *symbolRepository) GetIDByName(symbolName string) (ID int, err error) {
	hit, err := sr.engine.Table(server.Symbol{}).Select("id").Where("symbol=?", symbolName).Get(&ID)
	if err == nil && !hit {
		err = errors.New(fmt.Sprintf("invalid symbol: %s", symbolName))
	}
	return
}

func (sr *symbolRepository) UpdateByID(ID int, symbol *server.Symbol) error {
	_, err := sr.engine.Table(server.Symbol{}).Where("id=?", ID).AllCols().Update(symbol)
	return err
}

func (sr *symbolRepository) UpdateByName(symbolName string, symbol *server.Symbol) error {
	_, err := sr.engine.Table(server.Symbol{}).Where("symbol=?", symbolName).AllCols().Update(symbol)
	return err
}

func (sr *symbolRepository) NewTransaction() *xorm.Session {
	return sr.engine.NewSession()
}

func (sr *symbolRepository) TransactionDeleteByName(ss *xorm.Session, tableName interface{}, symbolName string) (num int64, err error) {
	return ss.Where("symbol=?", symbolName).NoAutoCondition(true).Delete(tableName)
}

func (sr *symbolRepository) TransactionDeleteSymbolByID(ss *xorm.Session, ID int) (num int64, err error) {
	return ss.Where("id=?", ID).NoAutoCondition(true).Delete(server.Symbol{})
}

func (sr *symbolRepository) TransactionDeleteSessionByID(ss *xorm.Session, ID int) error {
	_, err := ss.Where("symbol_id=?", ID).Delete(server.Session{})
	return err
}

func (sr *symbolRepository) GetAllSecuritySymbols() ([]map[string]string, error) {
	return sr.engine.QueryString("select security_id, group_concat(symbol Separator ',') as symbol from symbol group by security_id")
}

func (sr *symbolRepository) GetSecuritySymbols(securityID int) (symbols []string, err error) {
	err = sr.engine.Table(server.Symbol{}).Select("symbol").Where("security_id=?", securityID).Find(&symbols)
	return
}

func (sr *symbolRepository) UpdateSymbolSecurity(symbolID int, securityID int) (num int64, err error) {
	symbol := new(server.Symbol)
	symbol.SecurityID = securityID
	return sr.engine.Table(server.Symbol{}).Cols("security_id").Where("id=?", symbolID).Update(symbol)
}

func (sr *symbolRepository) ValidSymbolID(ID int) (valid bool, err error) {
	return sr.engine.Table(server.Symbol{}).Where("id=?", ID).Exist()
}

func (sr *symbolRepository) ValidSymbolName(symbolName string) (valid bool, err error) {
	return sr.engine.Table(server.Symbol{}).Where("symbol=?", symbolName).Exist()
}

func (sr *symbolRepository) ValidSymbolNameID(symbolName string, symbolID int) (valid bool, err error) {
	return sr.engine.Table(server.Symbol{}).Where("symbol=? and id=?", symbolName, symbolID).Exist()
}

func (sr *symbolRepository) ValidSymbolSecurity(symbolID int, securityID int) (valid bool, err error) {
	return sr.engine.Table(server.Symbol{}).Where("id=? and security_id=?", symbolID, securityID).Exist()
}

func (sr *symbolRepository) GetSymbolsNameBySecurityID(securityID int) (symbols []string, err error) {
	err = sr.engine.Table(server.Symbol{}).Select("symbol").Where("security_id=?", securityID).Find(&symbols)
	return
}

func (sr *symbolRepository) GetSymbolsName() (symbols []string, err error) {
	err = sr.engine.Table(server.Symbol{}).Select("symbol").Find(&symbols)
	return
}

func (sr *symbolRepository) SecurityHoldSymbols(securityID int) (hold bool, err error) {
	return sr.engine.Table(server.Symbol{}).Where("security_id=?", securityID).Exist()
}
