package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/labstack/echo/v4"
	"github.com/matthewyuh246/aws-cognito/internal/controller"
	"github.com/matthewyuh246/aws-cognito/internal/domain"
	"github.com/matthewyuh246/aws-cognito/internal/repository"
	"github.com/matthewyuh246/aws-cognito/internal/routes"
	"github.com/matthewyuh246/aws-cognito/internal/usecase"
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

	// リポジトリの初期化
	userRepo := repository.NewUserRepository(db)
	
	authConfig := repository.AuthConfig{
		CognitoDomain:    utils.GetEnv("COGNITO_DOMAIN_URL", ""),
		UserPoolClientID: config.UserPoolClientID,
		AllowedDomains: []string{
			"http://localhost:3000",
			"http://localhost:5173",
			utils.GetEnv("FE_URL", "http://localhost:5173"),
		},
	}
	authRepo := repository.NewAuthRepository(authConfig)

	// usecaseの初期化
	authUsecase := usecase.NewAuthUsecase(
		userRepo,
		authRepo,
		awsSession,
		config.UserPoolID,
		config.JWTSecret,
	)

	// controllerの初期化
	authController := controller.NewAuthController(authUsecase)

	// Echoサーバーの初期化
	e := echo.New()
	
	// ルート設定
	routes.SetupRoutes(e, authController)

	// サーバー起動（優雅な終了付き）
	port := utils.GetEnv("PORT", "8080")
	log.Printf("Server starting on port %s", port)

	// 優雅な終了の実装
	go func() {
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 終了シグナルの待機
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Server is shutting down...")

	// 10秒のタイムアウトでサーバーを停止
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}