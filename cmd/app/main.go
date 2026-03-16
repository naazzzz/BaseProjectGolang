package main

import (
	"BaseProjectGolang/internal/dependency/app"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

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
	container, _, err := app.InitializeContainer()
	if err != nil {
		log.Fatalf(err.Error())
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err = container.App.Run(ctx); err != nil {
		log.Fatalf(err.Error())
		return
	}
}
