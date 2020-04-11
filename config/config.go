package config

import (
	"fmt"
	"os"

	"github.com/monzo/terrors"
)

var (
	ConfigCORSAllowedOrigin    = getConfigFromOSEnv("CORS_ALLOWED_ORIGIN", "*", true)
	ConfigListenAddr           = getConfigFromOSEnv("LISTEN_ADDR", ":8080", true)
	ConfigGrafanaHost          = getConfigFromOSEnv("GRAFANA_HOST", "", false)
	ConfigGrafanaRenderTimeout = getConfigFromOSEnv("GRAFANA_RENDER_TIMEOUT", "1m", true)
	ConfigGrafanaMaxPeriod     = getConfigFromOSEnv("GRAFANA_MAX_PERIOD", "3h", true)
	ConfigGrafanaDefaultPeriod = getConfigFromOSEnv("GRAFANA_DEFAULT_PERIOD", "1h", true)
	ConfigImageCacheValidity   = getConfigFromOSEnv("IMAGE_CACHE_VALIDITY", "2m", true)

	// In mTLS mode, the reverse proxy in front of rimegate backend should check client
	// certificate, in which case the backend will not attempt to forward Grafana credentials,
	// but instead use a static read-only API token from Grafana to authenticate itself. If
	// this is set, username and password forwarded from client will be ignored.
	ConfigGrafanaAPIToken = getConfigFromOSEnv("GRAFANA_API_TOKEN", "", true)
)

// This is intended to run inside Kubernetes as a pod , so we just set service configurations from
// deployment Configuration.
func getConfigFromOSEnv(key, defaultValue string, canBeEmpty bool) string {
	envValue := os.Getenv(key)
	if envValue != "" {
		return envValue
	}

	if !canBeEmpty {
		panic(terrors.InternalService("invalid_Config", fmt.Sprintf("Config value cannot be empty: %s", key), nil))
	}

	return defaultValue
}
