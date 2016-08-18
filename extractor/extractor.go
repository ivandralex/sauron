package extractor

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/paulbellamy/ratecounter"
	"github.com/sauron/session"
	"github.com/sauron/stat"
)

//import _ "net/http/pprof"

//TODO: move to file
var config = struct {
	useDataHeader  bool
	sessionsPeriod int
	featuresPeriod int
	statPeriod     int
	//Number requests to a path for it to be added to dump
	pathsThreshold     int
	maxInactiveMinutes float64
	//Feature flags
	beholdStat        bool
	beholdFeatures    bool
	beholdSessionsEnd bool
}{
	useDataHeader:      true,
	beholdStat:         false,
	beholdSessionsEnd:  true,
	beholdFeatures:     true,
	sessionsPeriod:     9,
	statPeriod:         10,
	featuresPeriod:     5,
	pathsThreshold:     1,
	maxInactiveMinutes: 15.0}

var sessions = new(sstrg.SessionsTable)
var rpsCounter = ratecounter.NewRateCounter(10 * time.Second)

//PathVector vector of features inherited from http path
type PathVector struct {
	//Delay of the first request (to this path) in the session
	started float64
	last    time.Time
	//Total number of requests
	counter int
	//Delays between consecutive calls to this path
	delays []float64
	//Average delay
	averageDelay float64
	//Maximum delay
	maxDelay float64
	//Minimum delay
	minDelay float64
	//Delays after previous request to different path of the same content type
	chainDelays []float64
	//Has referrer been requested?
	validRef bool
}

func init() {
	fmt.Println("Extractor initialized")
	sessions.H = make(map[string]*sstrg.SessionHistory)
}

//FeatureVector feature representation of the session
type FeatureVector struct {
	//Vectors corresponding to every unique path
	pathVectors        map[string]*PathVector
	sessionDuration    float64
	sessionStartHour   int
	sessionStartMinute int
	clientTimeZone     int
}

//Start features extractor
func Start() {
	//Check if sessions active periodically
	if config.beholdSessionsEnd {
		go startSessionsBeholder(config.sessionsPeriod)
	}
	//Dump features periodically
	if config.beholdFeatures {
		go startFeaturesBeholder(config.featuresPeriod)
	}
	//Collect stat on sessions
	if config.beholdStat {
		go statutils.StartSessionsStatBeholder(sessions, config.statPeriod)
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
	var now = time.Now().UTC()

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

func startFeaturesBeholder(periodSec int) {
	// fire once per second
	t := time.NewTicker(time.Second * time.Duration(periodSec))

	for {
		dumpFeatures()
		<-t.C
	}
}

func dumpFeatures() {
	fmt.Fprintf(os.Stdout, "Dump features!\n")

	sessions.Lock()

	var pathsCounter = 0
	var allPaths = make(map[string]int)

	for _, s := range sessions.H {
		//We do not dump not-yet-finished sessions
		if s.Active {
			continue
		}

		for _, r := range s.Requests {
			allPaths[r.Path]++
			if allPaths[r.Path] == 1 {
				pathsCounter++
			}
		}
	}

	for key, s := range sessions.H {
		//We do not dump not-yet-finished sessions
		if s.Active {
			continue
		}

		var fv = extractFeatures(s)

		for path := range fv.pathVectors {
			fmt.Printf("%s: %d\n", path, allPaths[path])
			//TODO: get top X requests that are of dynamicContentType and are not referers
			if allPaths[path] <= config.pathsThreshold {

			}
		}

		//var fvDesc = fv.describe(allPaths)
		//fmt.Printf("%s\n", fvDesc)
		//TODO: save description

		//Delete session from storage
		delete(sessions.H, key)
	}

	fmt.Printf("Total paths: %d\n", pathsCounter)

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

//FeaturesHandler ouputs feature vector for specified session key
func FeaturesHandler(w http.ResponseWriter, r *http.Request) {
	var sessionKey = r.URL.Query().Get("key")

	if _, ok := sessions.H[sessionKey]; !ok {
		fmt.Fprintf(w, "Not Found")
		return
	}

	sessions.Lock()

	var fv = extractFeatures(sessions.H[sessionKey])

	sessions.Unlock()

	for path, pathVector := range fv.pathVectors {
		fmt.Fprintf(w, "%s:\nstarted in %.4f sec\tdelays:%v (%d total)\tmin: %.4f\taverage: %.4f\tmax: %.4f\tvalidRef:%v\n\n", path, pathVector.started, pathVector.delays, pathVector.counter, pathVector.minDelay, pathVector.averageDelay, pathVector.maxDelay, pathVector.validRef)
	}
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
		sessions.H[sessionKey] = new(sstrg.SessionHistory)
		sessions.H[sessionKey].Started = request.Time
		sessions.H[sessionKey].Active = true
	}

	sessions.H[sessionKey].Ended = request.Time

	//Save request to session
	sessions.H[sessionKey].Requests = append(sessions.H[sessionKey].Requests, request)

	sessions.Unlock()
}

func extractFeatures(s *sstrg.SessionHistory) *FeatureVector {
	//sstrg.SortRequestsByTime(s.Requests)

	var fv = new(FeatureVector)
	//TODO: init it above
	fv.pathVectors = make(map[string]*PathVector)

	var validRef bool
	//Build path vectors map from requests
	for _, r := range s.Requests {
		//fmt.Fprintf(os.Stdout, "%v\n", r.Time)

		//Have referrer of this request been requested?
		if _, ok := fv.pathVectors[r.Referer]; ok {
			validRef = true
		}

		if _, pv := fv.pathVectors[r.Path]; !pv {
			fv.pathVectors[r.Path] = new(PathVector)
			fv.pathVectors[r.Path].started = r.Time.Sub(s.Started).Seconds()
			fv.pathVectors[r.Path].minDelay = math.MaxFloat64
			fv.pathVectors[r.Path].validRef = validRef
		} else {
			//Delay after the last request with the same path
			var delay = r.Time.Sub(fv.pathVectors[r.Path].last).Seconds()
			fv.pathVectors[r.Path].delays = append(fv.pathVectors[r.Path].delays, delay)
			fv.pathVectors[r.Path].averageDelay += delay
			//Update max delay
			if delay > fv.pathVectors[r.Path].maxDelay {
				fv.pathVectors[r.Path].maxDelay = delay
			}
			//Update min delay
			if delay < fv.pathVectors[r.Path].minDelay {
				fv.pathVectors[r.Path].minDelay = delay
			}
		}

		//If validRef is false for a single request in a session it's gonna be false for corresponding PathVector
		fv.pathVectors[r.Path].validRef = fv.pathVectors[r.Path].validRef && validRef
		fv.pathVectors[r.Path].last = r.Time
		fv.pathVectors[r.Path].counter++
	}

	for _, pathVector := range fv.pathVectors {
		if pathVector.counter == 1 {
			pathVector.minDelay = 0
			pathVector.averageDelay = 0
		} else {
			pathVector.averageDelay /= float64(pathVector.counter - 1)
		}
	}

	return fv
}
