package memory

import (
	"import_symbol_config/symbol/server"
)

type holidayRepository struct {
}

var holidayRep *holidayRepository

func GetHolidayRepository() *holidayRepository {
	if holidayRep == nil {
		holidayRep = &holidayRepository{}
	}

	return holidayRep
}

func (hr *holidayRepository) Insert(holi *server.Holiday) error {
	return nil
}

func (hr *holidayRepository) UpdateByID(ID int, holi *server.Holiday) error {
	return nil
}

func (hr *holidayRepository) IsTradable(symbolName string) (bool, error) {
	return true, nil
}
