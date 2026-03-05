package pagination

import (
	"BaseProjectGolang/internal/constant"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// todo в будущем вынести в конфигурацию такие вещи
const (
	Limit              = "limit"
	Offset             = "offset"
	DefaultOffsetValue = 0
	DefaultLimitValue  = 10
	MaxPageSize        = 100
)

// Request PaginationRequest представляет запрос на пагинацию
type Request struct {
	Limit  int `json:"limit" query:"limit"`
	Offset int `json:"offset" query:"offset"`
	// Page     int `json:"page" query:"page"`
	// PageSize int `json:"page_size" query:"page_size"`
}

// Response Не используется, при необходимости добавить
// PaginationResponse представляет ответ с пагинацией
type Response struct {
	Data       interface{} `json:"data"`
	Pagination struct {
		CurrentPage int   `json:"current_page"`
		PageSize    int   `json:"page_size"`
		TotalItems  int64 `json:"total_items"`
		TotalPages  int   `json:"total_pages"`
		HasPrevPage bool  `json:"has_prev_page"`
		HasNextPage bool  `json:"has_next_page"`
	} `json:"pagination"`
}

// ParsePaginationFromQuery извлекает параметры пагинации из запроса
func ParsePaginationFromQuery(ctx fiber.Ctx) *Request {
	pagination := &Request{
		Limit:  DefaultLimitValue,
		Offset: DefaultOffsetValue,
		//Page:     1,
		//PageSize: DefaultPageSize,
	}

	fmt.Println(ctx.Query(Offset), ctx.Query(Limit))
	// Получаем номер страницы из запроса
	if offset := ctx.Query(Offset); offset != "" {
		if offsetNumber, err := strconv.Atoi(offset); err == nil && offsetNumber > 0 {
			pagination.Offset = offsetNumber
		}
	}

	// Получаем размер страницы из запроса
	if limit := ctx.Query(Limit); limit != "" {
		if limitValue, err := strconv.Atoi(limit); err == nil && limitValue > 0 {
			// Ограничиваем максимальный размер страницы
			if limitValue > MaxPageSize {
				limitValue = MaxPageSize
			}

			pagination.Limit = limitValue
		}
	}

	ctx.Locals(contextkeys.PaginationCtxKey, pagination)

	return pagination
}

// ApplyPagination применяет пагинацию к запросу GORM

func ApplyPagination(query *gorm.DB, pagination *Request) *gorm.DB {
	return query.Offset(pagination.Offset).Limit(pagination.Limit)
}

// Неплохой вариант реализации пагинации, если фронту понадобится информация о работе со страницами предпочтительнее реализовывать так
// CreatePaginationResponse создает ответ с пагинацией
// func CreatePaginationResponse(data abstraction{}, pagination Request, totalItems int64) Response {
//	response := Response{
//		Data: data,
//	}
//
//	// Заполняем информацию о пагинации
//	response.Pagination.CurrentPage = pagination.Page
//	response.Pagination.PageSize = pagination.PageSize
//	response.Pagination.TotalItems = totalItems
//
//	// Вычисляем общее количество страниц
//	totalPages := int(math.Ceil(float64(totalItems) / float64(pagination.PageSize)))
//	response.Pagination.TotalPages = totalPages
//
//	// Определяем, есть ли предыдущая и следующая страницы
//	response.Pagination.HasPrevPage = pagination.Page > 1
//	response.Pagination.HasNextPage = pagination.Page < totalPages
//
//	return response
//}
