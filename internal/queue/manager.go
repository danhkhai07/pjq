package queue

import (
	"context"
	"time"

	"pjq/internal/domain"
)

type JobManager struct {
	queue		*Queue
	jobCh 		chan domain.Job
	workerPool 	[]*worker
	numWorkers 	int
	registry  	*Registry
	store 		domain.JobStore
}

func NewJobManager(queue *Queue, numWorkers int, registry *Registry, store domain.JobStore) *JobManager {
	jm := JobManager{
		queue: queue,
		jobCh: make(chan domain.Job, numWorkers),
		workerPool: make([]*worker, numWorkers),
		numWorkers: numWorkers,
		registry: registry,
		store: store,
	}
	return &jm
}

func (jm *JobManager) PushJob(job domain.Job) {
	jm.queue.Push(job)
}

func (jm *JobManager) Run(ctx context.Context) {
	for i := range jm.workerPool {
		w := newWorker(i, jm.registry)
		jm.workerPool[i] = w
		go jm.RunWorker(ctx, w, jm.jobCh)
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

func (jm *JobManager) RunWorker(ctx context.Context, w *worker, jobCh chan domain.Job) (err error) {
	for {
		select {
		case job := <-jobCh:
			err = w.Process(ctx, &job)
			if err != nil {
				job.Error = err.Error()
				job.Logs = append(job.Logs, err.Error())
			}
			jm.store.Save(job)
		case <-ctx.Done():
			return nil
		}
	}
}
