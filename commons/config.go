package commons

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds the configuration values for the application.
type Config struct {
	// Database configuration
	DBHost         string `mapstructure:"POSTGRES_HOST"`
	DBUserName     string `mapstructure:"POSTGRES_USER"`
	DBUserPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBName         string `mapstructure:"POSTGRES_DB"`
	DBNameTest     string `mapstructure:"POSTGRES_DB_TEST"`
	DBPort         string `mapstructure:"POSTGRES_PORT"`

	// Redis
	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     string `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`

	// Server configuration
	AppEnv         string `mapstructure:"APP_ENV"`
	ServerProtocol string `mapstructure:"PROTOCOL"`
	ServerHost     string `mapstructure:"HOSt"`
	ServerPort     string `mapstructure:"PORT"`
	ClientOrigin   string `mapstructure:"CLIENT_ORIGIN"`

	// JWT Tokens configuration
	AccessTokenPrivateKey string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey  string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiresIn  time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	AccessTokenMaxAge     int           `mapstructure:"ACCESS_TOKEN_MAXAGE"`

	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	RefreshTokenExpiresIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
	RefreshTokenMaxAge     int           `mapstructure:"REFRESH_TOKEN_MAXAGE"`

	// Minio
	MinioEndpoint  string `mapstructure:"MINIO_ENDPOINT"`
	MinioAccessKey string `mapstructure:"MINIO_ACCESS_KEY"`
	MinioSecretKey string `mapstructure:"MINIO_SECRET_KEY"`
	MinioBucket    string `mapstructure:"MINIO_BUCKET"`
	MinioLocation  string `mapstructure:"MINIO_LOCATION"`
}

// LoadConfig loads configuration from the specified path.
// It reads environment variables and populates the Config struct.
// Returns the loaded config and an error if any.
func LoadConfig(path string) *Config {
	viper.AutomaticEnv()

	var err error

	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	// Read the .env file
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("load_config_err: read config: %v", err))
	}

	// Unmarshal the config into the Config struct
	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(fmt.Errorf("load_config_err: unmarshal: %v", err))
	}

	return &cfg
}
