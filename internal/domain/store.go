package domain

// Retriable: true - check if job can be retried | false - pass
type JobFilter struct {
	Status 		*Status
	Type 		*string
	Priority 	*int
	Retriable 	bool
}

func (jf JobFilter) Pass(job Job) bool {
	if jf.Status != nil && *jf.Status != job.Status {
		return false
	}
	if jf.Type != nil && *jf.Type != job.Type {
		return false
	}
	if jf.Priority != nil && *jf.Priority != job.Priority {
		return false
	}
	if jf.Retriable && job.Retries >= job.MaxRetries {
		return false
	} 
	return true
}

type JobStore interface {
	Save(job Job) error
	Get(id string) (Job, error)
	Delete(id string) error
	List(filter JobFilter) ([]Job, error)
}
