package detectors

import "github.com/sauron/session"

// Detector is an interface implemented by all Detectors
type Detector interface {
	Init(configPath string)
	GetLabel(s *sstrg.SessionData) string
}
