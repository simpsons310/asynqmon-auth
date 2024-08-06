package asynqmonauth

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
)

const (
	ServerShutdownTimeout = 5 * time.Second
)

func newServer(cfg *ServerConfig, handler http.Handler) *http.Server {
	addr := ":" + strconv.Itoa(cfg.Port)
	return &http.Server{
		Addr:    addr,
		Handler: handler,
	}
}

func (a *Application) StartServer(ctx context.Context) error {
	logger := a.logger

	errChan := make(chan error, 1)
	go func() {
		<-ctx.Done()

		logger.Println("context closed")
		shutdownCtx, done := context.WithTimeout(context.Background(), ServerShutdownTimeout)
		defer done()

		logger.Println("shutting down server")
		errChan <- a.server.Shutdown(shutdownCtx)
	}()

	logger.Println("server started at", a.server.Addr)
	if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server serve error: %w", err)
	}

	logger.Println("server shutdown gracefully")
	if err := <-errChan; err != nil {
		return fmt.Errorf("server shutdown error: %w", err)
	}

	return nil
}

func newAsynqMonHandler(cfg *AsynqConfig) (*asynqmon.HTTPHandler, error) {
	redisParsedOpt, err := asynq.ParseRedisURI(cfg.RedisDSN)
	if err != nil {
		return nil, err
	}

	redisOpt, ok := redisParsedOpt.(asynq.RedisClientOpt)
	if !ok {
		return nil, fmt.Errorf("invalid redis option type")
	}

	if cfg.RedisInSecureTLS {
		redisOpt.TLSConfig = &tls.Config{
			// TODO: add custom TLS configuration
			InsecureSkipVerify: true,
		}
	}

	return asynqmon.New(asynqmon.Options{
		RootPath:     cfg.MonRootPath,
		RedisConnOpt: redisOpt,
		// TODO: implement later
		// ResultFormatter ResultFormatter
		ReadOnly: true,
	}), nil
}
