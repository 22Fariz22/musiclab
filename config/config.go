package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// App config struct
type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	Logger   Logger
	Redis    RedisConfig
}

// Server config struct
type ServerConfig struct {
	AppVersion        string
	Port              string
	Mode              string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	CtxDefaultTimeout time.Duration
	Debug             bool
}

// Logger config
type Logger struct {
	Development       bool
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
}

// Postgresql config
type PostgresConfig struct {
	PostgresqlHost     string
	PostgresqlPort     string
	PostgresqlUser     string
	PostgresqlPassword string
	PostgresqlDbname   string
	PostgresqlSSLMode  bool
	PgDriver           string
}

type RedisConfig struct {
	Addr         string
	Password     string
	DB           int
	DefaultDB    int
	MinIdleConns int
	PoolSize     int
	PoolTimeout  int
}

// LoadConfig reads environment variables into a Config struct
func LoadConfig() (*Config, error) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. Falling back to environment variables.")
	}

	return &Config{
		Server: ServerConfig{
			AppVersion:        getEnv("APP_VERSION", "1.0.0"),
			Port:              getEnv("SERVER_PORT", ":8080"),
			Mode:              getEnv("MODE", "Development"),
			ReadTimeout:       getEnvAsDuration("READ_TIMEOUT", 10*time.Second),
			WriteTimeout:      getEnvAsDuration("WRITE_TIMEOUT", 10*time.Second),
			CtxDefaultTimeout: getEnvAsDuration("CTX_DEFAULT_TIMEOUT", 12*time.Second),
			Debug:             getEnvAsBool("DEBUG", false),
		},
		Logger: Logger{
			Development:       getEnvAsBool("LOGGER_DEVELOPMENT", true),
			DisableCaller:     getEnvAsBool("LOGGER_DISABLE_CALLER", false),
			DisableStacktrace: getEnvAsBool("LOGGER_DISABLE_STACKTRACE", false),
			Encoding:          getEnv("LOGGER_ENCODING", "console"),
			Level:             getEnv("LOGGER_LEVEL", "info"),
		},
		Postgres: PostgresConfig{
			PostgresqlHost:     getEnv("POSTGRES_HOST", "localhost"),
			PostgresqlPort:     getEnv("POSTGRES_PORT", "5432"),
			PostgresqlUser:     getEnv("POSTGRES_USER", "postgres"),
			PostgresqlPassword: getEnv("POSTGRES_PASSWORD", "postgres"),
			PostgresqlDbname:   getEnv("POSTGRES_DBNAME", "music_db"),
			PostgresqlSSLMode:  getEnvAsBool("POSTGRES_SSLMODE", false),
			PgDriver:           getEnv("POSTGRES_DRIVER", "pgx"),
		},
		Redis: RedisConfig{
			Addr:         getEnv("REDIS_ADDR", "redis:6379"),
			Password:     getEnv("REDIS_PASSWORD", ""),
			DB:           getEnvAsInt("REDIS_DB", 0),
			DefaultDB:    getEnvAsInt("REDIS_DEFAULT_DB", 0),
			MinIdleConns: getEnvAsInt("REDIS_MIN_IDLE_CONNS", 10),
			PoolSize:     getEnvAsInt("REDIS_POOL_SIZE", 500),
			PoolTimeout:  getEnvAsInt("REDIS_POOL_TIMEOUT", 30),
		},
	}, nil
}

// Helper functions to parse environment variables

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valStr := getEnv(key, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valStr := getEnv(key, "")
	if val, err := time.ParseDuration(valStr); err == nil {
		return val
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valStr := getEnv(key, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultValue
}
