package main

import (
	"log"
	"net/http"

	"github.com/sauron/app"
	"github.com/sauron/detectors"
)

//import _ "net/http/pprof"

func main() {

	//"../configs/human_paths.csv"

	var defaultDetector = new(detectors.BlackListDetector)
	defaultDetector.Init("configs/ip_black_list.csv")

	sauron.Start(defaultDetector)

	http.HandleFunc("/", sauron.RequestHandler) // each request calls handler
	http.HandleFunc("/raw", sauron.RawHandler)
	http.HandleFunc("/stat", sauron.StatHandler)
	http.HandleFunc("/check", sauron.SessionCheckHandler)

	log.Fatal(http.ListenAndServe("0.0.0.0:3000", nil))
	/*
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	*/
}
