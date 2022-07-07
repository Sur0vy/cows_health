package errors

import "errors"

var ErrExist = errors.New("entry already exist")
var ErrEmpty = errors.New("entry is missing")
