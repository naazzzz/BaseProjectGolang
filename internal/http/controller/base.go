package controller

import (
	"BaseProjectGolang/internal/config"

	"github.com/soner3/flora"
)

type BaseController struct {
	flora.Component
	Cfg *config.Config
}

func NewBaseController(cfg *config.Config) *BaseController {
	return &BaseController{Cfg: cfg}
}
