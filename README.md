# Base Golang Project
Базовый проект для развертки golang-сервисов, ориентированный на DDD.

## Перед началом работы

1. Скопировать .env.example в .env
2. Скопировать services.example.yml в services.yml

(Для локальной разработки)
3. Установить переменные окружения в текущую сборку Golang
      (самый быстрый и простой способ - установить Plugin для Goland - EnvFile, либо через environments в сборке передавать по одной)
      После чего перейти в Run/Debug Configuration -> Вкладка EnvFile -> + -> добавляем .env из корня проекта

## Автогенерация сваггера на основе fiber-swaggo
Ссылка для ознакомления https://github.com/swaggo/gin-swagger

~~~sh
swag init -g ./cmd/app/main.go -o ./api
~~~

_____
# Миграции

Файл atlas.hcl служит конфигурационным файлом для миграций,
в теории предполагается дополнительная настройка для mysql и
clickhouse миграций.

## Генерация миграций при помощи atlas и gorm

https://atlasgo.io/guides/orms./gorm/getting-started

~~~sh
atlas migrate diff migration --env {gorm_mysql или gorm_clickhouse}
~~~

P.S. Не пропускать name - по дефолту пусть будет migration, иначе при запуске контейнера они не пройдут

Для удобства в локальной среде можно использовать

Установка:
~~~sh
curl -sSf https://atlasgo.sh | sh
~~~

Запуск:
~~~sh 
sh local-migration-generate.sh
~~~

Пересборка hash
~~~sh
atlas migrate hash migration --env gorm_mysql
~~~
~~~sh
atlas migrate hash migration --env gorm_sqlite
~~~

____
Так как они полностью совместимы с golang-migrate можно использовать следующее:

> Установите пакет migrate https://github.com/golang-migrate/migrate/tree/v4.18.3/cmd/migrate

Создание файлов миграции(создает пустые файлы up и down миграции)
 ```sh
  migrate create -ext sql -dir ./internal/database/migrations/postgres migration_name 
  ```

> Применение всех миграций
```sh
 migrate -path ./internal/database/migrations/postgres -database "postgresql://postgres:12345@localhost:5434/listing_db?sslmode=disable" up 
 ```

> Применение последней миграции
```sh
 migrate -path ./internal/database/migrations/postgres -database "postgresql://postgres:12345@localhost:5434/listing_db?sslmode=disable" up 1 
```

> Откат последней миграции
```sh
 migrate -path ./internal/database/migrations/postgres -database "postgresql://postgres:12345@localhost:5434/listing_db?sslmode=disable" down 1 
 ```

> Если при применении миграции произошла ошибка, то перед тем как накатить миграцию заново нужно выполнить откат к предыдущей миграции
```sh
 migrate -path ./internal/database/migrations/postgres -database "postgresql://postgres:12345@localhost:5434/listing_db?sslmode=disable" force migration_number(20250801135242) 
 ```
_____
## DI

Install Wire by running: (для генерации файла зависимостей)
```sh
go install github.com/google/wire/cmd/wire@latest
```
and ensuring that `$GOPATH/bin` is added to your `$PATH`.

Для реализации DI и спользуется Wire - https://github.com/google/wire/tree/main.
В целом для корректной реализации требуется:

- Создть так называемые Provide-методы предоставляющий объект структуры (по сути конструктор структуры)
- Выстроить иерархию структур для необходимой логики с реализованными конструкторами
- В иерархическом порядке добавить их в файл wire.go
- Выполнить команду:

  (или же go run -mod=mod github.com/google/wire/cmd/wire)

~~~sh 
cd ./internal/dependency/cmd/app && wire
~~~
_____
## Стэк технологий

- Fiber - http-framework
- swaggo - в интеграции с fiber для удобной и быстрой генерации swagger
- gorm - ORM
- golang-migrate - для миграций (по-хорошему от atlas надо отходить)
- wire - dependency injection
- enums - go-enum (auto-generation)
- eris - трассировка ошибок для http_errors
- envConfig, viper, dotenv - для работы с переменными окружения и выстраивания конфигурации проекта
- logs - lumberjack для записи в файлы логов ("gopkg.in/natefinch/lumberjack.v2")

Тесты:
- моки http - в pkg/fast_http_mock
- Для реализации фабрик используются:
    - "github.com/Pallinder/go-randomdata"
    - "github.com/bluele/factory-go/factory"

____
### Структура проекта Golang

https://habr.com/ru/companies/inDrive/articles/690088/
____

# Тесты

1. Запустить контейнеры из docker-compose 

3. Запустить команду(либо то же самое кликом через IDE)

~~~sh
go test ./test/... 
 ~~~

### Запуск тестов для gitlab ci/cd

~~~sh
make quick-start-test
~~~

### Локальное отображение зон покрытия тестами в браузере

```shell
go test -v -covermode=atomic -coverprofile=/startuptv-api/coverage/coverage.out -coverpkg=./internal/... ./test/...
```

coverprofile можно игнорировать файл служит для отображения зон в ci/cd.

## Обработка ошибок
Для корректности отображения возвращаемых ошибок, а также trace по ошибке, все, что обрабатывается необходимо возвращаться через обертку Eris на самом низком уровне

~~~golang
	if err := c.ParamsParser(&se); err != nil {
		return eris.New(err.Error())
	}
~~~

# Работа с линтерами

Файл конфигурации
GolangCI-Lint ищет файлы конфигурации по следующим путям из текущего рабочего каталога:

```yaml
.golangci.yml
.golangci.yaml
.golangci.toml
.golangci.json
```

## Установка
проверить версию  - https://golangci-lint.run/welcome/install/
```sh
# binary will be $(go env GOPATH)/bin/golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.1.6

golangci-lint --version
```

Для локального запуска линтера

```sh
golangci-lint run ./config/... ./cmd/... ./internal/... ./pkg/... ./test/...
```

Запуск с автофиксом (обратить внимание, что автоматически фиксится не все, потому что не все линтеры поддерживают данную
функцию)

```sh
golangci-lint run ./config/... ./cmd/... ./internal/... ./pkg/... ./test/... --fix
```

## Сборка пайплайн с тестом и линтом для golang-сервисов
В настоящий момент, существуют следующие golang сервисы, в которые необходимо добавить стадии lint и test:
listing, auction, notification, services, background_deal_processor, background_chat_processor, background_auction_processor

Разберем на примере listing-сервиса:

1. Для stage lint и test существуют соответствующие Dockerfile(DockerfileForLint, DockerfileForTest), которые должны срабатывать на все мр-ры, идущие в дев и прод.
2. Для удобства запуска подготовлен соответсвующий Makefile с командами для запуска и поддержки
3. Команды из данного Makefile используются в .gitlab-ci.yml

Данные stage должны выполняться для каждого из вышеперечисленных сервисов.

Принцип работы lint:
Запускается автоанализатор кода на предмет синтаксиса, отсутпов и многого другого, в случае, если разработчик забыл(или не учел при разработке выполнить линтер) DockerfileForLint завершит свою работу с ошибкой отличной от 0.

Принцип работы test:
Запускает все тесты из директории test соответствующего golang-микросервиса, формирует артифакты junit-report.xml и coverage.out, которые использует gitlab для формирования удобного отображения пройденных и провалившихся тестов, а также процента покрытия кода тестами

Соответствующие ссылки:
- https://docs.gitlab.com/ci/testing/code_coverage/ - формирование отчета о покрытии кода тестами для gitlab
- https://docs.gitlab.com/ci/testing/unit_test_reports/ - формирование отчета о тестах (упавших и пройденных) для gitlab