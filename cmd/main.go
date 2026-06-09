package main

import (
	"context"
	"fmt"
	"os"
	"pjq/internal/api"
	"pjq/internal/application"
	"pjq/internal/infra"
	queuePackage "pjq/internal/queue"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load environment variables: %s\n", err)
	}
	DB_URL := os.Getenv("POSTGRES_URL")

	queue 			:= queuePackage.NewQueue()
	store, err 		:= infra.NewPSQLStore(DB_URL)
	if err != nil {
		fmt.Println("DB failed to start.")
		return
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
