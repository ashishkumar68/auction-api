package config

import (
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

var DBConfig *gorm.Config

func LoadDBConfig() {
	switch os.Getenv("APP_ENV") {
	case AppEnvProd:
		DBConfig = &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		}
	case AppEnvStaging:
		DBConfig = &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		}
	case AppEnvTest:
		DBConfig = &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
			PrepareStmt: false,
		}
	default:
		DBConfig = &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
			PrepareStmt: false,
		}
	}
}