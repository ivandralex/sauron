package statutils

import (
	"fmt"
	"os"
	"time"

	"github.com/sauron/session"
)

//StartSessionsStatBeholder collects and outputs statistics on sessions
func StartSessionsStatBeholder(sessions *sstrg.SessionsTable, periodSec int) {
	//Init vars
	var allPaths map[string]int
	var maxCounter int
	var maxPath string
	var timeOfMax time.Time
	allPaths = make(map[string]int)

	// fire once per second
	t := time.NewTicker(time.Second * time.Duration(periodSec))

	for {
		//calcDurationStat(sessions)
		calcPathsStat(sessions, allPaths, &maxCounter, &maxPath, &timeOfMax)
		<-t.C
	}
}

func calcPathsStat(sessions *sstrg.SessionsTable, allPaths map[string]int, maxCounter *int, maxPath *string, timeOfMax *time.Time) {
	fmt.Fprintf(os.Stdout, "----\n")

	sessions.Lock()

	for key, s := range sessions.H {
		//We do not dump not-yet-finished sessions
		if s.Active {
			continue
		}

		for _, r := range s.Requests {
			if r.ContentType != sstrg.DynamicContentType {
				continue
			}

			allPaths[r.Path]++

			if allPaths[r.Path] > *maxCounter {
				*maxCounter = allPaths[r.Path]
				*maxPath = r.Path
				*timeOfMax = r.Time
			}

			//We definitely know that r.Referer is a page
			if _, ok := allPaths[r.Referer]; ok {
				delete(allPaths, r.Referer)
			}
		}

		//Delete session
		//TODO: do not delete it here
		delete(sessions.H, key)
	}

	sessions.Unlock()

	var totalRequests int
	var topRequests int

	for path, counter := range allPaths {
		totalRequests++
		if counter > *maxCounter/10000 {
			fmt.Fprintf(os.Stdout, "%s,%d\n", path, counter)
			topRequests++
		}
	}

	fmt.Fprintf(os.Stdout, "----\n%s: %d Top: %d/%d Time: %v\n", *maxPath, *maxCounter, topRequests, totalRequests, *timeOfMax)
}

func calcDurationStat(sessions *sstrg.SessionsTable) {
	fmt.Fprintf(os.Stdout, "Sessions durations:\n")

	var dur float64
	var durations []float64

	sessions.RLock()

	for _, s := range sessions.H {
		dur = s.Ended.Sub(s.Started).Minutes()

		durations = append(durations, dur)
	}

	sessions.RUnlock()

	var length = len(sessions.H)

	var percentiles = getPercentile(durations)

	for i := 50; i < 100; i++ {
		var percentile = float64(i) * 0.01
		fmt.Fprintf(os.Stdout, "%d, %.2f, %.2f\n", length, percentile, percentiles[percentile])
	}
}
