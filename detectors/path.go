package detectors

import (
	"github.com/sauron/config"
	"github.com/sauron/session"
)

//PathDetector detects human by checking if client visited so called "human" paths
type PathDetector struct {
	paths map[string]bool
	label int
}

//SetLabel sets positive label for this detector
func (d *PathDetector) SetLabel(label int) {
	d.label = label
}

//Init initializes human path detector
func (d *PathDetector) Init(configPath string) {
	d.paths = make(map[string]bool)
	var paths = configutil.ReadPathsConfig(configPath)

	for _, path := range paths {
		d.paths[path] = true
	}
}

//GetLabel returns label for session by analyzing visited paths
func (d *PathDetector) GetLabel(s *sstrg.SessionData) int {
	for _, r := range s.Requests {
		if _, ok := d.paths[r.Path]; ok {
			return d.label
		}
	}

	return UnknownLabel
}
