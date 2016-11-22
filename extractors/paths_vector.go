package extractors

import (
	"math"
	"path/filepath"
	"strconv"
	"time"

	"github.com/sauron/config"
	"github.com/sauron/session"
)

//pathVector vector of features inherited from http path
type pathVector struct {
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
}

func (pv *pathVector) describe() []string {
	vector := []string{
		strconv.FormatInt(int64(pv.counter), 10)}
	//strconv.FormatFloat(pv.started, 'f', 2, 64),
	//strconv.FormatFloat(pv.averageDelay, 'f', 2, 64),
	//strconv.FormatFloat(pv.minDelay, 'f', 2, 64),
	//strconv.FormatFloat(pv.maxDelay, 'f', 2, 64)}

	return vector
}

//PathsVector feature representation of the session
type PathsVector struct {
	sessionDuration    float64
	sessionStartHour   int
	sessionStartMinute int
	clientTimeZone     int
	targetPaths        []string
}

//Init initializes extractor
func (fv *PathsVector) Init(configPath string) {
	absPath, _ := filepath.Abs(configPath)
	fv.targetPaths = configutil.ReadPathsConfig(absPath)
}

func (fv *PathsVector) describe(pathsFilter []string, pathVectors *map[string]*pathVector) []string {
	var finalVector []string
	for _, path := range pathsFilter {
		if pv, ok := (*pathVectors)[path]; ok {
			finalVector = append(finalVector, pv.describe()...)
		} else {
			//Add NaN vector if path was not visited
			finalVector = append(finalVector, "0" /*, "0", "0", "0", "0"*/)
		}
	}

	return finalVector
}

//ExtractFeatures extracts paths vector from session
func (fv *PathsVector) ExtractFeatures(s *sstrg.SessionData) []string {
	//sstrg.SortRequestsByTime(s.Requests)

	pathVectors := make(map[string]*pathVector)

	//Build path vectors map from requests
	for _, r := range s.Requests {
		//fmt.Fprintf(os.Stdout, "%v\n", r.Time)

		if _, pv := pathVectors[r.Path]; !pv {
			pathVectors[r.Path] = new(pathVector)
			pathVectors[r.Path].started = r.Time.Sub(s.Started).Seconds()
			pathVectors[r.Path].minDelay = math.MaxFloat64
		} else {
			//Delay after the last request with the same path
			var delay = r.Time.Sub(pathVectors[r.Path].last).Seconds()
			pathVectors[r.Path].delays = append(pathVectors[r.Path].delays, delay)
			pathVectors[r.Path].averageDelay += delay
			//Update max delay
			if delay > pathVectors[r.Path].maxDelay {
				pathVectors[r.Path].maxDelay = delay
			}
			//Update min delay
			if delay < pathVectors[r.Path].minDelay {
				pathVectors[r.Path].minDelay = delay
			}
		}

		pathVectors[r.Path].last = r.Time
		pathVectors[r.Path].counter++
	}

	for _, pathVector := range pathVectors {
		if pathVector.counter == 1 {
			pathVector.minDelay = 0
			pathVector.averageDelay = 0
		} else {
			pathVector.averageDelay /= float64(pathVector.counter - 1)
		}
	}

	return fv.describe(fv.targetPaths, &pathVectors)
}
