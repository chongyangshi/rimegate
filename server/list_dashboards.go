package server

import (
	"encoding/json"

	"github.com/monzo/slog"
	"github.com/monzo/terrors"
	"github.com/monzo/typhon"

	"github.com/icydoge/rimegate/apiclient"
	"github.com/icydoge/rimegate/types"
)

func serveListDashboards(req typhon.Request) typhon.Response {
	requestBytes, err := req.BodyBytes(false)
	if err != nil {
		slog.Error(req, "Error reading request bytes: %v", err)
	}

	request := types.ListDashboardsRequest{}
	err = json.Unmarshal(requestBytes, &request)
	if err != nil {
		slog.Error(req, "Error unmarshaling request: %v", err)
	}

	dashboards, err := apiclient.ListDashboards(req, request.Auth)
	if err != nil {
		slog.Error(req, "Error listing dashboards: %v", err)
		return typhon.Response{Error: terrors.InternalService("", "Error listing dashboards", nil)}
	}

	return req.Response(&types.ListDashboardsResponse{
		Dashboards: dashboards,
	})
}
