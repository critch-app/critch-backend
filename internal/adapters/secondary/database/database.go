package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Adapter struct {
	db *gorm.DB
}

func NewAdapter(dsn string) (*Adapter, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return &Adapter{db}, nil
}

func (dbA *Adapter) Migrate(models ...any) error {
	return dbA.db.AutoMigrate(models...)
}
