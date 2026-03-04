package database

import (
	"BaseProjectGolang/internal/config"
	"BaseProjectGolang/internal/infrastructure/database/driver"
	logInternal "BaseProjectGolang/pkg/log"

	"gorm.io/gorm"
)

type (
	Driver interface {
		MustGetGorm() *gorm.DB
	}

	DataBase struct {
		DatabaseDriver Driver
	}
)

func NewDataBase(
	cfg *config.Config,
	loggerInternal *logInternal.Logger,
) (database *DataBase, err error) {
	database = &DataBase{}

	if cfg.AppWorkMode == config.TestEnv {
		// Используем SQLite для тестов
		database.DatabaseDriver, err = driver.NewSQLiteForTestEnv(loggerInternal, cfg)
		if err != nil {
			return
		}
	} else {
		// Основной драйвер
		database.DatabaseDriver, err = driver.NewPSQLAndLoadMigrations(cfg, loggerInternal)
		if err != nil {
			return
		}
	}

	return
}
