package controller

import (
	"BaseProjectGolang/internal/config"
)

type BaseController struct {
	Cfg *config.Config
}

func NewBaseController(cfg *config.Config) *BaseController {
	return &BaseController{Cfg: cfg}
}
