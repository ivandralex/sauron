package extractors

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/sauron/config"
	"github.com/sauron/session"
)

//RequestsSequence feature representation of the session
type RequestsSequence struct {
	cardinality     int
	targetPathsMap  map[string]int
	targetPaths     []string
	temporalEnabled bool
}

//Init initializes extractor
func (fv *RequestsSequence) Init(configPath string) {
	absPath, _ := filepath.Abs(configPath)
	fv.targetPaths = configutil.ReadPathsConfig(absPath)

	fv.targetPathsMap = make(map[string]int)

	for index, path := range fv.targetPaths {
		fv.targetPathsMap[path] = index
	}
}

func (fv *RequestsSequence) getEmptyPathVector() []string {
	length := len(fv.targetPaths)

	if fv.temporalEnabled {
		length++
	}

	vector := make([]string, length, length)

	for i := range vector {
		vector[i] = "0"
	}

	return vector
}

//SetCardinality cardinality
func (fv *RequestsSequence) SetCardinality(cardinality int) {
	fv.cardinality = cardinality
}

//SetTemporalEnabled disables or enables temporal features
func (fv *RequestsSequence) SetTemporalEnabled(enabled bool) {
	fv.temporalEnabled = enabled
}

//ExtractFeatures extracts paths vector from session
func (fv *RequestsSequence) ExtractFeatures(s *sstrg.SessionData) []string {
	var features = []string{}
	var pathTimes = make(map[string]time.Time)

	requestsCounter := 0

	startTime := s.Requests[0].Time

	//Build path vectors map from requests
	for _, r := range s.Requests {
		if _, pv := fv.targetPathsMap[r.Path]; pv {
			vector := fv.getEmptyPathVector()

			//Setup one-hot vector for this request
			index := fv.targetPathsMap[r.Path]
			vector[index] = "1"

			if fv.temporalEnabled {
				//Init time for this request
				if _, ok := pathTimes[r.Path]; !ok {
					pathTimes[r.Path] = startTime
				}

				//Save delay of this request from previous request for the same path
				//(or from the beginning of the session if it's the first request for this path)
				requestDelay := r.Time.Sub(pathTimes[r.Path]).Seconds()

				if requestDelay > 100000 {
					fmt.Printf("Too long session of %f with key %s\n", requestDelay, s.IP)
				}
				vector[len(vector)-1] = strconv.FormatFloat(requestDelay, 'f', 1, 64)
				pathTimes[r.Path] = r.Time
			}

			features = append(features, vector...)

			requestsCounter++
			if requestsCounter == fv.cardinality {
				break
			}
		}
	}

	if requestsCounter != fv.cardinality {
		for requestsCounter < fv.cardinality {
			vector := fv.getEmptyPathVector()
			features = append(features, vector...)
			requestsCounter++
		}
	}

	return features
}

//GetFeaturesNames array of features names
func (fv *RequestsSequence) GetFeaturesNames() []string {
	head := []string{}

	for i := 0; i < fv.cardinality; i++ {
		index := strconv.FormatInt(int64(i), 10)
		for _, path := range fv.targetPaths {
			head = append(head, index+path)
		}
		if fv.temporalEnabled {
			head = append(head, index+"/delay")
		}
	}

	return head
}
