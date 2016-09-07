package detectors

import "github.com/sauron/session"

const (
	//UnknownLabel label for session without label
	UnknownLabel int = iota
	//BotLabel label for bot
	BotLabel int = iota
	//HumanLabel label for human
	HumanLabel int = iota
)

// Detector is an interface implemented by all Detectors
type Detector interface {
	Init(configPath string)
	GetLabel(s *sstrg.SessionData) int
}
