package sauron

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/paulbellamy/ratecounter"
	"github.com/sauron/detectors"
	"github.com/sauron/extractors"
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
	writeRelevantOnly bool
	beholdSessionsEnd bool
}{
	useDataHeader:      true,
	emulateTime:        true,
	beholdStat:         false,
	beholdSessionsEnd:  true,
	beholdFeatures:     true,
	writeRelevantOnly:  true,
	sessionsPeriod:     5,
	statPeriod:         5,
	featuresPeriod:     5,
	maxInactiveMinutes: 60.0}

var sessions = new(sstrg.SessionsTable)
var emulatedTime time.Time
var rpsCounter = ratecounter.NewRateCounter(10 * time.Second)

var defaultDetector detectors.Detector
var defaultExtractor extractors.Extractor

func init() {
	sessions.H = make(map[string]*sstrg.SessionData)
}

//Configure configures app
func Configure(detector detectors.Detector, extractor extractors.Extractor) {
	defaultDetector = detector
	defaultExtractor = extractor
}

//Start features extractor
func Start() {
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

	w, err := os.Create(absPath)

	if err != nil {
		log.Fatalln("Error creating features file:", err)
	}

	if err != nil {
		log.Fatalln("Could not open dump file for writing:", err)
	}

	columnNames := []string{"ip", "user_agent"}
	featureNames := defaultExtractor.GetFeaturesNames()
	columnNames = append(columnNames, featureNames...)
	columnNames = append(columnNames, "label")

	if err := printCSV(w, columnNames); err != nil {
		log.Fatalln("Error writing header to csv:", err)
	}

	for {
		dumpFeatures(w, sessions)
		<-t.C
	}
}

func printCSV(w io.Writer, row []string) error {
	sep := ""
	for _, cell := range row {
		_, err := fmt.Fprintf(w, `%s"%s"`, sep, strings.Replace(cell, `"`, `""`, -1))
		if err != nil {
			return err
		}
		sep = ","
	}
	_, err := fmt.Fprintf(w, "\n")
	if err != nil {
		return err
	}

	return nil
}

func dumpFeatures(w io.Writer, sessions *sstrg.SessionsTable) {
	sessions.Lock()

	for key, s := range sessions.H {
		//We do not dump not-yet-finished sessions
		if s.Active {
			continue
		}

		//Append label
		var label = defaultDetector.GetLabel(s)

		if config.writeRelevantOnly && label == detectors.IrrelevantLabel {
			continue
		}

		line := strings.Split(key, "|")

		var fvDesc = defaultExtractor.ExtractFeatures(s)
		line = append(line, fvDesc...)
		line = append(line, strconv.Itoa(label))

		if err := printCSV(w, line); err != nil {
			log.Fatalln("Error writing record to csv:", err)
		}

		//Delete session from storage
		delete(sessions.H, key)
	}

	sessions.Unlock()
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
