package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sauron/extractor"
	"github.com/sauron/session"
)

func main() {
	http.HandleFunc("/stat", extractor.StatHandler)

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	extractor.Start()

	readDump()
}

func readDump() {
	var dumpPath = "/home/andrew/repos/data-miner-utils/dump.list"

	f, err := os.Open(dumpPath)

	if err != nil {
		fmt.Println("error opening file ", err)
		os.Exit(1)
	}

	r := bufio.NewReader(f)
	for {
		str, err := r.ReadString(10)
		if err == io.EOF {
			break
		} else if err != nil {
			continue
		}

		pieces := strings.Split(str, "|")

		handleRequest(pieces)
	}

	f.Close()
}

func handleRequest(source []string) {
	if len(source) != 6 {
		return
	}

	//var args = source[5]

	var request = new(sstrg.RequestData)

	t, err := time.Parse(time.RFC3339Nano, source[3])

	if err == nil {
		request.Time = t
	} else {
		request.Time = time.Now().UTC()
	}

	request.Path = source[0]

	if strings.HasSuffix(request.Path, "/") {
		request.Path = request.Path[:len(request.Path)-1]
	}

	request.Method = "GET"
	request.Referer = source[1]

	request.Path, request.ContentType = sstrg.GetContentType(request.Path)

	var sessionKey = source[2] // source[4]

	extractor.HandleRequest(sessionKey, request)
}
