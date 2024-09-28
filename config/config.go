package config

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Log      *LogConfig
	Database *DatabaseConfig
	Redis    *RedisConfig
	Server   *ServerConfig
}

type LogConfig struct {
	Level  slog.Level
	Format string
}

type DatabaseConfig struct {
	ConnectionURL string
}

type RedisConfig struct {
	Host     string
	Password string
	DB       int
}

type ServerConfig struct {
	Addr            string
	WriteTimeout    time.Duration
	ReadTimeout     time.Duration
	IdleTimeout     time.Duration
	GracefulTimeout time.Duration
	SessionKey      string
}

func NewConfig() (*Config, error) {
	database, err := NewDatabaseConfig()
	if err != nil {
		return nil, err
	}

	server, err := NewServerConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		Log:      NewLogConfig(),
		Database: database,
		Server:   server,
	}, nil
}

func NewLogConfig() *LogConfig {
	var level slog.Level
	levelStr := os.Getenv("LOG_LEVEL")
	switch levelStr {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	format := os.Getenv("LOG_FORMAT")
	if format != "json" {
		format = "text"
	}

	return &LogConfig{
		Level:  level,
		Format: format,
	}
}

func NewDatabaseConfig() (*DatabaseConfig, error) {
	mode := os.Getenv("MODE")
	if mode == "" {
		log.Fatal("MODE environment variable not set")
	}
	var envFile string
	if mode == "production" {
		envFile = ".env.production"
	} else {
		envFile = ".env.development"
	}
	err := godotenv.Load(envFile)

	if err != nil {
		log.Println(err)
		log.Fatal(" loading .env file")
	}

	host := os.Getenv("DB_HOST")
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, errors.New("invalid DB_PORT")
	}
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")
	schema := os.Getenv("DB_SCHEMA")

	query := url.Values{
		"sslmode":  []string{"disable"},
		"timezone": []string{"utc"},
	}
	if schema != "" {
		query.Add("search_path", schema)
	}
	connURL := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(user, pass),
		Host:     host + ":" + strconv.Itoa(port),
		Path:     dbname,
		RawQuery: query.Encode(),
	}
	return &DatabaseConfig{
		ConnectionURL: connURL.String(),
	}, nil
}

func NewServerConfig() (*ServerConfig, error) {
	writeTimeout, err := time.ParseDuration(os.Getenv("SERVER_WRITE_TIMEOUT"))
	fmt.Println("SERVER_WRITE_TIMEOUT:", os.Getenv("SERVER_WRITE_TIMEOUT"))
	fmt.Println("SERVER_WRITE_TIMEOUT:", writeTimeout)

	if err != nil {
		return nil, errors.New("invalid SERVER_WRITE_TIMEOUT")
	}
	readTimeout, err := time.ParseDuration(os.Getenv("SERVER_READ_TIMEOUT"))
	if err != nil {
		return nil, errors.New("invalid SERVER_READ_TIMEOUT")
	}
	idleTimeout, err := time.ParseDuration(os.Getenv("SERVER_IDLE_TIMEOUT"))
	if err != nil {
		return nil, errors.New("invalid SERVER_IDLE_TIMEOUT")
	}
	gracefulTimeout, err := time.ParseDuration(os.Getenv("SERVER_GRACEFUL_TIMEOUT"))
	if err != nil {
		return nil, errors.New("invalid SERVER_GRACEFUL_TIMEOUT")
	}
	sessionKey := os.Getenv("SESSION_KEY")

	fmt.Printf("addr %s", os.Getenv("SERVER_ADDR"))
	fmt.Printf("port %s", os.Getenv("SERVER_PORT"))

	return &ServerConfig{
		Addr:            fmt.Sprintf("%s:%s", os.Getenv("SERVER_ADDR"), os.Getenv("SERVER_PORT")),
		GracefulTimeout: gracefulTimeout,
		WriteTimeout:    writeTimeout,
		ReadTimeout:     readTimeout,
		IdleTimeout:     idleTimeout,
		SessionKey:      sessionKey,
	}, nil
}
