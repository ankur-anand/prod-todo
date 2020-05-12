package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ankur-anand/prod-todo/pkg/apis/rest"
	logger2 "github.com/ankur-anand/prod-todo/pkg/logger"
)

type server struct {
	logger  *logger2.Logger
	handler *rest.MuxHandler
}

func main() {
	logger, _ := logger2.NewDevelopment()
	mh := rest.NewMuxHandler(logger)
	s := server{
		logger:  logger,
		handler: mh,
	}
	s.run("3000")
}

// run method to simply start the API Server
func (s server) run(port string) {

	msg := fmt.Sprintf("Starting RestServiceApplication with PID %d", os.Getpid())
	s.logger.Info(msg)
	msg = fmt.Sprintf("Server initializing with port(s): %s (http)", port)
	s.logger.Info(msg)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	addr := ":" + port

	// httpServer.
	h := &http.Server{Addr: addr, Handler: s.handler}

	go func() {
		err := h.ListenAndServe()
		s.logger.Fatal("err starting http rest server", logger2.Error(err))
	}()

	msg = fmt.Sprintf("Server started on port(s): %s (http)", port)
	s.logger.Info(msg)
	<-stop
	msg = fmt.Sprintf("Shutting down RestServiceApplication with PID %d", os.Getpid())
	s.logger.Info(msg)
	close(stop)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := h.Shutdown(ctx)
	if err != nil {
		s.logger.Fatal("error while gracefully stopping RestServiceApplication Server", logger2.Error(err))
	}
	msg = fmt.Sprintf("Server gracefully stopped on port(s): %s (http)", port)
	s.logger.Info(msg)
}
