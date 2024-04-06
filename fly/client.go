package fly

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/jnorman-us/mcfly/env"
	"github.com/jnorman-us/mcfly/fly/wirefmt"
	"gopkg.in/resty.v1"
)

type FlyClient struct {
	client *resty.Client
}

func NewFlyClient(cfg env.Config) *FlyClient {
	client := resty.New()
	client.SetAuthToken(cfg.FlyToken)
	client.SetHeader("Content-Type", "application/json")
	client.SetHostURL(fmt.Sprintf("https://api.machines.dev/v1/apps/%s/", cfg.FlyApp))

	return &FlyClient{
		client,
	}
}

var ErrorStrange = errors.New("response does not conform to spec")

func handleError(statusCode int, respErr any) error {
	if flyErr, ok := respErr.(*wirefmt.FlyError); ok {
		var err = *flyErr
		if statusCode == http.StatusBadRequest {
			return wirefmt.ErrorBadRequest{FlyError: err}
		}
		if statusCode == http.StatusRequestTimeout {
			return wirefmt.ErrorTimedOut{FlyError: err}
		}
	}
	return ErrorStrange
}
