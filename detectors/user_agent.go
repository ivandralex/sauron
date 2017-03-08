package detectors

import "github.com/sauron/session"

//UserAgentDetector checks session by checking user agent
type UserAgentDetector struct {
	ListDetector
}

//Init list detector
func (d *UserAgentDetector) Init(configPath string) {
	d.ListDetector.Init(configPath)
	d.keyGetter = d
}

func (d *UserAgentDetector) getKey(s *sstrg.SessionData) string {
	return s.Requests[0].Header.Get("User-Agent")
}
