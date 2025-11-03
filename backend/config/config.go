package config

import (
	"Inquiro/auth"
	jobpb "Inquiro/protos"
	"Inquiro/repositories"
	"Inquiro/utils/mailer"

	"github.com/alexedwards/scs/v2"
	"go.uber.org/zap"
)

type Config struct {
	Addr        string
	FrontendURL string
	DBConfig    DBConfig
	MailConfig  MailConfig
}

type Application struct {
	Config  Config
	Logger  *zap.SugaredLogger
	Store   repositories.Storage
	Auth    auth.Auth
	Mail    mailer.Client
	Session *scs.SessionManager
	Grpc    jobpb.JobServiceClient
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
