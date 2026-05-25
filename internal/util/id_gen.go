package util

import (
	"crypto/rand"

	"github.com/oklog/ulid"
)

func GenerateULID() string {
	id := ulid.MustNew(ulid.Now(), rand.Reader)
	return id.String()
}
