package main

import (
	"context"
	"fmt"
	"pjq/internal/api"
	"pjq/internal/application"
	"pjq/internal/infra"
	queuePackage "pjq/internal/queue"
)

func main() {
	queue 			:= queuePackage.NewQueue()
	store, err 		:= infra.NewPSQLStore(
		"postgres://postgres:gratefultobealive@localhost:5432/pjq_db?sslmode=disable",
	)
	if err != nil {
		fmt.Println("DB failed to start.")
	}
	registry		:= queuePackage.NewRegistry()
	primeHandler	:= infra.NewPrimeCalcHandler()
	registry.Register("prime", primeHandler)

	queueManager 	:= queuePackage.NewQueueManager(queue, 3, registry, store)
	jobService 		:= application.NewJobService(store, queueManager)

	svr := api.NewServer("0.0.0.0:8888", jobService)
	fmt.Println("App started!")
	svr.Run(context.Background())
}
