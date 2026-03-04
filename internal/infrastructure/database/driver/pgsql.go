package driver

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"BaseProjectGolang/internal/config"
	logInternal "BaseProjectGolang/pkg/log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file" //nolint:revive
	"github.com/rotisserie/eris"
	gormPgsql "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	ErrDBMysqlNil           = errors.New("dbConn is nil")
	ErrMigrationDirNotExist = errors.New("migration directory does not exist")
)

const SlowThreshold = 300 * time.Millisecond

type Pgsql struct {
	Gorm *gorm.DB
}

func NewPSQLAndLoadMigrations(
	cfg *config.Config,
	loggerInternal *logInternal.Logger,
) (pgsql *Pgsql, err error) {
	pgsql = &Pgsql{}

	conTimeEnv, maxConEnv, maxIdleConEnv := pgsql.сonvertEnvTypeToNeeded(cfg)

	var pgsqlConn *sql.DB

	if pgsqlConn, err = pgsql.CreateConnectionPgsql(conTimeEnv, maxConEnv, maxIdleConEnv, cfg); err != nil {
		return
	}

	if err = pgsql.CreatePgsqlSessionForGormUsage(loggerInternal, pgsqlConn); err != nil {
		return
	}

	if cfg.AppWorkMode == config.TestEnv || cfg.AppWorkMode == config.DebugEnv {
		if err = pgsql.LoadMigrations(cfg); err != nil {
			return
		}

		pgsql.Gorm = pgsql.Gorm.Debug()
	}

	return
}

func (pgsql *Pgsql) CreatePgsqlSessionForGormUsage(loggerInternal *logInternal.Logger, pgsqlConn *sql.DB) (err error) {
	pgsql.Gorm, err = gorm.Open(
		gormPgsql.New(gormPgsql.Config{
			Conn: pgsqlConn,
		}),
		&gorm.Config{
			Logger: logger.New(loggerInternal.Logger, logger.Config{
				SlowThreshold:             SlowThreshold,
				LogLevel:                  logger.Warn,
				IgnoreRecordNotFoundError: false,
				Colorful:                  false,
			}),
		})

	return
}

func (pgsql *Pgsql) сonvertEnvTypeToNeeded(
	cfg *config.Config,
) (conTimeEnv time.Duration, maxConEnv int, maxIdleConEnv int) {
	conTime, err := strconv.Atoi(cfg.Databases.ConnectionMaxLifeTime)
	if err != nil {
		log.Println(err)
	}

	conTimeEnv = time.Duration(conTime) * time.Second

	maxConEnv, err = strconv.Atoi(cfg.Databases.MaxOpenConnections)
	if err != nil {
		log.Println(err)
	}

	maxIdleConEnv, err = strconv.Atoi(cfg.Databases.MaxIdleConnections)
	if err != nil {
		log.Println(err)
	}

	return
}

func (pgsql *Pgsql) LoadMigrations(
	cfg *config.Config,
) (err error) {
	conn, err := pgsql.Gorm.DB()
	if err != nil {
		return err
	}

	if conn == nil {
		return eris.Wrap(ErrDBMysqlNil, "dbConn is nil")
	}

	var (
		driver    database.Driver
		sourceURL string
	)

	// Выбираем драйвер
	switch cfg.AppWorkMode {
	case config.TestEnv:
		driver, err = sqlite3.WithInstance(conn, &sqlite3.Config{})
		sourceURL = "file://internal/infrastructure/database/migrations/sqlite"
	default:
		driver, err = postgres.WithInstance(conn, &postgres.Config{})
		sourceURL = "file://internal/infrastructure/database/migrations/postgres"
	}

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
		log.Println("MySQL Migrations: " + err.Error())
	}

	return nil
}

func (pgsql *Pgsql) CreateConnectionPgsql(
	conTimeEnv time.Duration,
	maxConEnv int,
	maxIdleConEnv int,
	cfg *config.Config,
) (pgsqlConn *sql.DB, err error) {
	connStr := "user=%s password=%s host=%s port=%s dbname=%s sslmode=disable"

	path := fmt.Sprintf(connStr, cfg.Databases.Pgsql.User, cfg.Databases.Pgsql.Password, cfg.Databases.Pgsql.Host, cfg.Databases.Pgsql.Port, cfg.Databases.Pgsql.Database)

	if pgsqlConn, err = sql.Open("postgres", path); err != nil {
		return
	}

	pgsqlConn.SetConnMaxLifetime(conTimeEnv)
	pgsqlConn.SetMaxOpenConns(maxConEnv)
	pgsqlConn.SetMaxIdleConns(maxIdleConEnv)

	return
}

func (pgsql *Pgsql) MustGetGorm() *gorm.DB {
	return pgsql.Gorm
}
