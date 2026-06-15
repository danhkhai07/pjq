package main

import (
	"context"
	"fmt"
	"os"
	"pjq/internal/api"
	"pjq/internal/application"
	"pjq/internal/queue"
	"pjq/internal/util"

	queueinfra "pjq/internal/infra/back_queue"
	storeinfra "pjq/internal/infra/store"
	handlerinfra "pjq/internal/infra/handler"
	workerinfra "pjq/internal/infra/worker"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load environment variables: %s\n", err)
	}
	DB_URL := os.Getenv("POSTGRES_URL")

	fqueue 			:= queue.NewQueue()
	bqueue			:= queueinfra.NewInProcessBackQueue(5)
	store, err 		:= storeinfra.NewPSQLStore(DB_URL)
	if err != nil {
		fmt.Println("DB failed to start.")
		return
	}
	registry		:= util.NewRegistry()
	primeHandler	:= handlerinfra.NewPrimeCalcHandler()
	registry.Register("prime", primeHandler)

	queueManager 	:= queue.NewQueueManager(fqueue, bqueue, 3, registry, store)
	jobService 		:= application.NewJobService(store, queueManager)

	// spawn workers
	workers := make([]*workerinfra.Worker, 0)
	for i := 1; i <= 3; i++ {
		worker := workerinfra.NewInProcessWorker(
			registry,
			store,
			queueManager,
		)
		workers = append(workers, worker)
		go worker.RunWorker(context.Background())
	}

	svr := api.NewServer("0.0.0.0:8888", jobService)
	fmt.Println("App started!")
	svr.Run(context.Background())
}
