package detectors

import (
	"fmt"

	"github.com/sauron/config"
	"github.com/sauron/session"
)

//TODO: it's not a very good idea to use global var here
var humanPaths map[string]bool

//TODO: check guidelines if it's a good idea to read config here rather then get it as argument in GetLabel
func init() {
	humanPaths = make(map[string]bool)
	var paths = configutil.ReadPathsConfig("../configs/human_paths.csv")

	for _, path := range paths {
		humanPaths[path] = true
	}
}

//GetNaiveHumanLabel returns label for session by analyzing visited "human" paths
func GetNaiveHumanLabel(s *sstrg.SessionData) string {
	for _, r := range s.Requests {
		if _, ok := humanPaths[r.Path]; ok {
			fmt.Printf("Found human %s\n", s.IP)
			return "2"
		}
	}

	return "0"
}
