package detectors

import "github.com/sauron/session"

//PathDetector assigns specified label to session from enlisted ip
type PathDetector struct {
	ListDetector
}

//Init list detector
func (d *PathDetector) Init(configPath string) {
	d.ListDetector.Init(configPath)
	d.keyGetter = d
}

func (d *PathDetector) getKey(r *session.RequestData) string {
	return r.Path
}
