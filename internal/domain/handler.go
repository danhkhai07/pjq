package domain

import "context"

type JobHandler interface {
	Handle(ctx context.Context, job Job, log func(string)) error
}
