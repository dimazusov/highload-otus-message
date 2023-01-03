package main

import (
	"flag"
	"log"
	"math/rand"
	"message/internal/config"
	"message/internal/domain/message"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", config.DefaultConfigPath, "Path to configuration file")
	flag.Parse()
}

func main() {
	cfg, err := config.New(configFile)
	if err != nil {
		log.Fatalln(err)
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  cfg.DB.Postgres.Dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		log.Println(err)
	}

	for i := 0; i < 100; i++ {
		msg := &message.Message{
			FromUserID: uint(rand.Int()),
			ToUserID:   uint(rand.Int()),
			Text:       "test",
			CreatedAt:  time.Now(),
		}
		err := db.Exec("INSERT INTO message VALUES (nextval('serial_message_id'), ?,?,?,?,null);",
			msg.FromUserID,
			msg.ToUserID,
			msg.Text,
			msg.CreatedAt)
		if err != nil {
			log.Println(err)
		}
	}

	log.Println("Success!")
}
