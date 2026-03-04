package main

import (
	"BaseProjectGolang/internal/dependency"
	appDependency "BaseProjectGolang/internal/dependency/app"

	_ "github.com/lib/pq"
)

// @title						Monitoring service
// @version					1.0
// @description				Monitoring API
// @basePath					/api
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
func main() {
	db, cfg := dependency.InitServicesBeforeDi()

	fiberInstance, err := appDependency.InitializeApp(cfg, db)
	if err != nil {
		panic(err)
	}

	if err = fiberInstance.Run(); err != nil {
		panic(err)
	}
}
