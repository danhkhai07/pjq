package queue

import (
	"fmt"
	"pjq/internal/domain"
)

type Registry struct {
	handlers map[string]domain.JobHandler
}

func NewRegistry() *Registry {
	return &Registry{
		make(map[string]domain.JobHandler),
	}
}

func (r *Registry) Register(jobType string, handler domain.JobHandler) error {
	if _, ok := r.handlers[jobType]; ok == false {
		r.handlers[jobType] = handler
		return nil
	} 
	return fmt.Errorf("'%s' is already a registered job type", jobType)
}

func (r *Registry) Get(jobType string) (domain.JobHandler, error) {
	h, ok := r.handlers[jobType]
	if !ok {
		return nil, fmt.Errorf("cannot find handler for job type '%s'", jobType)
	}
	return h, nil
}
