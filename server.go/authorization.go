package server

import (
	"strings"

	"github.com/monzo/terrors"

	"github.com/monzo/typhon"

	"github.com/icydoge/rimegate/auth"
)

func checkAuthorization(req typhon.Request) error {
	token := req.Header.Get("Authorization")
	if token == "" || !strings.HasPrefix(token, "Bearer ") {
		return terrors.Unauthorized("missing_token", "401 Unauthorized (no token is presented in request)", nil)
	}

	validated, err := auth.VerifyToken(strings.TrimSpace(strings.TrimPrefix(token, "Bearer ")))
	if err != nil {
		return terrors.Unauthorized("bad_token", "401 Unauthorized (invalid token is presented in request)", nil)
	}

	if !validated {
		return terrors.Forbidden("bad_token", "403 Forbidden", nil)
	}

	return nil
}
