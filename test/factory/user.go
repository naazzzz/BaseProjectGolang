package user

import (
	"time"

	userModel "BaseProjectGolang/internal/infrastructure/database/orm/model/user"
	"BaseProjectGolang/test"

	"github.com/Pallinder/go-randomdata"
	"github.com/bluele/factory-go/factory"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NewUserFactory() *factory.Factory {
	timeNow := time.Now()

	return factory.NewFactory(
		&userModel.User{
			CreatedAt: &timeNow,
			UpdatedAt: &timeNow,
		},
	).
		SeqInt("ID", func(n int) (interface{}, error) {
			if n < 0 {
				return uint(0), nil
			}

			return uint(n), nil
		}).
		Attr("Username", func(_ factory.Args) (interface{}, error) {
			return randomdata.SillyName(), nil
		}).
		Attr("Active", func(_ factory.Args) (interface{}, error) {
			return true, nil
		}).
		Attr("Password", func(_ factory.Args) (interface{}, error) {
			// secret
			return "$2a$04$QGRxTBJo9sIeD6QKPF8EHOZNgPZevWI0C6NakuFFsCq.92LhZTfhO", nil
		}).
		OnCreate(func(args factory.Args) error {
			db := args.Context().Value(test.DBConst).(*gorm.DB)

			modelUser := args.Instance().(*userModel.User)

			if err := db.Omit(clause.Associations).Create(modelUser).Error; err != nil {
				return err
			}

			return nil
		})
}
