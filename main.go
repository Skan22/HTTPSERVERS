package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"unicode/utf8"
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
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write(fmt.Appendf(nil, `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
		</html>`, cfg.fileServerHits.Load()))
	}
}

func (cfg *apiConfig) resetMetricsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		cfg.fileServerHits.Swap(0)
	}

}

func validateChirpHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Body string `json:"body"`
		}
		type returnVals struct {
			Valid bool   `json:"valid"`
			Error string `json:"error"`
		}
		params := parameters{}
		respBody := returnVals{}

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&params)

		if err != nil {
			w.WriteHeader(500)
			respBody.Error = "Something went wrong"
		}

		if utf8.RuneCountInString(params.Body) > 140 {
			respBody.Valid = false
			respBody.Error = "Chirp is too long"
			w.WriteHeader(400)

		} else {
			respBody.Valid = true
			w.WriteHeader(200)
		}

		dat, err := json.Marshal(respBody)
		if err != nil {
			w.WriteHeader(500)
			respBody.Error = "Something went wrong"

		}
		
		w.Header().Set("Content-Type", "application/json")
		w.Write(dat)
	}
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
