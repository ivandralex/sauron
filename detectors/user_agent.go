package detectors

import (
	"fmt"

	"github.com/sauron/config"
	"github.com/sauron/session"
)

//UserAgentDetector checks session by checking user agent
type UserAgentDetector struct {
	userAgents map[string]bool
	label      int
}

//Init initializes human path detector
func (d *UserAgentDetector) Init(configPath string) {
	d.userAgents = configutil.ReadStringMap(configPath)
}

//SetLabel sets positive label for this detector
func (d *UserAgentDetector) SetLabel(label int) {
	d.label = label
}

//GetLabel returns label for session by checking
func (d *UserAgentDetector) GetLabel(s *sstrg.SessionData) int {
	userAgent := s.Requests[0].Header.Get("User-Agent")
	if _, ok := d.userAgents[userAgent]; ok {
		fmt.Printf("user_agent_detector: %s\n", userAgent)
		return d.label
	}

	return UnknownLabel
}
