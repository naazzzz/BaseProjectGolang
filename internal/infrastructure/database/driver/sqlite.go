package driver

import (
	"fmt"
	"log"
	"os"
	"time"

	"BaseProjectGolang/internal/config"
	logInternal "BaseProjectGolang/pkg/log"

	"github.com/glebarez/sqlite"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file" //nolint:revive
	"github.com/google/uuid"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Sqlite struct {
	Gorm *gorm.DB
}

func NewSQLiteForTestEnv(loggerInternal *logInternal.Logger, cfg *config.Config) (sqliteObj *Sqlite, err error) {
	sqliteObj = &Sqlite{}

	dsn := fmt.Sprintf("file:testdb_%d_%s?mode=memory&cache=shared", time.Now().Unix(), uuid.New().String())

	sqliteObj.Gorm, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.New(loggerInternal.Logger, logger.Config{
			SlowThreshold:             SlowThreshold,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: false,
			Colorful:                  false,
		}),
	})
	if err != nil {
		return
	}

	if err = sqliteObj.LoadMigrations(cfg); err != nil {
		return
	}

	return
}

func (sqlite *Sqlite) LoadMigrations(
	cfg *config.Config,
) (err error) {
	conn, err := sqlite.Gorm.DB()
	if err != nil {
		return err
	}

	if conn == nil {
		return eris.Wrap(ErrDBMysqlNil, "dbConn is nil")
	}

	driver, err := sqlite3.WithInstance(conn, &sqlite3.Config{})
	sourceURL := "file://internal/infrastructure/database/migrations/sqlite"

	if err != nil {
		return err
	}

	// Проверяем существование директории с миграциями
	if _, err = os.Stat(sourceURL[7:]); os.IsNotExist(err) {
		return eris.Wrapf(ErrMigrationDirNotExist, "path: %s", sourceURL[7:])
	}

	// Создаем экземпляр миграции
	instance, err := migrate.NewWithDatabaseInstance(
		sourceURL,
		cfg.Databases.Pgsql.Database,
		driver,
	)
	if err != nil {
		return err
	}

	if cfg.AppWorkMode != config.TestEnv {
		log.Println("Running migrations in debug mode...")
	}

	// Применяем миграции
	if err = instance.Up(); err != nil {
		log.Println("SQLite Migrations: " + err.Error())
	}

	return nil
}

func (sqlite *Sqlite) MustGetGorm() *gorm.DB {
	return sqlite.Gorm
}
