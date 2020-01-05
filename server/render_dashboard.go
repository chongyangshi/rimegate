package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/icydoge/rimegate/config"

	"github.com/monzo/slog"
	"github.com/monzo/terrors"
	"github.com/monzo/typhon"

	"github.com/icydoge/rimegate/apiclient"
	"github.com/icydoge/rimegate/types"
)

var dashboardURLRegexp = regexp.MustCompile(`[a-zA-Z0-9\-\/]+`)
var defaultPeriod time.Duration

func init() {
	var err error
	defaultPeriod, err = time.ParseDuration(config.ConfigGrafanaDefaultPeriod)
	if err != nil {
		defaultPeriod = time.Hour
	}
}

func serveRenderDashboard(req typhon.Request) typhon.Response {
	requestBytes, err := req.BodyBytes(false)
	if err != nil {
		slog.Error(req, "Error reading request bytes: %v", err)
	}

	request := types.RenderDashboardRequest{}
	err = json.Unmarshal(requestBytes, &request)
	if err != nil {
		slog.Error(req, "Error unmarshaling request: %v", err)
	}

	errParams := map[string]string{
		"start_time":    request.StartTime,
		"end_time":      request.EndTime,
		"height":        strconv.FormatInt(int64(request.Height), 10),
		"width":         strconv.FormatInt(int64(request.Width), 10),
		"dashboard_url": request.DashboardURL,
	}

	// Process start and end times specified, if unspecified or invalid, we use default.
	var startTime, endTime time.Time
	if request.StartTime == "" && request.EndTime == "" {
		endTime = time.Now()
		startTime = endTime.Add(-1 * defaultPeriod)
	}

	parsedStartTime, startTimeErr := time.Parse(time.RFC3339, request.StartTime)
	parsedEndTime, endTimeErr := time.Parse(time.RFC3339, request.EndTime)
	if startTimeErr != nil || endTimeErr != nil {
		endTime = time.Now()
		startTime = endTime.Add(-1 * defaultPeriod)
	}

	if parsedStartTime.After(parsedEndTime) {
		endTime = time.Now()
		startTime = endTime.Add(-1 * defaultPeriod)
	}

	startTime = parsedStartTime
	endTime = parsedEndTime

	errParams["start_time"] = startTime.Format(time.RFC3339)
	errParams["end_time"] = endTime.Format(time.RFC3339)

	switch {
	case !validateDashboardURL(request.DashboardURL):
		return typhon.Response{Error: terrors.BadRequest("invalid_url", fmt.Sprintf("Requested dashboard URL %s is invalid", request.DashboardURL), errParams)}
	case request.Height < 400:
		return typhon.Response{Error: terrors.BadRequest("bad_height", fmt.Sprintf("Requested height (%d) is below 400, which is unlikely to produce useful render", request.Height), errParams)}
	case request.Width < 400:
		return typhon.Response{Error: terrors.BadRequest("bad_width", fmt.Sprintf("Requested width (%d) is below 400, which is unlikely to produce useful render", request.Width), errParams)}
	}

	// Input validation done by apiclient
	render, timeRendered, err := apiclient.RenderDashboards(req, request.DashboardURL, startTime, endTime, request.Height, request.Width)
	if err != nil {
		slog.Error(req, "Error rendering dashboard: %v", err, errParams)
		return typhon.Response{Error: terrors.InternalService("", "Error rendering dashboard", errParams)}
	}

	rsp := typhon.NewResponse(req)
	rsp.StatusCode = http.StatusOK
	rsp.Header.Set("X-Time-Rendered", timeRendered.Format(time.RFC3339))
	rsp.Body = ioutil.NopCloser(bytes.NewReader(render))
	return rsp
}

func validateDashboardURL(dashboardURL string) bool {
	// Dashboard URLs in Grafana V5+ should look like /d/VNegG8BWz/multi-cluster-network-encapsulation-wylis
	return dashboardURLRegexp.MatchString(dashboardURL)
}
