package pathvector

import (
	"math"
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
	//Has referrer been requested?
	validRef bool
}

var targetPaths []string

func init() {
	targetPaths = configutil.ReadPathsConfig("../configs/target_paths.csv")
}

func (pv *PathVector) describe() []string {
	vector := []string{
		strconv.FormatInt(int64(pv.counter), 10),
		strconv.FormatFloat(pv.started, 'f', 2, 64),
		strconv.FormatFloat(pv.averageDelay, 'f', 2, 64),
		strconv.FormatFloat(pv.minDelay, 'f', 2, 64),
		strconv.FormatFloat(pv.maxDelay, 'f', 2, 64)}
	//strconv.FormatBool(pv.validRef)}

	return vector
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

func (fv *FeatureVector) describe(pathsFilter []string) []string {
	var finalVector []string
	for _, path := range pathsFilter {
		if pv, ok := fv.pathVectors[path]; ok {
			finalVector = append(finalVector, pv.describe()...)
		} else {
			//Add NaN vector if path was not visited
			finalVector = append(finalVector, "0", "0", "0", "0", "0")
		}
	}

	return finalVector
}

//ExtractFeatures extracts paths vector from session
func ExtractFeatures(s *sstrg.SessionData) []string {
	var fv = extractFeatureVector(s)

	return fv.describe(targetPaths)
}

//Extracts paths vector from session
func extractFeatureVector(s *sstrg.SessionData) *FeatureVector {
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
