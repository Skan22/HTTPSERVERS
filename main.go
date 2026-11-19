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
	Server := http.Server{
		Addr: ":" + port,
		Handler: RequestRouter,
	}

	RequestRouter.Handle("/",http.FileServer(http.Dir(rootdir)))
	
	fmt.Printf("Serving files from %s on port %s",rootdir,port)
	log.Fatal(Server.ListenAndServe())






}
