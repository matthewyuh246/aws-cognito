package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/matthewyuh246/aws-cognito/internal/domain"
	awsconfig "github.com/matthewyuh246/aws-cognito/pkg/aws"
	"github.com/matthewyuh246/aws-cognito/pkg/database"
	"github.com/matthewyuh246/aws-cognito/pkg/utils"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	dbConfig := database.NewConfig()
	return database.NewConnection(dbConfig)
}

func initAWS(region string) *session.Session {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		log.Fatalf("Failed to create AWS session: %v", err)
	}

	return sess
}

func migrateTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.User{},
	)
}

func main() {
	utils.LoadEnvFile()

	config := awsconfig.LoadCognitoConfig()

	db := initDB()
	defer database.Close(db)

	if err := migrateTables(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migration completed")

	awsSession := initAWS(config.AWSRegion)
}