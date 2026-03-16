package user

import (
	userModel "BaseProjectGolang/internal/infrastructure/database/orm/model/examplemdl"
	"BaseProjectGolang/test"

	"github.com/Pallinder/go-randomdata"
	"github.com/bluele/factory-go/factory"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NewExampleFactory() *factory.Factory {
	return factory.NewFactory(
		&userModel.ExampleModel{},
	).
		SeqInt("ID", func(n int) (interface{}, error) {
			if n < 0 {
				return uint(0), nil
			}

			return uint(n), nil
		}).
		Attr("Data", func(_ factory.Args) (interface{}, error) {
			return randomdata.Alphanumeric(50), nil
		}).
		OnCreate(func(args factory.Args) error {
			db := args.Context().Value(test.DBConst).(*gorm.DB)

			modelExample := args.Instance().(*userModel.ExampleModel)

			if err := db.Omit(clause.Associations).Create(modelExample).Error; err != nil {
				return err
			}

			return nil
		})
}
