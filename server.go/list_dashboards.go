package server

import "github.com/monzo/typhon"

type grafanaDashboardSearchResponse struct {
	UID         string `json:"uid"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	FolderUID   string `json:"folderUid"`
	FolderTitle string `json:"folderTitle"`
}

func serveListDashboards(req typhon.Request) typhon.Response {
	if err := checkAuthorization(req); err != nil {
		return typhon.Response{Error: err}
	}

	// Returns a plain 200 success response to show that
	// the server is still alive.
	return req.Response(healthCheckResponse{})
}
