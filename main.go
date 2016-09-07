package main

import (
	"log"
	"net/http"

	"github.com/sauron/app"
	"github.com/sauron/detectors"
	"github.com/sauron/replay"
)

//import _ "net/http/pprof"

func main() {
	http.HandleFunc("/", sauron.RequestHandler) // each request calls handler
	http.HandleFunc("/raw", sauron.RawHandler)
	http.HandleFunc("/stat", sauron.StatHandler)
	http.HandleFunc("/check", sauron.SessionCheckHandler)

	go func() {
		log.Println(http.ListenAndServe("localhost:3000", nil))
	}()

	/* Detectors */
	var compositeDetector = new(detectors.CompositeDetector)
	compositeDetector.Init("")

	var d1 = new(detectors.BlackListDetector)
	d1.Init("configs/ip_black_list.csv")
	compositeDetector.AddDetector(d1)

	var d2 = new(detectors.HumanPathDetector)
	d2.Init("configs/human_paths.csv")
	compositeDetector.AddDetector(d2)
	/* ~Detectors~ */

	sauron.Configure(compositeDetector)
	sauron.Start()

	replay.Start("/home/andrew/repos/data-miner-utils/dump.list")

	log.Println(http.ListenAndServe("localhost:6060", nil))
}
