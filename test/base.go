package test

import (
	"BaseProjectGolang/internal/bootstrap"
	"BaseProjectGolang/internal/config"
	common "BaseProjectGolang/internal/constant"
	"BaseProjectGolang/internal/dependency/app"
	"BaseProjectGolang/internal/infrastructure/database"
	"BaseProjectGolang/pkg/crypto"
	logUtil "BaseProjectGolang/pkg/log"
	"context"
	"crypto/rand"
	"encoding/json"
	"log"
	"math/big"
	"net/http"
	"os"
	"path"
	"runtime"
	"testing"

	factoryLib "github.com/bluele/factory-go/factory"
	"gorm.io/gorm"
)

type typeForTest string

const (
	DBConst      typeForTest = "db"
	dbContextKey typeForTest = "db"
	letterRunes  typeForTest = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func InitCfg(t *testing.T, envConfig ...string) (cfg *config.Config, err error) {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")

	if err = os.Chdir(dir); err != nil {
		t.Error(err)
	}

	cfg, err = config.NewConfig(true, envConfig...)
	if err != nil {
		t.Error(err)
	}

	return
}

// GetDefaultAppTest Инициализация дефолтного инстанса аpp для использования в тестах
func GetDefaultAppTest(t *testing.T, configWebServiceData *string) (newApp *bootstrap.App, cfg *config.Config) {
	var (
		configPath string
	)

	cfg, _ = InitCfg(t, "test.env", configPath)
	db := InitializeAndCleanDatabaseAfterTest(t, cfg)

	container, _, err := app.InitializeContainer()
	if err != nil {
		log.Fatalf(err.Error())
	}

	container.Config = cfg
	container.DataBase = db
	newApp = container.App

	return
}

// CreateObjectInTestDatabaseFromFactory Инициализация объекта фабрики
func CreateObjectInTestDatabaseFromFactory(
	factory *factoryLib.Factory,
	db *database.DataBase,
	options map[string]interface{},
) interface{} {
	db.DatabaseDriver.MustGetGorm().Statement.Context = context.WithValue(context.Background(), common.FactoryKey, true)
	ctx := context.WithValue(context.Background(), dbContextKey, db.DatabaseDriver.MustGetGorm())

	defer func(db *gorm.DB) {
		db.Statement.Context = context.Background()
	}(db.DatabaseDriver.MustGetGorm())

	return factory.MustCreateWithContextAndOption(ctx, options)
}

// InitializeAndCleanDatabaseAfterTest Откатывает транзакцию и возвращает базу в изначальное состояние
func InitializeAndCleanDatabaseAfterTest(t *testing.T, cfg *config.Config) *database.DataBase {
	logger := logUtil.NewLogger(cfg.Logs)

	DB, err := database.NewDataBase(cfg, logger)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		// Закрыть соединение с MySQL
		psqlConn, err := DB.DatabaseDriver.MustGetGorm().DB()
		if err != nil {
			t.Error(err)
		}

		if err := psqlConn.Close(); err != nil {
			t.Error(err)
		}
	})

	return DB
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	letterRunesLength := big.NewInt(int64(len(letterRunes)))

	for i := range b {
		num, err := rand.Int(rand.Reader, letterRunesLength)
		if err != nil {
			panic(err) // Handle error appropriately in production code
		}

		b[i] = rune(letterRunes[num.Int64()])
	}

	return string(b)
}

func FormSignatureForAuth(cfg *config.Config, username, password string, req *http.Request) {
	mapData := make(map[string]string)
	mapData["username"] = username
	mapData["password"] = password

	jsonData, err := json.Marshal(mapData)
	if err != nil {
		panic(err)
	}

	encodedData, err := crypto.Aes256cbcEncode(string(jsonData), cfg.Secure.AuthPrivateKey, "")
	if err != nil {
		return
	}

	q := req.URL.Query()
	q.Set("signature", encodedData)
	req.URL.RawQuery = q.Encode()
}
