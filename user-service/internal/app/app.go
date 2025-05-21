package app

import (
	"context"
	"fmt"
	"github.com/19parwiz/user-service/config"
	"github.com/19parwiz/user-service/internal/adapter/grpc"
	"github.com/19parwiz/user-service/internal/adapter/mail"
	"github.com/19parwiz/user-service/internal/adapter/mongo"
	"github.com/19parwiz/user-service/internal/usecase"
	"github.com/19parwiz/user-service/pkg/hashing"
	mongoConn "github.com/19parwiz/user-service/pkg/mongo"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const serviceName = "user-service"

type App struct {
	grpcServer *grpc.ServerAPI
}

func NewApp(ctx context.Context, cfg *config.Config) (*App, error) {
	log.Printf(fmt.Sprintf("Initializing %s service...", serviceName))

	log.Println("Connecting to DB:", cfg.Mongo.Database)
	mongoDB, err := mongoConn.NewDB(ctx, cfg.Mongo)
	if err != nil {
		return nil, fmt.Errorf("error connecting to DB: %v", err)
	}

	aiRepo := mongo.NewAutoInc(mongoDB.Conn)
	userRepo := mongo.NewUserRepo(mongoDB.Conn)

	hasher := hashing.NewBcryptHasher()

	// Initialize the mailer here
	mailer := mail.NewMailer(
		cfg.SMTP.Host,
		cfg.SMTP.Port,
		cfg.SMTP.Username,
		cfg.SMTP.Password,
	)

	userUsecase := usecase.NewUserUsecase(aiRepo, userRepo, hasher, mailer)

	grpcServer := grpc.New(cfg.Server, userUsecase)

	app := &App{
		grpcServer: grpcServer,
	}

	return app, nil
}

func (app *App) Start() error {
	errCh := make(chan error)

	app.grpcServer.Run(errCh)

	log.Printf(fmt.Sprintf("Starting %s service!", serviceName))

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case errRun := <-errCh:
		return errRun
	case sig := <-shutdownCh:
		log.Printf(fmt.Sprintf("Received shutdown signal: %s", sig.String()))
		app.Stop()
		log.Printf(fmt.Sprintf("Stopping %s service!", serviceName))
	}
	return nil
}

func (app *App) Stop() {
	err := app.grpcServer.Stop()
	if err != nil {
		log.Printf(fmt.Sprintf("Error stopping %s service: %v", serviceName, err))
	}
}
