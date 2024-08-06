package main

import (
	"context"
	"os"
	"os/signal"
	asynqmonauth "simpsons310/asynqmon-auth/internal"
	"strconv"
	"syscall"
)

func main() {
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer done()

	cfg, err := asynqmonauth.LoadEnv()
	if err != nil {
		panic(err)
	}

	// If a port number is provided as a command-line argument, use that.
	// In docker environment, the port is fixed (8080), reading from env can be override when running the container,
	// thus, providing a way to override the port number via command-line argument.
	if len(os.Args) == 2 {
		port, err := strconv.Atoi(os.Args[1])
		if err == nil && port > 0 {
			cfg.Server.Port = port
		}
	}

	app, err := asynqmonauth.NewApplication(cfg, nil)
	if err != nil {
		panic(err)
	}

	if err := app.StartServer(ctx); err != nil {
		panic(err)
	}
}
