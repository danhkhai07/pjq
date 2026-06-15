package infra

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync/atomic"
	"time"
	"math"

	"pjq/internal/domain"
	"pjq/internal/util"
	"pjq/internal/queue"
)

var (
	workerCounter int = 0

	ErrWorkerIsBusy error = errors.New("worker is busy")
)

const (
	BASE_RETRY_BACKOFF = 5 * time.Second
)

type Worker struct {
	id 				int
	job				*domain.Job
	registry		*util.Registry
	busy 			atomic.Bool
	store			domain.JobStore
	consumeJob		func () (*domain.Job, bool)
	pushJob			func (*domain.Job)
}

func NewInProcessWorker(
	registry 	*util.Registry,
	store		domain.JobStore,
	queue		*queue.QueueManager,
) *Worker {
	workerCounter++
	return &Worker{
		id: workerCounter,
		registry: registry,
		store: store,
		consumeJob: queue.PopJob,
		pushJob: queue.PushJob,
	}
}

func (w *Worker) GetID() int { return w.id }

func (w *Worker) IsBusy() bool { return w.busy.Load() }

func (w *Worker) process(ctx context.Context, job *domain.Job) error {
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

	err = handler.Handle(ctx, w.job, w.log)
	if err != nil {
		w.job.Status = domain.StatusFailed
		w.job.Error = err.Error()
		return err
	}
	w.job.Status = domain.StatusDone
	return nil
}

func (w *Worker) RunWorker(ctx context.Context) (err error) {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			job, ok := w.consumeJob()
			if !ok {
				continue
			}
			w.job = job
			now := time.Now()
			if job.StartedAt == nil {
				job.StartedAt = &now
			}
			w.changeStatus(domain.StatusRunning)
			err = w.process(ctx, job)
			job.FinishedAt = &now
			if err != nil {
				w.changeStatus(domain.StatusFailed)
				w.logError(err)
				if job.Retries < job.MaxRetries {
					w.changeStatus(domain.StatusRetrying)
					w.retry()
				}
			} else {
				w.changeStatus(domain.StatusDone)
			}
			w.store.Save(ctx, *job)
			time.Sleep(10 * time.Millisecond)
		}
		w.job = nil
	}
}

func (w *Worker) log(message string) {
	logTime := time.Now()
	if w.job != nil {
		w.job.Logs = append(w.job.Logs, logTime.Local().String() + " " + message)
		return
	} 
	fmt.Fprintf(os.Stderr, "error: cannot log with no job as worker id %d\n", w.id)
}

func (w *Worker) logError(err error) {
	w.job.Error = err.Error()
	w.log(err.Error())
}

func (w *Worker) changeStatus(status domain.Status) {
	w.job.Status = status
	w.log(fmt.Sprintf("Changed job status to '%s'.", status))
}

func (w *Worker) retry() {
	now := time.Now()
	w.job.Retries++
	// exponential back-off with power of 2: 5s, 10s, 20s, 40s,...
	runAt := now.Add(BASE_RETRY_BACKOFF * time.Duration(math.Pow(2, float64(w.job.Retries-1))))
	w.job.RunAt = &runAt
	w.pushJob(w.job)
}
