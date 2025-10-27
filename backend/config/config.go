package config

import (
	"Inquiro/auth"
	"Inquiro/repositories"
	"Inquiro/utils/mailer"

	"go.uber.org/zap"
)

type Config struct {
	Addr        string
	FrontendURL string
	DBConfig    DBConfig
	MailConfig  MailConfig
}

type Application struct {
	Config Config
	Logger *zap.SugaredLogger
	Store  repositories.Storage
	Auth   auth.Auth
	Mail   mailer.Client
}

type DBConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     int
}

type MailConfig struct {
	APIKey    string
	FromEmail string
}
