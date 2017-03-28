package errors

import (
	"errors"
	"fmt"
)

type EndPointError struct {
	error
}

func NewEndPointError(msg string, a ...interface{}) EndPointError {
	err := fmt.Sprintf("EndPointError: "+msg, a...)
	return EndPointError{error: errors.New(err)}
}
