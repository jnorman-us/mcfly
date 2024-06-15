package halter

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"github.com/madflojo/tasks"
)

const HaltWaitDuration = time.Minute * 2

type HalterQueue interface {
	Queue(machineID string) error
	Dequeue(machineID string)
}

type halterQueue struct {
	scheduler *tasks.Scheduler
	log       logr.Logger
	stopper   Stopper
}

func New(stopper Stopper, log logr.Logger) HalterQueue {
	return &halterQueue{
		scheduler: tasks.New(),
		log:       log,
		stopper:   stopper,
	}
}

func (h *halterQueue) Queue(machineID string) error {
	ctx := context.Background()
	log := h.log.WithValues("machineID", machineID)
	ctx = logr.NewContext(ctx, log)

	err := h.scheduler.AddWithID(machineID, &tasks.Task{
		RunOnce:     true,
		Interval:    HaltWaitDuration,
		TaskContext: tasks.TaskContext{Context: ctx},
		FuncWithTaskContext: func(tc tasks.TaskContext) error {
			return h.halt(tc.Context, machineID)
		},
		ErrFuncWithTaskContext: func(tc tasks.TaskContext, err error) {
			logr.FromContextOrDiscard(tc.Context).Error(err, "Problem stopping machine")
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (h *halterQueue) Dequeue(machineID string) {
	h.scheduler.Del(machineID)
}

func (h *halterQueue) halt(ctx context.Context, machineID string) error {
	log := logr.FromContextOrDiscard(ctx)
	stopper := h.stopper

	log.Info("Stopping machine")

	err := stopper.StopMachine(ctx, machineID)
	if err != nil {
		return err
	}
	return nil
}
