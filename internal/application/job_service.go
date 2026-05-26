package application

import (
	"context"
	"pjq/internal/domain"
	"pjq/internal/queue"
	"pjq/internal/util"
)

const (
	PRIORITY_DEFAULT = 1
	MAX_RETRIES_DEFAULT = 3
)

type JobService struct {
	store 			domain.JobStore
	queueManager	*queue.QueueManager
}

func NewJobService(
	store domain.JobStore,
	queueManager *queue.QueueManager,
) *JobService {
	return &JobService{
		store: store,
		queueManager: queueManager,
	}
}

func (js *JobService) Run(ctx context.Context) {
	js.queueManager.Run(ctx)
}

// Return new job id.
func (js *JobService) ProcessNewJob(
	jobType string,
	payload []byte,
) string {
	jobID := util.GenerateULID()
	job := domain.NewJob(
		jobID,
		jobType,
		payload,
		PRIORITY_DEFAULT,
		MAX_RETRIES_DEFAULT,
	)
	js.store.Save(job)
	js.queueManager.PushJob(job)
	return jobID
}

func (js *JobService) GetJobByID(id string) (domain.Job, error) {
	return js.store.Get(id)
}

func (js *JobService) ListJobsWithFilter(filter domain.JobFilter) ([]domain.Job, error) {
	return js.store.List(filter)
}
