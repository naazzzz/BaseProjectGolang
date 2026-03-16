package examplerep

import (
	common "BaseProjectGolang/internal/constant"
	"BaseProjectGolang/internal/domain/exampledmn"
	"BaseProjectGolang/internal/infrastructure/database/orm/model/examplemdl"
	"BaseProjectGolang/pkg/querybuilder"
	"context"

	"github.com/jinzhu/copier"
	"github.com/soner3/flora"
)

type ExampleRepository struct {
	flora.Component
}

func NewExampleRepository() *ExampleRepository {
	return &ExampleRepository{}
}

func (repository *ExampleRepository) Example(ctx context.Context, domain *exampledmn.ExampleDomain) error {
	model := &examplemdl.ExampleModel{}

	if err := copier.Copy(&model, &domain); err != nil {
		return err
	}

	//  Or inject db in repository struct

	qb := ctx.Value(common.QbCtxKey).(*querybuilder.Builder)
	return qb.Current.
		Create(model).
		Error
}
