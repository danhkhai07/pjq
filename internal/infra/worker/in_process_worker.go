package infra

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"pjq/internal/domain"
	"pjq/internal/util"
)

var (
	ErrWorkerIsBusy error = errors.New("worker is busy")
)

type Worker struct {
	id 			int
	job 		*domain.Job
	registry	*util.Registry
	busy 		atomic.Bool
}

func (w *Worker) IsBusy() bool {
	return w.busy.Load()
}

func (w *Worker) Log(message string) {
	logTime := time.Now()
	if w.job != nil {
		w.job.Logs = append(w.job.Logs, logTime.Local().String() + " " + message)
		return
	} 
	fmt.Fprintf(os.Stderr, "error: cannot log with no job as worker id %d\n", w.id)
}

func (w *Worker) Process(ctx context.Context, job *domain.Job) error {
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

