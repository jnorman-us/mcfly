package wirefmt

import "fmt"

type FlyError struct {
	ErrorString string `json:"error"`
	Status      string `json:"status"`
}

func (e FlyError) Error() string {
	return fmt.Sprintf("flyerr: %s (Status %s)", e.ErrorString, e.Status)
}

type ErrorBadRequest struct {
	FlyError
}

type ErrorTimedOut struct {
	FlyError
}
