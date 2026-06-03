package queue

import (
	"context"
	"fmt"
	"time"

	"pjq/internal/domain"
)

type QueueManager struct {
	queue		*Queue
	jobCh 		chan domain.Job
	workerPool 	[]*worker
	numWorkers 	int
	registry  	*Registry
	store 		domain.JobStore
}

func NewQueueManager(queue *Queue, numWorkers int, registry *Registry, store domain.JobStore) *QueueManager {
	qm := QueueManager{
		queue: queue,
		jobCh: make(chan domain.Job, numWorkers),
		workerPool: make([]*worker, numWorkers),
		numWorkers: numWorkers,
		registry: registry,
		store: store,
	}
	return &qm
}

func (qm *QueueManager) PushJob(job domain.Job) {
	qm.queue.Push(job)
}

func (qm *QueueManager) Run(ctx context.Context) {
	for i := range qm.workerPool {
		w := newWorker(i, qm.registry)
		qm.workerPool[i] = w
		go qm.RunWorker(ctx, w, qm.jobCh)
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			job, ok := qm.queue.Pop()
			if !ok {
				time.Sleep(10 * time.Millisecond)
				continue
			}
			qm.jobCh <- job
		}
	}
}

func (qm *QueueManager) RunWorker(ctx context.Context, w *worker, jobCh chan domain.Job) (err error) {
	for {
		select {
		case job := <-jobCh:
			now := time.Now()
			if job.StartedAt == nil {
				job.StartedAt = &now
			}
			changeStatus(&job, domain.StatusRunning)
			err = w.Process(ctx, &job)
			job.FinishedAt = &now
			if err != nil {
				changeStatus(&job, domain.StatusFailed)
				job.Error = err.Error()
				job.Logs = append(job.Logs, err.Error())
				if job.Retries < job.MaxRetries {
					changeStatus(&job, domain.StatusRetrying)
					qm.retry(job)
				}
			} else {
				changeStatus(&job, domain.StatusDone)
			}
			qm.store.Save(ctx, job)
			time.Sleep(10 * time.Millisecond)
		case <-ctx.Done():
			return nil
		}
	}
}

func changeStatus(job *domain.Job, status domain.Status) {
	job.Status = status
	job.Logs = append(
		job.Logs, 
		fmt.Sprintf("QueueManager: Changed job status to '%s'.", status),
	)
}

func (qm *QueueManager) retry(job domain.Job) {
	job.Retries++
	qm.PushJob(job)
}
