package test

import (
	common "BaseProjectGolang/internal/constant"
	"context"
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"path"
	"runtime"
	"strconv"
	"testing"

	"BaseProjectGolang/internal/bootstrap"
	"BaseProjectGolang/internal/config"
	appDependency "BaseProjectGolang/internal/dependency/app"
	"BaseProjectGolang/internal/http/middleware/auth"
	"BaseProjectGolang/internal/infrastructure/database"
	"BaseProjectGolang/pkg/crypto"
	logUtil "BaseProjectGolang/pkg/log"

	factoryLib "github.com/bluele/factory-go/factory"
	"github.com/dromara/carbon/v2"
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

	cfg, err = config.LoadConfig(true, envConfig...)
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

	newApp, err := appDependency.InitializeApp(cfg, db)
	if err != nil {
		t.Error(err)
	}

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
	logger := logUtil.InitLogger(cfg.Logs)

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

func SetAuthForRequest(
	t *testing.T,
	req *http.Request,
	queryValues url.Values,
	body *[]byte,
	time *carbon.Carbon,
	cfg *config.Config,
) {
	query, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		t.Fatal(err)
	}

	if time == nil {
		timeNow := carbon.Now()
		time = timeNow
	}

	signatureTimestamp := time.Timestamp()

	if body == nil {
		newBody := []byte("")
		body = &newBody
	}

	publicKey, err := auth.FormPublicKey(query, *body, int(signatureTimestamp))
	if err != nil {
		t.Fatal(err)
	}

	publicKeyMd5 := crypto.GetMD5Hash(publicKey)

	signature := crypto.EncodeHmacSha256(cfg.Secure.AuthPrivateKey, publicKeyMd5)

	queryValues.Set("signature", signature)
	queryValues.Set("signature_timestamp", strconv.FormatInt(signatureTimestamp, 10))

	req.URL.RawQuery = queryValues.Encode()
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
