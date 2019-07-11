package mysql

import (
	"github.com/go-xorm/xorm"
	"import_symbol_config/symbol/server"
)

type sessionRepository struct {
	engine *xorm.Engine
}

var sessionRep *sessionRepository

func GetSessionRepository() *sessionRepository {
	if sessionRep == nil {
		sessionRep = &sessionRepository{
			engine: xEngine,
		}
	}

	return sessionRep
}

// Session crud opeation

func (sr *sessionRepository) Insert(sess ...*server.Session) error {
	_, err := sr.engine.Table(server.Session{}).Insert(sess)
	return err
}

func (sr *sessionRepository) GetByName(symbolName string) (sess []*server.Session, err error) {
	err = sr.engine.Table(server.Session{}).Where("symbol=?", symbolName).OrderBy("type, weekday, time_span").Find(&sess)
	return
}

func (sr *sessionRepository) GetByID(symbolID int) (sess []*server.Session, err error) {
	err = sr.engine.Table(server.Session{}).Where("symbol_id=?", symbolID).OrderBy("type, weekday, time_span").Find(&sess)
	return
}

func (sr *sessionRepository) UpdateByID(sessionID int, sess *server.Session) error {
	_, err := sr.engine.Table(server.Session{}).AllCols().Where("id=?", sessionID).Update(sess)
	return err
}

func (sr *sessionRepository) DeleteByID(sessionID int) (int64, error) {
	return sr.engine.Table(server.Session{}).Where("id=?", sessionID).Delete(server.Session{})
}

func (sr *sessionRepository) ValidSessionID(sessionID int) (valid bool, err error) {
	return sr.engine.Table(server.Session{}).Where("id=?", sessionID).Exist()
}

func (sr *sessionRepository) GetSymbolIDBySessionID(sessionID int) (symbolID int, exist bool, err error) {
	exist, err = sr.engine.Table(server.Session{}).Select("symbol_id").Where("id=?", sessionID).Get(&symbolID)
	return
}
