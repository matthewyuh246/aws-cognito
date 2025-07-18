package awsconfig

import (
	"github.com/matthewyuh246/aws-cognito/pkg/utils"
)

type Config struct {
	Port string
	AWSRegion string
	UserPoolID string
	UserPoolClientID string
	JWTSecret string
}

func LoadCognitoConfig() *Config {
	return &Config{
		Port: utils.GetEnv("PORT", "8080"),
		AWSRegion: utils.GetEnv("AWS_REGION", ""),
		UserPoolID: utils.GetEnv("USER_POOL_ID", ""),
		UserPoolClientID: utils.GetEnv("USER_POOL_CLIENT_ID", ""),
		JWTSecret: utils.GetEnv("JWT_SECRET", "your-secret-key"),
	}
}