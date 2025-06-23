package application

import (
	"context"
	"go.uber.org/fx"
	"log"
	"net/http"
)

const port = ":8080"

func NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, world!"))
	})
	return mux
}

func StartServer(lc fx.Lifecycle, mux *http.ServeMux) {
	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Println("Starting server on " + port)
				if err := server.ListenAndServe(); err != http.ErrServerClosed {
					log.Fatalf("ListenAndServe: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Shutting down server...")
			return server.Shutdown(ctx)
		},
	})
}

var Module = fx.Options(
	fx.Provide(NewMux),
	fx.Invoke(StartServer),
)
