package asynqmonauth

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

const (
	AuthModeNone  = "none"
	AuthModeBasic = "basic"
	AuthModeHttp  = "http"
)

type Config struct {
	Server *ServerConfig `env:", prefix=SERVER_"`
	Asynq  *AsynqConfig  `env:", prefix=ASYNQ_"`
}

type ServerConfig struct {
	// In Docker environment, the port is set to 8080
	Port      int        `env:"PORT"`
	AuthMode  string     `env:"AUTH_MODE"`
	AuthBasic *AuthBasic `env:", prefix=AUTH_BASIC_"`
}

type AsynqConfig struct {
	MonRootPath      string `env:"MON_ROOT_PATH"`
	ReadOnly         bool   `env:"READ_ONLY"`
	RedisDSN         string `env:"REDIS_DSN"`
	RedisInSecureTLS bool   `env:"REDIS_INSECURE_TLS"`
}

type AuthBasic struct {
	Username string `env:"USERNAME"`
	Password string `env:"PASSWORD"`
}

func LoadEnv() (*Config, error) {
	port := 8080
	if portEnv := os.Getenv("SERVER_PORT"); portEnv != "" {
		portInt, err := strconv.Atoi(portEnv)
		if err != nil {
			return nil, fmt.Errorf("invalid port number: %s", portEnv)
		}
		port = portInt
	}

	authMode := AuthModeNone
	if mode := os.Getenv("SERVER_AUTH_MODE"); mode != "" {
		authMode = mode
	}

	rootPath := os.Getenv("ASYNQ_MON_ROOT_PATH")
	if rootPath == "" {
		rootPath = "/"
	}

	asynqReadOnly := false
	if os.Getenv("ASYNQ_REDIS_DSN") == "true" {
		asynqReadOnly = true
	}

	redisDSN := os.Getenv("ASYNQ_REDIS_DSN")
	if redisDSN == "" {
		redisDSN = "redis://127.0.0.1:6379/0"
	}

	redisInSecureTLS := false
	if os.Getenv("ASYNQ_REDIS_INSECURE_TLS") == "true" {
		redisInSecureTLS = true
	}

	authBasicUsername := os.Getenv("SERVER_AUTH_BASIC_USERNAME")
	authBasicPassword := os.Getenv("SERVER_AUTH_BASIC_PASSWORD")
	if authMode == AuthModeBasic && (authBasicUsername == "" || authBasicPassword == "") {
		return nil, fmt.Errorf("basic auth requires username and password")
	}

	return &Config{
		Server: &ServerConfig{
			Port:     port,
			AuthMode: authMode,
			AuthBasic: &AuthBasic{
				Username: authBasicUsername,
				Password: authBasicPassword,
			},
		},
		Asynq: &AsynqConfig{
			MonRootPath:      rootPath,
			ReadOnly:         asynqReadOnly,
			RedisDSN:         redisDSN,
			RedisInSecureTLS: redisInSecureTLS,
		},
	}, nil
}

type Application struct {
	server *http.Server
	config *Config
	logger Logger
}

func NewApplication(cfg *Config, logger Logger) (*Application, error) {
	if logger == nil {
		logger = NewLogger()
	}

	asynqHandler, err := newAsynqMonHandler(cfg.Asynq)
	if err != nil {
		return nil, err
	}

	// Add auth middleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		asynqHandler.ServeHTTP(w, r)
	})
	logger.Println("auth mode:", cfg.Server.AuthMode)

	if cfg.Server.AuthMode == AuthModeBasic {
		handler = basicAuthHandler(handler, cfg.Server.AuthBasic)
	} else if cfg.Server.AuthMode == AuthModeHttp {
		return nil, fmt.Errorf("http auth not implemented yet")
	}

	// Register handler to server
	logger.Println("root path:", asynqHandler.RootPath())
	if asynqHandler.RootPath() != "/" {
		redirectToRootHandler(asynqHandler.RootPath())
	}
	http.Handle(asynqHandler.RootPath()+"/", handler)
	srv := newServer(cfg.Server, handler)

	return &Application{
		logger: logger,
		config: cfg,
		server: srv,
	}, nil
}
