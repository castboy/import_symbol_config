package server

import (
	"github.com/juju/errors"
	"sync"
	"time"
)

type Holiday struct {
	ID          int             `json:"id" xorm:"id"`
	Enable      bool            `json:"enable" xorm:"enable"`
	Date        string          `json:"date" xorm:"date"`
	From        string          `json:"from" xorm:"from"`
	To          string          `json:"to" xorm:"to"`
	Category    HolidayCategory `json:"category" xorm:"category"`
	Symbol      string          `json:"symbol" xorm:"symbol"`
	Description string          `json:"description" xorm:"description"`
}

type HolidayCategory int

const (
	HolidayAll HolidayCategory = iota
	HolidaySecurity
	HolidaySymbol
)

type HolidayRepository interface {
	Insert(holi *Holiday) (num int64, err error)
	UpdateByID(ID int, holi *Holiday) error
	GetByDate(date string) (holidays []*Holiday, err error)
	ValidHolidayID(ID int) (valid bool, err error)
}

type holidayOperator struct {
	holidayRepo HolidayRepository
}

var holidayOp *holidayOperator

func GetHolidayOperator() *holidayOperator {
	return holidayOp
}

func InitHolidayOperator(holidayRepo HolidayRepository) *holidayOperator {
	if holidayOp == nil {
		holidayOp = &holidayOperator{
			holidayRepo,
		}
	}
	return holidayOp
}

func holidayFormatCheck(holi *Holiday) error {
	_, err := time.ParseInLocation("2006-01-02", holi.Date, time.UTC)
	if err != nil {
		return errors.NotValidf("date %s", holi.Date)
	}

	_, err = time.ParseInLocation("15:04:05", holi.From, time.UTC)
	if err != nil {
		return errors.NotValidf("from %s", holi.From)
	}

	_, err = time.ParseInLocation("15:04:05", holi.To, time.UTC)
	if err != nil {
		return errors.NotValidf("to %s", holi.To)
	}

	switch holi.Category {
	case HolidayAll:

	case HolidaySecurity:
		valid, err := GetSecurityOperator().ValidSecurityName(holi.Symbol)
		if err != nil {
			return errors.Annotate(err, "sql exec")
		}
		if !valid {
			return errors.NotValidf("security %s", holi.Symbol)
		}

	case HolidaySymbol:
		valid, err := GetSymbolOperator().ValidSymbolName(holi.Symbol)
		if err != nil {
			return errors.Annotate(err, "sql exec")
		}
		if !valid {
			return errors.NotValidf("symbol %s", holi.Symbol)
		}

	default:
		return errors.NotValidf("category %d", holi.Category)
	}

	return nil
}

func (ho *holidayOperator) InsertHoliday(holi *Holiday) error {
	if err := holidayFormatCheck(holi); err != nil {
		return errors.NewNotValid(err, "validation holiday")
	}

	if _, err := ho.holidayRepo.Insert(holi); err != nil {
		return errors.Annotate(err, "sql exec")
	}

	return nil
}

func (ho *holidayOperator) UpdateHolidayByID(ID int, holi *Holiday) error {
	valid, err := ho.holidayRepo.ValidHolidayID(ID)
	if err != nil {
		return errors.Annotate(err, "sql exec")
	}
	if !valid {
		return errors.NotValidf("holiday id %d", ID)
	}

	if err := holidayFormatCheck(holi); err != nil {
		return errors.NewNotValid(err, "validation holiday")
	}

	if err := ho.holidayRepo.UpdateByID(ID, holi); err != nil {
		return errors.Annotate(err, "sql exec")
	}

	return nil
}

func (ho *holidayOperator) IsTradable(symbol string) bool {
	holiCaches.RLock()
	defer holiCaches.RUnlock()

	date := time.Now().UTC().Format("2006-01-02")
	if _, ok := holiCaches.info[date]; !ok {
		return true
	}

	if !holiCaches.info[date].todaySymbol[symbol] {
		return true
	}

	for _, v := range holiCaches.info[date].symbolHolidays[symbol] {
		if v.from == "00:00:00" && v.to == "00:00:00" {
			return false
		}

		now := time.Now().UTC().Format("15:04:05")
		if now >= v.from && now <= v.to {
			return true
		}
	}

	return false
}

