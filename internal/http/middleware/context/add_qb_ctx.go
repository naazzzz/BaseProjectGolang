package context

import (
	common "BaseProjectGolang/internal/constant"
	"BaseProjectGolang/internal/infrastructure/database"
	"BaseProjectGolang/internal/infrastructure/database/query"

	"github.com/gofiber/fiber/v3"
)

type SetupCtxQB struct {
	db *database.DataBase
}

// MiddlewareOptions определяет, какие модификаторы применять
type MiddlewareOptions struct {
	UseSort       bool
	UsePagination bool
	UseFilter     bool
}

func NewSetupCtxQB(db *database.DataBase) *SetupCtxQB {
	return &SetupCtxQB{
		db: db,
	}
}

// DefaultQueryBuilderMiddleware - middleware с использованием всех модификаторов
// Примеры инициализации в роутинге
// // Использование всех модификаторов
//
//	app.Get("/items/full", func(c fiber.Ctx) error {
//	    builder := ctx.Locals(common.QbCtxKey).(*query.Builder)
//	    return builder.DefaultQueryBuilderMiddleware(c)
//	}, handler)
//
//	// Использование только сортировки
//	app.Get("/items/sort", func(c fiber.Ctx) error {
//	    builder := ctx.Locals(common.QbCtxKey).(*query.Builder)
//	    return builder.CustomQueryBuilderMiddleware(c, MiddlewareOptions{
//	        UseSort: true,
//	    })
//	}, handler)
func (m *SetupCtxQB) DefaultQueryBuilderMiddleware(ctx fiber.Ctx) error {
	return m.CustomQueryBuilderMiddleware(ctx, MiddlewareOptions{
		UseSort:       true,
		UsePagination: true,
		UseFilter:     true,
	})
}

// CustomQueryBuilderMiddleware - настраиваемый middleware
// Примеры инициализации в роутинге
// // Использование всех модификаторов
//
//	app.Get("/items/full", func(c fiber.Ctx) error {
//	    builder := NewBuilder(db)
//	    return builder.DefaultQueryBuilderMiddleware(c)
//	}, handler)
//
//	// Использование только сортировки
//	app.Get("/items/sort", func(c fiber.Ctx) error {
//	    builder := NewBuilder(db)
//	    return builder.CustomQueryBuilderMiddleware(c, MiddlewareOptions{
//	        UseSort: true,
//	    })
//	}, handler)
func (m *SetupCtxQB) CustomQueryBuilderMiddleware(ctx fiber.Ctx, opts MiddlewareOptions) error {
	qb := query.NewBuilder(m.db)

	if ctx.Method() != fiber.MethodGet {
		ctx.Locals(common.QbCtxKey, qb)

		return ctx.Next()
	}

	// Применяем модификаторы в зависимости от опций
	if opts.UseSort {
		qb.WithSort(ctx)
	}

	if opts.UsePagination {
		qb.WithPagination(ctx)
	}

	if opts.UseFilter {
		qb.WithFilter(ctx)
	}

	// Сохраняем qb в контексте
	ctx.Locals(common.QbCtxKey, qb)

	return ctx.Next()
}
