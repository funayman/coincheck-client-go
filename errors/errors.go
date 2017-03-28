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

type GenericError struct {
	error
}

func NewGenericError(msg string, a ...interface{}) GenericError {
	err := fmt.Sprintf("EndPointError: "+msg, a...)
	return GenericError{error: errors.New(err)}
}
