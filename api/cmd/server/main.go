package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sei-ri/pingpong.io/api/internal/env"
	"github.com/sei-ri/pingpong.io/api/internal/gateway"
	"github.com/sei-ri/pingpong.io/api/internal/swaggerui"
	"github.com/sei-ri/pingpong.io/protos"
)

const (
	DefaultPort       = "3000"
	DefaultSwaggerDir = "apis/v1"
)

func main() {

	port := env.Get("PORT", DefaultPort)
	swaggerDir := env.Get("SWAGGER_DIR", DefaultSwaggerDir)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	opts := []runtime.ServeMuxOption{}
	gw, err := gateway.NewGateway(ctx, opts, []gateway.Service{
		{
			Name:    "ping",
			Address: env.Get("PING_SERVER_ADDR", ":10010"),
			Handler: protos.RegisterServiceHandler,
		},
		{
			Name:    "ping2",
			Address: env.Get("PING2_SERVER_ADDR", ":10011"),
			Handler: protos.RegisterServiceHandler,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/swagger/", gateway.SwaggerServer(swaggerDir))
	swaggerui.Serve(mux)

	mux.Handle("/", gw)

	srv := &http.Server{
		Addr:         net.JoinHostPort("", port),
		Handler:      gateway.AllowCORS(mux),
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	go func() {
		<-ctx.Done()
		log.Println("Shutting down the http server")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("Failed to shutdown http server: %v", err)
		}
	}()

	log.Printf("http server listening at %v", net.JoinHostPort("", port))
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Failed to listen and serve: %v", err)
	}
}
