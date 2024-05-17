package fly

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jnorman-us/mcfly/fly/wirefmt"
)

func (c *FlyClient) ListMachines(ctx context.Context) (wirefmt.ListMachinesOutput, error) {
	resp, err := c.client.R().
		SetContext(ctx).
		SetResult(&wirefmt.ListMachinesOutput{}).
		SetError(&wirefmt.FlyError{}).
		Get("/machines")

	if err != nil {
		return []wirefmt.Machine{}, err
	}

	if resp.StatusCode() == http.StatusOK {
		if output, ok := resp.Result().(*wirefmt.ListMachinesOutput); ok {
			return *output, nil
		}
	}
	return []wirefmt.Machine{}, handleError(resp.StatusCode(), resp.Error())
}

func (c *FlyClient) GetMachine(ctx context.Context, id string) (*wirefmt.GetMachineOutput, error) {
	resp, err := c.client.R().
		SetContext(ctx).
		SetResult(&wirefmt.GetMachineOutput{}).
		SetError(&wirefmt.FlyError{}).
		Get(fmt.Sprintf("/machines/%s", id))

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() == http.StatusOK {
		if output, ok := resp.Result().(*wirefmt.GetMachineOutput); ok {
			return output, nil
		}
	}
	return nil, handleError(resp.StatusCode(), resp.Error())
}

func (c *FlyClient) CreateMachine(ctx context.Context, input wirefmt.CreateMachineInput) (*wirefmt.CreateMachineOutput, error) {
	resp, err := c.client.R().
		SetContext(ctx).
		SetBody(input).
		SetResult(&wirefmt.CreateMachineOutput{}).
		SetError(&wirefmt.FlyError{}).
		Post("/machines")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() == http.StatusOK {
		if output, ok := resp.Result().(*wirefmt.CreateMachineOutput); ok {
			return output, nil
		}
	}

	return nil, handleError(resp.StatusCode(), resp.Error())
}

func (c *FlyClient) StartMachine(ctx context.Context, id string) error {
	resp, err := c.client.R().
		SetContext(ctx).
		Post(fmt.Sprintf("/machines/%s/start", id))

	if err != nil {
		return err
	}

	if resp.StatusCode() == http.StatusOK {
		return nil
	}

	return handleError(resp.StatusCode(), nil)
}

func (c *FlyClient) StopMachine(ctx context.Context, id string) error {
	resp, err := c.client.R().
		SetContext(ctx).
		Post(fmt.Sprintf("/machines/%s/stop", id))

	if err != nil {
		return err
	}

	if resp.StatusCode() == http.StatusOK {
		return nil
	}

	return handleError(resp.StatusCode(), nil)
}
