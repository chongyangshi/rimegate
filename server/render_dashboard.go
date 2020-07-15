package server

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/monzo/slog"
	"github.com/monzo/terrors"
	"github.com/monzo/typhon"

	"github.com/chongyangshi/rimegate/apiclient"
	"github.com/chongyangshi/rimegate/config"
	"github.com/chongyangshi/rimegate/types"
)

var dashboardURLRegexp = regexp.MustCompile(`[a-zA-Z0-9\-\/]+`)
var defaultPeriod, maxPeriod time.Duration

func init() {
	var err error
	defaultPeriod, err = time.ParseDuration(config.ConfigGrafanaDefaultPeriod)
	if err != nil {
		defaultPeriod = time.Hour
	}
	maxPeriod, err = time.ParseDuration(config.ConfigGrafanaMaxPeriod)
	if err != nil {
		maxPeriod = time.Hour * 3
	}
}

func serveRenderDashboard(req typhon.Request) typhon.Response {
	requestBytes, err := req.BodyBytes(false)
	if err != nil {
		slog.Error(req, "Error reading request bytes: %v", err)
		return typhon.Response{Error: err}
	}

	request := types.RenderDashboardRequest{}
	err = json.Unmarshal(requestBytes, &request)
	if err != nil {
		slog.Error(req, "Error unmarshaling request: %v", err)
		return typhon.Response{Error: err}
	}

	errParams := map[string]string{
		"org_id":        strconv.FormatInt(int64(request.OrgID), 10),
		"start_time":    request.StartTime,
		"end_time":      request.EndTime,
		"height":        strconv.FormatInt(int64(request.Height), 10),
		"width":         strconv.FormatInt(int64(request.Width), 10),
		"dashboard_url": request.DashboardURL,
	}

	// Process start and end times specified, if unspecified or invalid, we use default.
	startTime, endTime := validateStartEndTimes(req, request.StartTime, request.EndTime)
	errParams["start_time"] = startTime.Format(time.RFC3339)
	errParams["end_time"] = endTime.Format(time.RFC3339)

	// If org ID not specified (or 0), find the current org ID.
	orgID := request.OrgID
	if orgID < 1 {
		org, err := apiclient.GetCurrentOrganization(req, request.Auth)
		if err != nil {
			slog.Error(req, "Cannot find current Grafana organization: %v", err, errParams)
			return typhon.Response{Error: terrors.PreconditionFailed("cannot_find_org", "Cannot find current Grafana organization", errParams)}
		}
		orgID = org.ID
	}

	switch {
	case !validateDashboardURL(request.DashboardURL):
		return typhon.Response{Error: terrors.BadRequest("invalid_url", fmt.Sprintf("Requested dashboard URL %s is invalid", request.DashboardURL), errParams)}
	case request.Height < 400:
		return typhon.Response{Error: terrors.BadRequest("bad_height", fmt.Sprintf("Requested height (%d) is below 400, which is unlikely to produce useful render", request.Height), errParams)}
	case request.Width < 400:
		return typhon.Response{Error: terrors.BadRequest("bad_width", fmt.Sprintf("Requested width (%d) is below 400, which is unlikely to produce useful render", request.Width), errParams)}
	}

	render, timeRendered, err := apiclient.RenderDashboards(req, request.Auth, request.DashboardURL, startTime, endTime, request.Height, request.Width, orgID, request.AutoFitPanel)
	if err != nil {
		// Proxy Unauthorized responses if credentials supplied are invalid.
		if terrors.PrefixMatches(err, "grafana_401") {
			return typhon.Response{Error: terrors.Unauthorized("", "Grafana username or password incorrect", nil)}
		}

		slog.Error(req, "Error rendering dashboard: %v", err, errParams)
		return typhon.Response{Error: terrors.InternalService("", "Error rendering dashboard", errParams)}
	}

	h, m, s := timeRendered.Clock()

	return req.Response(&types.RenderDashboardResponse{
		Payload:      base64.StdEncoding.EncodeToString(render),
		RenderedTime: timeRendered.Format(time.RFC3339),
		UTCWallClock: fmt.Sprintf("%02d:%02d:%02d", h, m, s),
	})
}

func validateDashboardURL(dashboardURL string) bool {
	// Dashboard URLs in Grafana V5+ should look like /d/VNegG8BWz/multi-cluster-network-encapsulation-wylis
	return dashboardURLRegexp.MatchString(dashboardURL)
}

func validateStartEndTimes(ctx context.Context, startTimeStr, endTimeStr string) (time.Time, time.Time) {
	defaultEndTime := time.Now()
	defaultStartTime := defaultEndTime.Add(-1 * defaultPeriod)

	if startTimeStr == "" || endTimeStr == "" {
		// Default period
		return defaultStartTime, defaultEndTime
	}

	parsedStartTime, startTimeErr := time.Parse(time.RFC3339, startTimeStr)
	parsedEndTime, endTimeErr := time.Parse(time.RFC3339, endTimeStr)
	if startTimeErr != nil || endTimeErr != nil {
		slog.Warn(ctx, "Error parsing start and/or end times (%v, %v), using default.", startTimeErr, endTimeErr)
		return defaultStartTime, defaultEndTime
	}

	if parsedStartTime.After(parsedEndTime) {
		slog.Warn(ctx, "Start time %v is after end time %v, using default start time interval.", parsedStartTime, parsedEndTime)
		startTime := parsedEndTime.Add(-1 * defaultPeriod)
		return startTime, parsedEndTime
	}

	if parsedEndTime.Sub(parsedStartTime) > maxPeriod {
		slog.Warn(ctx, "Start time %v is too far from end time %v, using max start time interval.", parsedStartTime, parsedEndTime)
		startTime := parsedEndTime.Add(-1 * maxPeriod)
		return startTime, parsedEndTime
	}

	return parsedStartTime, parsedEndTime
}
