package queue

import (
	"pjq/internal/domain"
)

type BackQueue interface {
	Push(*domain.Job) error
	Pop() (*domain.Job, bool)

	GetName() string
}
