/*******************************************************************************
 * // Copyright AnchyTec Corp. All Rights Reserved.
 * // SPDX-License-Identifier: Apache-2.0
 * // Author: shaozhiming
 ******************************************************************************/

package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var (
	SymbolConfig  *viper.Viper
	AccountConfig *viper.Viper
	PriceConfig   *viper.Viper
	OrderConfig   *viper.Viper
	TradeConfig   *viper.Viper

	appMode       string
	appConfigPath string
)

func init() {
	// `APP_MODE` mapping config-path：`dev` -> `dev`, `prod` -> `prod`, `test` -> `test`
	appMode = os.Getenv("APP_MODE")
	if appMode == "" {
		appMode = "dev"
	}

	workPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	appConfigPath = filepath.Join(workPath, "config/resource")

	SymbolConfig = viper.New()
	initConfigFromFiles(SymbolConfig, "symbol")
}

func initConfigFromFiles(config *viper.Viper, fileName string) {
	config.SetConfigName(fileName)
	config.AddConfigPath(".")
	config.AddConfigPath(appConfigPath)
	// for modules test
	config.AddConfigPath("../config/resource")
	config.AddConfigPath("../../config/resource")
	config.AddConfigPath("../../../config/resource")
	config.AddConfigPath("../../../../config/resource")

	err := config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file [%s]: %s \n", fileName, err))
	}
}

func GetConfigService(moduleName string) *viper.Viper {
	switch moduleName {
	case "symbol":
		return SymbolConfig.Sub(appMode)
	case "price":
		return PriceConfig.Sub(appMode)
	case "order":
		return OrderConfig.Sub(appMode)
	case "account":
		return AccountConfig.Sub(appMode)
	case "trade":
		return TradeConfig.Sub(appMode)
	default:
		panic(fmt.Errorf("Unsupported configuration module : %s\n", moduleName))
	}
}
