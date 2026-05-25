package queue

import (
	"context"
	"time"

	"pjq/internal/domain"
)

type JobManager struct {
	queue		*Queue
	jobCh 		chan domain.Job
	workerPool 	[]*Worker
	numWorkers 	int
	registry  	*Registry
}

func NewJobManager(queue *Queue, numWorkers int, registry *Registry) *JobManager {
	jm := JobManager{
		queue: queue,
		jobCh: make(chan domain.Job, numWorkers),
		workerPool: make([]*Worker, numWorkers),
		numWorkers: numWorkers,
		registry: registry,
	}
	return &jm
}

func (jm *JobManager) Run(ctx context.Context) {
	for i := range jm.workerPool {
		w := NewWorker(i, jm.registry)
		jm.workerPool[i] = w
		go RunWorker(ctx, w, jm.jobCh)
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			job, ok := jm.queue.Pop()
			if !ok {
				time.Sleep(10 * time.Millisecond)
				continue
			}
			jm.jobCh <- job
		}
	}
}

func RunWorker(ctx context.Context, w *Worker, jobCh chan domain.Job) (err error) {
	for {
		select {
		case job := <-jobCh:
			err = w.Process(ctx, &job)
			if err != nil {
				job.Error = err.Error()
				job.Logs = append(job.Logs, err.Error())
			}
			// store job
		case <-ctx.Done():
			return nil
		}
	}
}
