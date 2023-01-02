package apperror

import (
	"github.com/pkg/errors"
)

var ErrNotFound = errors.New("Not found")
var ErrTokenExpired = errors.New("Not found")
var ErrInternal = errors.New("Internal error")
