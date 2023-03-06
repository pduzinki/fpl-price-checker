package storage

import (
	"errors"
)

var ErrDataAlreadyExists = errors.New("data already exists")
var ErrDataNotFound = errors.New("data not found")
