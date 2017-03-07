package sauron

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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
	beholdStat:         true,
	beholdSessionsEnd:  true,
	beholdFeatures:     false,
	writeRelevantOnly:  true,
	sessionsPeriod:     5,
	statPeriod:         5,
	featuresPeriod:     5,
	maxInactiveMinutes: 20.0}

var sessions = new(sstrg.SessionsTable)
var emulatedTime time.Time

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
			fmt.Fprintf(os.Stdout, "Closed session: %s cauze inactive for %f\n", k, dur)
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

	//TODO: use composite serializable key
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
