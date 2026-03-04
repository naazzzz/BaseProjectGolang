package main

import (
	"fmt"
	"io"
	"os"

	"BaseProjectGolang/internal/infrastructure/database/orm/model/user"

	"ariga.io/atlas-provider-gorm/gormschema"
)

func main() {
	// TODO что-то с этим сделать (подумать над вариантом через рефлексию цеплять структуры через файлы в папке models)
	models := []any{
		user.OAuthAccessToken{},
		user.User{},
	}

	stmts, err := gormschema.New("sqlite").Load(models...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}

	if _, err = io.WriteString(os.Stdout, stmts); err != nil {
		panic(err)
	}
}
