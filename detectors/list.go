package detectors

import (
	"regexp"

	"github.com/sauron/config"
	"github.com/sauron/session"
)

//RequestKeyGetter interface for session key getter
type RequestKeyGetter interface {
	getKey(s *session.RequestData) string
}

//ListDetector returns label for session by matching key against list of regexps
type ListDetector struct {
	expList   []*regexp.Regexp
	label     int
	keyGetter RequestKeyGetter
}

//Init list detector
func (d *ListDetector) Init(configPath string) {
	expStrs := configutil.ReadPathsConfig(configPath)
	d.expList = make([]*regexp.Regexp, len(expStrs), len(expStrs))

	for i, expStr := range expStrs {
		d.expList[i] = regexp.MustCompile(expStr)
	}
}

//SetLabel sets positive label for this detector
func (d *ListDetector) SetLabel(label int) {
	d.label = label
}

//GetLabel returns label for session by matching key against list of regexps
func (d *ListDetector) GetLabel(s *session.SessionData) int {
	for _, r := range s.Requests {
		key := d.keyGetter.getKey(r)

		for _, re := range d.expList {
			if re.MatchString(key) {
				return d.label
			}
		}
	}

	return UnknownLabel
}
