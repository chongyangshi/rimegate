package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/monzo/slog"
	"github.com/monzo/typhon"

	"github.com/icydoge/rimegate/config"
	"github.com/icydoge/rimegate/server"
)

func main() {
	initContext := context.Background()

	// Initialise server for incoming requests
	svc := server.Service()
	srv, err := typhon.Listen(svc, fmt.Sprintf("%s:%s", config.ConfigListenAddr, config.ConfigListenPort))
	if err != nil {
		panic(err)
	}
	slog.Info(initContext, "Rimegate incoming listening on %v", srv.Listener().Addr())

	// Log termination gracefully
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
	slog.Info(initContext, "Rimegate shutting down")
	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Stop(c)
}
