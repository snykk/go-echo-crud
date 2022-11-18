package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/snykk/go-echo-crud/entity"
)

type Config struct {
	Host     string
	Port     string
	Password string
	User     string
	DBName   string
	SSLMode  string
}

func InitializeDatabase(config *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return db, err
	}

	if err := AutoMigrate(db); err != nil {
		return nil, err
	}

	log.Println("success connecting to postgres and running migration")
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(&entity.Movie{})
	return err
}
