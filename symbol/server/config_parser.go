package server

import (
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"reflect"
	"time"
	"import_symbol_config/config"
)

var iViper = config.GetConfigService("symbol")
var decodeConfigTagName = iViper.GetString("decoder_config_tagname")

func symbolDecodeHook() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() == reflect.Float64 && t == reflect.TypeOf(decimal.Decimal{}) {
			return decimal.NewFromFloat(data.(float64)), nil
		}

		return data, nil
	}
}

func holidayDecodeHook() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() == reflect.String && t == reflect.TypeOf(time.Time{}) {
			return time.ParseInLocation("2006-01-02 15:04:05", data.(string), time.UTC)
		}

		if f.Kind() == reflect.String && t.Kind() == reflect.Int {
			v := data.(string)
			switch v {
			case "all":
				return 0, nil
			case "security":
				return 1, nil
			case "symbol":
				return 2, nil
			default:
				panic("invalid category: " + v)
			}
		}

		return data, nil
	}
}

func decodeHookWithTag(hook mapstructure.DecodeHookFunc, tagName string) viper.DecoderConfigOption {
	return func(c *mapstructure.DecoderConfig) {
		c.DecodeHook = hook
		c.TagName = tagName
	}
}

func parseSymbols() (symbols []Symbol, err error) {
	//iViper := config.GetConfigService("symbol")

	decodeHookFunc := symbolDecodeHook()
	decodeConfigTagName := iViper.GetString("decoder_config_tagname")
	decoderConfigOption := decodeHookWithTag(decodeHookFunc, decodeConfigTagName)

	err = iViper.UnmarshalKey("symbols", &symbols, decoderConfigOption)

	return
}

func parseSecurity() (securities []Security, err error) {
	err = iViper.UnmarshalKey("security", &securities, func(c *mapstructure.DecoderConfig) {
		c.TagName = decodeConfigTagName
	})
	return
}

func parseHoliday() (holidays []Holiday, err error) {
	err = iViper.UnmarshalKey("holidays", &holidays,
		func(c *mapstructure.DecoderConfig) {
			c.DecodeHook = holidayDecodeHook()
			c.TagName = decodeConfigTagName
		})

	return
}
