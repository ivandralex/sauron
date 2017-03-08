package detectors

import (
	"fmt"
	"regexp"

	"github.com/sauron/config"
	"github.com/sauron/session"
)

//SessionKeyGetter interface for session key getter
type SessionKeyGetter interface {
	getKey(s *sstrg.RequestData) string
}

//ListDetector checks session by checking user agent
type ListDetector struct {
	expList   []*regexp.Regexp
	label     int
	keyGetter SessionKeyGetter
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

//GetLabel returns label for session by checking
func (d *ListDetector) GetLabel(s *sstrg.SessionData) int {
	for _, r := range s.Requests {
		key := d.keyGetter.getKey(r)

		for _, re := range d.expList {
			if re.MatchString(key) {
				fmt.Printf("list_detector3: %s\n", key)
				return d.label
			}
		}
	}

	return UnknownLabel
}
