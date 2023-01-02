package test_data

import (
	"database/sql"
	"log"
	"math/rand"
	"message/internal/pkg/hasher"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/cheggaaa/pb/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"message/internal/domain/user"
)

const countBuildings = 1000000
const batchSize = 5000
const minAge = 18
const maxAge = 30

type generator struct {
	db *gorm.DB
}

func NewGenerator(sqlDB *sql.DB) (*generator, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.Logger = newLogger()
	return &generator{
		db: db,
	}, nil
}

func newLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
}

func (m generator) GenerateTestData() (err error) {
	log.Println("generate buildings")
	if err := m.generateUsers(); err != nil {
		return err
	}

	return nil
}
func New() *boolgen {
	return &boolgen{src: rand.NewSource(time.Now().UnixNano())}
}
func (m generator) generateUsers() error {
	// Сгенерировать любым способ 1,000,000 анкет.
	// Имена и Фамилии должны быть реальными (чтобы учитывать селективность индекса)
	users := make([]user.User, 0, batchSize)
	boolGen := newBoolGenerator()

	bar := pb.StartNew(countBuildings)
	for i := 0; i < countBuildings; i += 5000 {
		for j := 0; j < batchSize; j++ {
			age := uint(rand.Intn(maxAge-minAge) + minAge)
			pass := gofakeit.Password(true, false, true, false, false, 0)
			hashedPass, err := hasher.New().GetHashFromStruct(pass)
			if err != nil {
				return err
			}
			users = append(users, user.User{
				Name:     gofakeit.Name(),
				Email:    gofakeit.Email(),
				Password: hashedPass,
				Age:      age,
				Sex:      boolGen.Bool(),
				City:     gofakeit.City(),
				Interest: gofakeit.Phrase(),
			})

			err = m.db.Create(&users).Error
			if err != nil {
				return err
			}

			bar.Add(len(users))

			users = make([]user.User, 0, batchSize)
		}
	}
	bar.Finish()

	return nil
}
