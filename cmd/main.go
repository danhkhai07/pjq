package main

import (
	"context"
	"fmt"
	"pjq/internal/api"
	"pjq/internal/application"
	"pjq/internal/domain"
	"pjq/internal/infra"
	queuePackage "pjq/internal/queue"
)

func main() {
	fmt.Println("App started!")
	queue 			:= queuePackage.NewQueue()
	store 			:= infra.NewInMemoryStore(map[string]domain.Job{})
	registry		:= queuePackage.NewRegistry()
	primeHandler	:= infra.NewPrimeCalcHandler()
	registry.Register("prime", primeHandler)

	queueManager 	:= queuePackage.NewQueueManager(queue, 3, registry, store)
	jobService 		:= application.NewJobService(store, queueManager)

	svr := api.NewServer("0.0.0.0:8888", jobService)
	svr.Run(context.Background())
}
