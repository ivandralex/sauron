package detectors

import (
	"github.com/sauron/config"
	"github.com/sauron/session"
)

//HumanPathDetector detects human by checking if client visited so called "human" paths
type HumanPathDetector struct {
	humanPaths map[string]bool
}

//Init initializes human path detector
func (d *HumanPathDetector) Init(configPath string) {
	d.humanPaths = make(map[string]bool)
	var paths = configutil.ReadPathsConfig(configPath)

	for _, path := range paths {
		d.humanPaths[path] = true
	}
}

//GetLabel returns label for session by analyzing visited "human" paths
func (d *HumanPathDetector) GetLabel(s *sstrg.SessionData) int {
	for _, r := range s.Requests {
		if _, ok := d.humanPaths[r.Path]; ok {
			return HumanLabel
		}
	}

	return UnknownLabel
}
