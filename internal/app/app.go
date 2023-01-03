package app

import (
	"context"
	"message/internal/domain/message"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"message/internal/cache"
	"message/internal/config"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type Domain struct {
	Message DomainMessage
}

type DomainMessage struct {
	Service message.Service
}

type App struct {
	cfg    *config.Config
	db     *gorm.DB
	Cache  cache.Cache
	Domain Domain
}

func New(config *config.Config) *App {
	return &App{cfg: config}
}

func (m *App) DB() *gorm.DB {
	return m.db
}

func (m *App) Init() error {
	if err := m.initDB(); err != nil {
		return err
	}
	if err := m.initServices(); err != nil {
		return err
	}

	return nil
}

func (m *App) initDB() error {
	switch m.cfg.Repository.Type {
	case "postgres":
		db, err := gorm.Open(postgres.Open(m.cfg.DB.Postgres.Dsn), &gorm.Config{})
		if err != nil {
			return err
		}
		m.db = db
	}
	return nil
}

func (m *App) initCache() (err error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     m.cfg.Redis.Address,
		Password: m.cfg.Redis.Password,
	})
	err = rdb.Ping(context.Background()).Err()
	if err != nil {
		return errors.New("cannot connect to redis")
	}

	m.Cache = cache.New(rdb)

	return nil
}

func (m *App) initServices() (err error) {
	rep := message.NewRepository(m.db)
	m.Domain.Message.Service = *message.New(rep)
	return nil
}
