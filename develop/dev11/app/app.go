package app

import (
	"dev11/config"
	"dev11/storage"
	"github.com/sirupsen/logrus"
)

type App struct {
	MyLogger     *logrus.Logger
	Config       *config.Conf
	EventStorage *storage.EventStorage
}

func NewApp(MyLogger *logrus.Logger, Conf *config.Conf, EventStorage *storage.EventStorage) *App {
	return &App{
		MyLogger:     MyLogger,
		Config:       Conf,
		EventStorage: EventStorage,
	}
}
