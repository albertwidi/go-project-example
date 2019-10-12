package booking

import "errors"

// list of booking errors
var (
	ErrInvalidType = errors.New("booking: type is not valid")
)
