package pushmessage

// Message for push message
type Message struct {
	Token   string
	Title   string
	Body    string
	Data    map[string]string
	Android *AndroidConfig
	IOS     *IOSConfig
}

// AndroidConfig push message specific
// TODO: add more things that is available in the SDK
type AndroidConfig struct {
	Icon string
	// in hex format
	Color         string
	Sound         string
	Tag           string
	ClickAction   string
	BodyLockKey   string
	BodyLockArgs  string
	TitleLockKey  string
	TItleLockArgs string
}

// IOSConfig push message specific
// TODO: add more things that is available in the SDK
type IOSConfig struct {
	MutableContent   bool
	ContentAvailable bool
	Category         string
	ThreadID         string
	Badge            *int
	CustomData       map[string]interface{}
}
