package main

import (
	"context"
	"pjq/internal/api"
)

func main() {
	svr := api.NewServer("0.0.0.0:8888")
	svr.Run(context.Background())
}
