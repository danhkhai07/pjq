package queue

import (
	"context"
	"errors"
	"log"
	"time"

	"pjq/internal/domain"
	"pjq/internal/util"
)

var (
	ErrFrontQueueNonPoppable = errors.New("front queue is non-poppable")
)

const (
	DISPATCH_RETRY_BACKOFF_SECONDS = 5
)

type QueueManager struct {
	fqueue		*FrontQueue
	bqueue 		BackQueue
	workerPool 	[]Worker
	wakeup		chan struct{}
	numWorkers 	int
	registry  	*util.Registry
	store 		domain.JobStore
}

func NewQueueManager(
	fqueue *FrontQueue,
	bqueue BackQueue,
	numWorkers int,
	registry *util.Registry,
	store domain.JobStore,
) *QueueManager {
	qm := QueueManager{
		fqueue: fqueue,
		bqueue: bqueue,
		workerPool: make([]Worker, numWorkers),
		wakeup: make(chan struct{}, 1),
		numWorkers: numWorkers,
		registry: registry,
		store: store,
	}
	return &qm
}

func (qm *QueueManager) PushJob(job *domain.Job) {
	qm.fqueue.Push(*job)
	qm.WakeUp()
}

func (qm *QueueManager) PopJob() (*domain.Job, bool) {
	return qm.bqueue.Pop()
}

func (qm *QueueManager) WakeUp() {
	select {
	case qm.wakeup <- struct{}{}:
	default:
	}
}

func (qm *QueueManager) Run(ctx context.Context) {
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
				qm.dispatchToBackQueue(0, nil)
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
				qm.dispatchToBackQueue(0, nil)
				continue
			}
		}
	}
}

// Default job param to nil so that front queue is popped
func (qm *QueueManager) dispatchToBackQueue(i int, job *domain.Job) {
	var err error
	if job == nil {
		j, ok := qm.fqueue.Pop()
		if !ok {
			err = ErrFrontQueueNonPoppable
		}
		job = &j
	}
	if !isProcessable(job) {
		return
	}
	if err == nil { err = qm.bqueue.Push(job) }
	if err != nil {
		log.Println(err.Error())
		if i < 5 {
			i++
			log.Println(
				"WARNING: Dispatching to %s failed. Retrying attempt number %d in %s seconds.",
				qm.bqueue.GetName(),
				i,
				DISPATCH_RETRY_BACKOFF_SECONDS,
			)
			time.Sleep(DISPATCH_RETRY_BACKOFF_SECONDS * time.Second)
			qm.dispatchToBackQueue(i, job)
		}
	}
}

func isProcessable(job *domain.Job) bool {
	return job.Status == domain.StatusPending || job.Status == domain.StatusRetrying
}
