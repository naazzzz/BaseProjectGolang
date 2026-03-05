package trait

import (
	"BaseProjectGolang/internal/config"
	"BaseProjectGolang/internal/infrastructure/database"
	"BaseProjectGolang/internal/infrastructure/database/orm/model/userModel"
	"BaseProjectGolang/test"
	user2 "BaseProjectGolang/test/factory"

	factoryLib "github.com/bluele/factory-go/factory"
)

func CreateUserWithServiceInfo(db *database.DataBase, _ *config.Config, userAttributes map[string]interface{}) *userModel.User {
	userFactory := user2.NewUserFactory()

	for key, value := range userAttributes {
		userFactory = userFactory.Attr(key, func(_ factoryLib.Args) (interface{}, error) {
			return value, nil
		})
	}

	return test.CreateObjectInTestDatabaseFromFactory(userFactory, db, nil).(*userModel.User)
}
