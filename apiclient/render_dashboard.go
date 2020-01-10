package apiclient

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/monzo/slog"
	"github.com/monzo/terrors"
	"github.com/monzo/typhon"

	"github.com/icydoge/rimegate/config"
	"github.com/icydoge/rimegate/types"
)

func RenderDashboards(ctx context.Context, auth *types.Auth, dashboardURL string, startTime, endTime time.Time, height, width, orgID int, fitPanel bool) ([]byte, *time.Time, error) {
	// The formats of dashboard URL and other input params should have already been validated by the caller.
	// In case dashboard URL starts with a forward slash, strip it.
	if strings.HasPrefix(dashboardURL, "/") {
		dashboardURL = strings.TrimPrefix(dashboardURL, "/")
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
		slog.Debug(ctx, "Rendered %s from cache", dashboardURL, errParams)
		return entry.payload, &entry.timeRendered, nil
	}

	// Render in kiosk mode to hide the header panel, which is of no use to the client
	requestURL := fmt.Sprintf(
		"%s/render/%s?orgId=%d&from=%d&to=%d&width=%d&height=%d&kiosk",
		config.ConfigGrafanaHost,
		dashboardURL,
		orgID,
		unixTimeWithMilliseconds(startTime),
		unixTimeWithMilliseconds(endTime),
		width,
		height,
	)

	req := typhon.NewRequest(ctx, http.MethodGet, requestURL, nil)
	req.SetBasicAuth(auth.GrafanaUsername, auth.GrafanaPassword)

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

	if rsp.StatusCode >= 400 {
		err := terrors.WrapWithCode(fmt.Errorf("%s", string(rspBytes[:])), errParams, fmt.Sprintf("grafana_%d", rsp.StatusCode))
		slog.Error(ctx, "Grafana returned error: %s", err, errParams)
		return nil, nil, err
	}

	slog.Debug(ctx, "Rendered %s without cache", requestURL, errParams)

	// Update the render cache
	entry = cacheRender(dashboardURL, height, width, rspBytes)
	return entry.payload, &entry.timeRendered, nil
}

// Annoying, Grafana timestamp query params are neither Unix() or UnixNano(),
// but represents milliseconds as three trailling digits in the integer.
func unixTimeWithMilliseconds(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}
