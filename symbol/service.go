/*******************************************************************************
 * Copyright (c) 2019. AnchyTec Corp. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 * Author: shaozhiming
 ******************************************************************************/

package symbol

import (
	"import_symbol_config/symbol/server"
	"import_symbol_config/symbol/server/repository/mysql"
)

var symbolService SymbolService

//go:generate mockgen -destination ../cmd/trading-core/mock/symbol/mock_service.go -package symbol import_symbol_config/symbol SymbolService
type SymbolService interface {
	SymbolOperater
	SessionOperater
	HolidayOperater
	SecurityOperater
}

type SymbolOperater interface {
	GetSymbolInfoByName(symbolName string) (symbol *server.Symbol)
	GetSymbolInfoByID(ID int) (symbol *server.Symbol)
	GetSymbols() (symbols []server.Symbol)

	// session can be inserted at the same time, if session field is not null.
	InsertSymbol(symbol *server.Symbol) error
	// only update symbol not include session below two methods.
	UpdateSymbolByID(ID int, symbol *server.Symbol) error
	UpdateSymbolByName(symbolName string, symbol *server.Symbol) error
	// delete related session at the same time below two methods.
	DeleteSymbolByID(ID int) error
	DeleteSymbolByName(symbolName string) error

	SetSymbolSecurity(symbolID int, securityID int) error
	UpdateSymbolSecurity(symbolID int, oldSecurityID int, newSecurityID int) error

	Start()
}

type SessionOperater interface {
	InsertSessions(sess ...*server.Session) error
	GetSessionsByName(symbolName string) (sess []*server.Session, err error)
	GetSessionsByID(symbolID int) (sess []*server.Session, err error)
	UpdateSessionByID(sessionID int, sess *server.Session) error
	DeleteSessionByID(sessionID int) error
}

type HolidayOperater interface {
	InsertHoliday(holi *server.Holiday) error
	UpdateHolidayByID(ID int, holi *server.Holiday) error
	IsTradable(symbolName string) bool
}

type SecurityOperater interface {
	GetSecurityInfo(id int) (*server.Security, error)
	GetAllSecuritiesInfos() ([]*server.Security, error)
	UpdateSecurityInfo(id int, info *server.Security) error
	InsertSecurityInfo(info *server.Security) error
	DeleteSecurityInfo(id int) error
}

type Operators struct {
	SymbolOperater
	SessionOperater
	HolidayOperater
	SecurityOperater
}

func NewOperators(syo SymbolOperater, seo SessionOperater, ho HolidayOperater, se SecurityOperater) *Operators {
	return &Operators{
		SymbolOperater:   syo,
		SessionOperater:  seo,
		HolidayOperater:  ho,
		SecurityOperater: se,
	}
}

func GetSymbolService() SymbolService {
	if symbolService == nil {
		syr := mysql.GetSymbolRepository()
		syo := server.InitSymbolOperator(syr)

		ser := mysql.GetSessionRepository()
		seo := server.InitSessionOperator(ser)

		hr := mysql.GetHolidayRepository()
		ho := server.InitHolidayOperator(hr)

		secr := mysql.GetSecurityRepository()
		seco := server.InitSecurityOperator(secr)

		symbolService = NewOperators(syo, seo, ho, seco)
	}
	return symbolService
}
