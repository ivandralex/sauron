package main

import (
	"log"
	"net/http"

	"github.com/sauron/detectors"
	"github.com/sauron/extractor"
)

//import _ "net/http/pprof"

func main() {

	//"../configs/human_paths.csv"

	var defaultDetector = new(detectors.BlackListDetector)
	defaultDetector.Init("configs/ip_black_list.csv")

	extractor.Start(defaultDetector)

	http.HandleFunc("/", extractor.RequestHandler) // each request calls handler
	http.HandleFunc("/raw", extractor.RawHandler)
	http.HandleFunc("/stat", extractor.StatHandler)
	http.HandleFunc("/check", extractor.SessionCheckHandler)

	log.Fatal(http.ListenAndServe("0.0.0.0:3000", nil))
	/*
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	*/
}
