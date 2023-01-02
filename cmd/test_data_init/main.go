package main

import (
	"flag"
	"log"
	"message/internal/app"
	"message/internal/config"
	"message/internal/test_data"
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

	application := app.New(cfg)
	err = application.Init()
	if err != nil {
		log.Fatalln(err)
	}

	generator, err := test_data.NewGenerator(application.DB())
	if err != nil {
		log.Fatalln(err)
	}
	err = generator.GenerateTestData()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("generate successful")
}
