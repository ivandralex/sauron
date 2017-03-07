package detectors

import (
	"fmt"
	"regexp"

	"github.com/sauron/config"
	"github.com/sauron/session"
)

//UserAgentDetector checks session by checking user agent
type UserAgentDetector struct {
	userAgentsExp []*regexp.Regexp
	label         int
}

//Init initializes human path detector
func (d *UserAgentDetector) Init(configPath string) {
	userAgents := configutil.ReadPathsConfig(configPath)
	d.userAgentsExp = make([]*regexp.Regexp, len(userAgents), len(userAgents))

	for i, expStr := range userAgents {
		d.userAgentsExp[i] = regexp.MustCompile(expStr)
	}
}

//SetLabel sets positive label for this detector
func (d *UserAgentDetector) SetLabel(label int) {
	d.label = label
}

//GetLabel returns label for session by checking
func (d *UserAgentDetector) GetLabel(s *sstrg.SessionData) int {
	userAgent := s.Requests[0].Header.Get("User-Agent")

	for _, re := range d.userAgentsExp {
		if re.MatchString(userAgent) {
			fmt.Printf("user_agent_detector2: %s\n", userAgent)
			return d.label
		}
	}

	return UnknownLabel
}
