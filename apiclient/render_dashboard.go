package apiclient

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/monzo/slog"
	"github.com/monzo/typhon"

	"github.com/icydoge/rimegate/config"
)

func RenderDashboards(ctx context.Context, dashboardURL string, startTime, endTime time.Time, height, width, orgID int) ([]byte, *time.Time, error) {
	// The formats of dashboard URL and other input params should have already been validated by the caller.
	// In case dashboard URL does not start with a forward slash, append it.
	if !strings.HasPrefix(dashboardURL, "/") {
		dashboardURL = fmt.Sprintf("/%s", dashboardURL)
	}

	errParams := map[string]string{
		"start_time":    startTime.Format(time.RFC3339),
		"end_time":      endTime.Format(time.RFC3339),
		"height":        strconv.FormatInt(int64(height), 10),
		"width":         strconv.FormatInt(int64(width), 10),
		"org_id":        strconv.FormatInt(int64(orgID), 10),
		"dashboard_url": dashboardURL,
	}

	// Try cache first, if hit, return
	ok, entry := getCachedRender(dashboardURL, height, width)
	if ok {
		return entry.payload, &entry.timeRendered, nil
	}

	// Render in kiosk mode to hide the header panel, which is of no use to the client
	requestURL := fmt.Sprintf(
		"%s/render/%s?orgId=%d&from=%d&to=%d&width=%d&height=%d&kiosk",
		config.ConfigGrafanaHost,
		dashboardURL,
		orgID,
		startTime.Unix(),
		endTime.Unix(),
		width,
		height,
	)

	req := typhon.NewRequest(ctx, http.MethodGet, requestURL, nil)
	req.Header.Set("Authorization", fmt.Sprint("Bearer: %s", config.ConfigGrafanaAPIKey))

	rsp := req.Send().Response()
	if rsp.Error != nil {
		slog.Error(ctx, "Grafana returned error: %v", rsp.Error, errParams)
		return nil, nil, rsp.Error
	}

	rspBytes, err := rsp.BodyBytes(true)
	if err != nil {
		slog.Error(ctx, "Error reading response: %v", err, errParams)
		return nil, nil, err
	}

	// Update the render cache
	entry = cacheRender(dashboardURL, height, width, rspBytes)
	return entry.payload, &entry.timeRendered, nil
}
