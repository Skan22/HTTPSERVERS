package main

import (
	"fmt"
	"log"
	"net/http"
)



func main(){
	const port = "8080"
	const rootdir = "."

	RequestRouter := http.NewServeMux()
	Server := &http.Server{
		Addr: ":" + port,
		Handler: RequestRouter,
	}
	
	RequestRouter.Handle("/app/",http.StripPrefix("/app",http.FileServer(http.Dir(rootdir))))

	RequestRouter.Handle("/healthz",healthHandler())
	fmt.Printf("Serving files from %s on port %s",rootdir,port)
	log.Fatal(Server.ListenAndServe())

}

func healthHandler()http.HandlerFunc{
		return func(w http.ResponseWriter, r *http.Request) {

			w.Header().Set("Content-Type","text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(http.StatusText(http.StatusOK)))
		}
	}