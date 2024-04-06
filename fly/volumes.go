package fly

import (
	"context"
	"net/http"

	"github.com/jnorman-us/mcfly/fly/wirefmt"
)

func (c *FlyClient) ListVolumes(ctx context.Context) (wirefmt.ListVolumesOutput, error) {
	resp, err := c.client.R().
		SetContext(ctx).
		SetResult(&wirefmt.ListVolumesOutput{}).
		SetError(&wirefmt.FlyError{}).
		Get("/volumes")

	if err != nil {
		return []wirefmt.Volume{}, err
	}

	if resp.StatusCode() == http.StatusOK {
		if output, ok := resp.Result().(*wirefmt.ListVolumesOutput); ok {
			return *output, nil
		}
	}
	return []wirefmt.Volume{}, handleError(resp.StatusCode(), resp.Error())
}

func (c *FlyClient) CreateVolume(ctx context.Context, input wirefmt.CreateVolumeInput) (*wirefmt.CreateVolumeOutput, error) {
	input.Region = c.region

	resp, err := c.client.R().
		SetContext(ctx).
		SetBody(input).
		SetResult(&wirefmt.CreateVolumeOutput{}).
		SetError(&wirefmt.FlyError{}).
		Post("/volumes")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() == http.StatusCreated {
		if output, ok := resp.Result().(*wirefmt.CreateVolumeOutput); ok {
			return output, nil
		}
	}

	return nil, handleError(resp.StatusCode(), resp.Error())
}
