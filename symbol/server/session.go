package server

import (
	"github.com/juju/errors"
	"strings"
	"time"
)

type Session struct {
	ID       int          `xorm:"id"`
	SymbolID int          `xorm:"symbol_id"`
	Type     SessionType  `xorm:"type"`
	Weekday  time.Weekday `xorm:"weekday"`
	TimeSpan string       `xorm:"time_span"`
}

type SessionType int

const (
	Quote SessionType = iota
	Trade
)

func NewSession(SymbolID int, Type SessionType, Weekday time.Weekday, TimeSpan string) *Session {
	return &Session{
		SymbolID: SymbolID,
		Type:     Type,
		Weekday:  Weekday,
		TimeSpan: TimeSpan,
	}
}

type SessionRepository interface {
	Insert(sess ...*Session) error
	GetByName(symbolName string) (sess []*Session, err error)
	GetByID(symbolID int) (sess []*Session, err error)
	UpdateByID(sessionID int, sess *Session) error
	DeleteByID(sessionID int) (int64, error)
	ValidSessionID(sessionID int) (valid bool, err error)
	GetSymbolIDBySessionID(sessionID int) (symbolID int, exist bool, err error)
}

type sessionOperator struct {
	sessionRepo SessionRepository
}

var sessionOp *sessionOperator

func GetSessionOperator() *sessionOperator {
	return sessionOp
}

func InitSessionOperator(sessionRepo SessionRepository) *sessionOperator {
	if sessionOp == nil {
		sessionOp = &sessionOperator{
			sessionRepo,
		}
	}
	return sessionOp
}

func sessionFormatCheck(sess *Session) error {
	valid, err := GetSymbolOperator().ValidSymbolID(sess.SymbolID)
	if err != nil {
		err = errors.Annotatef(err, "sql exec")
		return err
	}
	if !valid {
		return errors.NotFoundf("symbol id %d", sess.SymbolID)
	}

	if sess.Type != Quote && sess.Type != Trade {
		return errors.NotValidf("type %s", sess.Type)
	}

	if sess.Weekday != time.Sunday && sess.Weekday != time.Monday && sess.Weekday != time.Tuesday && sess.Weekday != time.Wednesday &&
		sess.Weekday != time.Tuesday && sess.Weekday != time.Friday && sess.Weekday != time.Friday {
		return errors.NotValidf("weekday %s", sess.Weekday)
	}

	// TODO sess.TimeSpan check
	return nil
}

func (so *sessionOperator) InsertSessions(sess ...*Session) error {
	//for i := range sess {
	//	if err := sessionFormatCheck(sess[i]); err != nil {
	//		return errors.NewNotValid(err, "validation session")
	//	}
	//}

	err := so.sessionRepo.Insert(sess...)
	if err != nil {
		return errors.Annotatef(err, "sql exec")
	}

	// delete cache to re-cache.

	return nil
}

func (so *sessionOperator) GetSessionsByName(symbolName string) (sess []*Session, err error) {
	return so.sessionRepo.GetByName(symbolName)
}

func (so *sessionOperator) GetSessionsByID(symbolID int) (sess []*Session, err error) {
	return so.sessionRepo.GetByID(symbolID)
}

func (so *sessionOperator) UpdateSessionByID(sessionID int, sess *Session) error {
	if err := sessionFormatCheck(sess); err != nil {
		return errors.NewNotValid(err, "validation session")
	}

	valid, err := so.sessionRepo.ValidSessionID(sessionID)
	if err != nil {
		return errors.Annotatef(err, "sql exec")
	}
	if !valid {
		return errors.NotValidf("session id %d", sessionID)
	}

	err = so.sessionRepo.UpdateByID(sessionID, sess)
	if err != nil {
		return errors.Annotatef(err, "sql exec")
	}

	// delete cache to re-cache.

	return nil
}

func (so *sessionOperator) DeleteSessionByID(sessionID int) error {
	_, exist, err := so.sessionRepo.GetSymbolIDBySessionID(sessionID)
	if err != nil {
		return errors.Annotatef(err, "sql exec")
	}

	if !exist {
		return errors.NotFoundf("session id %d", sessionID)
	}

	hit, err := so.sessionRepo.DeleteByID(sessionID)
	if err != nil {
		return errors.Annotatef(err, "sql exec")
	}

	if hit == 0 {
		return errors.NotFoundf("session id %d", sessionID)
	}

	// delete cache to re-cache.

	return nil
}

// encode or decode session.
func EncodeSession(sess []*Session) (qt map[SessionType]map[time.Weekday]string, err error) {
	qt = map[SessionType]map[time.Weekday]string{
		Quote: map[time.Weekday]string{time.Saturday: "", time.Sunday: ""},
		Trade: map[time.Weekday]string{time.Saturday: "", time.Sunday: ""},
	}

	for _, session := range sess {
		if qt[session.Type][session.Weekday] == "" {
			qt[session.Type][session.Weekday] = session.TimeSpan
		} else {
			qt[session.Type][session.Weekday] += "," + session.TimeSpan
		}
	}

	for session, _ := range qt {
		for weekday, _ := range qt[session] {
			L := len(qt[session][weekday])
			switch {
			case L < len("00:00-00:00"):
				qt[session][weekday] = "00:00-00:00,00:00-00:00,00:00-00:00"
			case L < len("00:00-00:00,00:00-00:00"):
				qt[session][weekday] += ",00:00-00:00,00:00-00:00"
			default:
				qt[session][weekday] += ",00:00-00:00"
			}
		}
	}

	return
}

func DecodeSession(symbolID int, sessionType SessionType, sessions map[time.Weekday]string) []*Session {
	sess := make([]*Session, 0)
	for weekday, se := range sessions {
		sl := strings.Split(se, ",")
		for _, time := range sl {
			if time == "00:00-00:00" {
				continue
			}
			sess = append(sess, NewSession(symbolID, sessionType, weekday, time))
		}
	}

	return sess
}
