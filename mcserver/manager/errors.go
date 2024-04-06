package manager

import "errors"

var ErrorServerNotRegistered = errors.New("server does not exist")
var ErrorNotAuthorized = errors.New("user not authorized")
var ErrorCloud = errors.New("problem with cloud")
