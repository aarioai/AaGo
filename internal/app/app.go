package app

import (
	"github.com/aarioai/AaGo/internal/app/config"
	"github.com/aarioai/AaGo/internal/app/helper"
	"github.com/aarioai/AaGo/internal/app/logger"
)

type App struct {
	Config *config.Config
	Log    logger.LogInterface

	Time *helper.Time
}

func New(cfgPath string, logger logger.LogInterface) *App {
	c := config.New(cfgPath)
	return &App{
		Config: c,
		Log:    logger,
		Time:   helper.NewTime(c.TimeLocation),
	}
}
