package server

import (
	"github.com/monzo/typhon"
)

type livenessResponse struct{}

func serveLiveness(req typhon.Request) typhon.Response {
	// Returns a plain 200 success response to show that
	// the server is still alive.
	return req.Response(livenessResponse{})
}
