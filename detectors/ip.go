package detectors

import "github.com/sauron/session"

//IPListDetector assigns specified label to session from enlisted ip
type IPListDetector struct {
	ListDetector
}

//Init list detector
func (d *IPListDetector) Init(configPath string) {
	d.ListDetector.Init(configPath)
	d.keyGetter = d
}

func (d *IPListDetector) getKey(r *session.RequestData) string {
	return r.IP
}
