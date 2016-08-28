package detectors

import (
	"github.com/sauron/config"
	"github.com/sauron/session"
)

//TODO: it's not a very good idea to use global var here
var humanPaths map[string]bool

//TODO: check guidelines if it's a good idea to read config here rather than get it as argument in GetLabel
func init() {
	humanPaths = make(map[string]bool)
	var paths = configutil.ReadPathsConfig("../configs/human_paths.csv")

	for _, path := range paths {
		humanPaths[path] = true
	}
}

//GetLabel returns label for session by analyzing visited "human" paths
func GetLabel(s *sstrg.SessionData) string {
	for _, r := range s.Requests {
		if _, ok := humanPaths[r.Path]; ok {
			return "0"
		}
	}

	return "1"
}
