package queue

import (
	"context"
	"errors"
	"pjq/internal/domain"
)

var (
	ErrWorkerIsBusy error = errors.New("worker is busy")
)

type Worker interface {
	IsAlive() bool
	IsBusy() bool
	Log(message string)
	Process(ctx context.Context, job *domain.Job) error
}
