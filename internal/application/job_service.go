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
	ctx context.Context,
	jobType string,
	payload []byte,
) (string, error) {
	jobID := util.GenerateULID()
	job := domain.NewJob(
		jobID,
		jobType,
		payload,
		PRIORITY_DEFAULT,
		MAX_RETRIES_DEFAULT,
	)
	err := js.store.Save(ctx, job)
	if err != nil {
		return "", err
	}
	js.queueManager.PushJob(job)
	return jobID, nil
}

func (js *JobService) GetJobByID(ctx context.Context, id string) (domain.Job, error) {
	return js.store.Get(ctx, id)
}

func (js *JobService) ListJobsWithFilter(ctx context.Context, filter domain.JobFilter) ([]domain.Job, error) {
	return js.store.List(ctx, filter)
}
