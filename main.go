package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/monzo/slog"
	"github.com/monzo/typhon"

	"github.com/icydoge/rimegate/apiclient"
	"github.com/icydoge/rimegate/config"
	"github.com/icydoge/rimegate/server"
)

func main() {
	initContext := context.Background()

	// Initialise client for forwarding requests to Grafana
	err := apiclient.Init(initContext)
	if err != nil {
		panic(err)
	}
	slog.Info(initContext, "Rimegate client initialised for Grafana server at %s", config.ConfigGrafanaHost)

	// Initialise server for incoming requests
	svc := server.Service()
	srv, err := typhon.Listen(svc, config.ConfigListenAddr)
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
