package aws

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/matthewyuh246/aws-cognito/pkg/utils"
)

type Config struct {
	Port string
	AWSRegion string
	UserPoolID string
	UserPoolClientID string
	JWTSecret string
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

func loadCognitoConfig() *Config {
	return &Config{
		Port: utils.GetEnv("PORT", "8080"),
		AWSRegion: utils.GetEnv("AWS_REGION", ""),
		UserPoolID: utils.GetEnv("USER_POOL_ID", ""),
		UserPoolClientID: utils.GetEnv("USER_POOL_CLIENT_ID", ""),
		JWTSecret: utils.GetEnv("JWT_SECRET", "your-secret-key"),
	}
}