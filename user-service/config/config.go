package config

import (
	"github.com/19parwiz/user-service/pkg/mongo"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"log"

	"time"
)

type (
	Config struct {
		Mongo   mongo.Config
		Server  Server
		SMTP    SMTPConfig //  Added smptconfig here
		Version string     `env:"VERSION"`
	}

	Server struct {
		HTTPServer HTTPServer
		GRPCServer GRPCServer
	}

	HTTPServer struct {
		Port           int           `env:"HTTP_PORT,required"`
		ReadTimeout    time.Duration `env:"HTTP_READ_TIMEOUT" envDefault:"30s"`
		WriteTimeout   time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"30s"`
		IdleTimeout    time.Duration `env:"HTTP_IDLE_TIMEOUT" envDefault:"60s"`
		MaxHeaderBytes int           `env:"HTTP_MAX_HEADER_BYTES" envDefault:"1048576"` // 1 MB
		Mode           string        `env:"GIN_MODE" envDefault:"release"`              // Can be: release, debug, test
	}

	GRPCServer struct {
		Port    int           `env:"GRPC_PORT,required"`
		Timeout time.Duration `env:"GRPC_TIMEOUT" envDefault:"30s"`
	}

	SMTPConfig struct { // <-- Added SMTPConfig struct
		Host     string `env:"SMTP_HOST,required"`
		Port     int    `env:"SMTP_PORT,required"`
		Username string `env:"SMTP_USERNAME,required"`
		Password string `env:"SMTP_PASSWORD,required"`
	}
)

func New() (*Config, error) {
	//Loading local .env file for private configuration
	if err := godotenv.Load("local.env"); err != nil {
		log.Printf("Error loading local.env file")
	}

	var cfg Config
	err := env.Parse(&cfg)

	return &cfg, err
}
