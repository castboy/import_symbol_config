package server

import (
	"fmt"
	"sync"
	"time"
	"errors"
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
		return errors.New(fmt.Sprintf("invalid date: %s", holi.Date))
	}

	_, err = time.ParseInLocation("15:04:05", holi.From, time.UTC)
	if err != nil {
		return errors.New(fmt.Sprintf("invalid from: %s", holi.From))
	}

	_, err = time.ParseInLocation("15:04:05", holi.To, time.UTC)
	if err != nil {
		return errors.New(fmt.Sprintf("invalid to: %s", holi.To))
	}

	switch holi.Category {
	case HolidayAll:

	case HolidaySecurity:
		valid, err := GetSecurityOperator().ValidSecurityName(holi.Symbol)
		if err != nil {
			return err
		}
		if !valid {
			return errors.New(fmt.Sprintf("invalid security: %s", holi.Symbol))
		}

	case HolidaySymbol:
		valid, err := GetSymbolOperator().ValidSymbolName(holi.Symbol)
		if err != nil {
			return err
		}
		if !valid {
			return errors.New(fmt.Sprintf("invalid symbol: %s", holi.Symbol))
		}

	default:
		return errors.New(fmt.Sprintf("invalid category: %d", holi.Category))
	}

	return nil
}

func (ho *holidayOperator) InsertHoliday(holi *Holiday) error {
	err := holidayFormatCheck(holi)
	if err != nil {
		return err
	}
	_, err = ho.holidayRepo.Insert(holi)

	return err
}

func (ho *holidayOperator) UpdateHolidayByID(ID int, holi *Holiday) error {
	valid, err := ho.holidayRepo.ValidHolidayID(ID)
	if err != nil {
		return err
	}
	if !valid {
		return errors.New(fmt.Sprintf("invalid holiday id: %d", ID))
	}

	if err := holidayFormatCheck(holi); err != nil {
		return err
	}

	return ho.holidayRepo.UpdateByID(ID, holi)
}

func (ho *holidayOperator) IsTradable(symbol string) bool {
	if holiCache.notHoliday {
		return true
	}

	holiCache.RLock()
	defer holiCache.RUnlock()

	if !holiCache.todaySymbol[symbol] {
		return true
	}

	for _, v := range holiCache.symbolHolidays[symbol] {
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

type holidayCache struct {
	notHoliday     bool
	todaySymbol    map[string]bool
	symbolHolidays map[string][]holidayTime
	sync.RWMutex
}

var holiCache holidayCache

func LoadHolidayCacheByDate(date string) error {
	holidays, err := GetHolidayOperator().holidayRepo.GetByDate(date)
	if err != nil {
		panic(err)
	}

	if holidays == nil {
		holiCache.notHoliday = true
		return nil
	}

	holiCache.todaySymbol = make(map[string]bool)
	holiCache.symbolHolidays = make(map[string][]holidayTime)

	holiCache.Lock()
	defer holiCache.Unlock()

	for i := range holidays {
		switch holidays[i].Category {
		case HolidayAll:
			symbols, err := GetSymbolOperator().GetSymbolsName()
			if err != nil {
				panic(err)
			}

			for j := range symbols {
				holiCache.todaySymbol[symbols[j]] = true

				if holiCache.symbolHolidays[symbols[j]] == nil {
					holiCache.symbolHolidays[symbols[j]] = make([]holidayTime, 0)
				}

				holiCache.symbolHolidays[symbols[j]] = append(holiCache.symbolHolidays[symbols[j]],
					holidayTime{from: holidays[i].From, to: holidays[i].To})
			}

		case HolidaySecurity:
			securityID, err := GetSecurityOperator().GetIDByName(holidays[i].Symbol)
			if err != nil {
				panic(err)
			}

			symbols, err := GetSymbolOperator().GetSymbolsNameBySecurityID(securityID)
			if err != nil {
				panic(err)
			}

			for j := range symbols {
				holiCache.todaySymbol[symbols[j]] = true

				if holiCache.symbolHolidays[symbols[j]] == nil {
					holiCache.symbolHolidays[symbols[j]] = make([]holidayTime, 0)
				}

				holiCache.symbolHolidays[symbols[j]] = append(holiCache.symbolHolidays[symbols[j]],
					holidayTime{from: holidays[i].From, to: holidays[i].To})
			}

		case HolidaySymbol:
			holiCache.todaySymbol[holidays[i].Symbol] = true

			if holiCache.symbolHolidays[holidays[i].Symbol] == nil {
				holiCache.symbolHolidays[holidays[i].Symbol] = make([]holidayTime, 0)
			}

			holiCache.symbolHolidays[holidays[i].Symbol] = append(holiCache.symbolHolidays[holidays[i].Symbol],
				holidayTime{from: holidays[i].From, to: holidays[i].To})

		default:
			panic(fmt.Sprintf("invalid holiday category: %d", holidays[i].Category))
		}
	}

	return nil
}
