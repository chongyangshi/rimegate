package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/monzo/typhon"

	"github.com/icydoge/rimegate/config"
)

// N.B. The service itself does not enforce authentication. It is recommended that
// basic or certificate-based authentication is applied at the reverse proxy serving
// this service, if served over the internet.
func Service() typhon.Service {
	router := typhon.Router{}
	router.POST("/render", serveRenderDashboard)
	router.POST("/dashboards", serveListDashboards)
	router.GET("/grafana-credentials-required", serveGrafanaCredentialsRequired)
	router.GET("/healthz", serveLiveness)
	router.GET("/ping", serveLiveness)

	svc := router.Serve().Filter(typhon.ErrorFilter).Filter(typhon.H2cFilter).Filter(ClientErrorFilter).Filter(CORSFilter)

	return svc
}

// ClientErrorFilter strips sensitive error info before returning error to client, leaving
// only code and message; on a best-effort basis.
func ClientErrorFilter(req typhon.Request, svc typhon.Service) typhon.Response {
	rsp := svc(req)
	if rsp.Error != nil {
		var basicErr = basicError{}
		bodyBytes, err := rsp.BodyBytes(false)
		if err != nil {
			rsp.Body = ioutil.NopCloser(bytes.NewReader(basicErr.toFailbackBytes(rsp.StatusCode)))
			return rsp
		}

		err = json.Unmarshal(bodyBytes, &basicErr)
		if err != nil {
			rsp.Body = ioutil.NopCloser(bytes.NewReader(basicErr.toFailbackBytes(rsp.StatusCode)))
			return rsp
		}

		seralized, err := basicErr.toSerialized(rsp.StatusCode)
		if err != nil {
			rsp.Body = ioutil.NopCloser(bytes.NewReader(basicErr.toFailbackBytes(rsp.StatusCode)))
			return rsp
		}

		rsp.Body = ioutil.NopCloser(bytes.NewReader(seralized))
		return rsp
	}

	return rsp
}

func CORSFilter(req typhon.Request, svc typhon.Service) typhon.Response {
	var rsp typhon.Response
	if req.Method == http.MethodOptions {
		rsp = typhon.NewResponse(req)
		rsp.Body = ioutil.NopCloser(bytes.NewReader([]byte("ok")))
		rsp.StatusCode = http.StatusOK
	} else {
		rsp = svc(req)
	}

	rsp.Header.Set("Access-Control-Allow-Origin", config.ConfigCORSAllowedOrigin)
	rsp.Header.Set("Access-Control-Allow-Methods", "GET, PUT, POST")
	rsp.Header.Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Accept-Langauge, Content-Language")

	return rsp
}

type basicError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// For 4xx errors, message can be disclosed to the client; while for 5xx errors this shouldn't be disclosed.
func (b basicError) toFailbackBytes(statusCode int) []byte {
	switch {
	case statusCode >= 400 && statusCode < 500:
		return []byte(fmt.Sprintf("Error (%s): %s", b.Code, b.Message))
	default:
		return []byte(fmt.Sprintf("Error (%s): %d", b.Code, statusCode))
	}
}

func (b basicError) toSerialized(statusCode int) ([]byte, error) {
	exportedError := b
	switch {
	case statusCode >= 400 && statusCode < 500:
		// Keep the error message
	default:
		exportedError.Message = strconv.FormatInt(int64(statusCode), 10)
	}

	seralized, err := json.Marshal(exportedError)
	return seralized, err
}
