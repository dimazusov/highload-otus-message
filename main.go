package main

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Company struct {
	ID        uint      `gorm:"id"`
	Name      string    `gorm:"name"`
	ImageUrl  string    `gorm:"image_url"`
	CreatedAt time.Time `gorm:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at"`
}

func main() {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=localhost user=postgres password= dbname=postgres port=5432 sslmode=disable",
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		log.Println(err)
	}

	for i := 0; i < 100; i++ {
		err := db.Create(&Company{
			Name:      "test",
			ImageUrl:  "test",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		if err != nil {
			log.Println(err)
		}
	}
}
