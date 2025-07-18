package main

import (
	"flag"
	"log"

	"github.com/matthewyuh246/aws-cognito/internal/domain"
	"github.com/matthewyuh246/aws-cognito/pkg/database"
	"gorm.io/gorm"
)

func runMigrationsUp(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.User{},
	)
}

func runMigrationsDown(db *gorm.DB) error {
	return db.Migrator().DropTable(
		&domain.User{},
	)
}

func main() {
	var (
		up = flag.Bool("up", false, "Run migrations up")
		down = flag.Bool("down", false, "Run migrations down")
	)
	flag.Parse()

	dbConfig := database.NewConfig()
	
	db := database.NewConnection(dbConfig)
	defer database.Close(db)

	if *up {
		log.Println("Running migrations up...")
		if err := runMigrationsUp(db); err != nil {
			log.Fatalf("Failed to run migrations up: %v", err)
		}
		log.Println("Migrations up completed successfully")
	} else if *down {
		log.Println("Running migrations down...")
		if err := runMigrationsDown(db); err != nil {
			log.Fatalf("Failed to run migrations down: %v", err)
		}
		log.Println("Migrations down completed successfully")
	} else {
		log.Println("Please specify -up or -down flag")
		flag.Usage()
	}
}
