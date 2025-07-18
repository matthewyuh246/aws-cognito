package database

import (
	"fmt"
	"log"

	"github.com/matthewyuh246/aws-cognito/pkg/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host string
	Port string
	Database string
	User string
	Password string
	SSLMode string
}

func NewConfig() *Config {
	return &Config{
		Host: utils.GetEnv("POSTGRES_HOST", "localhost"),
		Port: utils.GetEnv("POSTGRES_PORT", "5445"),
		Database: utils.GetEnv("POSTGRES_DB", ""),
		User: utils.GetEnv("POSTGRES_USER", ""),
		Password: utils.GetEnv("POSTGRES_PW", ""),
		SSLMode: utils.GetEnv("POSTGRES_SSLMODE", "disable"),
	}
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode)
}

func NewConnection(config *Config) *gorm.DB {
	dsn := config.GetDSN()
	log.Printf("Connecting to database: host=%s port=%s dbname=%s user=%s",
		config.Host, config.Port, config.Database, config.User)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return db
}

func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}