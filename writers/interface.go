package writers

//SessionDumpWriter is an interface for session dump writers used to persist sessions or their features
type SessionDumpWriter interface {
	Init(path string)
	WriteHead(featureNames []string)
	WriteSession(key string, features []string, label string)
}
