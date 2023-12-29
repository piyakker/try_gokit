package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/piyakker/try_gokit/model"
)

var (
	httpAddr    = ":8080"
	httpTimeout = 60 * time.Second
)

func main() {
	service := NewService()
	handler := MakeHandler(service)
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(httpTimeout))
	r.Mount("/api/v1", handler)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		log.Printf("listening on %s", httpAddr)
		errs <- http.ListenAndServe(httpAddr, r)
	}()

	log.Println("exit", <-errs)
}
