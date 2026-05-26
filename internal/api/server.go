package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"pjq/internal/application"
)

type Server struct {
	addr 			string
	httpServer 		*http.Server
	jobService		*application.JobService
}

func (svr *Server) Run(
	ctx context.Context,
	// args []string,
	// getenv func(string) string,
) {
	ctx, osCancel := signal.NotifyContext(ctx, os.Interrupt)
	defer osCancel()

	go func() {
		err := svr.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			fmt.Println("error http listen and serve")
			log.Fatal(err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutDownCtx := context.Background()
		shutDownCtx, shutDownCancel := context.WithTimeout(shutDownCtx, 10*time.Second)
		defer shutDownCancel()
		
		if err := svr.httpServer.Shutdown(shutDownCtx); err != nil {
			fmt.Println("error http shutdown")
			log.Fatal(err)
		}
	}()
	wg.Wait()
}

func NewServer(
	addr string,
	jobService *application.JobService,
) (svr *Server) {
	svr = &Server{
		addr: addr,
		httpServer: &http.Server{
			Addr: addr,
			Handler: NewHttpHandler(svr),
		},
		jobService: jobService,
	}
	return svr
}

func NewHttpHandler(
	svr *Server,
) http.Handler {
	mux := http.NewServeMux()
	addRoutes(
		svr,
		mux,
	)
	return mux
}
