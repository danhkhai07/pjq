package queue

import (
	"pjq/internal/domain"
)

type BackQueue interface {
	Push(job *domain.Job) error
}
