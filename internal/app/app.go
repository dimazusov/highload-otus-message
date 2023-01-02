package app

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"message/internal/cache"
	"message/internal/config"
	"message/internal/domain/auth_token"
	"message/internal/domain/user"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type Domain struct {
	User      DomainUser
	AuthToken DomainAuthToken
}

type DomainAuthToken struct {
	Service auth_token.Service
}

type DomainUser struct {
	Repository user.Repository
	Service    user.Service
}

type Storage struct {
	writerNodes  []*sql.DB
	readersNodes []*sql.DB
}

type App struct {
	cfg     *config.Config
	storage *Storage
	Cache   cache.Cache
	Domain  Domain
}

func New(config *config.Config) *App {
	return &App{cfg: config}
}

func (m *App) DB() *Storage {
	return m.storage
}

func (m *App) Init() error {
	if err := m.initDB(); err != nil {
		return err
	}
	//if err := m.initCache(); err != nil {
	//	return err
	//}
	if err := m.initRepositories(); err != nil {
		return err
	}
	if err := m.initServices(); err != nil {
		return err
	}

	return nil
}

func (m *App) initDB() error {
	switch m.cfg.Repository.Type {
	case "mysql":
		m.storage = &Storage{
			readersNodes: make([]*sql.DB, 0, len(m.cfg.DB.Mysql.ReaderNodes)),
			writerNodes:  make([]*sql.DB, 0, len(m.cfg.DB.Mysql.WriterNodes)),
		}

		for _, wNode := range m.cfg.DB.Mysql.WriterNodes {
			db, err := sql.Open(m.cfg.DB.Mysql.Dialect, wNode)
			if err != nil {
				return errors.Wrap(err, "cannot connect to db")
			}
			db.SetConnMaxLifetime(time.Minute * 3)
			db.SetMaxOpenConns(m.cfg.DB.Mysql.MaxConn)
			db.SetMaxIdleConns(m.cfg.DB.Mysql.MaxConn)
			m.storage.writerNodes = append(m.storage.writerNodes, db)
		}

		for _, rNode := range m.cfg.DB.Mysql.ReaderNodes {
			db, err := sql.Open(m.cfg.DB.Mysql.Dialect, rNode)
			if err != nil {
				return errors.Wrap(err, "cannot connect to db")
			}
			db.SetConnMaxLifetime(time.Minute * 3)
			db.SetMaxOpenConns(m.cfg.DB.Mysql.MaxConn)
			db.SetMaxIdleConns(m.cfg.DB.Mysql.MaxConn)
			m.storage.readersNodes = append(m.storage.readersNodes, db)
		}
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

func (m *App) initRepositories() (err error) {
	m.Domain.User.Repository = user.NewRepository(m.storage.readersNodes, m.storage.writerNodes)

	return nil
}

func (m *App) initServices() (err error) {
	m.Domain.User.Service = user.NewService(m.Domain.User.Repository)
	m.Domain.AuthToken.Service = auth_token.NewService()

	return nil
}
