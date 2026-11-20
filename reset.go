package main

import "net/http"

func (cfg *apiConfig) resetMetricsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Store(0)
		w.WriteHeader(http.StatusOK)
	}

}
