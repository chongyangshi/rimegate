package types

type GrafanaDashboard struct {
	UID         string `json:"uid"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	FolderUID   string `json:"folderUid"`
	FolderTitle string `json:"folderTitle"`
}

type ListDashboardsResponse struct {
	Dashboards map[string][]GrafanaDashboard `json:"dashboards"`
}
