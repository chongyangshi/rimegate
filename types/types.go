package types

type Organization struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type GrafanaDashboard struct {
	UID         string `json:"uid"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	FolderUID   string `json:"folderUid"`
	FolderTitle string `json:"folderTitle"`
}

// Auth is the passthrough Grafana username and password, we don't maintain our own auth
// as it is unnecessary when Grafana is already supporting basic auth and LDAP.
// Empty usernames and passwords are accepted, as some Grafana dashboards could be
// unauthenticated.
type Auth struct {
	GrafanaUsername string `json:"grafana_username"`
	GrafanaPassword string `json:"grafana_password"`
}

type ListDashboardsRequest struct {
	*Auth
}

type ListDashboardsResponse struct {
	Dashboards map[string][]GrafanaDashboard `json:"dashboards"`
}

type RenderDashboardRequest struct {
	*Auth
	OrgID        int    `json:"org_id"`
	DashboardURL string `json:"dashboard_url"`
	Height       int    `json:"height"`
	Width        int    `json:"width"`
	StartTime    string `json:"start_time"` // RFC3339
	EndTime      string `json:"end_time"`   // RFC3339
	AutoFitPanel bool   `json:"auto_fit_panel"`
}

type RenderDashboardResponse struct {
	Payload      string `json:"payload"`        // Base64
	RenderedTime string `json:"rendered_time"`  // RFC3339
	UTCWallClock string `json:"utc_wall_clock"` // A backend-supplied UTC wall clock for clients unable to parse RFC3339 in-browser.
}
