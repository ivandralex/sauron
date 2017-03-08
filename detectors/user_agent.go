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

func (d *UserAgentDetector) getKey(r *sstrg.RequestData) string {
	return r.Header.Get("User-Agent")
}
