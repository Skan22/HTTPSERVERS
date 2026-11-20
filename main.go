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
	RequestRouter.Handle("GET /admin/metrics", apiCfg.metricsHandler())
	RequestRouter.Handle("POST /admin/reset", apiCfg.resetMetricsHandler())
	RequestRouter.Handle("POST /api/validate_chirp", validateChirpHandler())
	fmt.Printf("Serving files from :  %s \n on port  : %s", rootdir, port)
	log.Fatal(Server.ListenAndServe())

}
