package featureflag

var _ff = FeatureFlag{}

// FeatureFlag struct
type FeatureFlag struct {
	backend Backend
}

// Backend of feature flag
type Backend interface {
	GetString(key string) (string, error)
	GetBoolean(key string) (bool, error)
	GetInt(key string) (int, error)
}

// SetBackend for feature flag
// so it can be used globally
func SetBackend(backend Backend) {
	_ff.backend = backend
	_ff.work()
}

// StopUpdate the feature flag
// means the checking to backend will be stopped
func StopUpdate() {
	_ff.stop()
}

// regularly check
func (ff *FeatureFlag) work() {

}

func (ff *FeatureFlag) stop() {

}