//cache
type holidayTime struct {
	from string
	to   string
}

type holidayCaches struct {
	info map[string]*holidayCache
	sync.RWMutex
}

type holidayCache struct {
	todaySymbol    map[string]bool
	symbolHolidays map[string][]holidayTime
}

var holiCaches = holidayCaches{info: make(map[string]*holidayCache)}

func LoadHolidayCacheByDate(date string) error {
	holidays, err := GetHolidayOperator().holidayRepo.GetByDate(date)
	if err != nil {
		return errors.Annotate(err, "sql exec")
	}

	if len(holidays) == 0 {
		return nil
	}

	holiCaches.Lock()
	defer holiCaches.Unlock()

	if _, ok := holiCaches.info[date]; ok {
		return nil
	}

	holiCaches.info[date] = &holidayCache{}
	holiCaches.info[date].todaySymbol = make(map[string]bool)
	holiCaches.info[date].symbolHolidays = make(map[string][]holidayTime)

	for i := range holidays {
		switch holidays[i].Category {
		case HolidayAll:
			symbols, err := GetSymbolOperator().GetSymbolsName()
			if err != nil {
				return errors.Annotate(err, "sql exec")
			}

			for j := range symbols {
				holiCaches.info[date].todaySymbol[symbols[j]] = true

				if holiCaches.info[date].symbolHolidays[symbols[j]] == nil {
					holiCaches.info[date].symbolHolidays[symbols[j]] = make([]holidayTime, 0)
				}

				holiCaches.info[date].symbolHolidays[symbols[j]] = append(holiCaches.info[date].symbolHolidays[symbols[j]],
					holidayTime{from: holidays[i].From, to: holidays[i].To})
			}

		case HolidaySecurity:
			securityID, exist, err := GetSecurityOperator().GetIDByName(holidays[i].Symbol)
			if err != nil {
				return errors.Annotate(err, "sql exec")
			}

			if !exist {
				return errors.NotFoundf("security name %s", holidays[i].Symbol)
			}

			symbols, err := GetSymbolOperator().GetSymbolsNameBySecurityID(securityID)
			if err != nil {
				return errors.Annotate(err, "sql exec")
			}

			for j := range symbols {
				holiCaches.info[date].todaySymbol[symbols[j]] = true

				if holiCaches.info[date].symbolHolidays[symbols[j]] == nil {
					holiCaches.info[date].symbolHolidays[symbols[j]] = make([]holidayTime, 0)
				}

				holiCaches.info[date].symbolHolidays[symbols[j]] = append(holiCaches.info[date].symbolHolidays[symbols[j]],
					holidayTime{from: holidays[i].From, to: holidays[i].To})
			}

		case HolidaySymbol:
			holiCaches.info[date].todaySymbol[holidays[i].Symbol] = true

			if holiCaches.info[date].symbolHolidays[holidays[i].Symbol] == nil {
				holiCaches.info[date].symbolHolidays[holidays[i].Symbol] = make([]holidayTime, 0)
			}

			holiCaches.info[date].symbolHolidays[holidays[i].Symbol] = append(holiCaches.info[date].symbolHolidays[holidays[i].Symbol],
				holidayTime{from: holidays[i].From, to: holidays[i].To})

		default:
			err = errors.NotValidf("category %d", holidays[i].Category)
			panic(err)
		}
	}

	return nil
}

func RmHolidayCacheBeforeToday() {
	holiCaches.Lock()
	defer holiCaches.Unlock()

	today := time.Now().UTC().Format("2006-01-02")

	for day := range holiCaches.info {
		if day < today {
			delete(holiCaches.info, day)
		}
	}
}
