package device

// Type of device
type Type int

// type of device
const (
	TypeAndroid Type = 1
	TypeIOS     Type = 2

	TypeAndroidString = "android"
	TypeIOSString     = "ios"
)

// Device struct
type Device struct {
	// FCMToken for firebase cloud messaging
	// changed when applicaion uninstalled
	FCMToken string
	// ACMToken for apple cloud messaging
	// changed when application uninstalled
	ACMTOken string
	// ID of device
	// only get changed when phone got formatted
	DeviceID string
}
