package queue

import (
	"pjq/internal/domain"
)

type BQueue interface {
	Push(job *domain.Job) error
}
