package order

import (
	"fmt"
	"strings"

	"BaseProjectGolang/pkg/condition"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// SortDirection представляет направление сортировки
type SortDirection string

const (
	Asc  SortDirection = "asc"  // по возрастанию
	Desc SortDirection = "desc" // по убыванию
)

// SortOption представляет опцию сортировки
type SortOption struct {
	Field     string        `json:"field"`
	Direction SortDirection `json:"direction"`
}

// SortRequest представляет запрос на сортировку
type SortRequest struct {
	Sort []SortOption `json:"sort"`
}

// ApplySort применяет сортировку к запросу GORM
func ApplySort(query *gorm.DB, sortOptions []SortOption) *gorm.DB {
	for _, sort := range sortOptions {
		direction := "ASC"
		if sort.Direction == Desc {
			direction = "DESC"
		}

		query = query.Order(fmt.Sprintf("%s %s", sort.Field, direction))
	}

	return query
}

// ParseSortFromQuery извлекает параметры сортировки из запроса
func ParseSortFromQuery(ctx fiber.Ctx) []SortOption {
	var sortOptions []SortOption

	// Получаем все параметры запроса
	params := ctx.Queries()

	for key, value := range params {
		// Проверяем, является ли параметр фильтром
		if strings.HasPrefix(key, "order[") {
			parts := strings.Split(key, "[")
			if len(parts) != 2 { //nolint:mnd
				continue
			}

			field := strings.TrimSuffix(parts[1], "]")

			sortOptions = append(sortOptions, SortOption{
				Field:     field,
				Direction: getSortDirection(value),
			})
		}
	}

	return sortOptions
}

func getSortDirection(directionStr string) SortDirection {
	directionStr = strings.ToLower(directionStr)
	return condition.If[SortDirection](directionStr == "desc", Desc, Asc)
}
