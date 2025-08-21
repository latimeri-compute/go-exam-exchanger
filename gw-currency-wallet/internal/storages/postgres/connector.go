package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBOptions struct {
	DBName     string
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     int
}

func NewConnection(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}
