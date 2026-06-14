package queue

import (
	"context"
	"fmt"
	"math"
	"time"

	"pjq/internal/domain"
	"pjq/internal/util"
)

const (
	BASE_RETRY_BACKOFF = 5 * time.Second
)

type QueueManager struct {
	fqueue		*FrontQueue
	bqueue 		chan domain.Job
	workerPool 	[]Worker
	wakeup		chan struct{}
	numWorkers 	int
	registry  	*util.Registry
	store 		domain.JobStore
}

func NewQueueManager(fqueue *FrontQueue, numWorkers int, registry *util.Registry, store domain.JobStore) *QueueManager {
	qm := QueueManager{
		fqueue: fqueue,
		bqueue: make(chan domain.Job, numWorkers),
		workerPool: make([]Worker, numWorkers),
		wakeup: make(chan struct{}, 1),
		numWorkers: numWorkers,
		registry: registry,
		store: store,
	}
	return &qm
}

func (qm *QueueManager) PushJob(job domain.Job) {
	qm.fqueue.Push(job)
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
		go qm.RunWorker(ctx, w, qm.bqueue)
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			next := qm.fqueue.Peek()

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
				job, ok := qm.fqueue.Pop()
				if ok {
					qm.bqueue <- job
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
				job, ok := qm.fqueue.Pop()
				if ok {
					qm.bqueue <- job
				}
				continue
			}
		}
	}
}

func (qm *QueueManager) RunWorker(ctx context.Context, w Worker, bqueue chan domain.Job) (err error) {
	for {
		select {
		case job := <-bqueue:
			now := time.Now()
			if job.StartedAt == nil {
				job.StartedAt = &now
			}
			changeStatus(&job, domain.StatusRunning)
			err = w.Process(ctx, &job)
			job.FinishedAt = &now
			if err != nil {
				changeStatus(&job, domain.StatusFailed)
				logError(&job, err)
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

func logError(job *domain.Job, err error) {
	logTime := time.Now().Local().Local().String()
	job.Error = err.Error()
	job.Logs = append(job.Logs, logTime + " " + err.Error())
}

func changeStatus(job *domain.Job, status domain.Status) {
	logTime := time.Now().Local().Local().String()
	job.Status = status
	job.Logs = append(
		job.Logs, 
		fmt.Sprintf("%s QueueManager: Changed job status to '%s'.", logTime, status),
	)
}

func (qm *QueueManager) retry(job domain.Job) {
	now := time.Now()
	job.Retries++
	// exponential back-off with power of 2: 5s, 10s, 20s, 40s,...
	runAt := now.Add(BASE_RETRY_BACKOFF * time.Duration(math.Pow(2, float64(job.Retries-1))))
	job.RunAt = &runAt
	qm.PushJob(job)
}
