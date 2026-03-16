package exampledmn

import (
	"context"
)

type IExampleRepository interface {
	Example(ctx context.Context, domain *ExampleDomain) error
}
