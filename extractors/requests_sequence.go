package extractors

import (
	"path/filepath"
	"strconv"

	"github.com/sauron/config"
	"github.com/sauron/session"
)

//RequestsSequence feature representation of the session
type RequestsSequence struct {
	cardinality    int
	targetPathsMap map[string]int
	targetPaths    []string
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

//ExtractFeatures extracts paths vector from session
func (fv *RequestsSequence) ExtractFeatures(s *sstrg.SessionData) []string {
	var features = []string{}

	requestsCounter := 0

	//Build path vectors map from requests
	for _, r := range s.Requests {
		if _, pv := fv.targetPathsMap[r.Path]; pv {
			vector := fv.getEmptyPathVector()

			index := fv.targetPathsMap[r.Path]

			vector[index] = "1"
			requestsCounter++

			features = append(features, vector...)

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
		for _, path := range fv.targetPaths {
			head = append(head, strconv.FormatInt(int64(i), 10)+path)
		}
	}

	return head
}
