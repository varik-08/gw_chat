package config

import (
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	instance *App
	onceApp  sync.Once
)

type App struct {
	Config       *Cfg
	DB           *pgxpool.Pool
	Repositories *Repository
	Services     *Service
}

func GetApp() (*App, error) {
	var err error

	onceApp.Do(func() {
		instance, err = initApp()
	})

	return instance, err
}

func initApp() (*App, error) {
	conf, err := GetConfig()
	if err != nil {
		return nil, err
	}

	db := initDB(conf.DB)

	repos := newRepository(db)
	services := newService(conf, repos)

	return &App{
		Config:       conf,
		DB:           db,
		Repositories: repos,
		Services:     services,
	}, nil
}
