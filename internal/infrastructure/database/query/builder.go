package query

import (
	"BaseProjectGolang/internal/infrastructure/database"
	"BaseProjectGolang/internal/infrastructure/database/query/filter"
	"BaseProjectGolang/internal/infrastructure/database/query/order"
	"BaseProjectGolang/internal/infrastructure/database/query/pagination"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// Builder wraps a GORM DB instance with query building capabilities
type Builder struct {
	// Original is the clean, unmodified DB instance
	Original *gorm.DB

	// Current is the DB instance with applied filters, pagination, etc.
	Current *gorm.DB
}

// NewBuilder creates a new query builder
func NewBuilder(db *database.DataBase) *Builder {
	dbGorm := db.DatabaseDriver.MustGetGorm()

	return &Builder{
		Original: dbGorm,
		Current:  dbGorm,
	}
}

// WithSort применяет сортировку
func (qb *Builder) WithSort(ctx fiber.Ctx) *Builder {
	if orderConditions := order.ParseSortFromQuery(ctx); orderConditions != nil {
		qb.Current = order.ApplySort(qb.Current, orderConditions)
	}

	return qb
}

// WithPagination применяет пагинацию
func (qb *Builder) WithPagination(ctx fiber.Ctx) *Builder {
	if paginationConditions := pagination.ParsePaginationFromQuery(ctx); paginationConditions != nil {
		qb.Current = pagination.ApplyPagination(qb.Current, paginationConditions)
	}

	return qb
}

// WithFilter применяет фильтры
func (qb *Builder) WithFilter(ctx fiber.Ctx) *Builder {
	if filterConditions := filter.ParseFiltersFromQuery(ctx); filterConditions != nil {
		qb.Current = filter.ApplyFilters(qb.Current, filterConditions)
	}

	return qb
}

// Reset resets the Current DB instance to match the Original
func (qb *Builder) Reset() *Builder {
	qb.Current = qb.Original

	return qb
}
