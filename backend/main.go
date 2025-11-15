package main

import (
	"Inquiro/auth"
	"Inquiro/config"
	"Inquiro/config/env"
	"Inquiro/controller"
	"Inquiro/db"
	jobpb "Inquiro/protos"
	"Inquiro/repositories"
	"Inquiro/routes"
	"Inquiro/services"
	"Inquiro/utils/mailer"
	_ "Inquiro/utils/token"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	PythonServerAddress = "localhost:50051"
)

func main() {
	logger := zap.Must(zap.NewProduction()).Sugar()
	configuration := config.Config{
		Addr:        "localhost:8080",
		FrontendURL: "http://localhost:3000",
		DBConfig: config.DBConfig{
			Host:     env.GetString("DB_HOST", "postgres"),
			User:     env.GetString("DB_USER", "admin"),
			Password: env.GetString("DB_PASSWORD", "1234"),
			DBName:   env.GetString("DB_NAME", "admin"),
			Port:     env.GetInt("DB_PORT", 5432),
		},
		MailConfig: config.MailConfig{
			APIKey:    env.GetString("RESEND_API", "re_2fo8WcM7_6uNEbMPou98kjNKoMZpoFsxw"),
			FromEmail: env.GetString("RESEND_FROM_EMAIL", "support@bloggerspot.xyz"),
		},
	}
	conn, err := grpc.NewClient(PythonServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	logger.Infow("Connecting to python service", "address : ", PythonServerAddress)
	if err != nil {
		logger.Fatalf("failed to connect to job service: %v", err.Error())
	}
	defer conn.Close()
	grpcClient := jobpb.NewJobServiceClient(conn)
	mailer := mailer.NewResendClient(configuration.MailConfig.APIKey, configuration.MailConfig.FromEmail, logger)

	cfg := config.Application{
		Config: configuration,
		Mail:   mailer,
		Logger: logger,
		Grpc:   grpcClient,
	}
	defer logger.Sync()

	db_conn, err := db.NewDB(cfg.Logger,
		cfg.Config.DBConfig.Host,
		cfg.Config.DBConfig.User,
		cfg.Config.DBConfig.Password,
		cfg.Config.DBConfig.DBName,
		cfg.Config.DBConfig.Port)
	if err != nil {
		logger.Fatalf("failed to connect to db: %v", err.Error())
	}
	defer db_conn.Close()
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any major browsers
	}))
	// Registering session manager
	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	r.Use(sessionManager.LoadAndSave)
	cfg.Session = sessionManager
	cfg.Auth = auth.NewAuth(cfg.Store, sessionManager)

	apiRouter := chi.NewRouter()
	cfg.Store = repositories.NewStorage(db_conn)

	// Handling users
	logger.Infof("registering user routes")
	srv := services.NewService(
		cfg.Store,
		cfg.Logger,
		cfg.Mail,
	)
	userController := controller.NewController(srv, cfg)
	userRoutes := routes.NewUserRoutes(userController)
	userRoutes.RegisterUserRoutes(apiRouter)

	// Handling resumes
	logger.Infof("regiter resume routes")
	resumeController := controller.NewController(srv, cfg)
	resumeRoutes := routes.NewResumeRoutes(resumeController)
	resumeRoutes.RegisterResumeRoutes(apiRouter)

	r.Mount("/api", apiRouter)

	Run(cfg, r)

}
