package config

import (
	"BaseProjectGolang/pkg/mail"
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

const (
	ProdEnv  = "prod"
	DevEnv   = "dev"
	TestEnv  = "test"
	DebugEnv = "debug"
)

type Config struct {
	AppPort     string `envconfig:"APP_PORT" default:"8080" validate:"required"`
	SwaggerHost string `envconfig:"SWAGGER_HOST" default:"127.0.0.1:8080" validate:"required"`
	AppWorkMode string `envconfig:"APP_WORK_MODE" default:"debug" validate:"required"`
	ImagesTag   string `envconfig:"IMAGES_TAG" default:"v1.0.0" validate:"required"`
	Databases   *Databases
	Logs        *Logs
	App         *App
	Secure      *Secure
	Mail        *mail.GomailServiceConfig
}

func NewConfig(isTest bool, envFile ...string) (*Config, error) {
	var (
		config Config
		err    error
	)

	// Initialize Viper for configuration
	viperConfig := viper.New()
	viperConfig.SetConfigType("yml")

	// Handle test configuration
	if isTest {
		if len(envFile) > 0 && envFile[0] != "" {
			// Load environment variables from .env file if provided
			if err := godotenv.Load(envFile[0]); err != nil {
				log.Printf("Warning: error loading .env file: %viperConfig", err)
			}
		} else {
			log.Println("No test env file specified, using default environment")
		}

		// If a config file path is provided as second argument, load it directly
		if len(envFile) > 1 && envFile[1] != "" {
			viperConfig.SetConfigFile(envFile[1])

			if err = viperConfig.ReadInConfig(); err != nil {
				return nil, fmt.Errorf("error reading config file: %w", err)
			}
		}
	}

	// Process environment variables
	if err := envconfig.Process("", &config); err != nil {
		log.Printf("Warning: error processing environment variables: %viperConfig", err)
	}

	// Initialize nested structs with default values if they weren't set
	if config.Databases == nil {
		config.Databases = &Databases{
			Pgsql: &Pgsql{},
		}
	}

	if config.Logs == nil {
		config.Logs = &Logs{}
	}

	if config.App == nil {
		config.App = &App{}
	}

	// Set derived values
	if config.App != nil && config.App.Host != "" {
		config.App.URL = "http://" + config.App.Host
	}

	// Unmarshal
	if err = viperConfig.UnmarshalKey("notification_recipients", &config.Mail.From); err != nil {
		return nil, fmt.Errorf("error unmarshaling services config: %w", err)
	}

	// Validate configuration
	if err = validator.New().Struct(config); err != nil {
		return &config, fmt.Errorf("configuration validation failed: %w", err)
	}

	return &config, nil
}
