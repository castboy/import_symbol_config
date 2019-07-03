package server

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"errors"
)

type Session struct {
	ID       int    `xorm:"id"`
	SymbolID int    `xorm:"symbol_id"`
	Symbol   string `xorm:"symbol"`
	Type     string `xorm:"type"`
	Weekday  string `xorm:"weekday"`
	TimeSpan string `xorm:"time_span"`
}

func NewSession(SymbolID int, Symbol, Type, Weekday, TimeSpan string) *Session {
	return &Session{
		SymbolID: SymbolID,
		Symbol:   Symbol,
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
	DeleteByID(sessionID int) error
	ValidSessionID(sessionID int) (valid bool, err error)
	GetSymbolIDBySessionID(sessionID int) (symbolID int, err error)
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
	valid, err := GetSymbolOperator().ValidSymbolNameID(sess.Symbol, sess.SymbolID)
	if err != nil {
		return err
	}
	if !valid {
		return errors.New(fmt.Sprintf("invalid symbol id: %d or symbol name: %s", sess.SymbolID, sess.Symbol))
	}

	if sess.Type != "quote" && sess.Type != "trade" {
		return errors.New(fmt.Sprintf("invalid session type: %s", sess.Type))
	}

	if sess.Weekday != "0" && sess.Weekday != "1" && sess.Weekday != "2" && sess.Weekday != "3" &&
		sess.Weekday != "4" && sess.Weekday != "5" && sess.Weekday != "6" {
		return errors.New(fmt.Sprintf("invalid session weekday: %s", sess.Weekday))
	}

	// TODO sess.TimeSpan check
	return nil
}

func (so *sessionOperator) InsertSessions(sess ...*Session) error {
	for i := range sess {
		if err := sessionFormatCheck(sess[i]); err != nil {
			return err
		}
	}

	err := so.sessionRepo.Insert(sess...)
	if err != nil {
		return err
	}

	// update cache.
	symbolName := sess[0].Symbol
	deleteSymbolCacheByName(symbolName)
	symbol := GetSymbolOperator().GetSymbolInfoByName(symbolName)
	if symbol == nil {
		return errors.New("update symbol cache failed by name")
	}
	insertSymbolToCache(symbol)

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
		return err
	}

	valid, err := so.sessionRepo.ValidSessionID(sessionID)
	if err != nil {
		return err
	}
	if !valid {
		return errors.New(fmt.Sprintf("invalid session id: %d", sessionID))
	}

	err = so.sessionRepo.UpdateByID(sessionID, sess)
	if err != nil {
		return err
	}

	// update cache.
	symbolName := sess.Symbol
	deleteSymbolCacheByName(symbolName)
	symbol := GetSymbolOperator().GetSymbolInfoByName(symbolName)
	if symbol == nil {
		return errors.New("get symbol info failed by name, when UpdateSessionByID()")
	}
	insertSymbolToCache(symbol)

	return nil
}

func (so *sessionOperator) DeleteSessionByID(sessionID int) error {
	symbolID, err := so.sessionRepo.GetSymbolIDBySessionID(sessionID)
	if err != nil {
		return err
	}

	err = so.sessionRepo.DeleteByID(sessionID)
	if err != nil {
		return err
	}

	// update cache.
	deleteSymbolCacheByID(symbolID)
	symbol := GetSymbolOperator().GetSymbolInfoByID(symbolID)
	if symbol == nil {
		return errors.New("get symbol info failed by id, when DeleteSessionByID()")
	}
	insertSymbolToCache(symbol)

	return nil
}

// encode or decode session.
func EncodeSession(sess []*Session) (qt map[string]map[time.Weekday]string, err error) {
	qt = map[string]map[time.Weekday]string{
		"quote": map[time.Weekday]string{time.Saturday: "", time.Sunday: ""},
		"trade": map[time.Weekday]string{time.Saturday: "", time.Sunday: ""},
	}

	for _, session := range sess {
		weekday, err := strconv.Atoi(session.Weekday)
		if err != nil {
			return nil, err
		}

		if qt[session.Type][time.Weekday(weekday)] == "" {
			qt[session.Type][time.Weekday(weekday)] = session.TimeSpan
		} else {
			qt[session.Type][time.Weekday(weekday)] += "," + session.TimeSpan
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

func DecodeSession(symbolID int, symbolName, sessionType string, sessions map[time.Weekday]string) []*Session {
	sess := make([]*Session, 0)
	for weekday, se := range sessions {
		sl := strings.Split(se, ",")
		for _, time := range sl {
			if time == "00:00-00:00" {
				continue
			}
			sess = append(sess, NewSession(symbolID, symbolName, sessionType, strconv.Itoa(int(weekday)), time))
		}
	}

	return sess
}
