package config

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/spf13/viper"
)

const DefaultConfigPath = "configs/config.yaml"

type Http struct {
	Host string `config:"host"`
	Port string `config:"port"`
}

type Server struct {
	Http Http `config:"http"`
}

type Logger struct {
	Path  string `config:"path"`
	Level string `config:"level"`
}

type Postgres struct {
	Dialect               string `config:"dialect"`
	Dsn                   string `config:"dsn"`
	MaxConnLifetimeMinute int    `config:"maxConnLifetimeMinute"`
	MaxConn               int    `config:"maxConn"`
}

type DB struct {
	Postgres Postgres `config:"msyql"`
}

type Repository struct {
	Type string `config:"type"`
}

type Redis struct {
	Address  string `config:"address"`
	Password string `config:"password"`
}

type Kafka struct {
	Address       string `config:"address"`
	BatchTimeout  string `config:"batchtimeout"`
	DialerTimeout string `config:"dialertimeout"`
	Group         string `config:"group"`
}

type Config struct {
	Server     Server     `config:"server"`
	Logger     Logger     `config:"logger"`
	DB         DB         `config:"db"`
	Repository Repository `config:"repository"`
	Redis      Redis      `config:"redis"`
	Kafka      Kafka      `config:"kafka"`
}

func New(filePath string) (*Config, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	err = viper.ReadConfig(bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Server: Server{
			Http: Http{
				Host: viper.GetString("server.http.host"),
				Port: viper.GetString("server.http.port"),
			},
		},
		Logger: Logger{
			Path:  viper.GetString("logger.path"),
			Level: viper.GetString("logger.level"),
		},
		DB: DB{
			Postgres: Postgres{
				Dialect:               viper.GetString("db.postgres.dialect"),
				Dsn:                   viper.GetString("db.postgres.dsn"),
				MaxConnLifetimeMinute: viper.GetInt("db.postgres.maxConnLifetimeMinute"),
				MaxConn:               viper.GetInt("db.postgres.maxConn"),
			},
		},
		Repository: Repository{
			Type: viper.GetString("repository.type"),
		},
	}

	return cfg, nil
}
