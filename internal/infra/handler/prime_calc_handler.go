package infra

import (
	"context"
	"encoding/json"
	"errors"

	"pjq/internal/domain"
	dto "pjq/internal/dto/handler"
)

type PrimeCalcHandler struct {
	flag []bool
	nthPrime []int
}

func NewPrimeCalcHandler() *PrimeCalcHandler {
	h := PrimeCalcHandler{
		flag: make([]bool, 10001),
		nthPrime: make([]int, 0),
	}

	h.nthPrime = append(h.nthPrime, 0) // offset
	for i := 2; i <= 10000; i++ {
		if !h.flag[i] {
			h.nthPrime = append(h.nthPrime, i)
			for j := 2; j*i <= 10000; j++ {
				h.flag[j*i] = true
			}
		}
	}
	h.flag = nil
	return &h
}

func (h *PrimeCalcHandler) Handle(ctx context.Context, job *domain.Job, log func(string)) error {
	log("Handler: Job started processing.")
	payload := dto.PrimeCalcHandlerPayload{}
	err := json.Unmarshal(job.Payload, &payload)
	if err != nil {
		errMsg := "Handler: Error: Bad payload."
		job.Error = errMsg
		return errors.New(errMsg)
	}

	if payload.N < 1 || payload.N > 1000 {
		errMsg := "Handler: Error: N out of range [1, 1000]."
		job.Error = errMsg
		return errors.New(errMsg)
	}

	job.Result = h.nthPrime[payload.N]
	log("Handler: Job finished processing.")
	return nil
}
