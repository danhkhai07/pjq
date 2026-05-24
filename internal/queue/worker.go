package queue

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync/atomic"

	"pjq/internal/domain"
)

var (
	ErrWorkerIsBusy error = errors.New("worker is busy")
)

type Worker struct {
	ID 			int
	Job 		*domain.Job
	Registry	*Registry
	busy 		atomic.Bool
}

func NewWorker(id int, registry *Registry) Worker {
	return Worker{
		ID: id,
		Registry: registry,
	}
}

func (w *Worker) IsBusy() bool {
	return w.busy.Load()
}

func (w *Worker) Log(message string) {
	if w.Job != nil {
		w.Job.Logs = append(w.Job.Logs, message)
		return
	} 
	fmt.Fprintf(os.Stderr, "error: cannot log with no job as worker id %d\n", w.ID)
}

func (w *Worker) Process(ctx context.Context, job *domain.Job) error {
	if w.IsBusy() {
		return fmt.Errorf("error: worker id %d is busy\n", w.ID)
	}

	w.busy.Store(true)
	w.Job = job
	defer w.busy.Store(false)
	defer func() { w.Job = nil }()

	handler, err := w.Registry.Get(w.Job.Type)
	if err != nil {
		return err
	}

	err = handler.Handle(ctx, w.Job, w.Log)
	if err != nil {
		w.Job.Status = domain.Failed
		w.Job.Error = err.Error()
		return err
	}
	w.Job.Status = domain.Done
	return nil
}
