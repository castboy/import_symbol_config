package mysql

import (
	"github.com/go-xorm/xorm"
	"import_symbol_config/symbol/server"
)

type holidayRepository struct {
	engine *xorm.Engine
}

var holidayRep *holidayRepository

func GetHolidayRepository() *holidayRepository {
	if holidayRep == nil {
		holidayRep = &holidayRepository{
			engine: xEngine,
		}
	}

	return holidayRep
}

func (hr *holidayRepository) GetByDate(date string) (holidays []*server.Holiday, err error) {
	err = hr.engine.Table(server.Holiday{}).Where("date=? and enable=1", date).Find(&holidays)
	return
}

func (hr *holidayRepository) Insert(holi *server.Holiday) (num int64, err error) {
	return hr.engine.Table(server.Holiday{}).Insert(holi)
}

func (hr *holidayRepository) UpdateByID(ID int, holi *server.Holiday) error {
	_, err := hr.engine.Table(server.Holiday{}).AllCols().Where("id=?", ID).Update(holi)
	return err
}

func (hr *holidayRepository) ValidHolidayID(ID int) (valid bool, err error) {
	return hr.engine.Table(server.Holiday{}).Where("id=?", ID).Exist()
}
