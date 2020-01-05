package server

import (
	"github.com/monzo/slog"
	"github.com/monzo/terrors"
	"github.com/monzo/typhon"

	"github.com/icydoge/rimegate/apiclient"
	"github.com/icydoge/rimegate/types"
)

func serveListDashboards(req typhon.Request) typhon.Response {
	dashboards, err := apiclient.ListDashboards(req)
	if err != nil {
		slog.Error(req, "Error listing dashboards: %v", err)
		return typhon.Response{Error: terrors.InternalService("", "Error listing dashboards", nil)}
	}

	return req.Response(&types.ListDashboardsResponse{
		Dashboards: dashboards,
	})
}
