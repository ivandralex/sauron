package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/sauron/app"
	"github.com/sauron/detectors"
	"github.com/sauron/extractors"
	"github.com/sauron/replay"
	"github.com/sauron/writers"
)

//import _ "net/http/pprof"

const replayMode = "replay"
const detectMode = "detect"

func main() {
	//Define flags
	mode := flag.String("mode", detectMode, "Application mode: replay|detect")
	//replayFrom := flag.String("replay-from", "", "Replay s=tart date in ISO8601 format")
	outputFile := flag.String("output", "", "Path to output file for features")
	dumpFile := flag.String("dump", "", "Path to dump file")
	flag.Parse()

	if *mode == replayMode && *dumpFile == "" {
		log.Fatal("Dump file not specified")
	}

	fmt.Printf("mode: %s\n", *mode)
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

	var d1 = new(detectors.IPListDetector)
	d1.Init("configs/ip_black_list.csv")
	compositeDetector.AddDetector(d1)

	var d2 = new(detectors.PathDetector)
	d2.Init("configs/human_paths.csv")
	compositeDetector.AddDetector(d2)
	/* ~Detectors~ */

	var extractor = new(extractors.RequestsSequence)
	extractor.Init("configs/requests_sequence.csv")

	var writer = new(writers.CSVWriter)
	writer.Init(*outputFile)

	sauron.Configure(compositeDetector, extractor, writer)
	sauron.Start()

	if *mode == replayMode {
		fmt.Println("Gonna replay traffic")
		replay.Start(*dumpFile)
	} else {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}
}
