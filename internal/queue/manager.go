package queue

import (
	"context"
	"time"

	"pjq/internal/domain"
	"pjq/internal/util"
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
					qm.bqueue.Push(&job)
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
					qm.bqueue.Push(&job)
				}
				continue
			}
		}
	}
}
