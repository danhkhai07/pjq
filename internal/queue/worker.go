package queue

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"pjq/internal/domain"
)

var (
	ErrWorkerIsBusy error = errors.New("worker is busy")
)

type worker struct {
	id 			int
	job 		*domain.Job
	registry	*Registry
	busy 		atomic.Bool
}

func newWorker(id int, registry *Registry) *worker {
	return &worker{
		id: id,
		registry: registry,
	}
}

func (w *worker) IsBusy() bool {
	return w.busy.Load()
}

func (w *worker) Log(message string) {
	logTime := time.Now()
	if w.job != nil {
		w.job.Logs = append(w.job.Logs, logTime.Local().String() + " " + message)
		return
	} 
	fmt.Fprintf(os.Stderr, "error: cannot log with no job as worker id %d\n", w.id)
}

func (w *worker) Process(ctx context.Context, job *domain.Job) error {
	if w.IsBusy() {
		return fmt.Errorf("error: worker id %d is busy\n", w.id)
	}

	w.busy.Store(true)
	w.job = job
	defer w.busy.Store(false)
	defer func() { w.job = nil }()

	handler, err := w.registry.Get(w.job.Type)
	if err != nil {
		return err
	}

	err = handler.Handle(ctx, w.job, w.Log)
	if err != nil {
		w.job.Status = domain.StatusFailed
		w.job.Error = err.Error()
		return err
	}
	w.job.Status = domain.StatusDone
	return nil
}
