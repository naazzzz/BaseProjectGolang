package exampleusc

import (
	"BaseProjectGolang/internal/config"
	"BaseProjectGolang/internal/domain/exampledmn"
	"context"

	"github.com/jinzhu/copier"
	"github.com/soner3/flora"
)

type ExampleHandler struct {
	flora.Component
	exampleRepository exampledmn.IExampleRepository
	cfg               *config.Config
}

func NewExampleHandler(
	cfg *config.Config,
	exampleRepository exampledmn.IExampleRepository,

) *ExampleHandler {
	tokenService := &ExampleHandler{
		cfg:               cfg,
		exampleRepository: exampleRepository,
	}

	return tokenService
}

func (handler *ExampleHandler) Execute(
	ctx context.Context,
	exampleCmd *ExampleCommand,
) (*ExampleResult, error) {
	exampleDomain := &exampledmn.ExampleDomain{}

	if err := copier.Copy(&exampleDomain, &exampleCmd); err != nil {
		return nil, err
	}

	if err := handler.exampleRepository.Example(ctx, exampleDomain); err != nil {
		return nil, err
	}

	return &ExampleResult{
		Data: "example result",
	}, nil
}
