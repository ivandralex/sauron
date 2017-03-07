package sauron

import (
	"fmt"
	"net/http"
	"time"

	"github.com/paulbellamy/ratecounter"
	"github.com/sauron/session"
)

var rpsCounter = ratecounter.NewRateCounter(10 * time.Second)

//RequestHandler handles incoming request
func RequestHandler(w http.ResponseWriter, r *http.Request) {
	//Form session key
	var sessionKey = sstrg.GetSessionKey(r)
	fmt.Println(sessionKey)
	//Extract useful data from request
	var request = sstrg.GetRequestData(r, config.useDataHeader)

	HandleRequest(sessionKey, request)

	fmt.Fprintf(w, "OK")
}

//HandleRequest handles parsed request
func HandleRequest(sessionKey string, request *sstrg.RequestData) {
	rpsCounter.Incr(1)
	//Create stat struct for new session
	sessions.Lock()

	//Save new session to storage
	if _, ok := sessions.H[sessionKey]; !ok {
		sessions.H[sessionKey] = new(sstrg.SessionData)
		sessions.H[sessionKey].Started = request.Time
		sessions.H[sessionKey].IP = request.IP
	}

	//If session was inactive and was not deleted and we received request with same session key we re-activate session
	sessions.H[sessionKey].Active = true
	sessions.H[sessionKey].Ended = request.Time

	//Update emulated value
	if config.emulateTime {
		emulatedTime = request.Time
	}

	//Save request to session
	sessions.H[sessionKey].Requests = append(sessions.H[sessionKey].Requests, request)

	sessions.Unlock()
}

//SessionCheckHandler checks session with default detector
func SessionCheckHandler(w http.ResponseWriter, r *http.Request) {
	var sessionKey = r.URL.Query().Get("key")

	if sessionKey == "" {
		sessionKey = sstrg.GetSessionKey(r)
	}

	sessions.RLock()

	if _, ok := sessions.H[sessionKey]; !ok {
		w.WriteHeader(404)
		fmt.Fprint(w, "Key Not Found: "+sessionKey)
		sessions.RUnlock()
		return
	}

	var session = sessions.H[sessionKey]
	var label = defaultDetector.GetLabel(session)

	sessions.RUnlock()

	fmt.Fprint(w, label)
}

//RawHandler outputs raw requests data for specified session key
func RawHandler(w http.ResponseWriter, r *http.Request) {
	var sessionKey = r.URL.Query().Get("key")

	if _, ok := sessions.H[sessionKey]; !ok {
		fmt.Fprint(w, "Not Found")
		return
	}

	sessions.RLock()

	var session = sessions.H[sessionKey]

	fmt.Fprintf(w, "Started: %v\nLast: %v\nRequests: %d\nActive: %v\n\n", session.Started, session.Ended, len(session.Requests), session.Active)

	for _, r := range sessions.H[sessionKey].Requests {

		fmt.Fprintf(w, "%s %s\n", r.Method, r.Path)

		for _, c := range r.Cookies {
			fmt.Fprintf(w, "%s=%s;", c.Name, c.Path)
		}

		fmt.Fprintf(w, "\n")

		for k, h := range r.Header {
			fmt.Fprintf(w, "%s: %s\n", k, h[0])
		}

		fmt.Fprintf(w, "\n\n\n\n")
	}

	sessions.RUnlock()
}

//StatHandler outputs current RPS
func StatHandler(w http.ResponseWriter, r *http.Request) {
	rps := rpsCounter.Rate() / 10

	fmt.Fprintf(w, "RPS: %d", rps)
}
