package time

import "time"

var (
	_defaultLocation *time.Location
)

// SetDefaultLocation will set the default location of timeutil library to the defined location
func SetDefaultLocation(location string) error {
	var err error

	_defaultLocation, err = time.LoadLocation(location)
	return err
}

// Now return the current time based on _defaultLocation
func Now() time.Time {
	return time.Now().In(_defaultLocation)
}

// Time return the Go time struct
func Time() time.Time {
	return time.Time{}
}
