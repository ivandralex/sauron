package sauron

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/sauron/detectors"
	"github.com/sauron/extractors"
	"github.com/sauron/session"
	"github.com/sauron/stat"
	"github.com/sauron/writers"
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
	beholdStat:         true,
	beholdSessionsEnd:  true,
	beholdFeatures:     true,
	writeRelevantOnly:  true,
	sessionsPeriod:     5,
	statPeriod:         5,
	featuresPeriod:     5,
	maxInactiveMinutes: 20.0}

var sessions = new(session.SessionsTable)
var emulatedTime time.Time

var defaultDetector detectors.Detector
var defaultExtractor extractors.Extractor
var defaultWriter writers.SessionDumpWriter

func init() {
	sessions.H = make(map[string]*session.SessionData)
}

//Configure configures app
func Configure(detector detectors.Detector, extractor extractors.Extractor, writer writers.SessionDumpWriter) {
	defaultDetector = detector
	defaultExtractor = extractor
	defaultWriter = writer
}

//Start features extractor
func Start() {
	//Check if sessions active periodically
	if config.beholdSessionsEnd {
		go startSessionsBeholder(config.sessionsPeriod)
	}
	//Collect stat on sessions
	if config.beholdStat {
		go statutils.StartSessionsStatBeholder(config.statPeriod, sessions, &defaultDetector)
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
	fmt.Fprintf(os.Stdout, "Closing sessions!\n")

	var dur float64
	var now time.Time
	if config.emulateTime {
		now = emulatedTime
	} else {
		now = time.Now().UTC()
	}

	sessions.RLock()

	for k, s := range sessions.H {
		if !s.Active {
			continue
		}

		dur = now.Sub(s.Ended).Minutes()

		if dur > config.maxInactiveMinutes {
			fmt.Fprintf(os.Stdout, "Closed session: %s\n", k)
			//Mark session as inactive. It will be deleted after the next dump
			s.Active = false
		}
	}

	sessions.RUnlock()
}

func startFeaturesBeholder(sessions *session.SessionsTable, periodSec int) {
	// fire once per second
	t := time.NewTicker(time.Second * time.Duration(periodSec))

	featureNames := defaultExtractor.GetFeaturesNames()

	defaultWriter.WriteHead(featureNames)

	for {
		dumpFeatures(sessions)
		<-t.C
	}
}

func dumpFeatures(sessions *session.SessionsTable) {
	sessions.Lock()

	for key, s := range sessions.H {
		//We do not dump not-yet-finished sessions
		if s.Active {
			continue
		}

		var label = defaultDetector.GetLabel(s)

		if !config.writeRelevantOnly || label != detectors.IrrelevantLabel {

			var fvDesc = defaultExtractor.ExtractFeatures(s)

			defaultWriter.WriteSession(key, fvDesc, strconv.Itoa(label))
		}

		//TODO: use listeners counter for session
		//Delete session from storage
		delete(sessions.H, key)
	}

	sessions.Unlock()
}
