package server

import (
	"github.com/monzo/typhon"
)

type livenessResponse struct {
	Status string `json:"status"`
}

func serveLiveness(req typhon.Request) typhon.Response {
	// Returns a plain 200 success response to show that
	// the server is still alive.
	return req.Response(&livenessResponse{Status: "ok"})
}
