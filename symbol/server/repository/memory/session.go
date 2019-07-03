package memory

import (
	"import_symbol_config/symbol/server"
)

type sessionRepository struct {
}

var sessionRep *sessionRepository

func GetSessionRepository() *sessionRepository {
	if sessionRep == nil {
		sessionRep = &sessionRepository{}
	}

	return sessionRep
}

func (sr *sessionRepository) Insert(sess ...*server.Session) error {
	return nil
}

func (sr *sessionRepository) GetByName(symbolName string) (sess []*server.Session, err error) {
	return nil, nil
}

func (sr *sessionRepository) GetByID(ID int) (sess []*server.Session, err error) {
	return nil, nil
}

func (sr *sessionRepository) Update(sess *server.Session) error {
	return nil
}

func (sr *sessionRepository) Delete(sessionID int) error {
	return nil
}
