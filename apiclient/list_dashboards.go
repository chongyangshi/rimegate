package apiclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/monzo/slog"
	"github.com/monzo/typhon"

	"github.com/icydoge/rimegate/config"
	"github.com/icydoge/rimegate/types"
)

func ListDashboards(ctx context.Context) (map[string][]types.GrafanaDashboard, error) {
	// Perform an empty search to list all dashboards the API token has read access to
	requestURL := fmt.Sprintf("%s/api/%s", config.ConfigGrafanaHost, "search?query=")
	errParams := map[string]string{
		"request_url": requestURL,
	}

	req := typhon.NewRequest(ctx, http.MethodGet, requestURL, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.ConfigGrafanaAPIKey))

	rsp := req.Send().Response()
	if rsp.Error != nil {
		slog.Error(ctx, "Grafana returned error: %v", rsp.Error, errParams)
		return nil, rsp.Error
	}

	rspBytes, err := rsp.BodyBytes(true)
	if err != nil {
		slog.Error(ctx, "Error reading response: %v", err, errParams)
		return nil, err
	}

	dashboards := []types.GrafanaDashboard{}
	err = json.Unmarshal(rspBytes, &dashboards)
	if err != nil {
		slog.Error(ctx, "Error unmarshaling response: %v", err, errParams)
		return nil, err
	}

	categorisedDashboards := map[string][]types.GrafanaDashboard{}
	for _, dashboard := range dashboards {
		folderSeen := false
		for folder := range categorisedDashboards {
			if folder == dashboard.FolderUID {
				folderSeen = true
				break
			}
		}

		if !folderSeen {
			categorisedDashboards[dashboard.FolderUID] = []types.GrafanaDashboard{}
		}

		categorisedDashboards[dashboard.FolderUID] = append(categorisedDashboards[dashboard.FolderUID], dashboard)
	}

	return categorisedDashboards, nil
}
