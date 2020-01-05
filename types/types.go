package types

type GrafanaDashboard struct {
	UID         string `json:"uid"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	FolderUID   string `json:"folderUid"`
	FolderTitle string `json:"folderTitle"`
}

type ListDashboardsResponse struct {
	Dashboards map[string][]*GrafanaDashboard `json:"dashboards"`
}

type RenderDashboardRequest struct {
	DashboardURL string `json:"dashboard_url"`
	Height       int    `json:"height"`
	Width        int    `json:"width"`
	StartTime    string `json:"start_time"` // RFC3339
	EndTime      string `json:"end_time"`   // RFC3339
}
