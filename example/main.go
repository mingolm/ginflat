package main

import (
	"context"
	"github.com/mingolm/ginflat/example/middleware"
	"net/http"
	"time"

	"github.com/mingolm/ginflat"
	"github.com/mingolm/ginflat/example/controller"
)

func main() {
	handler := ginflat.NewHandler()
	handler.Use(middleware.ErrHandler())
	handler.RegisterController("localhost.test", &controller.User{})
	handler.RegisterController("admin.localhost.test", &controller.Admin{})

	srv := newHttpServer(handler)

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

func newHttpServer(handler *ginflat.Handler) *http.Server {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	if err := handler.InitControllers(ctx); err != nil {
		panic(err)
	}
	cancel()

	return &http.Server{
		Addr:         ":23380",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  150 * time.Second,
	}
}
