package server

import (
	"encoding/json"

	"github.com/monzo/slog"
	"github.com/monzo/terrors"
	"github.com/monzo/typhon"

	"github.com/chongyangshi/rimegate/apiclient"
	"github.com/chongyangshi/rimegate/types"
)

func serveListDashboards(req typhon.Request) typhon.Response {
	requestBytes, err := req.BodyBytes(false)
	if err != nil {
		slog.Error(req, "Error reading request bytes: %v", err)
		return typhon.Response{Error: err}
	}

	request := types.ListDashboardsRequest{}
	err = json.Unmarshal(requestBytes, &request)
	if err != nil {
		slog.Error(req, "Error unmarshaling request: %v", err)
		return typhon.Response{Error: err}
	}

	dashboards, err := apiclient.ListDashboards(req, request.Auth)
	if err != nil {
		// Proxy Unauthorized responses if credentials supplied are invalid.
		if terrors.PrefixMatches(err, "grafana_401") {
			return typhon.Response{Error: terrors.Unauthorized("", "Grafana username or password incorrect", nil)}
		}

		slog.Error(req, "Error listing dashboards: %v", err)
		return typhon.Response{Error: terrors.InternalService("", "Error listing dashboards", nil)}
	}

	return req.Response(&types.ListDashboardsResponse{
		Dashboards: dashboards,
	})
}
