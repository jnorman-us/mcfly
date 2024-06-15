package halter

import "errors"

var ErrSchedulingTask = errors.New("problem scheduling halt task")
var ErrCloud = errors.New("problem with cloud")
