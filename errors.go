package cfg

import (
	"errors"
)

var ErrInvalidSpecification = errors.New("specification must be a struct pointer")
var ErrInvalidArgPointerRequired = errors.New("cfg requires struct pointers")
var ErrInvalidArgMapParseSpec = errors.New("map argument requires pairs")
var ErrIgnoreTag = errors.New("this Tag isn't in use")
