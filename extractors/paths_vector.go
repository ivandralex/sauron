package extractors

import (
	"math"
	"path/filepath"
	"strconv"
	"time"

	"github.com/sauron/config"
	"github.com/sauron/session"
)

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
}

func (pv *PathVector) describe() []string {
	vector := []string{
		strconv.FormatInt(int64(pv.counter), 10)}
	//strconv.FormatFloat(pv.started, 'f', 2, 64),
	//strconv.FormatFloat(pv.averageDelay, 'f', 2, 64),
	//strconv.FormatFloat(pv.minDelay, 'f', 2, 64),
	//strconv.FormatFloat(pv.maxDelay, 'f', 2, 64)}

	return vector
}

var targetPaths []string

//FeatureVector feature representation of the session
type FeatureVector struct {
	//Vectors corresponding to every unique path
	PathVectors        map[string]*PathVector
	sessionDuration    float64
	sessionStartHour   int
	sessionStartMinute int
	clientTimeZone     int
}

//Init initializes extractor
func (fv *FeatureVector) Init(configPath string) {
	absPath, _ := filepath.Abs(configPath)
	targetPaths = configutil.ReadPathsConfig(absPath)
}

func (fv *FeatureVector) describe(pathsFilter []string) []string {
	var finalVector []string
	for _, path := range pathsFilter {
		if pv, ok := fv.PathVectors[path]; ok {
			finalVector = append(finalVector, pv.describe()...)
		} else {
			//Add NaN vector if path was not visited
			finalVector = append(finalVector, "0" /*, "0", "0", "0", "0"*/)
		}
	}

	return finalVector
}

//ExtractFeatures extracts paths vector from session
func (fv *FeatureVector) ExtractFeatures(s *sstrg.SessionData) []string {
	//sstrg.SortRequestsByTime(s.Requests)

	//TODO: init it above
	fv.PathVectors = make(map[string]*PathVector)

	//Build path vectors map from requests
	for _, r := range s.Requests {
		//fmt.Fprintf(os.Stdout, "%v\n", r.Time)

		if _, pv := fv.PathVectors[r.Path]; !pv {
			fv.PathVectors[r.Path] = new(PathVector)
			fv.PathVectors[r.Path].started = r.Time.Sub(s.Started).Seconds()
			fv.PathVectors[r.Path].minDelay = math.MaxFloat64
		} else {
			//Delay after the last request with the same path
			var delay = r.Time.Sub(fv.PathVectors[r.Path].last).Seconds()
			fv.PathVectors[r.Path].delays = append(fv.PathVectors[r.Path].delays, delay)
			fv.PathVectors[r.Path].averageDelay += delay
			//Update max delay
			if delay > fv.PathVectors[r.Path].maxDelay {
				fv.PathVectors[r.Path].maxDelay = delay
			}
			//Update min delay
			if delay < fv.PathVectors[r.Path].minDelay {
				fv.PathVectors[r.Path].minDelay = delay
			}
		}

		fv.PathVectors[r.Path].last = r.Time
		fv.PathVectors[r.Path].counter++
	}

	for _, pathVector := range fv.PathVectors {
		if pathVector.counter == 1 {
			pathVector.minDelay = 0
			pathVector.averageDelay = 0
		} else {
			pathVector.averageDelay /= float64(pathVector.counter - 1)
		}
	}

	return fv.describe(targetPaths)
}
