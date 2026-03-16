package filter

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

// Operator FilterOperator представляет оператор для фильтрации
type Operator string

// todo будущем вынести в конфигурацию
var allowedOperations = map[Operator]bool{ //nolint:gochecknoglobals
	Eq:         true,
	Neq:        true,
	Gt:         true,
	Gte:        true,
	Lt:         true,
	Lte:        true,
	Like:       true,
	NotLike:    true,
	In:         true,
	NotIn:      true,
	IsNull:     true,
	IsNotNull:  true,
	DateBefore: true,
	DateAfter:  true,
	Search:     true,
}

const (
	Eq         Operator = "eq"     // равно
	Neq        Operator = "neq"    // не равно
	Gt         Operator = "gt"     // больше
	Gte        Operator = "gte"    // больше или равно
	Lt         Operator = "lt"     // меньше
	Lte        Operator = "lte"    // меньше или равно
	Like       Operator = "like"   // LIKE %value%
	NotLike    Operator = "nlike"  // NOT LIKE %value%
	In         Operator = "in"     // IN (values)
	NotIn      Operator = "nin"    // NOT IN (values)
	IsNull     Operator = "null"   // IS NULL
	IsNotNull  Operator = "nnull"  // IS NOT NULL
	DateBefore Operator = "before" // Дата до
	DateAfter  Operator = "after"  // Дата после
	Search     Operator = "search" // Дата после
)

// Condition FilterCondition представляет условие фильтрации
type Condition struct {
	Field    string      `json:"field"`
	Operator Operator    `json:"operator"`
	Value    interface{} `json:"value"`
}

// Request FilterRequest представляет запрос на фильтрацию
type Request struct {
	Filters []Condition `json:"filters"`
}

// ApplyFilters применяет фильтры к запросу GORM

func ApplyFilters(query *gorm.DB, filters []Condition) *gorm.DB {
	for _, filter := range filters {
		query = applyFilter(query, filter)
	}

	return query
}

// applyFilter применяет одно условие фильтрации
func applyFilter(query *gorm.DB, filter Condition) *gorm.DB {
	switch filter.Operator {
	case Eq:
		return query.Where(fmt.Sprintf("%s = ?", filter.Field), filter.Value)
	case Neq:
		return query.Where(fmt.Sprintf("%s <> ?", filter.Field), filter.Value)
	case Gt:
		return query.Where(fmt.Sprintf("%s > ?", filter.Field), filter.Value)
	case Gte:
		return query.Where(fmt.Sprintf("%s >= ?", filter.Field), filter.Value)
	case Lt:
		return query.Where(fmt.Sprintf("%s < ?", filter.Field), filter.Value)
	case Lte:
		return query.Where(fmt.Sprintf("%s <= ?", filter.Field), filter.Value)
	case Like:
		return query.Where(fmt.Sprintf("%s LIKE ?", filter.Field), fmt.Sprintf("%%%v%%", filter.Value))
	case NotLike:
		return query.Where(fmt.Sprintf("%s NOT LIKE ?", filter.Field), fmt.Sprintf("%%%v%%", filter.Value))
	case In:
		return query.Where(fmt.Sprintf("%s IN (?)", filter.Field), filter.Value)
	case NotIn:
		return query.Where(fmt.Sprintf("%s NOT IN ?", filter.Field), filter.Value)
	case IsNull:
		return query.Where(fmt.Sprintf("%s IS NULL", filter.Field))
	case IsNotNull:
		return query.Where(fmt.Sprintf("%s IS NOT NULL", filter.Field))
	case DateBefore:
		if dateStr, ok := filter.Value.(string); ok {
			date, err := time.Parse(time.RFC3339, dateStr)
			if err == nil {
				return query.Where(fmt.Sprintf("%s < ?", filter.Field), date)
			}
		}

		return query
	case DateAfter:
		if dateStr, ok := filter.Value.(string); ok {
			date, err := time.Parse(time.RFC3339, dateStr)
			if err == nil {
				return query.Where(fmt.Sprintf("%s > ?", filter.Field), date)
			}
		}

		return query

	case Search:
		return query.Where(
			fmt.Sprintf("to_tsvector('simple', %s) @@ websearch_to_tsquery('simple', ?)", filter.Field),
			filter.Value,
		)
	default:
		return query
	}
}

// ParseFiltersFromQuery извлекает фильтры из параметров запроса
func ParseFiltersFromQuery(ctx fiber.Ctx) []Condition {
	var filters []Condition

	// Получаем все параметры запроса
	params := ctx.Queries()

	for key, value := range params {
		if strings.HasPrefix(key, "filter[") {
			parts := strings.Split(key, "[")
			if len(parts) != 3 { //nolint:mnd
				continue
			}

			field := strings.TrimSuffix(parts[1], "]")
			operatorStr := strings.TrimSuffix(parts[2], "]")
			operator := Operator(operatorStr)

			// Проверяем, разрешён ли оператор
			if !allowedOperations[operator] {
				continue
			}

			parsedValue := ParseArrayValue(value)
			if !reflect.ValueOf(parsedValue).IsValid() {
				continue
			}

			// Создаем условие фильтрации
			condition := Condition{
				Field:    field,
				Operator: operator,
				Value:    parsedValue,
			}

			filters = append(filters, condition)
		}
	}

	return filters
}

func ParseArrayValue(query string) interface{} {
	if query == "" {
		return nil
	}

	query = strings.TrimSpace(query)

	// Split the querybuilder to get the first element for type inference
	valuesStr := strings.Split(query, ",")
	if len(valuesStr) == 0 {
		return nil
	}

	// Infer the type from the first element
	firstValue := valuesStr[0]

	var inferredType reflect.Type

	if _, err := strconv.ParseUint(firstValue, 10, 64); err == nil {
		inferredType = reflect.TypeOf(uint(0))
	} else if _, err := strconv.Atoi(firstValue); err == nil {
		inferredType = reflect.TypeOf(0)
	} else if _, err := strconv.ParseFloat(firstValue, 64); err == nil {
		inferredType = reflect.TypeOf(0.0)
	} else if _, err := strconv.ParseBool(firstValue); err == nil {
		inferredType = reflect.TypeOf(true)
	} else {
		inferredType = reflect.TypeOf("")
	}

	// Example a slice of the inferred type
	values := reflect.MakeSlice(reflect.SliceOf(inferredType), 0, 0)

	for _, value := range valuesStr {
		parsedValue, err := ParseValue(value, inferredType)
		if err != nil {
			continue
		}

		values = reflect.Append(values, reflect.ValueOf(parsedValue))
	}

	return values.Interface()
}

func ParseValue(value string, inferredType reflect.Type) (interface{}, error) {
	switch inferredType.Kind() {
	case reflect.Uint:
		u, err := strconv.ParseUint(value, 10, 64)
		return uint(u), err
	case reflect.Int:
		return strconv.Atoi(value)
	case reflect.Float64:
		return strconv.ParseFloat(value, 64)
	case reflect.Bool:
		return strconv.ParseBool(value)
	case reflect.String:
		return value, nil
	default:
		return nil, eris.New("unsupported type conversion")
	}
}
