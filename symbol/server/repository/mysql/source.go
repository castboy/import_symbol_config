package mysql

import (
_ "github.com/go-sql-driver/mysql"
"github.com/go-xorm/xorm"
"import_symbol_config/symbol/server"
)

type sourceRepository struct {
	engine *xorm.Engine
}

var sourceRep *sourceRepository

func GetSourceRepository() *sourceRepository {
	if sourceRep == nil {
		sourceRep = &sourceRepository{
			engine: xEngine,
		}
	}

	return sourceRep
}

func (sr *sourceRepository) GetIDByName(source string) (ID int, exist bool, err error) {
	exist, err = sr.engine.Table(server.Source{}).Select("id").Where("source=?", source).Get(&ID)
	return
}

func (sr *sourceRepository) InsertSource(sec *server.Source) (err error) {
	_, err = sr.engine.Table(server.Source{}).Omit("id").Insert(sec)
	return
}