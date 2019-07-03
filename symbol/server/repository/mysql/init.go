package mysql

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"time"
	"import_symbol_config/config"
)

var xEngine *xorm.Engine

func init() {
	cnf := config.GetConfigService("symbol")
	cnfRoot := "mysql"
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		cnf.GetString(cnfRoot+".user"), cnf.GetString(cnfRoot+".pwd"), cnf.GetString(cnfRoot+".host"),
		cnf.GetString(cnfRoot+".port"), cnf.GetString(cnfRoot+".dbname"))

	engine, err := xorm.NewEngine("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}

	engine.DatabaseTZ = time.UTC

	xEngine = engine
}
