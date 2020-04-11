package apiclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/monzo/slog"
	"github.com/monzo/terrors"
	"github.com/monzo/typhon"

	"github.com/icydoge/rimegate/config"
	"github.com/icydoge/rimegate/types"
)

func GetCurrentOrganization(ctx context.Context, auth *types.Auth) (*types.Organization, error) {
	requestURL := fmt.Sprintf("%s/api/org", config.ConfigGrafanaHost)
	errParams := map[string]string{
		"request_url": requestURL,
	}

	req := typhon.NewRequest(ctx, http.MethodGet, requestURL, nil)
	setAuthenticationCredentials(auth.GrafanaUsername, auth.GrafanaPassword)

	rsp := req.Send().Response()
	if rsp.Error != nil {
		slog.Error(ctx, "Grafana returned error: %v", rsp.Error, errParams)
		return nil, rsp.Error
	}

	rspBytes, err := rsp.BodyBytes(true)
	if err != nil {
		slog.Error(ctx, "Error reading response: %v", err, errParams)
		return nil, err
	}

	if rsp.StatusCode >= 400 {
		err := terrors.WrapWithCode(fmt.Errorf("%s", string(rspBytes[:])), errParams, fmt.Sprintf("grafana_%d", rsp.StatusCode))
		slog.Error(ctx, "Grafana returned error: %s", err, errParams)
		return nil, err
	}

	org := &types.Organization{}
	err = json.Unmarshal(rspBytes, org)
	if err != nil {
		slog.Error(ctx, "Error unmarshaling response: %v", err, errParams)
		return nil, err
	}

	return org, nil
}
