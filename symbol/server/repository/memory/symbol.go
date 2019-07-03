/*******************************************************************************
 * // Copyright AnchyTec Corp. All Rights Reserved.
 * // SPDX-License-Identifier: Apache-2.0
 * // Author: shaozhiming
 ******************************************************************************/

package memory

import (
	"sync"
	"import_symbol_config/symbol/server"
)

type symbolRepository struct {
	repository sync.Map
}

var symbolRep *symbolRepository

func GetSymbolRepository() *symbolRepository {
	if symbolRep == nil {
		symbolRep = &symbolRepository{}
	}

	return symbolRep
}

func (sr *symbolRepository) GetSymbolInfoByName(symbolName string) (symbol *server.Symbol) {
	if symbol, ok := sr.repository.Load(symbolName); ok {
		return symbol.(*server.Symbol)
	} else {
		return nil
	}
}

func (sr *symbolRepository) GetSymbolInfoByID(ID int) (symbol *server.Symbol) {
	//TODO
	return nil
}

func (sr *symbolRepository) GetSymbols() (symbols []server.Symbol) {
	symbols = make([]server.Symbol, 0)

	sr.repository.Range(func(key, value interface{}) bool {
		symbolV := value.(*server.Symbol)
		symbols = append(symbols, *symbolV)
		return true
	})

	return
}

// create new record, update existing record, delete existing record.

func (sr *symbolRepository) Insert(symbol *server.Symbol) error {
	sr.repository.Store(symbol.Symbol, symbol)
	return nil
}

// the following code is just to avoid reporting errors and is not implemented.
// we will switch to the mysql implementation later.

func (sr *symbolRepository) UpdateByID(ID int, symbol *server.Symbol) error {
	return nil
}

func (sr *symbolRepository) UpdateByName(symbolName string, symbol *server.Symbol) error {
	return nil
}

func (sr *symbolRepository) DeleteByID(ID int) error {
	return nil
}

func (sr *symbolRepository) DeleteByName(symbolName string) error {
	return nil
}

func (sr *symbolRepository) InsertSessions(sess ...*server.Session) error {
	return nil
}

func (sr *symbolRepository) GetSessionsByName(symbolName string) (sess []*server.Session, err error) {
	return
}

func (sr *symbolRepository) GetSessionsByID(sessionID int) (sess []*server.Session, err error) {
	return
}

func (sr *symbolRepository) UpdateSession(sess *server.Session) error {
	return nil
}

func (sr *symbolRepository) DeleteSession(sessionID int) error {
	return nil
}
