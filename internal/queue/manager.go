package queue

import (
	"context"
	"fmt"
	"math"
	"time"

	"pjq/internal/domain"
)

const (
	BASE_RETRY_BACKOFF = 5 * time.Second
)

type QueueManager struct {
	queue		*Queue
	jobCh 		chan domain.Job
	workerPool 	[]*worker
	wakeup		chan struct{}
	numWorkers 	int
	registry  	*Registry
	store 		domain.JobStore
}

func NewQueueManager(queue *Queue, numWorkers int, registry *Registry, store domain.JobStore) *QueueManager {
	qm := QueueManager{
		queue: queue,
		jobCh: make(chan domain.Job, numWorkers),
		workerPool: make([]*worker, numWorkers),
		wakeup: make(chan struct{}, 1),
		numWorkers: numWorkers,
		registry: registry,
		store: store,
	}
	return &qm
}

func (qm *QueueManager) PushJob(job domain.Job) {
	qm.queue.Push(job)
	qm.WakeUp()
}

func (qm *QueueManager) WakeUp() {
	select {
	case qm.wakeup <- struct{}{}:
	default:
	}
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
			next := qm.queue.Peek()

			if next == nil {
				select {
				case <-ctx.Done():
					return
				case <-qm.wakeup:
					continue
				}
			}

			wait := time.Until(*next.RunAt)
			if wait <= 0 {
				job, ok := qm.queue.Pop()
				if ok {
					qm.jobCh <- job
				}
				continue
			}
			timer := time.NewTimer(wait)

			select {
			case <-ctx.Done():
				if !timer.Stop() {
					<-timer.C
				}
				return
			case <-qm.wakeup:
				if !timer.Stop() {
					<-timer.C
				}
				continue
			case <-timer.C:
				job, ok := qm.queue.Pop()
				if ok {
					qm.jobCh <- job
				}
				continue
			}
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
	now := time.Now()
	job.Retries++
	// exponential back-off with power of 2
	runAt := now.Add(BASE_RETRY_BACKOFF * time.Duration(math.Pow(2, float64(job.Retries-1))))
	job.RunAt = &runAt
	qm.PushJob(job)
}
