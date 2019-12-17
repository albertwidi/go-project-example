package image

import "errors"

// list of errors
var (
	ErrTooManyTags            = errors.New("image: cannot have more than 5 tags")
	ErrTempPathNotFound       = errors.New("image: temporary path not found")
	ErrInvalidAccessAttribute = errors.New("image: invalid access attribute")
)
