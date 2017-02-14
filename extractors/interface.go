package extractors

import "github.com/sauron/session"

//Extractor is an interface implemented by all Extractors
type Extractor interface {
	Init(configPath string)
	ExtractFeatures(s *sstrg.SessionData) []string
	GetFeaturesNames() []string
}
