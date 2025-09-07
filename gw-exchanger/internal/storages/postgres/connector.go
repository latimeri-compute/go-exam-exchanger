package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func NewConnection(dsn string) (*DB, error) {
	cfg := []gorm.Option{
		&gorm.Config{TranslateError: true},
	}
	db, err := gorm.Open(postgres.Open(dsn), cfg...)
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}
