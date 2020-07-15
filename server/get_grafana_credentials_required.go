package server

import (
	"github.com/monzo/typhon"

	"github.com/chongyangshi/rimegate/config"
	"github.com/chongyangshi/rimegate/types"
)

// This endpoint tells the client whether it needs to collect Grafana username/password
// from the end user, based on whether the server operates with mTLS and therefore relies
// on a static Grafana API token instead.
func serveGrafanaCredentialsRequired(req typhon.Request) typhon.Response {
	rsp := types.GrafanaCredentialsRequiredResponse{
		Required: config.ConfigGrafanaAPIToken == "",
	}

	return req.Response(&rsp)
}
