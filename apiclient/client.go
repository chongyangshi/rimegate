package apiclient

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/monzo/slog"
	"github.com/monzo/typhon"

	"github.com/icydoge/rimegate/config"
	"github.com/icydoge/rimegate/types"
)

var cacheValidity time.Duration

func Init(ctx context.Context) error {
	timeOut, err := time.ParseDuration(config.ConfigGrafanaRenderTimeout)
	if err != nil {
		slog.Error(ctx, "Error parsing timeout from environment config: %v", err)
		return err
	}

	slog.Info(ctx, "Grafana request timeout: %v.", timeOut)

	if _, err = url.Parse(config.ConfigGrafanaHost); err != nil {
		slog.Error(ctx, "Invalid Grafana host URL %s from environment config: %v", config.ConfigGrafanaHost, err)
		return err
	}

	cacheValidity, err = time.ParseDuration(config.ConfigImageCacheValidity)
	if err != nil {
		slog.Error(ctx, "Invalid cache validity duration %s from environment config: %v", config.ConfigImageCacheValidity, err)
		return err
	}

	// Configure timeout
	roundTripper := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   timeOut,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   timeOut,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := func(req typhon.Request) typhon.Response {
		return typhon.HttpService(roundTripper)(req)
	}

	typhon.Client = client

	return nil
}

func setAuthenticationCredentials(req *typhon.Request, auth *types.Auth) (*typhon.Request, error) {
	switch {
	case config.ConfigGrafanaAPIToken != "":
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.ConfigGrafanaAPIToken))
	case auth != nil:
		// Empty username / password supported.
		req.SetBasicAuth(auth.GrafanaUsername, auth.GrafanaPassword)
	}

	// No credential is supported, depending on Grafana authentication mode
	return req
}