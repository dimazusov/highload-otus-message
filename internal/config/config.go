package config

import (
	"bytes"
	"io/ioutil"
	"log"
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

type Mysql struct {
	Dialect               string   `config:"dialect"`
	ReaderNodes           []string `config:"reader_nodes"`
	WriterNodes           []string `config:"writer_nodes"`
	MaxConnLifetimeMinute int      `config:"maxConnLifetimeMinute"`
	MaxConn               int      `config:"maxConn"`
}

type DB struct {
	Mysql Mysql `config:"msyql"`
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

	res := viper.GetStringSlice("db.mysql.reader_nodes")

	log.Println(res[0])

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
			Mysql: Mysql{
				Dialect:               viper.GetString("db.mysql.dialect"),
				ReaderNodes:           viper.GetStringSlice("db.mysql.reader_nodes"),
				WriterNodes:           viper.GetStringSlice("db.mysql.writer_nodes"),
				MaxConnLifetimeMinute: viper.GetInt("db.mysql.maxConnLifetimeMinute"),
				MaxConn:               viper.GetInt("db.mysql.maxConn"),
			},
		},
		Repository: Repository{
			Type: viper.GetString("repository.type"),
		},
	}

	return cfg, nil
}
