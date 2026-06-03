package domain

import "context"

// Retriable: true - check if job can be retried | false - pass
type JobFilter struct {
	Status 		*Status
	Type 		*string
	Retriable 	bool
}

func (jf JobFilter) Pass(job Job) bool {
	if jf.Status != nil && *jf.Status != job.Status {
		return false
	}
	if jf.Type != nil && *jf.Type != job.Type {
		return false
	}
	if jf.Retriable && job.Retries >= job.MaxRetries {
		return false
	} 
	return true
}

type JobStore interface {
	Save(ctx context.Context, job Job) error
	Get(ctx context.Context, id string) (Job, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter JobFilter) ([]Job, error)
	Recover(ctx context.Context) ([]Job, error)
}
