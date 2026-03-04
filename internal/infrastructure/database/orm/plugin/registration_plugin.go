package plugin

import (
	"BaseProjectGolang/internal/infrastructure/database"
)

type IRegistrationPlugin interface {
	RegisterPluginsInGorm() (err error)
}

type RegistrationPlugin struct {
	db *database.DataBase
	// add new plugin
}

func NewRegistrationPlugin(
	db *database.DataBase,
) *RegistrationPlugin {
	return &RegistrationPlugin{
		db: db,
	}
}

func (p *RegistrationPlugin) RegisterPluginsInGorm() (err error) {
	// use new plugin
	//if err = p.db.Pgsql.Gorm.Use(p.messagePlugin); err != nil {
	//	return
	//}

	return
}
