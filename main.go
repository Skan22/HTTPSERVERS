package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

const port = "8080"
const rootdir = "."

type apiConfig struct {
	fileServerHits atomic.Int32
}

func main() {
	apiCfg := apiConfig{}

	RequestRouter := http.NewServeMux()
	Server := &http.Server{
		Addr:    ":" + port,
		Handler: RequestRouter,
	}

	RequestRouter.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(rootdir)))))
	RequestRouter.Handle("GET /api/healthz", healthHandler())
	RequestRouter.Handle("GET /api/metrics", apiCfg.metricsHandler())
	RequestRouter.Handle("POST /api/reset", apiCfg.resetMetricsHandler())

	fmt.Printf("Serving files from %s on port %s", rootdir, port)
	log.Fatal(Server.ListenAndServe())

}

func healthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(http.StatusOK)))
	}
}

func (cfg *apiConfig) metricsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(fmt.Appendf(nil, "Hits: %d", cfg.fileServerHits.Load()))
	}
}

func (cfg *apiConfig) resetMetricsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		cfg.fileServerHits.Swap(0)
	}

}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
