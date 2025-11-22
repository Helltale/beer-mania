package config

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
)

//nolint:golines // long struct tags with metadata
type DatabaseConfig struct {
	Host         string `env:"POSTGRES_HOST" env-default:"localhost" validate:"required"`
	Port         string `env:"POSTGRES_PORT" env-default:"5432" validate:"required"`
	User         string `env:"POSTGRES_USER" env-default:"beermania_user" validate:"required"`
	Password     string `env:"POSTGRES_PASSWORD" env-default:"beermania_password" validate:"required"`
	DBName       string `env:"POSTGRES_DB" env-default:"beermania_db" validate:"required"`
	SSLMode      string `env:"POSTGRES_SSLMODE" env-default:"disable" validate:"oneof=disable require verify-full verify-ca"`
	MaxIdleConns int    `env:"POSTGRES_MAX_IDLE_CONNS" env-default:"10" validate:"min=1,max=1000"`
	MaxOpenConns int    `env:"POSTGRES_MAX_OPEN_CONNS" env-default:"100" validate:"min=1,max=10000"`
}

//nolint:golines // long struct tags with metadata
type RabbitMQConfig struct {
	Host     string `env:"RABBITMQ_HOST" env-default:"localhost" validate:"required"`
	Port     string `env:"RABBITMQ_PORT" env-default:"5672" validate:"required"`
	User     string `env:"RABBITMQ_USER" env-default:"beermania_user" validate:"required"`
	Password string `env:"RABBITMQ_PASSWORD" env-default:"beermania_password" validate:"required"`
	VHost    string `env:"RABBITMQ_VHOST" env-default:"/" validate:"required"`
}

//nolint:golines // long struct tags with metadata
type MinIOConfig struct {
	Endpoint                    string `env:"MINIO_ENDPOINT" env-default:"localhost:9000" validate:"required"`
	AccessKey                   string `env:"MINIO_ACCESS_KEY" env-default:"minioadmin" validate:"required"`
	SecretKey                   string `env:"MINIO_SECRET_KEY" env-default:"minioadmin" validate:"required"`
	UseSSL                      bool   `env:"MINIO_USE_SSL" env-default:"false"`
	BucketUploads               string `env:"MINIO_BUCKET_UPLOADS" env-default:"uploads" validate:"required"`
	BucketProcessed             string `env:"MINIO_BUCKET_PROCESSED" env-default:"processed" validate:"required"`
	PresignedURLExpirationHours int    `env:"MINIO_PRESIGNED_URL_EXPIRATION_HOURS" env-default:"168" validate:"min=1,max=8760"` // Default: 7 days (168 hours), max: 1 year
}

func (c *MinIOConfig) PresignedURLExpiration() time.Duration {
	return time.Duration(c.PresignedURLExpirationHours) * time.Hour
}

//nolint:golines // long struct tags with metadata
type BackendConfig struct {
	Port     string `env:"BACKEND_PORT" env-default:"8080" validate:"required"`
	LogLevel string `env:"BACKEND_LOG_LEVEL" env-default:"info" validate:"oneof=debug info warn error"`
	Env      string `env:"BACKEND_ENV" env-default:"development" validate:"oneof=development production staging"`
}

type Config struct {
	Database DatabaseConfig
	RabbitMQ RabbitMQConfig
	MinIO    MinIOConfig
	Backend  BackendConfig
}

func Load() (*Config, error) {
	cfg := &Config{}

	// Load each nested config using cleanenv
	if err := cleanenv.ReadEnv(&cfg.Database); err != nil {
		return nil, fmt.Errorf("failed to load database configuration: %w", err)
	}

	if err := cleanenv.ReadEnv(&cfg.RabbitMQ); err != nil {
		return nil, fmt.Errorf("failed to load rabbitmq configuration: %w", err)
	}

	if err := cleanenv.ReadEnv(&cfg.MinIO); err != nil {
		return nil, fmt.Errorf("failed to load minio configuration: %w", err)
	}

	if err := cleanenv.ReadEnv(&cfg.Backend); err != nil {
		return nil, fmt.Errorf("failed to load backend configuration: %w", err)
	}

	// Validate configuration using validator
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	validate := validator.New()

	// Validate nested structs
	if err := validate.Struct(c.Database); err != nil {
		return fmt.Errorf("database config validation failed: %w", err)
	}

	if err := validate.Struct(c.RabbitMQ); err != nil {
		return fmt.Errorf("rabbitmq config validation failed: %w", err)
	}

	if err := validate.Struct(c.MinIO); err != nil {
		return fmt.Errorf("minio config validation failed: %w", err)
	}

	if err := validate.Struct(c.Backend); err != nil {
		return fmt.Errorf("backend config validation failed: %w", err)
	}

	return nil
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}
