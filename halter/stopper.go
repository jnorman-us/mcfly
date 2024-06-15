package halter

import "context"

type Stopper interface {
	StopMachine(context.Context, string) error
}
