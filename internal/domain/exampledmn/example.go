package exampledmn

import (
	"time"
)

type ExampleDomain struct {
	ID        string
	Data      string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
