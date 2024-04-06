package wirefmt

import (
	"errors"
	"fmt"
)

type FlyError struct {
	ErrorString string `json:"error"`
	Status      string `json:"status"`
}

func (e FlyError) Error() string {
	return fmt.Sprintf("flyerr: %s (Status %s)", e.ErrorString, e.Status)
}

var ErrorBadRequest = errors.New("err 400")
var ErrorTimedOut = errors.New("err 408")
