package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/a-h/templ"
	"github.com/danielmmetz/templ/templates"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	if err := mainE(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		// Only exit non-zero if our initial context has yet to be canceled.
		// Otherwise it's very likely that the error we're seeing is a result of our attempt at graceful shutdown.
		if ctx.Err() == nil {
			os.Exit(1)
		}
	}
}

func mainE(ctx context.Context) error {
	http.Handle("/", templ.Handler(templates.Index()))
	itemsHandler := ItemsHandler{items: []string{"apples", "bananas", "carrots"}}
	http.Handle("/list", Handle(itemsHandler.List))
	http.Handle("/add-item", Handle(itemsHandler.AddItem))
	http.Handle("/delete-item", Handle(itemsHandler.DeleteItem))

	s := http.Server{Addr: "localhost:8080"}
	var eg errgroup.Group
	eg.Go(func() error {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		return s.Shutdown(ctx)
	})

	switch err := s.ListenAndServe(); err {
	case http.ErrServerClosed:
		return eg.Wait()
	default:
		return err
	}
}
