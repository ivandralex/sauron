package extractor

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/paulbellamy/ratecounter"
	"github.com/sauron/detectors"
	"github.com/sauron/features"
	"github.com/sauron/session"
	"github.com/sauron/stat"
)

//import _ "net/http/pprof"

//TODO: move to file
var config = struct {
	useDataHeader      bool
	emulateTime        bool
	sessionsPeriod     int
	featuresPeriod     int
	statPeriod         int
	maxInactiveMinutes float64
	//Feature flags
	beholdStat        bool
	beholdFeatures    bool
	beholdSessionsEnd bool
}{
	useDataHeader:      true,
	emulateTime:        true,
	beholdStat:         false,
	beholdSessionsEnd:  true,
	beholdFeatures:     false,
	sessionsPeriod:     5,
	statPeriod:         5,
	featuresPeriod:     5,
	maxInactiveMinutes: 60.0}

var sessions = new(sstrg.SessionsTable)
var emulatedTime time.Time
var rpsCounter = ratecounter.NewRateCounter(10 * time.Second)

var defaultDetector detectors.Detector

func init() {
	sessions.H = make(map[string]*sstrg.SessionData)
}

//Start features extractor
func Start(detector detectors.Detector) {
	defaultDetector = detector

	//Check if sessions active periodically
	if config.beholdSessionsEnd {
		go startSessionsBeholder(config.sessionsPeriod)
	}
	//Collect stat on sessions
	if config.beholdStat {
		go statutils.StartSessionsStatBeholder(sessions, config.statPeriod)
	}

	//Dump features periodically
	if config.beholdFeatures {
		go startFeaturesBeholder(sessions, config.featuresPeriod)
	}
}

func startSessionsBeholder(periodSec int) {
	// fire once per second
	t := time.NewTicker(time.Second * time.Duration(periodSec))

	for {
		closeSessions()
		<-t.C
	}
}

func closeSessions() {
	//fmt.Fprintf(os.Stdout, "Closing sessions!\n")

	var dur float64
	var now time.Time
	if config.emulateTime {
		now = emulatedTime
	} else {
		now = time.Now().UTC()
	}

	sessions.RLock()

	for _, s := range sessions.H {
		if !s.Active {
			continue
		}

		dur = now.Sub(s.Ended).Minutes()

		if dur > config.maxInactiveMinutes {
			//fmt.Fprintf(os.Stdout, "Closed session: %s cauze inactive for %f\n", k, dur)
			//Mark session as inactive. It will be deleted after the next dump
			s.Active = false
		}
	}

	sessions.RUnlock()
}

func startFeaturesBeholder(sessions *sstrg.SessionsTable, periodSec int) {
	// fire once per second
	t := time.NewTicker(time.Second * time.Duration(periodSec))

	absPath, _ := filepath.Abs("output/features/new.csv")
	os.Remove(absPath)

	f, err := os.Create(absPath)

	if err != nil {
		log.Fatalln("Error creating features file:", err)
	}

	w := csv.NewWriter(f)

	for {
		dumpFeatures(w, sessions)
		<-t.C
	}
}

func dumpFeatures(w *csv.Writer, sessions *sstrg.SessionsTable) {
	sessions.Lock()

	for key, s := range sessions.H {
		//We do not dump not-yet-finished sessions
		if s.Active {
			continue
		}

		var fvDesc = pathvector.ExtractFeatures(s)
		//Append label
		var label = defaultDetector.GetLabel(s)

		//If bot was not detected check if it's a human
		if label == "0" {
			label = defaultDetector.GetLabel(s)
		} else {
			fmt.Printf("Found bot %s\n", s.IP)
		}

		fvDesc = append(fvDesc, label)

		if err := w.Write(fvDesc); err != nil {
			log.Fatalln("Error writing record to csv:", err)
		}

		//Delete session from storage
		delete(sessions.H, key)
	}

	sessions.Unlock()

	w.Flush()
}

//StatHandler outputs current RPS
func StatHandler(w http.ResponseWriter, r *http.Request) {
	rps := rpsCounter.Rate() / 10

	fmt.Fprintf(w, "RPS: %d", rps)
}

//RequestHandler handles incoming request
func RequestHandler(w http.ResponseWriter, r *http.Request) {
	//Form session key
	var sessionKey = sstrg.GetSessionKey(r)
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

//SessionCheckHandler checks session with default detector
func SessionCheckHandler(w http.ResponseWriter, r *http.Request) {
	var sessionKey = r.URL.Query().Get("key")

	if sessionKey == "" {
		sessionKey = sstrg.GetSessionKey(r)
	}

	fmt.Fprint(os.Stdout, sessionKey)

	sessions.RLock()

	if _, ok := sessions.H[sessionKey]; !ok {
		fmt.Fprint(w, "Key Not Found: "+sessionKey)
		sessions.RUnlock()
		return
	}

	var session = sessions.H[sessionKey]
	var label = defaultDetector.GetLabel(session)

	sessions.RUnlock()

	fmt.Fprint(w, label)
}
