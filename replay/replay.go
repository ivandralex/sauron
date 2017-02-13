package replay

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/sauron/app"
	"github.com/sauron/session"
)

//Start starts replay of the traffic from the dump file
func Start(dumpPath string) {
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

	fmt.Println("Replay finished")

	f.Close()
}

func handleRequest(source []string) {
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
	request.IP = strings.Split(source[2], ",")[0]
	request.Path, request.ContentType = sstrg.GetContentType(request.Path)
	request.Cookies = sstrg.GetCookiesFromCookiesString(source[6])

	var sessionKey = request.IP // source[4]

	sauron.HandleRequest(sessionKey, request)
}
