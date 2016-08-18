package main

import (
	"log"
	"net/http"

	"github.com/sauron/extractor"
)

//import _ "net/http/pprof"

func main() {
	extractor.Start()

	http.HandleFunc("/", extractor.RequestHandler) // each request calls handler
	http.HandleFunc("/raw", extractor.RawHandler)
	http.HandleFunc("/stat", extractor.StatHandler)
	http.HandleFunc("/features", extractor.FeaturesHandler)

	log.Fatal(http.ListenAndServe("0.0.0.0:3000", nil))
	/*
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	*/
}
