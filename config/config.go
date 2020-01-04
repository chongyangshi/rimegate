package config

import (
	"fmt"
	"os"

	"github.com/monzo/terrors"
)

var (
	ConfigCORSAllowedOrigin        = getConfigFromOSEnv("CORS_ALLOWED_ORIGIN", "*", true)
	ConfigListenAddr               = getConfigFromOSEnv("LISTEN_ADDR", "", true)
	ConfigListenPort               = getConfigFromOSEnv("LISTEN_PORT", "8080", true)
	ConfigGrafanaAPIKey            = getConfigFromOSEnv("GRAFANA_API_KEY", "", false)
	ConfigGrafanaHost              = getConfigFromOSEnv("GRAFANA_HOST", "", false)
	ConfigImageCacheDirectory      = getConfigFromOSEnv("IMAGE_CACHE_DIRECTORY", "/cache", true)
	ConfigImageCacheValidity       = getConfigFromOSEnv("IMAGE_CACHE_VALIDITY", "1m", true)
	ConfigAuthenticationSigningKey = getConfigFromOSEnv("AUTHENTICATION_SIGNING_KEY", "", false)
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
