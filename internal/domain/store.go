package domain

type JobFilter struct {
	Status 		*Status
	Type 		*string
	Priority 	*int
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
	return true
}

type JobStore interface {
	Save(job Job) error
	Get(id string) (Job, error)
	Delete(id string) error
	List(filter JobFilter) ([]Job, error)
}
